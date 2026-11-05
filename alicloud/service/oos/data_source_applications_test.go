package oos_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudOOSApplicationsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.OOSApplicationSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudOosApplicationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_oos_application.default.application_name}"]`,
		}),
		FakeConfig: testAccCheckAlicloudOosApplicationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_oos_application.default.application_name}_fake"]`,
		}),
	}
	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudOosApplicationsDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_oos_application.default.application_name}"]`,
			"tags": `{Created = "TF"}`,
		}),
		FakeConfig: testAccCheckAlicloudOosApplicationsDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_oos_application.default.application_name}_fake"]`,
			"tags": `{Created = "TF_FAKE"}`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudOosApplicationsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_oos_application.default.application_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudOosApplicationsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_oos_application.default.application_name}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudOosApplicationsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_oos_application.default.application_name}"]`,
			"name_regex": `"${alicloud_oos_application.default.application_name}"`,
			"tags":       `{Created = "TF"}`,
		}),
		FakeConfig: testAccCheckAlicloudOosApplicationsDataSourceName(rand, map[string]string{
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
			"applications.0.create_time":       acctest.CHECKSET,
			"applications.0.description":       acctest.CHECKSET,
			"applications.0.update_time":       acctest.CHECKSET,
			"applications.0.resource_group_id": acctest.CHECKSET,
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
	var alicloudOosApplicationsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_oos_applications.default",
		ExistMapFunc: existAlicloudOosApplicationsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudOosApplicationsDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudOosApplicationsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, tagsConf, nameRegexConf, allConf)
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
