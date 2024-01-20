package lister

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegionsAll(t *testing.T) {
	assert.Equal(t,
		[]string{
			"us-east-2",
			"us-east-1",
			"us-west-1",
			"us-west-2",
			"af-south-1",
			"ap-east-1",
			"ap-south-2",
			"ap-southeast-3",
			"ap-southeast-4",
			"ap-south-1",
			"ap-northeast-3",
			"ap-northeast-2",
			"ap-southeast-1",
			"ap-southeast-2",
			"ap-northeast-1",
			"ca-central-1",
			"ca-west-1",
			"eu-central-1",
			"eu-west-1",
			"eu-west-2",
			"eu-south-1",
			"eu-west-3",
			"eu-south-2",
			"eu-north-1",
			"eu-central-2",
			"il-central-1",
			"me-south-1",
			"me-central-1",
			"sa-east-1",
		},
		allRegions())
}

func TestSanitizeRegions(t *testing.T) {
	assert.Equal(t, "", sanitizeRegions("\t\t\t   ; : \n\n  \r"))
	assert.Equal(t, "asdfghghjkl", sanitizeRegions("ASDFGH GHJKL"))
}

func TestParseRegions(t *testing.T) {
	assert.Equal(t,
		[]string{
			"us-east-2",
			"us-east-1",
			"us-west-1",
			"us-west-2",
			"af-south-1",
			"ap-east-1",
			"ap-south-2",
			"ap-southeast-3",
			"ap-southeast-4",
			"ap-south-1",
			"ap-northeast-3",
			"ap-northeast-2",
			"ap-southeast-1",
			"ap-southeast-2",
			"ap-northeast-1",
			"ca-central-1",
			"ca-west-1",
			"eu-central-1",
			"eu-west-1",
			"eu-west-2",
			"eu-south-1",
			"eu-west-3",
			"eu-south-2",
			"eu-north-1",
			"eu-central-2",
			"il-central-1",
			"me-south-1",
			"me-central-1",
			"sa-east-1",
		},
		parseRegions(""),
	)
}