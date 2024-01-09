package aws_connector

type ListRDSParams struct {
	region string
}

type ListRDSResult struct {
	Name string
}

type RDSLister interface {
	ListRDS(p ListRDSParams) []ListRDSResult
}
