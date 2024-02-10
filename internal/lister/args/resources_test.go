package args

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResourcesAll(t *testing.T) {
	assert.Equal(t,
		[]string{
			"s3",
			"rds",
		},
		allResources())
}

func TestSanitizeResourceArgs(t *testing.T) {
	assert.Equal(t, "", sanitizeResourceArgs("\t\t\t   ; : \n\n  \r"))
	assert.Equal(t, "asdfghghjkl", sanitizeResourceArgs("ASDFGH GHJKL"))
}

func TestParseResourceArgs(t *testing.T) {
	assert.Equal(t,
		[]string{
			"s3",
			"rds",
		},
		ParseResources(""),
	)

	assert.Equal(t,
		[]string{
			"s3",
			"rds",
		},
		ParseResources("all"),
	)

	assert.Equal(t,
		[]string{
			"s3",
			"rds",
		},
		ParseResources("rds,s3"),
	)

	assert.Equal(t,
		[]string{
			"s3",
		},
		ParseResources("s3,something"),
	)
}
