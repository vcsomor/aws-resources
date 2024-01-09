package aws_connector

type ListS3Params struct {
	region string
}

type ListS3Result struct {
	Name string
}

type S3Lister interface {
	ListS3(p ListS3Params) []ListS3Result
}
