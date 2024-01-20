package lister

import "strings"

const regionSeparator = ","

const regionAll string = "all"
const regionUsEast2 string = "us-east-2"
const regionUsEast1 string = "us-east-1"
const regionUsWest1 string = "us-west-1"
const regionUsWest2 string = "us-west-2"
const regionAfSouth1 string = "af-south-1"
const regionApEast1 string = "ap-east-1"
const regionApSouth2 string = "ap-south-2"
const regionApSoutheast3 string = "ap-southeast-3"
const regionApSoutheast4 string = "ap-southeast-4"
const regionApSouth1 string = "ap-south-1"
const regionApNortheast3 string = "ap-northeast-3"
const regionApNortheast2 string = "ap-northeast-2"
const regionApSoutheast1 string = "ap-southeast-1"
const regionApSoutheast2 string = "ap-southeast-2"
const regionApNortheast1 string = "ap-northeast-1"
const regionCaCentral1 string = "ca-central-1"
const regionCaWest1 string = "ca-west-1"
const regionEuCentral1 string = "eu-central-1"
const regionEuWest1 string = "eu-west-1"
const regionEuWest2 string = "eu-west-2"
const regionEuSouth1 string = "eu-south-1"
const regionEuWest3 string = "eu-west-3"
const regionEuSouth2 string = "eu-south-2"
const regionEuNorth1 string = "eu-north-1"
const regionEuCentral2 string = "eu-central-2"
const regionIlCentral1 string = "il-central-1"
const regionMeSouth1 string = "me-south-1"
const regionMeCentral1 string = "me-central-1"
const regionSaEast1 string = "sa-east-1"

func allRegions() []string {
	return []string{
		regionUsEast2,
		regionUsEast1,
		regionUsWest1,
		regionUsWest2,
		regionAfSouth1,
		regionApEast1,
		regionApSouth2,
		regionApSoutheast3,
		regionApSoutheast4,
		regionApSouth1,
		regionApNortheast3,
		regionApNortheast2,
		regionApSoutheast1,
		regionApSoutheast2,
		regionApNortheast1,
		regionCaCentral1,
		regionCaWest1,
		regionEuCentral1,
		regionEuWest1,
		regionEuWest2,
		regionEuSouth1,
		regionEuWest3,
		regionEuSouth2,
		regionEuNorth1,
		regionEuCentral2,
		regionIlCentral1,
		regionMeSouth1,
		regionMeCentral1,
		regionSaEast1,
	}
}

func parseRegions(regions string) []string {
	r := sanitizeRegions(regions)

	if len(r) == 0 {
		return allRegions()
	}

	var result []string
	for _, region := range strings.Split(r, regionSeparator) {
		if region == regionAll {
			return allRegions()
		}
		result = append(result, region)
	}
	return result
}

func sanitizeRegions(regions string) string {
	r := regions
	r = strings.ReplaceAll(r, " ", "")
	r = strings.ReplaceAll(r, ";", "")
	r = strings.ReplaceAll(r, ":", "")
	r = strings.ReplaceAll(r, "\n", "")
	r = strings.ReplaceAll(r, "\t", "")
	r = strings.ReplaceAll(r, "\r", "")
	return strings.ToLower(r)

}
