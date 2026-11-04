package elasticsearch

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func elasticsearchEnablePublicDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOkExists("enable_public"); ok && v.(bool) == true {
		return false
	}
	return true
}

func elasticsearchEnableKibanaPublicDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOkExists("enable_kibana_public_network"); ok && v.(bool) == true {
		return false
	}
	return true
}

func elasticsearchEnableKibanaPrivateDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOkExists("enable_kibana_private_network"); ok && v.(bool) == true {
		return false
	}
	return true
}

func esVersionDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {

	oldVersion := strings.Split(strings.Split(old, "_")[0], ".")
	newVersion := strings.Split(strings.Split(new, "_")[0], ".")

	if len(oldVersion) >= 2 && len(newVersion) >= 2 {
		if oldVersion[0] == newVersion[0] && oldVersion[1] == newVersion[1] {
			return true
		}
	}

	return false
}

func esDataNodeDiskPerformanceLevelDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("data_node_disk_type"); ok && v.(string) != "cloud_essd" {
		return true
	}
	return false
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
