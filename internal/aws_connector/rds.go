package aws_connector

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/sirupsen/logrus"
	"time"
)

type ListRDSParams struct {
}

type ListRDSResult struct {
	Arn        string
	ID         string
	CreateTime *time.Time
}

type RDSOperations interface {
	ListRDS(ctx context.Context, p ListRDSParams) ([]ListRDSResult, error)
}

type defaultRDSOperations struct {
	logger *logrus.Entry
	client *rds.Client
}

var _ RDSOperations = (*defaultRDSOperations)(nil)

func NewDefaultRDSOperations(logger *logrus.Entry, client *rds.Client) RDSOperations {
	return &defaultRDSOperations{
		logger: logger,
		client: client,
	}
}

func (op *defaultRDSOperations) ListRDS(ctx context.Context, _ ListRDSParams) ([]ListRDSResult, error) {
	describeResult, err := op.client.DescribeDBInstances(ctx, nil)
	if err != nil {
		return nil, err
	}

	var res []ListRDSResult
	for _, r := range describeResult.DBInstances {
		res = append(res, ListRDSResult{
			Arn:        *r.DBInstanceArn,
			ID:         *r.DBInstanceIdentifier,
			CreateTime: r.InstanceCreateTime,
		})
	}

	return res, nil
}
