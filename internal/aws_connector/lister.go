package aws_connector

import "context"

type awsLister struct {
	clientFactory ClientFactory
}

var _ Lister = (*awsLister)(nil)

func NewLister(factory ClientFactory) Lister {
	return &awsLister{
		clientFactory: factory,
	}
}

func (l *awsLister) ListS3(ctx context.Context, _ ListS3Params) []ListS3Result {
	_, _ = l.clientFactory.S3Client(ctx)

	return []ListS3Result{
		{
			Name: "my-bucket",
		},
	}
}

func (l *awsLister) ListRDS(ctx context.Context, _ ListRDSParams) []ListRDSResult {
	_, _ = l.clientFactory.RDSClient(ctx)

	return []ListRDSResult{
		{
			Name: "my-rds-instance",
		},
	}
}
