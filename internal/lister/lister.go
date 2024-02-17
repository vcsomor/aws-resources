package lister

import (
	"context"
	"github.com/sirupsen/logrus"
	conn "github.com/vcsomor/aws-resources/internal/aws_connector"
	"github.com/vcsomor/aws-resources/internal/executor"
	"slices"
)

type Lister interface {
	List(ctx context.Context) []Result
}

type taskBasedLister struct {
	clientFactory                 conn.ClientFactory
	executor                      executor.SynchronousExecutor
	logger                        *logrus.Logger
	individualResultWriterFactory IndividualResultWriterFactory
	summarizedResultWriterFactory SummarizedResultWriterFactory

	regions   []string
	resources []string
}

var _ Lister = (*taskBasedLister)(nil)

func (l *taskBasedLister) List(ctx context.Context) (res []Result) {
	if slices.Contains(l.resources, "s3") {
		res = append(res, l.listS3(ctx)...)
	}
	if slices.Contains(l.resources, "rds") {
		res = append(res, l.listRDS(ctx)...)
	}

	l.logger.WithField(logKeyResourceCount, len(res)).
		Debug("resources listed")

	// TODO vcsomor this is temporary here
	if wf := l.individualResultWriterFactory; wf != nil {
		for _, r := range res {
			if err := wf(r).Write(r); err != nil {
				l.logger.WithError(err).
					Errorf("error withing the result for %s", r.Arn)
			}
		}
	}

	if wf := l.summarizedResultWriterFactory; wf != nil {
		if err := wf(res).Write(res); err != nil {
			l.logger.WithError(err).
				Errorf("error withing the results")
		}
	}

	return res
}
