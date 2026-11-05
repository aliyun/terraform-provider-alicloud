package smartag

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Diff-suppress helpers used by the smartag resources. The flat-layout root
// package owned the originals; the generic kmsDiffSuppressFunc is still shared
// from the root by other products, while the sag-specific helpers are private
// to this package. They are copied here verbatim so the migrated resources can
// keep calling them unqualified within package smartag.

// kmsDiffSuppressFunc suppresses diffs for sensitive fields when a password is
// already configured, mirroring the root package helper reused by
// resource_alicloud_sag_client_user.
func kmsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("password"); ok && v.(string) != "" {
		return true
	}
	if v, ok := d.GetOk("account_password"); ok && v.(string) != "" {
		return true
	}
	return false
}

// sagDnatEntryTypeDiffSuppressFunc suppresses the type diff for non-Intranet
// sag dnat entries.
func sagDnatEntryTypeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("type").(string) != "Intranet" {
		return true
	}
	return false
}

// sagClientUserPasswordSuppressFunc suppresses the password diff when no user
// name is set on the sag client user.
func sagClientUserPasswordSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("user_name").(string) == "" {
		return true
	}
	return false
}
