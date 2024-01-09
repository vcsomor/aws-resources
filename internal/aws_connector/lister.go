package aws_connector

type awsLister struct {
}

var _ Lister = (*awsLister)(nil)

func NewLister() Lister {
	return &awsLister{}
}

func (l *awsLister) ListS3(_ ListS3Params) []ListS3Result {
	return []ListS3Result{
		{
			Name: "my-bucket",
		},
	}
}

func (l *awsLister) ListRDS(_ ListRDSParams) []ListRDSResult {
	return []ListRDSResult{
		{
			Name: "my-rds-instance",
		},
	}
}
