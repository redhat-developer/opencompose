package util

import (
	"fmt"
	"regexp"
)

const (
	maxResourceNameLength    = 253
	resourceNameRegexpString = "([A-Za-z0-9][-A-Za-z0-9_.]*)?"
)

var (
	resourceNameRegexp = regexp.MustCompile(resourceNameRegexpString)
)

func ValidateResourceName(name string) error {
	if len(name) > maxResourceNameLength {
		return fmt.Errorf("invalid resource name: length %d is bigger than maximum %d", len(name), maxResourceNameLength)
	}

	if !resourceNameRegexp.MatchString(name) {
		return fmt.Errorf("invalid resource name: it must match regexp '%s'", resourceNameRegexpString)
	}

	return nil
}
