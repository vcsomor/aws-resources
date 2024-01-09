package aws_connector

type Lister interface {
	S3Lister
	RDSLister
}
