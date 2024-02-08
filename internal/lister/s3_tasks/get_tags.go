package s3_tasks

import (
	"context"
	"github.com/sirupsen/logrus"
	conn "github.com/vcsomor/aws-resources/internal/aws_connector"
	"github.com/vcsomor/aws-resources/internal/executor"
)

type getTagsTask struct {
	ctx        context.Context
	logger     *logrus.Entry
	client     conn.S3Client
	bucketName string
}

type GetTagsResult struct {
	BucketName string
	Tags       map[string]*string
	Error      error
}

var _ executor.Task = (*getTagsTask)(nil)

func NewS3GetTagsTask(
	ctx context.Context,
	logger *logrus.Entry,
	client conn.S3Client,
	bucketName string,
) executor.Task {
	return &getTagsTask{
		ctx:        ctx,
		logger:     logger,
		client:     client,
		bucketName: bucketName,
	}
}

func (t *getTagsTask) Execute() any {
	res, err := t.client.GetTags(t.ctx, conn.NewGetS3BucketTagsParams(t.bucketName))
	if err != nil {
		t.logger.WithError(err).
			Error("unable to get tags")
		return GetTagsResult{
			BucketName: t.bucketName,
			Error:      err,
		}
	}

	t.logger.Debugf("bucket tags fetched")

	return GetTagsResult{
		BucketName: t.bucketName,
		Tags:       res.Tags,
		Error:      nil,
	}
}

func FindGetTagsResult(bucketName string, results []GetTagsResult) GetTagsResult {
	for _, result := range results {
		if result.BucketName == bucketName {
			return result
		}

	}
	return GetTagsResult{}
}
