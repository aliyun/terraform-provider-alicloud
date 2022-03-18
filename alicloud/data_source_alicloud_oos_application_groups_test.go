package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudOOSApplicationGroupsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.OOSApplicationSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosApplicationGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_oos_application_group.default.application_group_name}"]`,
		}),
		fakeConfig: testAccCheckAlicloudOosApplicationGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_oos_application_group.default.application_group_name}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosApplicationGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_oos_application_group.default.application_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudOosApplicationGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_oos_application_group.default.application_group_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosApplicationGroupsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_oos_application_group.default.application_group_name}"]`,
			"name_regex": `"${alicloud_oos_application_group.default.application_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudOosApplicationGroupsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_oos_application_group.default.application_group_name}_fake"]`,
			"name_regex": `"${alicloud_oos_application_group.default.application_group_name}_fake"`,
		}),
	}
	var existAlicloudOosApplicationGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"names.#":                         "1",
			"groups.#":                        "1",
			"groups.0.application_group_name": fmt.Sprintf("tf-testAccApplicationGroup-%d", rand),
			"groups.0.application_name":       CHECKSET,
			"groups.0.deploy_region_id":       os.Getenv("ALICLOUD_REGION"),
			"groups.0.description":            fmt.Sprintf("tf-testAccApplicationGroup-%d", rand),
			"groups.0.import_tag_key":         fmt.Sprintf("tf-testAccApplicationGroup-%d", rand),
			"groups.0.import_tag_value":       fmt.Sprintf("tf-testAccApplicationGroup-%d", rand),
			"groups.0.create_time":            CHECKSET,
			"groups.0.update_time":            CHECKSET,
			"groups.0.cms_group_id":           "",
			"groups.0.id":                     CHECKSET,
		}
	}
	var fakeAlicloudOosApplicationGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudOosApplicationGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_oos_application_groups.default",
		existMapFunc: existAlicloudOosApplicationGroupsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudOosApplicationGroupsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudOosApplicationGroupsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudOosApplicationGroupsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "tf-testAccApplicationGroup-%d"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_oos_application" "default" {
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  application_name  = var.name
  description       = var.name
  tags              = {
    Created = "TF"
  }
}

resource "alicloud_oos_application_group" "default" {
  application_group_name = var.name
  application_name       = alicloud_oos_application.default.id
  deploy_region_id       = "%s"
  description            = var.name
  import_tag_key         = var.name
  import_tag_value       = var.name
}

data "alicloud_oos_application_groups" "default" {
  application_name = alicloud_oos_application.default.id
  %s
}
`, rand, os.Getenv("ALICLOUD_REGION"), strings.Join(pairs, " \n "))
	return config
}
