package alicloud

import "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"

// Helpers still used by the PolarDB resources and services that remain in the
// flat package after alicloud_polardb_cluster moved to service/polardb.

func convertPolarDBSubCategoryUpdateRequest(source string) string {
	switch source {
	case "Exclusive":
		return "normal_exclusive"
	}
	return "normal_general"
}

// IsContain keeps the flat package's argument order on top of helper.InArray,
// which takes the target first.
func IsContain(items []string, item string) bool {
	return helper.InArray(item, items)
}
