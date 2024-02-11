package aws_connector

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/ptr"
	"time"
)

type ListS3Params struct {
}

type ListS3Result struct {
	Name         string
	CreationTime *time.Time
}

type GetS3RegionParams struct {
	name string
}

type GetS3RegionResult struct {
	Name   string
	Region string
}

func NewGetS3RegionParams(name string) GetS3RegionParams {
	return GetS3RegionParams{name: name}
}

type GetS3BucketTagsParams struct {
	name string
}

type GetS3BucketTagsResult struct {
	Name string
	Tags map[string]*string
}

func NewGetS3BucketTagsParams(name string) GetS3BucketTagsParams {
	return GetS3BucketTagsParams{name: name}
}

type S3Client interface {
	List(ctx context.Context, p ListS3Params) ([]ListS3Result, error)
	GetRegion(ctx context.Context, p GetS3RegionParams) (GetS3RegionResult, error)
	GetTags(ctx context.Context, p GetS3BucketTagsParams) (GetS3BucketTagsResult, error)
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

func (c *s3Client) List(ctx context.Context, _ ListS3Params) ([]ListS3Result, error) {
	buckets, err := c.client.ListBuckets(ctx, nil)
	if err != nil {
		return nil, err
	}

	var res []ListS3Result
	for _, b := range buckets.Buckets {
		res = append(res, ListS3Result{
			Name:         *b.Name,
			CreationTime: b.CreationDate,
		})
	}

	return res, nil
}

func (c *s3Client) GetRegion(ctx context.Context, p GetS3RegionParams) (GetS3RegionResult, error) {
	loc, err := c.client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{Bucket: ptr.String(p.name)})
	if err != nil {
		return GetS3RegionResult{}, err
	}

	return GetS3RegionResult{
		Name:   p.name,
		Region: string(loc.LocationConstraint),
	}, nil
}

func (c *s3Client) GetTags(ctx context.Context, p GetS3BucketTagsParams) (GetS3BucketTagsResult, error) {
	tags, err := c.client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{Bucket: ptr.String(p.name)})
	if err != nil {
		return GetS3BucketTagsResult{}, err
	}

	res := GetS3BucketTagsResult{
		Name: p.name,
		Tags: map[string]*string{},
	}

	for _, t := range tags.TagSet {
		if t.Key == nil {
			continue
		}
		res.Tags[*t.Key] = t.Value
	}

	return res, nil
}
