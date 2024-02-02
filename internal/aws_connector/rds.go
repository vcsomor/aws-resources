package aws_connector

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"time"
)

type ListRDSParams struct {
}

type ListRDSResult struct {
	Arn        string
	ID         string
	CreateTime *time.Time
}

type RDSClient interface {
	List(ctx context.Context, p ListRDSParams) ([]ListRDSResult, error)
}

type rdsClient struct {
	client *rds.Client
}

var _ RDSClient = (*rdsClient)(nil)

func newRDSClient(client *rds.Client) RDSClient {
	return &rdsClient{
		client: client,
	}
}

func (c *rdsClient) List(ctx context.Context, _ ListRDSParams) ([]ListRDSResult, error) {
	listResult, err := c.client.DescribeDBInstances(ctx, nil)
	if err != nil {
		return nil, err
	}

	var res []ListRDSResult
	for _, r := range listResult.DBInstances {
		res = append(res, ListRDSResult{
			Arn:        *r.DBInstanceArn,
			ID:         *r.DBInstanceIdentifier,
			CreateTime: r.InstanceCreateTime,
		})
	}

	return res, nil
}

func transformRDSListTags(tags []types.Tag) map[string]*string {
	res := map[string]*string{}
	for _, t := range tags {
		if t.Key == nil {
			continue
		}
		res[*t.Key] = t.Value
	}
	return res
}
