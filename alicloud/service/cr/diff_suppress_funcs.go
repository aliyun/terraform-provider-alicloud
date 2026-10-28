package cr

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// kmsDiffSuppressFunc suppresses the KMS-encryption diff when a password or
// account_password is configured. It is migrated verbatim from the root
// alicloud package's diff_suppress_funcs.go so the cr subpackage keeps the
// original production behavior without importing the monolith.
func kmsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("password"); ok && v.(string) != "" {
		return true
	}
	if v, ok := d.GetOk("account_password"); ok && v.(string) != "" {
		return true
	}
	return false
}
