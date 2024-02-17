package args

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSanitizeOutputArgs(t *testing.T) {
	assert.Equal(t, "", sanitizeOutputArgs("\t\t\t   ; : \n\n  \r"))
	assert.Equal(t, "asdfghghjkl", sanitizeOutputArgs("ASDFGH GHJKL"))
}

func TestParseOutputArgs(t *testing.T) {
	assert.Equal(t,
		[]string{
			"stdout",
		},
		ParseOutputs(""),
	)

	assert.Equal(t,
		[]string{
			"file",
			"stdout",
		},
		ParseOutputs("file,stdout"),
	)

	assert.Equal(t,
		[]string{
			"file",
		},
		ParseOutputs("file,something"),
	)
}
