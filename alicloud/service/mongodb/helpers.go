package mongodb

import (
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// The following helpers were copied verbatim from the root alicloud package
// (diff_suppress_funcs.go, resource_alicloud_db_instance.go and
// extension_rds.go) so that the mongodb product package stays self-contained
// after the directory migration. Behavior is unchanged from the originals.

// PostPaidDiffSuppressFunc mirrors the root alicloud diff_suppress_funcs.go
// implementation.
func PostPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	// payment_type is the instance_charge_type's replacement.
	// If both instance_charge_type and payment_type are "", it means hiding a default "PostPaid"
	if v, ok := d.GetOk("instance_charge_type"); ok && strings.ToLower(v.(string)) == "prepaid" {
		return false
	}
	if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
		return false
	}
	return true
}

// kmsDiffSuppressFunc mirrors the root alicloud diff_suppress_funcs.go
// implementation.
func kmsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("password"); ok && v.(string) != "" {
		return true
	}
	if v, ok := d.GetOk("account_password"); ok && v.(string) != "" {
		return true
	}
	return false
}

// parameterToHash mirrors the root alicloud resource_alicloud_db_instance.go
// implementation.
func parameterToHash(v interface{}) int {
	m := v.(map[string]interface{})
	return helper.Hashcode(m["name"].(string) + "|" + m["value"].(string))
}

// LOCAL_HOST_IP mirrors the root alicloud common.go declaration.
const LOCAL_HOST_IP = "127.0.0.1"

// BACKUP_TIME mirrors the root alicloud extension_rds.go declaration.
var BACKUP_TIME = []string{
	"00:00Z-01:00Z", "01:00Z-02:00Z", "02:00Z-03:00Z", "03:00Z-04:00Z", "04:00Z-05:00Z",
	"05:00Z-06:00Z", "06:00Z-07:00Z", "07:00Z-08:00Z", "08:00Z-09:00Z", "09:00Z-10:00Z",
	"10:00Z-11:00Z", "11:00Z-12:00Z", "12:00Z-13:00Z", "13:00Z-14:00Z", "14:00Z-15:00Z",
	"15:00Z-16:00Z", "16:00Z-17:00Z", "17:00Z-18:00Z", "18:00Z-19:00Z", "19:00Z-20:00Z",
	"20:00Z-21:00Z", "21:00Z-22:00Z", "22:00Z-23:00Z", "23:00Z-24:00Z",
}
