package rkvstore

import (
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Local private copies of product-specific diff-suppress helpers, per
// stock-method-disposition. Root-package copies are retained for other
// consumers; these keep behavior identical for the migrated R-kvstore types.

func kmsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("password"); ok && v.(string) != "" {
		return true
	}
	if v, ok := d.GetOk("account_password"); ok && v.(string) != "" {
		return true
	}
	return false
}

func redisPostPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("payment_type"); ok && v.(string) == "PrePaid" {
		return false
	}

	if v, ok := d.GetOk("instance_charge_type"); ok && v.(string) == "PrePaid" {
		return false
	}

	return true
}

func redisPostPaidAndRenewDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("payment_type"); ok && v.(string) == "PrePaid" && d.Get("auto_renew").(bool) {
		return false
	}

	if v, ok := d.GetOk("instance_charge_type"); ok && v.(string) == "PrePaid" && d.Get("auto_renew").(bool) {
		return false
	}

	return true
}

func redisSecurityGroupIdDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	oldArray := strings.Split(old, ",")
	newArray := strings.Split(new, ",")
	if len(oldArray) != len(newArray) {
		return false
	}
	sort.Strings(oldArray)
	sort.Strings(newArray)
	for i := range newArray {
		if newArray[i] != oldArray[i] {
			return false
		}
	}
	return true
}
