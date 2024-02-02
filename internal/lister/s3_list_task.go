package lister

import (
	"context"
	"github.com/sirupsen/logrus"
	conn "github.com/vcsomor/aws-resources/internal/aws_connector"
	"github.com/vcsomor/aws-resources/internal/threads"
	"time"
)

type s3ListTask struct {
	ctx    context.Context
	logger *logrus.Entry
	client conn.S3Client
}

type s3ListTaskResultItem struct {
	name    string
	created *time.Time
}

type s3ListTaskResult struct {
	buckets []s3ListTaskResultItem
	error   error
}

var _ threads.Task = (*s3ListTask)(nil)

func newS3ListTask(
	ctx context.Context,
	logger *logrus.Entry,
	client conn.S3Client,
) threads.Task {
	return &s3ListTask{
		ctx:    ctx,
		logger: logger,
		client: client,
	}
}

func (t *s3ListTask) Execute() any {
	buckets, err := t.client.List(t.ctx, conn.ListS3Params{})
	if err != nil {
		t.logger.WithError(err).
			Error("unable to list buckets")
		return s3ListTaskResult{error: err}
	}

	var results []s3ListTaskResultItem

	for _, bucket := range buckets {
		results = append(results, s3ListTaskResultItem{
			name:    bucket.Name,
			created: bucket.CreationTime,
		})
	}

	t.logger.WithField(logKeyResourceCount, len(buckets)).
		Debugf("resouces listed")

	return s3ListTaskResult{buckets: results}
}
