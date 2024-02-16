package lister

import (
	"github.com/vcsomor/aws-resources/internal/lister/writer"
	"time"
)

const (
	logKeyRegion        = "region"
	logKeyResourceType  = "resource-type"
	logKeyResourceCount = "resource-count"

	s3ResourceType  = "S3"
	rdsResourceType = "RDS"
)

type Result struct {
	Arn          string     `json:"arn"`
	ID           string     `json:"id"`
	CreationTime *time.Time `json:"creationTime"`

	Data any `json:"data"`
}

type S3Data struct {
	LocationConstraint string             `json:"locationConstraint"`
	Tags               map[string]*string `json:"tags"`
}

type RDSData struct {
	Tags map[string]*string `json:"tags"`
}

type ResultBasedWriterFactory func(d Result) writer.Writer
