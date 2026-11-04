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

// dbInstanceVswitchIdsOrderDiffSuppressFunc suppresses order-only changes on
// `vswitch_ids`. For a multi-AZ RDS instance the OpenAPI VSwitchId field is a
// comma-joined string whose member order is bound to the zone order
// (ZoneId/ZoneIdSlave1/ZoneIdSlave2): reordering vswitch_ids while the zones
// stay fixed is rejected by CreateDBInstance/ModifyDBInstance, so member order
// is not an independently applicable part of the attribute's semantics.
// Suppressing an order-only diff also keeps a ForceNew vswitch_ids reorder
// from destroying the instance to apply a config the API would then reject.
// Membership changes (adding/removing a vswitch) remain ForceNew and are not
// suppressed.
//
// The SDK invokes the suppressor once per changed list element, where the
// element values alone cannot decide reorder-ness, so the full old and new
// lists are read from the resource data. The unit net calls this function with
// a nil ResourceData and the member lists serialized as comma-joined strings.
func dbInstanceVswitchIdsOrderDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d == nil {
		return dbInstanceVswitchIdsSameMembers(strings.Split(old, ","), strings.Split(new, ","))
	}
	listKey := k
	if i := strings.LastIndexByte(k, '.'); i > 0 {
		listKey = k[:i]
	}
	o, n := d.GetChange(listKey)
	oldList, ok := o.([]interface{})
	if !ok {
		return false
	}
	newList, ok := n.([]interface{})
	if !ok {
		return false
	}
	return dbInstanceVswitchIdsSameMembers(dbInstanceVswitchIdsToStringSlice(oldList), dbInstanceVswitchIdsToStringSlice(newList))
}

// dbInstanceVswitchIdsSameMembers reports whether two vswitch id slices hold
// the same members in any order. An empty operand never suppresses on its own.
func dbInstanceVswitchIdsSameMembers(old, new []string) bool {
	if len(old) != len(new) {
		return false
	}
	if len(old) == 0 {
		return false
	}
	counts := make(map[string]int, len(old))
	for _, v := range old {
		counts[v]++
	}
	for _, v := range new {
		counts[v]--
		if counts[v] < 0 {
			return false
		}
	}
	return true
}

// dbInstanceVswitchIdsToStringSlice converts a TypeList element slice to plain
// strings; a non-string member never suppresses.
func dbInstanceVswitchIdsToStringSlice(list []interface{}) []string {
	out := make([]string, len(list))
	for i, v := range list {
		s, ok := v.(string)
		if !ok {
			return nil
		}
		out[i] = s
	}
	return out
}

// dbInstanceVswitchIdDiffSuppressFunc suppresses `vswitch_id` diffs that are
// not independently applicable. The list form `vswitch_ids` and the string
// form `vswitch_id` are mutually exclusive (ConflictsWith); when the user has
// configured `vswitch_ids`, the `vswitch_id` diff is suppressed to avoid
// perpetual drift between the two forms.
//
// Whether `vswitch_ids` is configured is read from the raw config, not from
// d.GetOk: `vswitch_ids` is Optional+Computed and Read always backfills the
// API value into state, so GetOk would report it as set even when the user
// only configured `vswitch_id`, silently suppressing a single-vswitch
// ForceNew change and drifting config from the instance (the F-1 regression).
// A multi-zone `vswitch_id` ("vsw-a,vsw-b") is also suppressed, matching the
// prior behavior for single-vswitch changes.
func dbInstanceVswitchIdDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d == nil {
		return dbInstanceVswitchIdDecide(true, false, new)
	}
	return dbInstanceVswitchIdDecide(d.Id() == "", dbInstanceVswitchIdsConfigured(d), new)
}

// dbInstanceVswitchIdsConfigured reports whether the user configured the list
// form `vswitch_ids` in their HCL config. It reads the raw config so only
// attributes the user actually set are considered, ignoring the Computed state
// value Read backfills (see dbInstanceVswitchIdDiffSuppressFunc).
func dbInstanceVswitchIdsConfigured(d *schema.ResourceData) bool {
	raw := d.GetRawConfig()
	if raw.IsNull() {
		return false
	}
	v := raw.GetAttr("vswitch_ids")
	return !v.IsNull() && v.IsKnown() && v.LengthInt() > 0
}

// dbInstanceVswitchIdDecide is the pure decision core of
// dbInstanceVswitchIdDiffSuppressFunc, extracted for unit testing. isNew
// reports a new resource (d.Id() == ""). vswitchIdsConfigured reports whether
// the user configured the list form `vswitch_ids` (never the Computed state
// value Read backfills). newVswitchId is the configured `vswitch_id`.
func dbInstanceVswitchIdDecide(isNew, vswitchIdsConfigured bool, newVswitchId string) bool {
	if isNew {
		return false
	}
	if vswitchIdsConfigured {
		return true
	}
	return len(strings.Split(newVswitchId, ",")) > 1
}
