package rds

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// These diff-suppress funcs are the package-rds copies of the root
// alicloud/diff_suppress_funcs.go helpers. They are root-native (not leaf
// re-exports). Four are RDS-specific (no other product references them); the
// rest are shared across products. Per stock-method-disposition rule #4 the
// root keeps its copies and RDS keeps these behavior-preserving local copies
// rather than importing the root package. Bodies are verbatim from the root.

func logRetentionPeriodDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("enable_backup_log").(bool) {
		return false
	}
	if d.Get("log_backup").(bool) {
		return false
	}
	if v, err := strconv.Atoi(new); err != nil && v > d.Get("backup_retention_period").(int) {
		return false
	}
	if v, err := strconv.Atoi(new); err != nil && v > d.Get("retention_period").(int) {
		return false
	}
	return true
}

func enableBackupLogDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("enable_backup_log").(bool) {
		return false
	}
	if d.Get("log_backup").(bool) {
		return false
	}

	return true
}

func archiveBackupPeriodDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("enable_backup_log").(bool) {
		return false
	}
	if d.Get("log_backup").(bool) {
		return false
	}
	if v, err := strconv.Atoi(new); err != nil && v+730 >= d.Get("backup_retention_period").(int) {
		return false
	}
	if v, err := strconv.Atoi(new); err != nil && v+730 >= d.Get("retention_period").(int) {
		return false
	}

	return true
}

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

func PostPaidAndRenewDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("instance_charge_type"); ok && strings.ToLower(v.(string)) == "prepaid" && d.Get("auto_renew").(bool) {
		return false
	}
	if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" && d.Get("auto_renew").(bool) {
		return false
	}
	return true
}

func kmsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("password"); ok && v.(string) != "" {
		return true
	}
	if v, ok := d.GetOk("account_password"); ok && v.(string) != "" {
		return true
	}
	return false
}

func sslEnabledDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("ssl_enabled"); ok && v.(int) == 1 {
		return false
	}
	return true
}

func sslActionDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("ssl_action"); ok && (v.(string) == "Open" || v.(string) == "Update") {
		return false
	}
	return true
}

func securityIpsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("security_ips"); ok && len(v.(*schema.Set).List()) > 0 {
		return false
	}
	return true
}

func kernelVersionDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("upgrade_db_instance_kernel_version"); ok && v.(bool) == true {
		return false
	}
	return true
}

func kernelSmallVersionDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.HasChange("target_minor_version") {
		if v, ok := d.GetOk("target_minor_version"); ok && v.(string) != "" {
			return false
		}
	}
	return true
}

func StorageAutoScaleDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("storage_auto_scale"); ok && strings.ToLower(v.(string)) == "enable" {
		return false
	}
	return true
}
