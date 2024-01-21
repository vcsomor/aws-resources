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
	"sync"
	"time"
)

const (
	logKeyRegion        = "region"
	logKeyResourceType  = "resource-type"
	logKeyResourceCount = "resource-count"

	s3ResourceType  = "S3"
	rdsResourceType = "RDS"
)

type Result struct {
	Arn          string     `json:"arn"`
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

	taskMtx     sync.Mutex
	taskFutures []*threads.TaskFuture
	results     [][]Result
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

	threadpool := threads.NewThreadpool(threadCount)

	l := NewDefaultLister(
		logger,
		conn.NewClientFactory(logger),
		threadpool,
		regions,
	)
	resources := l.List(context.TODO())

	threadpool.Shutdown()

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
	l.startListS3(ctx)
	l.startListRDSInRegions(ctx, l.regions)

	return l.fetchResults()
}

func (l *defaultLister) startListS3(ctx context.Context) {
	logger := l.logger.WithField(logKeyResourceType, s3ResourceType)

	client, err := l.clientFactory.S3Client(ctx)
	if err != nil {
		logger.WithError(err).
			Error("unable to create the client")
		return
	}

	if err = l.startTask(newS3Task(ctx, logger, client)); err != nil {
		logger.WithError(err).
			Error("unable to start task")
	}
}

func (l *defaultLister) startListRDSInRegions(ctx context.Context, regions []string) {
	if regions == nil {
		l.startListRdsInRegion(ctx, nil)
	}

	for _, region := range regions {
		currRegion := region
		l.startListRdsInRegion(ctx, &currRegion)
	}
}

func (l *defaultLister) startListRdsInRegion(ctx context.Context, region *string) {
	logger := l.logger.
		WithField(logKeyResourceType, rdsResourceType).
		WithField(logKeyRegion, derefRegion(region))

	client, err := l.clientFactory.RDSClient(ctx, region)
	if err != nil {
		logger.WithError(err).
			Error("unable to create the client")
		return
	}

	if err = l.startTask(newRDSTask(ctx, logger, client)); err != nil {
		logger.WithError(err).
			Error("unable to start task")
	}
}

func (l *defaultLister) startTask(task threads.Task) error {
	l.taskMtx.Lock()
	defer l.taskMtx.Unlock()
	future, err := l.threadpool.SubmitTask(task)
	if err != nil {
		return err
	}
	l.taskFutures = append(l.taskFutures, &future)
	return nil
}

func (l *defaultLister) fetchResults() (results []Result) {
	for _, tf := range l.taskFutures {
		f := (*tf)
		rawResult := f.Get()
		if taskResult, ok := rawResult.(s3TaskResult); ok {
			if err := taskResult.error; err == nil {
				results = append(results, taskResult.results...)
			} else {
				l.logger.WithError(err).
					Error("not appending results")
			}
			continue
		}

		if taskResult, ok := rawResult.(rdsTaskResult); ok {
			if err := taskResult.error; err == nil {
				results = append(results, taskResult.results...)
			} else {
				l.logger.WithError(err).
					Error("not appending results")
			}
			continue
		}

		l.logger.Warnf("unhandled type %t", rawResult)
	}

	return
}

func writeResult(res []Result) {
	js, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to display data %s", err)
		return
	}
	fmt.Printf("%s\n", js)
}

func derefRegion(r *string) string {
	if r != nil {
		return *r
	}
	return "default"
}
