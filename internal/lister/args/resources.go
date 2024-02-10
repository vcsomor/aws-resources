package args

import (
	"slices"
	"strings"
)

const resourcesSeparator = ","

const resourcesAll string = "all"
const resourcesS3 string = "s3"
const resourcesRDS string = "rds"

func allResources() []string {
	return []string{
		resourcesS3,
		resourcesRDS,
	}
}

func ParseResources(res string) []string {
	r := sanitizeResourceArgs(res)

	if len(r) == 0 {
		return allResources()
	}

	if r == resourcesAll {
		return allResources()
	}

	desiredRes := strings.Split(r, resourcesSeparator)
	if len(desiredRes) == 0 {
		return allResources()
	}

	var result []string
	for _, resource := range allResources() {
		if !slices.Contains(desiredRes, resource) {
			continue
		}
		result = append(result, resource)
	}
	return result
}

func sanitizeResourceArgs(regions string) string {
	r := regions
	r = strings.ReplaceAll(r, " ", "")
	r = strings.ReplaceAll(r, ";", "")
	r = strings.ReplaceAll(r, ":", "")
	r = strings.ReplaceAll(r, "\n", "")
	r = strings.ReplaceAll(r, "\t", "")
	r = strings.ReplaceAll(r, "\r", "")
	return strings.ToLower(r)
}
