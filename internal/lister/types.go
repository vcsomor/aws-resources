package lister

import "time"

const (
	logKeyRegion        = "region"
	logKeyResourceType  = "resource-type"
	logKeyResourceCount = "resource-count"

	s3ResourceType  = "S3"
	rdsResourceType = "RDS"
)

type ResultDataType interface {
	S3Data | RDSData
}

type Result[T ResultDataType] struct {
	Arn          string     `json:"arn"`
	ID           string     `json:"id"`
	CreationTime *time.Time `json:"creationTime"`

	Data T `json:"data"`
}

type S3Data struct {
	LocationConstraint string             `json:"locationConstraint"`
	Tags               map[string]*string `json:"tags"`
}

type RDSData struct {
	Tags map[string]*string `json:"tags"`
}
