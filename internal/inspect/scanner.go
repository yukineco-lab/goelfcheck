package inspect

import (
	"os"
	"regexp"
	"sort"
	"strings"
)

const maxPathFindings = 100

var absolutePathPatterns = []*regexp.Regexp{
	regexp.MustCompile(`/(home|Users|tmp)/[A-Za-z0-9._@%+\-/]+`),
	regexp.MustCompile(`[A-Za-z]:\\[A-Za-z0-9._@%+\- ]{2,}(?:\\[A-Za-z0-9._@%+\- ]+)+`),
	regexp.MustCompile(`\\\\[A-Za-z0-9._-]{2,}\\[A-Za-z0-9._@%+\- ]{2,}(?:\\[A-Za-z0-9._@%+\- ]{2,})+`),
}

func scanStrings(path string) (StringScan, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return StringScan{}, err
	}

	found := map[string]struct{}{}
	for _, s := range printableStrings(data, 4) {
		for _, re := range absolutePathPatterns {
			for _, match := range re.FindAllString(s, -1) {
				match = strings.Trim(match, "\x00\r\n\t ")
				if match != "" {
					found[match] = struct{}{}
				}
			}
		}
	}

	paths := make([]string, 0, len(found))
	for p := range found {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	if len(paths) > maxPathFindings {
		paths = paths[:maxPathFindings]
	}
	return StringScan{AbsolutePaths: paths}, nil
}

func printableStrings(data []byte, minLen int) []string {
	var out []string
	start := -1
	for i, b := range data {
		if isPrintableASCII(b) {
			if start == -1 {
				start = i
			}
			continue
		}
		if start != -1 && i-start >= minLen {
			out = append(out, string(data[start:i]))
		}
		start = -1
	}
	if start != -1 && len(data)-start >= minLen {
		out = append(out, string(data[start:]))
	}
	return out
}

func isPrintableASCII(b byte) bool {
	return b >= 0x20 && b <= 0x7e
}
