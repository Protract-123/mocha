package manifest

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	versionTokenPattern    = regexp.MustCompile(`[0-9]+|[A-Za-z]+`)      // a run of digits and/or letters
	prereleaseTokenPattern = regexp.MustCompile(`(?i)alpha|beta|rc|pre`) // prerelease tags
)

func tokenizeVersion(version string) []string {
	return versionTokenPattern.FindAllString(version, -1)
}

func CompareVersions(reference string, difference string) int {
	reference = strings.ReplaceAll(reference, "+", "-")
	difference = strings.ReplaceAll(difference, "+", "-")

	if reference == difference {
		return 0
	}

	referenceParts := tokenizeVersion(reference)
	differenceParts := tokenizeVersion(difference)

	if referenceParts[0] == "nightly" && differenceParts[0] == "nightly" {
		return 0
	}

	maxLen := max(len(referenceParts), len(differenceParts))

	for i := 0; i < maxLen; i++ {
		if i >= len(referenceParts) {
			if prereleaseTokenPattern.MatchString(differenceParts[i]) {
				return -1
			}
			return 1
		}
		if i >= len(differenceParts) {
			if prereleaseTokenPattern.MatchString(referenceParts[i]) {
				return 1
			}
			return -1
		}

		if result := comparePart(referenceParts[i], differenceParts[i]); result != 0 {
			return result
		}
	}

	return 0
}

func comparePart(part1 string, part2 string) int {
	number1, err1 := strconv.ParseInt(part1, 10, 64)
	number2, err2 := strconv.ParseInt(part2, 10, 64)

	if err1 == nil && err2 == nil {
		switch {
		case number2 > number1:
			return 1
		case number2 < number1:
			return -1
		default:
			return 0
		}
	}

	return strings.Compare(strings.ToLower(part2), strings.ToLower(part1))
}
