package lister

import (
	"context"
	"github.com/sirupsen/logrus"
	conn "github.com/vcsomor/aws-resources/internal/aws_connector"
	"github.com/vcsomor/aws-resources/internal/executor"
	"slices"
)

type Lister interface {
	List(ctx context.Context) []any
}

type taskBasedLister struct {
	clientFactory conn.ClientFactory
	executor      executor.SynchronousExecutor
	logger        *logrus.Logger

	regions   []string
	resources []string
}

var _ Lister = (*taskBasedLister)(nil)

func (l *taskBasedLister) List(ctx context.Context) []any {
	var res []any

	if slices.Contains(l.resources, "s3") {
		res = append(res, l.listS3(ctx)...)
	}
	if slices.Contains(l.resources, "rds") {
		res = append(res, l.listRDS(ctx)...)
	}

	l.logger.WithField(logKeyResourceCount, len(res)).
		Debug("resources listed")

	return res
}
