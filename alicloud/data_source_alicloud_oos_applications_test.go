package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudOOSApplicationsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.OOSApplicationSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosApplicationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_oos_application.default.application_name}"]`,
		}),
		fakeConfig: testAccCheckAlicloudOosApplicationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_oos_application.default.application_name}_fake"]`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosApplicationsDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_oos_application.default.application_name}"]`,
			"tags": `{Created = "TF"}`,
		}),
		fakeConfig: testAccCheckAlicloudOosApplicationsDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_oos_application.default.application_name}_fake"]`,
			"tags": `{Created = "TF_FAKE"}`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosApplicationsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_oos_application.default.application_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudOosApplicationsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_oos_application.default.application_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosApplicationsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_oos_application.default.application_name}"]`,
			"name_regex": `"${alicloud_oos_application.default.application_name}"`,
			"tags":       `{Created = "TF"}`,
		}),
		fakeConfig: testAccCheckAlicloudOosApplicationsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_oos_application.default.application_name}_fake"]`,
			"name_regex": `"${alicloud_oos_application.default.application_name}_fake"`,
			"tags":       `{Created = "TF_FAKE"}`,
		}),
	}
	var existAlicloudOosApplicationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"names.#":                          "1",
			"applications.#":                   "1",
			"applications.0.id":                fmt.Sprintf("tf-testAccApplication%d", rand),
			"applications.0.application_name":  fmt.Sprintf("tf-testAccApplication%d", rand),
			"applications.0.create_time":       CHECKSET,
			"applications.0.description":       CHECKSET,
			"applications.0.update_time":       CHECKSET,
			"applications.0.resource_group_id": CHECKSET,
			"applications.0.tags.%":            "1",
			"applications.0.tags.Created":      "TF",
		}
	}
	var fakeAlicloudOosApplicationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"names.#":        "0",
			"applications.#": "0",
		}
	}
	var alicloudOosApplicationsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_oos_applications.default",
		existMapFunc: existAlicloudOosApplicationsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudOosApplicationsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudOosApplicationsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, tagsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudOosApplicationsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccApplication%d"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_oos_application" "default" {
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  application_name  = var.name
  description       = var.name
  tags = {
    Created = "TF"
  }
}

data "alicloud_oos_applications" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
