package adb

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// Product-local diff suppressors moved from the root alicloud package; product
// subpackages cannot import the root diff_suppress_funcs.go definitions. Behavior
// is preserved 1:1.

func kmsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("password"); ok && v.(string) != "" {
		return true
	}
	if v, ok := d.GetOk("account_password"); ok && v.(string) != "" {
		return true
	}
	return false
}

func adbPostPaidAndRenewDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("pay_type"); ok && v.(string) == "PrePaid" && d.Get("renewal_status").(string) != string(RenewNotRenewal) {
		return false
	}
	if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" && d.Get("renewal_status").(string) != string(RenewNotRenewal) {
		return false
	}
	return true
}

func adbPostPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("pay_type"); ok && v.(string) == "PrePaid" {
		return false
	}
	if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
		return false
	}
	return true
}
