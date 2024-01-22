package aws_connector

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"time"
)

type ListS3Params struct {
}

type ListS3Result struct {
	Arn          string
	Name         string
	CreationTime *time.Time
}

type S3Client interface {
	ListS3(ctx context.Context, p ListS3Params) ([]ListS3Result, error)
}

type s3Client struct {
	client *s3.Client
}

var _ S3Client = (*s3Client)(nil)

func newS3Client(client *s3.Client) S3Client {
	return &s3Client{
		client: client,
	}
}

func (c *s3Client) ListS3(ctx context.Context, _ ListS3Params) ([]ListS3Result, error) {
	buckets, err := c.client.ListBuckets(ctx, nil)
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
