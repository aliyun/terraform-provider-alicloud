package ons

import (
	"fmt"
	"regexp"
	"strings"
)

func validateOnsGroupId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !(strings.HasPrefix(value, "GID-") || strings.HasPrefix(value, "GID_")) {
		errors = append(errors, fmt.Errorf("%q is invalid, it must start with 'GID-' or 'GID_'", k))
	}
	if reg := regexp.MustCompile(`^[\w\-]{7,64}$`); !reg.MatchString(value) {
		errors = append(errors, fmt.Errorf("%q length is limited to 7-64 and only characters such as letters, digits, '_' and '-' are allowed", k))
	}
	return
}
