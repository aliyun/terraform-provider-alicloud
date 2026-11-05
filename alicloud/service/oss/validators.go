package oss

import (
	"fmt"
	"time"
)

// validateOssBucketDateTimestamp mirrors the OSS-specific validator kept in
// the root package (alicloud/validators.go). It migrates with the OSS product
// so the service package stays self-contained without importing the root
// package; the value, semantics and signature are identical to the root
// definition.
func validateOssBucketDateTimestamp(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := time.Parse("2006-01-02", value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q cannot be parsed as date YYYY-MM-DD Format", value))
	}
	return
}
