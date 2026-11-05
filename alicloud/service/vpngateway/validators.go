package vpngateway

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// validateSslVpnPortValue was previously defined in the root package alicloud
// (validators.go). It is reproduced here 1:1 so that the migrated package does
// not import the root package.
func validateSslVpnPortValue(is []int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		ws, errors = validation.IntBetween(1, 65535)(v, k)
		if errors != nil {
			return
		}

		value := v.(int)
		for _, i := range is {
			if i == value {
				errors = append(errors, fmt.Errorf(
					"%q must contain a valid int value should not be in array %#v, got %q",
					k, is, value))
				return
			}
		}
		return

	}
}
