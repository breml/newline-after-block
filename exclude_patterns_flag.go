package newlineafterblock

import (
	"fmt"
	"regexp"
	"strings"
)

// excludePatterns is a custom flag type that holds regex patterns for excluding files.
type excludePatterns struct {
	patterns []*regexp.Regexp
	raw      []string
}

// String returns a string representation of the exclude patterns.
func (e *excludePatterns) String() string {
	return strings.Join(e.raw, ",")
}

// Set adds a new exclude pattern, validating it as a regex.
func (e *excludePatterns) Set(value string) error {
	re, err := regexp.Compile(value)
	if err != nil {
		return fmt.Errorf("invalid regex pattern %q: %w", value, err)
	}

	e.patterns = append(e.patterns, re)
	e.raw = append(e.raw, value)
	return nil
}

// matches checks if a given path matches any of the exclude patterns.
func (e *excludePatterns) matches(path string) bool {
	for _, re := range e.patterns {
		if re.MatchString(path) {
			return true
		}
	}

	return false
}
