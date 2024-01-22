package lister

import (
	"context"
	"github.com/sirupsen/logrus"
	conn "github.com/vcsomor/aws-resources/internal/aws_connector"
	"github.com/vcsomor/aws-resources/internal/threads"
)

type s3Task struct {
	ctx    context.Context
	logger *logrus.Entry
	client conn.S3Client
}

type s3TaskResult struct {
	results []Result
	error   error
}

var _ threads.Task = (*s3Task)(nil)

func newS3Task(
	ctx context.Context,
	logger *logrus.Entry,
	client conn.S3Client,
) threads.Task {
	return &s3Task{
		ctx:    ctx,
		logger: logger,
		client: client,
	}
}

func (t *s3Task) Execute() any {
	resources, err := t.client.ListS3(t.ctx, conn.ListS3Params{})
	if err != nil {
		t.logger.WithError(err).
			Error("unable to list resources")
		return s3TaskResult{error: err}
	}

	var results []Result

	for _, res := range resources {
		results = append(results, Result{
			Arn:          res.Arn,
			ID:           res.Name,
			CreationTime: res.CreationTime,
		})
	}

	t.logger.WithField(logKeyResourceCount, len(resources)).
		Debugf("resouces listed")

	return s3TaskResult{results: results}
}
