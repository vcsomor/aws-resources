package args

import (
	"slices"
	"strings"
)

const (
	outputSeparator = ","

	OutputStdout string = "stdout"
	OutputFile   string = "file"

	OutputDefault = OutputStdout
)

func ParseOutputs(arg string) []string {
	r := sanitizeOutputArgs(arg)

	if len(r) == 0 {
		return []string{OutputDefault}
	}

	desiredOutputs := strings.Split(r, outputSeparator)
	if len(desiredOutputs) == 0 {
		return []string{OutputDefault}
	}

	var result []string
	for _, output := range []string{OutputFile, OutputStdout} {
		if !slices.Contains(desiredOutputs, output) {
			continue
		}
		result = append(result, output)
	}
	return result
}

func sanitizeOutputArgs(regions string) string {
	r := regions
	r = strings.ReplaceAll(r, " ", "")
	r = strings.ReplaceAll(r, ";", "")
	r = strings.ReplaceAll(r, ":", "")
	r = strings.ReplaceAll(r, "\n", "")
	r = strings.ReplaceAll(r, "\t", "")
	r = strings.ReplaceAll(r, "\r", "")
	return strings.ToLower(r)
}
