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

	InstanceType     *string
	AvailabilityZone *string
	AllocatedStorage *int32
	Engine           *string
	EngineVersion    *string
	ReplicaMode      string
	Status           *string
	MultiAz          *bool
	MultiTenant      *bool

	Tags map[string]*string
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
	desc, err := c.client.DescribeDBInstances(ctx, nil)
	if err != nil {
		return nil, err
	}

	var res []ListRDSResult
	for _, r := range desc.DBInstances {
		res = append(res, ListRDSResult{
			Arn:        safeDeref(r.DBInstanceArn, ""),
			ID:         safeDeref(r.DBInstanceIdentifier, ""),
			CreateTime: r.InstanceCreateTime,

			InstanceType:     r.DBInstanceClass,
			AvailabilityZone: r.AvailabilityZone,
			AllocatedStorage: r.AllocatedStorage,
			Engine:           r.Engine,
			EngineVersion:    r.EngineVersion,
			ReplicaMode:      string(r.ReplicaMode),
			Status:           r.DBInstanceStatus,
			MultiAz:          r.MultiAZ,
			MultiTenant:      r.MultiTenant,

			Tags: transformTagsList(r.TagList),
		})
	}

	return res, nil
}

func transformTagsList(tags []types.Tag) map[string]*string {
	res := map[string]*string{}
	for _, t := range tags {
		if t.Key == nil {
			continue
		}
		res[*t.Key] = t.Value
	}
	return res
}

func safeDeref[T any](p *T, def T) T {
	if p == nil {
		return def
	}
	return *p
}
