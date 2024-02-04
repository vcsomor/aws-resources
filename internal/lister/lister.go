package lister

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vcsomor/aws-resources/config"
	conn "github.com/vcsomor/aws-resources/internal/aws_connector"
	"github.com/vcsomor/aws-resources/internal/executor"
	"github.com/vcsomor/aws-resources/internal/lister/rds_tasks"
	"github.com/vcsomor/aws-resources/internal/lister/s3_tasks"
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
	executor      executor.SynchronousExecutor
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

	threadpool, err := executor.NewThreadpool(threadCount)
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
		executor.NewSynchronousExecutor(threadpool),
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
	executor executor.SynchronousExecutor,
	regions []string,
) Lister {
	return &defaultLister{
		clientFactory: clientFactory,
		executor:      executor,
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

	listResult, err := l.executor.Execute(s3_tasks.NewListTask(ctx, logger, client))
	if err != nil {
		logger.WithError(err).
			Error("unable to start listing task")
		return
	}

	s3Buckets := listResult.(s3_tasks.ListTaskResult)
	if err = s3Buckets.Error; err != nil {
		logger.WithError(err).
			Error("bucket fetch error")
		return
	}

	var getRegionTasks []executor.Task
	for _, b := range s3Buckets.Buckets {
		getRegionTasks = append(getRegionTasks, s3_tasks.NewS3GetRegionTask(ctx, logger, client, b.Name))
	}

	getRegionResults := l.executor.ExecuteAll(getRegionTasks)

	for _, r := range getRegionResults {
		if err = r.Error; err != nil {
			logger.WithError(err).
				Error("error while fetching the region")
			continue
		}

		getRegionResult := r.Outcome.(s3_tasks.GetRegionResult)
		if err = getRegionResult.Error; err != nil {
			logger.WithError(err).
				Error("region fetch error")
			continue
		}

		if regionFiler(getRegionResult.Region, regions) {
			results = append(results, aResultOfBucketOperations(getRegionResult, s3Buckets.Buckets))
		}
	}
	return
}

func (l *defaultLister) listRDS(ctx context.Context, regions []string) (results []Result) {
	logger := l.logger.
		WithField(logKeyResourceType, rdsResourceType)

	var tasks []executor.Task
	if regions == nil {
		t, err := l.rdsTaskInRegion(ctx, logger.WithField(logKeyRegion, "default"), nil)
		if err != nil {
			logger.WithError(err).
				Error("unable to add RDS list task in default region")
		} else {
			tasks = append(tasks, t)
		}
	} else {
		for _, region := range regions {
			currRegion := region
			t, err := l.rdsTaskInRegion(ctx, logger.WithField(logKeyRegion, region), &currRegion)
			if err != nil {
				logger.WithError(err).
					Errorf("unable to add RDS list task in region")
				continue
			}
			tasks = append(tasks, t)
		}
	}

	execResults := l.executor.ExecuteAll(tasks)

	for _, r := range execResults {
		if err := r.Error; err != nil {
			logger.WithError(err).
				Error("error while fetching the RDS instances")
			continue
		}

		listResult := r.Outcome.(rds_tasks.ListResult)
		if err := listResult.Error; err != nil {
			logger.WithError(err).
				Error("rds list task error")
			continue
		}

		for _, rds := range listResult.RDSInstances {
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

func (l *defaultLister) rdsTaskInRegion(ctx context.Context, logger *logrus.Entry, region *string) (executor.Task, error) {
	client, err := l.clientFactory.RDSClient(ctx, region)
	if err != nil {
		logger.WithError(err).
			Error("unable to create the client")
		return nil, err
	}

	return rds_tasks.NewListTask(ctx, logger, client), nil
}

func writeResult(res []Result) {
	js, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to display data %s", err)
		return
	}
	fmt.Printf("%s\n", js)
}

func aResultOfBucketOperations(res s3_tasks.GetRegionResult, buckets []s3_tasks.ListTaskBucketData) Result {
	for _, b := range buckets {
		name := res.BucketName
		if name == b.Name {
			return Result{
				Arn:          fmt.Sprintf("arn:aws:s3:::%s", name),
				Region:       res.Region,
				ID:           name,
				CreationTime: b.Created,
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
