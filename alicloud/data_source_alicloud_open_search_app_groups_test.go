package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudOpenSearchAppGroupDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOpenSearchAppGroupDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_open_search_app_group.default.id}"]`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudOpenSearchAppGroupDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_open_search_app_group.default.id}_fake"]`,
			"enable_details": "true",
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOpenSearchAppGroupDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_open_search_app_group.default.app_group_name}"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudOpenSearchAppGroupDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_open_search_app_group.default.app_group_name}_fake"`,
			"enable_details": "true",
		}),
	}
	instanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOpenSearchAppGroupDataSourceName(rand, map[string]string{
			"instance_id":    `"${alicloud_open_search_app_group.default.instance_id}"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudOpenSearchAppGroupDataSourceName(rand, map[string]string{
			"instance_id":    `"${alicloud_open_search_app_group.default.instance_id}_fake"`,
			"enable_details": "true",
		}),
	}
	typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOpenSearchAppGroupDataSourceName(rand, map[string]string{
			"type":           `"${alicloud_open_search_app_group.default.type}"`,
			"instance_id":    `"${alicloud_open_search_app_group.default.instance_id}"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudOpenSearchAppGroupDataSourceName(rand, map[string]string{
			"type":           `"enhanced"`,
			"enable_details": "true",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOpenSearchAppGroupDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_open_search_app_group.default.id}"]`,
			"name_regex":     `"${alicloud_open_search_app_group.default.app_group_name}"`,
			"instance_id":    `"${alicloud_open_search_app_group.default.instance_id}"`,
			"type":           `"${alicloud_open_search_app_group.default.type}"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudOpenSearchAppGroupDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_open_search_app_group.default.id}_fake"]`,
			"name_regex":     `"${alicloud_open_search_app_group.default.app_group_name}_fake"`,
			"type":           `"enhanced"`,
			"enable_details": "true",
		}),
	}
	var existAlicloudOpenSearchAppGroupDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"names.#":                 "1",
			"groups.#":                "1",
			"groups.0.payment_type":   "PayAsYouGo",
			"groups.0.app_group_name": fmt.Sprintf("tf_testacc_%d", rand),
		}
	}
	var fakeAlicloudOpenSearchAppGroupDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudOpenSearchAppGroupCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_open_search_app_groups.default",
		existMapFunc: existAlicloudOpenSearchAppGroupDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudOpenSearchAppGroupDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.OpenSearchSupportRegions)
	}
	alicloudOpenSearchAppGroupCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, instanceIdConf, typeConf, allConf)
}
func testAccCheckAlicloudOpenSearchAppGroupDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf_testacc_%d"
}

resource "alicloud_open_search_app_group" "default" {
  app_group_name = var.name
  payment_type   = "PayAsYouGo"
  type           = "standard"
  quota {
    doc_size         = 1
    compute_resource = 20
    spec             = "opensearch.share.common"
  }
}

data "alicloud_open_search_app_groups" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
