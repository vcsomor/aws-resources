package s3_tasks

import (
	"context"
	"github.com/sirupsen/logrus"
	conn "github.com/vcsomor/aws-resources/internal/aws_connector"
	"github.com/vcsomor/aws-resources/internal/executor"
	"time"
)

type listTask struct {
	ctx    context.Context
	logger *logrus.Entry
	client conn.S3Client
}

type ListTaskBucketData struct {
	Name    string
	Created *time.Time
}

type ListTaskResult struct {
	Buckets []ListTaskBucketData
	Error   error
}

var _ executor.Task = (*listTask)(nil)

func NewListTask(
	ctx context.Context,
	logger *logrus.Entry,
	client conn.S3Client,
) executor.Task {
	return &listTask{
		ctx:    ctx,
		logger: logger,
		client: client,
	}
}

func (t *listTask) Execute() any {
	buckets, err := t.client.List(t.ctx, conn.ListS3Params{})
	if err != nil {
		t.logger.WithError(err).
			Error("unable to list buckets")
		return ListTaskResult{Error: err}
	}

	var results []ListTaskBucketData

	for _, bucket := range buckets {
		results = append(results, ListTaskBucketData{
			Name:    bucket.Name,
			Created: bucket.CreationTime,
		})
	}

	t.logger.Debugf("resouces listed")

	return ListTaskResult{Buckets: results}
}
