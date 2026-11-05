package alidns

import (
	"fmt"
	"strings"
)

// validateRR is a local 1:1 copy of the DNS RR-specific definition in
// alicloud/validators.go. The root copy is retained for the root-package test
// consumer in alicloud/validators_test.go; alidns cannot import the root
// package. The generic validators leaf package (alicloud/validators) has no RR
// validator, so a product-local copy is the behavior-preserving disposition.
func validateRR(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if strings.HasPrefix(value, "-") || strings.HasSuffix(value, "-") {
		errors = append(errors, fmt.Errorf("RR is invalid, it can not starts or ends with '-'"))
	}

	if len(value) > 253 {
		errors = append(errors, fmt.Errorf("RR can not longer than 253 characters."))
	}

	for _, part := range strings.Split(value, ".") {
		if len(part) > 63 {
			errors = append(errors, fmt.Errorf("Each part of RR split with . can not longer than 63 characters."))
			return
		}
	}
	return
}
