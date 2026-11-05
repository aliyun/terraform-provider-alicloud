package polardb

import "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"

// parameterToHash hashes a parameter set block for the schema.Set of the
// PolarDB cluster parameters resource. Migrated verbatim from the root
// package (alicloud/resource_alicloud_db_instance.go).
func parameterToHash(v interface{}) int {
	m := v.(map[string]interface{})
	return helper.Hashcode(m["name"].(string) + "|" + m["value"].(string))
}
