package redis

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// validateConfig validates the redis/tair "config" map attribute: it must
// contain at least one key. This is the product-specific validator that
// travelled with alicloud_redis_tair_instance when it migrated out of the flat
// root package. The still-flat alicloud_kvstore_instance keeps its own copy in
// the root package, so this leaf copy is intentionally local and does not touch
// the shared validators package.
func validateConfig(v interface{}, k string) (ws []string, errors []error) {
	value, _ := v.(map[string]interface{})

	if len(value) < 1 {
		errors = append(errors, fmt.Errorf("invalid value for %s (%s)", k, value))
	}

	return
}

// Compile-time assertion that validateConfig satisfies the SDK validator
// signature, the same way the root package's copy is wired into SchemaValidateFunc.
var _ schema.SchemaValidateFunc = validateConfig
