package polardb

import (
	"sort"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

func getChargeType(paymentType string) string {
	if strings.ToLower(paymentType) == strings.ToLower(string(PostPaid)) {
		return string(PostPaid)
	}
	return string(PrePaid)
}

func parameterToHash(v interface{}) int {
	m := v.(map[string]interface{})
	return helper.Hashcode(m["name"].(string) + "|" + m["value"].(string))
}

func convertIpsSetListToString(arr1 *schema.Set) []string {
	var ips []string
	for _, v := range arr1.List() {
		ips = append(ips, v.(string))
	}
	return ips
}

func convertIpsSetToString(sourceIps string) []string {
	ipsMap := make(map[string]string)

	for _, ip := range strings.Split(sourceIps, helper.CommaSeparated) {
		ipsMap[ip] = ip
	}
	var ips []string
	if len(ipsMap) > 0 {
		for key := range ipsMap {
			ips = append(ips, key)
		}
	}
	return ips
}

func arrValueEqual(arr1, arr2 []string) bool {
	if len(arr1) != len(arr2) {
		return false
	}
	sort.Strings(arr2)
	sort.Strings(arr1)
	for i, v := range arr1 {
		if v != arr2[i] {
			return false
		}
	}
	return true
}
