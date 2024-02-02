package lister

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vcsomor/aws-resources/config"
	conn "github.com/vcsomor/aws-resources/internal/aws_connector"
	"github.com/vcsomor/aws-resources/internal/threads"
	"github.com/vcsomor/aws-resources/log"
	"os"
	"strconv"
	"time"
)

type Result struct {
	Arn          string     `json:"arn"`
	Region       string     `json:"region"`
	ID           string     `json:"id"`
	CreationTime *time.Time `json:"creationTime"`
}

type Lister interface {
	List(ctx context.Context) []Result
}

type defaultLister struct {
	clientFactory conn.ClientFactory
	threadpool    threads.Threadpool
	logger        *logrus.Logger
	regions       []string
}

var _ Lister = (*defaultLister)(nil)

// CmdListResources is the command entry point
func CmdListResources(command *cobra.Command, _ []string) {
	logger := log.NewLogger(config.Config())

	th := command.Flag("threads").
		Value.
		String()

	threadCount, err := strconv.Atoi(th)
	if err != nil {
		logger.WithError(err).
			Error("invalid thread count")
		return
	}

	regions := parseRegions(command.Flag("regions").
		Value.
		String())

	logger.Debugf("regions: %v", regions)

	threadpool, err := threads.NewThreadpool(threadCount)
	if err != nil {
		logger.WithError(err).
			Error("threadpool error")
		return
	}

	defer func() {
		if threadpool != nil {
			threadpool.Shutdown()
		}
	}()

	l := NewDefaultLister(
		logger,
		conn.NewClientFactory(logger),
		threadpool,
		regions,
	)
	resources := l.List(context.TODO())

	logger.WithField(logKeyResourceCount, len(resources)).
		Debug("resources listed")

	// TODO vcsomor do the write
	writeResult(resources)
}

func NewDefaultLister(
	logger *logrus.Logger,
	clientFactory conn.ClientFactory,
	threadpoolFactory threads.Threadpool,
	regions []string,
) Lister {
	return &defaultLister{
		clientFactory: clientFactory,
		threadpool:    threadpoolFactory,
		logger:        logger,
		regions:       regions,
	}
}

func (l *defaultLister) List(ctx context.Context) []Result {
	res := []Result{}

	res = append(res, l.listS3(ctx, l.regions)...)
	res = append(res, l.listRDS(ctx, l.regions)...)

	return res
}

func (l *defaultLister) listS3(ctx context.Context, regions []string) (results []Result) {
	logger := l.logger.WithField(logKeyResourceType, s3ResourceType)

	client, err := l.clientFactory.S3Client(ctx, nil)
	if err != nil {
		logger.WithError(err).
			Error("unable to create the client")
		return
	}

	listResult, err := l.threadpool.SubmitTask(newS3ListTask(ctx, logger, client))
	if err != nil {
		logger.WithError(err).
			Error("unable to start listing task")
	}

	s3Buckets, ok := listResult.GetWait().(s3ListTaskResult)
	if !ok {
		logger.Error("unable to get s3 list result, invalid type")
		return
	}

	if err = s3Buckets.error; err != nil {
		logger.WithError(err).
			Error("bucket fetch error")
		return
	}

	var tasks []threads.Future
	for _, b := range s3Buckets.buckets {
		f, errSubmit := l.threadpool.SubmitTask(newS3GetRegionTask(ctx, logger, client, b.name))
		if errSubmit != nil {
			logger.WithError(errSubmit).
				WithField("bucket-name", b.name).
				Error("unable to submit task")
			continue
		}
		tasks = append(tasks, f)
	}

	for _, t := range tasks {
		getResult, regionOK := t.GetWait().(s3GetRegionResult)
		if !regionOK {
			logger.Error("unable to assert return value")
			continue
		}

		if err = getResult.error; err != nil {
			logger.WithError(err).
				Error("error while fetching the region")
			continue
		}

		region := getResult.region
		if regionFiler(region, regions) {
			results = append(results, aResultOfBucketOperations(getResult.bucketName, s3Buckets.buckets, region))
		}
	}
	return
}

func (l *defaultLister) listRDS(ctx context.Context, regions []string) (results []Result) {
	logger := l.logger.
		WithField(logKeyResourceType, rdsResourceType)

	var tasks []threads.Future
	if regions == nil {
		f, err := l.startListRdsInRegion(ctx, logger.WithField(logKeyRegion, "default"), nil)
		if err != nil {
			logger.WithError(err).
				Error("error while fetching the rds instances")
		} else {
			tasks = append(tasks, f)
		}
	} else {
		for _, region := range regions {
			currRegion := region
			f, err := l.startListRdsInRegion(ctx, logger.WithField(logKeyRegion, region), &currRegion)
			if err != nil {
				logger.WithError(err).
					Error("error while fetching the rds instances")
				continue
			}
			tasks = append(tasks, f)
		}
	}

	for _, t := range tasks {
		taskResult, ok := t.GetWait().(rdsTaskResult)
		if !ok {
			logger.Error("unable to assert return value")
			continue
		}

		if err := taskResult.error; err != nil {
			logger.WithError(err).
				Error("error while fetching the RDS instances")
			continue
		}

		for _, rds := range taskResult.results {
			results = append(results, Result{
				Arn:          rds.Arn,
				Region:       rds.Region,
				ID:           rds.ID,
				CreationTime: rds.CreationTime,
			})
		}
	}
	return
}

func (l *defaultLister) startListRdsInRegion(ctx context.Context, logger *logrus.Entry, region *string) (threads.Future, error) {
	client, err := l.clientFactory.RDSClient(ctx, region)
	if err != nil {
		logger.WithError(err).
			Error("unable to create the client")
		return nil, err
	}

	return l.threadpool.SubmitTask(newRDSTask(ctx, logger, client))
}

func writeResult(res []Result) {
	js, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to display data %s", err)
		return
	}
	fmt.Printf("%s\n", js)
}

func aResultOfBucketOperations(name string, buckets []s3ListTaskResultItem, region string) Result {
	for _, b := range buckets {
		if name == b.name {
			return Result{
				Arn:          fmt.Sprintf("arn:aws:s3:::%s", name),
				Region:       region,
				ID:           name,
				CreationTime: b.created,
			}
		}
	}
	return Result{}
}

func regionFiler(region string, regions []string) bool {
	if len(regions) == 0 {
		return true
	}
	for _, r := range regions {
		if region == r {
			return true
		}
	}
	return false
}
