package rds

import (
	"fmt"
	"strconv"
)

// validateDBConnectionPort is the package-rds copy of the root helper. It is
// DB-connection specific and migrates with RDS; the root package keeps its own
// copy for non-RDS products (e.g. gpdb) that still reference the root helper.
func validateDBConnectionPort(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		port, err := strconv.Atoi(value)
		if err != nil {
			errors = append(errors, err)
		}
		if port < 1000 || port > 5999 {
			errors = append(errors, fmt.Errorf("%q cannot be less than 3001 and larger than 3999.", k))
		}
	}
	return
}
