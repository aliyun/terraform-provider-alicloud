package polardb

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// PolarDB resource-specific diff-suppress helpers, migrated verbatim from the
// root package (alicloud/diff_suppress_funcs.go). Package-private, kept beside
// the resources that own them; the product package does not import root
// validators or suppressors. Behavior is unchanged.

func kmsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("password"); ok && v.(string) != "" {
		return true
	}
	if v, ok := d.GetOk("account_password"); ok && v.(string) != "" {
		return true
	}
	return false
}

func polardbPostPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("pay_type").(string) == "PrePaid" {
		return false
	}
	return true
}

func polardbDBClusterVersionDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	var latestVersion interface{}
	latestVersion, ok := d.GetOk("cluster_latest_version")
	if !ok {
		return true
	}
	if old == latestVersion.(string) {
		return true
	}
	return false
}

func polardbTDEAndEnabledDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("tde_status"); ok && v.(string) == "Enabled" && old != "" && new != "" && old != new {
		return true
	}
	return false
}

func polardbStorageTypeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("creation_category"); ok && v.(string) == "SENormal" {
		if w, ok := d.GetOk("storage_type"); ok && w.(string) == "ESSDAUTOPL" {
			return false
		}
	}
	return true
}

func polardbServrelessTypeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("db_type").(string) == "MySQL" && (d.Get("db_version").(string) == "8.0" || d.Get("db_version").(string) == "5.7") {
		if d.Get("serverless_type").(string) == "AgileServerless" || (d.Get("serverless_type").(string) == "SteadyServerless" && d.Get("serverless_steady_switch").(string) == "ON") {
			return false
		}
	} else if d.Get("db_type").(string) == "PostgreSQL" || d.Get("db_type").(string) == "Oracle" {
		if d.Get("serverless_type").(string) == "AgileServerless" || (d.Get("serverless_type").(string) == "SteadyServerless" && d.Get("serverless_steady_switch").(string) == "ON") {
			return false
		}
	}
	return true
}

func polardbPostPaidAndRenewDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("pay_type").(string) == "PrePaid" && d.Get("renewal_status").(string) != string(RenewNotRenewal) {
		return false
	}
	return true
}

func polardbProxyClassDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("creation_category"); ok && v.(string) == "SENormal" {
		return false
	}
	return true
}

func polardbProxyTypeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("creation_category").(string) == "SENormal" || d.Get("creation_category").(string) == "Normal" {
		return false
	}
	return true
}

func polardbDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	creationOption, optionOk := d.GetOk("creation_option")
	if dbType, ok := d.GetOk("db_type"); ok && dbType.(string) == "MySQL" && (creationOption == "Normal" || !optionOk) {
		return false
	}
	return true
}

func polardbXengineDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	creationOption, optionOk := d.GetOk("creation_option")
	if dbType, ok := d.GetOk("db_type"); ok && dbType.(string) == "MySQL" && d.Get("db_version").(string) == "8.0" && (creationOption != "CreateGdnStandby" || !optionOk) {
		return false
	}
	return true
}

func polardbAndCreationDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	creationCategory, categoryOk := d.GetOk("creation_category")
	if creationCategory == "Normal" || creationCategory == "NormalMultimaster" || !categoryOk {
		return false
	}
	return true
}

func polardbCompressStorageDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	creationOption, optionOk := d.GetOk("creation_option")
	if dbType, ok := d.GetOk("db_type"); ok && dbType.(string) == "MySQL" && (creationOption == "Normal" || !optionOk) && d.Get("storage_type").(string) == "PSL4" {
		return false
	}
	return true
}

func polardbStandbyAzDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	creationOption, optionOk := d.GetOk("creation_option")
	if (creationOption == "Normal" || !optionOk) && d.Get("hot_standby_cluster").(string) != "OFF" && old != "auto" && new != "auto" {
		return false
	}
	return true
}
