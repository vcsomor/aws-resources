package aws_connector

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
	"time"
)

type ListS3Params struct {
	region string
}

type ListS3Result struct {
	Arn          string
	Name         string
	CreationTime *time.Time
}

type S3Operations interface {
	ListS3(ctx context.Context, p ListS3Params) ([]ListS3Result, error)
}

type defaultS3Operations struct {
	logger *logrus.Logger
	client *s3.Client
}

var _ S3Operations = (*defaultS3Operations)(nil)

func NewDefaultS3Operations(logger *logrus.Logger, client *s3.Client) S3Operations {
	return &defaultS3Operations{
		logger: logger,
		client: client,
	}
}

func (op *defaultS3Operations) ListS3(ctx context.Context, _ ListS3Params) ([]ListS3Result, error) {
	buckets, err := op.client.ListBuckets(ctx, nil)
	if err != nil {
		return nil, err
	}

	var res []ListS3Result
	for _, b := range buckets.Buckets {
		res = append(res, ListS3Result{
			Arn:          fmt.Sprintf("arn:aws:s3:::%s", *b.Name),
			Name:         *b.Name,
			CreationTime: b.CreationDate,
		})
	}

	return res, nil
}
