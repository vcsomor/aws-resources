package aws_connector

import (
	"context"
	"github.com/sirupsen/logrus"
)

type awsLister struct {
	clientFactory ClientFactory
	logger        *logrus.Logger
}

var _ Lister = (*awsLister)(nil)

func NewLister(logger *logrus.Logger, factory ClientFactory) Lister {
	return &awsLister{
		clientFactory: factory,
		logger:        logger,
	}
}

func (l *awsLister) ListS3(ctx context.Context, _ ListS3Params) ([]ListS3Result, error) {
	c, err := l.clientFactory.S3Client(ctx)
	if err != nil {
		return nil, err
	}

	buckets, err := c.ListBuckets(ctx, nil)
	if err != nil {
		return nil, err
	}

	var res []ListS3Result
	for _, b := range buckets.Buckets {
		res = append(res, ListS3Result{
			Name: *b.Name,
		})
	}

	return res, nil
}

func (l *awsLister) ListRDS(ctx context.Context, _ ListRDSParams) ([]ListRDSResult, error) {
	c, err := l.clientFactory.RDSClient(ctx)
	if err != nil {
		return nil, err
	}

	rds, err := c.DescribeDBInstances(ctx, nil)
	if err != nil {
		return nil, err
	}

	var res []ListRDSResult
	for _, r := range rds.DBInstances {
		res = append(res, ListRDSResult{
			Name: *r.DBInstanceIdentifier,
		})
	}

	return res, nil
}
