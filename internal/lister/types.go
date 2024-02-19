package lister

import (
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

	Properties any `json:"properties"`
}

type S3Data struct {
	LocationConstraint string             `json:"locationConstraint"`
	Tags               map[string]*string `json:"tags"`
}

type RDSData struct {
	InstanceType     *string            `json:"instanceType"`
	AvailabilityZone *string            `json:"availabilityZone"`
	AllocatedStorage *int32             `json:"allocatedStorage"`
	Engine           *string            `json:"engine"`
	EngineVersion    *string            `json:"engineVersion"`
	ReplicaMode      string             `json:"replicaMode"`
	Status           *string            `json:"status"`
	MultiAz          *bool              `json:"multiAz"`
	MultiTenant      *bool              `json:"multiTenant"`
	Tags             map[string]*string `json:"tags"`
}
