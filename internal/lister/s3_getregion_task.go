package lister

import (
	"context"
	"github.com/sirupsen/logrus"
	conn "github.com/vcsomor/aws-resources/internal/aws_connector"
	"github.com/vcsomor/aws-resources/internal/threads"
)

type s3GetRegionTask struct {
	ctx        context.Context
	logger     *logrus.Entry
	client     conn.S3Client
	bucketName string
}

type s3GetRegionResult struct {
	bucketName string
	region     string
	error      error
}

var _ threads.Task = (*s3GetRegionTask)(nil)

func newS3GetRegionTask(
	ctx context.Context,
	logger *logrus.Entry,
	client conn.S3Client,
	bucketName string,
) threads.Task {
	return &s3GetRegionTask{
		ctx:        ctx,
		logger:     logger,
		client:     client,
		bucketName: bucketName,
	}
}

func (t *s3GetRegionTask) Execute() any {
	res, err := t.client.GetRegion(t.ctx, conn.NewGetS3RegionParams(t.bucketName))
	if err != nil {
		t.logger.WithError(err).
			Error("unable to get region")
		return s3GetRegionResult{
			bucketName: t.bucketName,
			error:      err,
		}
	}

	t.logger.Debugf("bucket region fetched")

	return s3GetRegionResult{
		bucketName: res.Name,
		region:     res.Region,
		error:      nil,
	}
}
