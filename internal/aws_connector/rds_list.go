package aws_connector

import "context"

type ListRDSParams struct {
	region string
}

type ListRDSResult struct {
	Name string
}

type RDSLister interface {
	ListRDS(ctx context.Context, p ListRDSParams) ([]ListRDSResult, error)
}
