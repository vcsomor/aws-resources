package s3_tasks

import (
	"context"
	"github.com/sirupsen/logrus"
	conn "github.com/vcsomor/aws-resources/internal/aws_connector"
	"github.com/vcsomor/aws-resources/internal/executor"
)

type getRegionTask struct {
	ctx        context.Context
	logger     *logrus.Entry
	client     conn.S3Client
	bucketName string
}

type GetRegionResult struct {
	BucketName string
	Region     string
	Error      error
}

var _ executor.Task = (*getRegionTask)(nil)

func NewS3GetRegionTask(
	ctx context.Context,
	logger *logrus.Entry,
	client conn.S3Client,
	bucketName string,
) executor.Task {
	return &getRegionTask{
		ctx:        ctx,
		logger:     logger,
		client:     client,
		bucketName: bucketName,
	}
}

func (t *getRegionTask) Execute() any {
	res, err := t.client.GetRegion(t.ctx, conn.NewGetS3RegionParams(t.bucketName))
	if err != nil {
		t.logger.WithError(err).
			Error("unable to get region")
		return GetRegionResult{
			BucketName: t.bucketName,
			Error:      err,
		}
	}

	t.logger.Debugf("bucket region fetched")

	return GetRegionResult{
		BucketName: res.Name,
		Region:     res.Region,
		Error:      nil,
	}
}
