package polardb

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// parameterToHash hashes a parameter set block for the schema.Set of the
// PolarDB cluster parameters resource. Migrated verbatim from the root
// package (alicloud/resource_alicloud_db_instance.go).
func parameterToHash(v interface{}) int {
	m := v.(map[string]interface{})
	return helper.Hashcode(m["name"].(string) + "|" + m["value"].(string))
}

// tdeAutomaticRotationRequestValue reports the EnableAutomaticRotation value
// that an update has to submit and whether it has to be submitted at all.
// The argument is input only (the API never returns it), so the merged update
// state cannot tell "explicitly false" apart from "absent or removed". The raw
// configuration is authoritative for what the user wrote: null or unknown
// means the argument is not configured, and only an explicit removal performed
// by this update is submitted through the planned value.
func tdeAutomaticRotationRequestValue(d *schema.ResourceData) (bool, bool) {
	rawConfig := d.GetRawConfig()
	if !rawConfig.IsNull() {
		value, diags := d.GetRawConfigAt(cty.GetAttrPath("enable_automatic_rotation"))
		if !diags.HasError() && value.Type().Equals(cty.Bool) && value.IsKnown() && !value.IsNull() {
			return value.True(), true
		}
	}
	if !d.IsNewResource() && d.HasChange("enable_automatic_rotation") {
		return d.Get("enable_automatic_rotation").(bool), true
	}
	return false, false
}
