package rkvstore

import "fmt"

// Product-specific validator moved from the root package (validators.go),
// per stock-method-disposition: product-specific validators relocate with
// their resource. Behavior preserved verbatim.

func validateRedisConfig(v interface{}, k string) (ws []string, errors []error) {
	value, _ := v.(map[string]interface{})

	if len(value) < 1 {
		errors = append(errors, fmt.Errorf("invalid value for %s (%s)", k, value))
	}

	return
}
