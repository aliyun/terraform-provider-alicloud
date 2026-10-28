package oos_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudOOSStateConfigurationsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.OOSSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudOosStateConfigurationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_oos_state_configuration.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudOosStateConfigurationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_oos_state_configuration.default.id}_fake"]`,
		}),
	}
	resourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudOosStateConfigurationsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_oos_state_configuration.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
		}),
		FakeConfig: "",
	}
	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudOosStateConfigurationsDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_oos_state_configuration.default.id}"]`,
			"tags": `{Created = "TF"}`,
		}),
		FakeConfig: testAccCheckAlicloudOosStateConfigurationsDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_oos_state_configuration.default.id}"]`,
			"tags": `{Created-Fake = "TF"}`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudOosStateConfigurationsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_oos_state_configuration.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
			"tags":              `{Created = "TF"}`,
		}),
		FakeConfig: testAccCheckAlicloudOosStateConfigurationsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_oos_state_configuration.default.id}_fake"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
			"tags":              `{Created-Fake = "TF"}`,
		}),
	}
	var existAlicloudOosStateConfigurationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                   "1",
			"configurations.#":                        "1",
			"configurations.0.configure_mode":         "ApplyOnly",
			"configurations.0.create_time":            acctest.CHECKSET,
			"configurations.0.update_time":            acctest.CHECKSET,
			"configurations.0.description":            fmt.Sprintf("tf-testAccStateConfiguration-%d", rand),
			"configurations.0.parameters":             acctest.CHECKSET,
			"configurations.0.resource_group_id":      acctest.CHECKSET,
			"configurations.0.schedule_expression":    acctest.CHECKSET,
			"configurations.0.schedule_type":          "rate",
			"configurations.0.id":                     acctest.CHECKSET,
			"configurations.0.state_configuration_id": acctest.CHECKSET,
			"configurations.0.tags.%":                 "1",
			"configurations.0.tags.Created":           "TF",
			"configurations.0.targets":                acctest.CHECKSET,
			"configurations.0.template_id":            acctest.CHECKSET,
			"configurations.0.template_name":          "ACS-ECS-InventoryDataCollection",
			"configurations.0.template_version":       acctest.CHECKSET,
		}
	}
	var fakeAlicloudOosStateConfigurationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudOosStateConfigurationsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_oos_state_configurations.default",
		ExistMapFunc: existAlicloudOosStateConfigurationsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudOosStateConfigurationsDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudOosStateConfigurationsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, resourceGroupIdConf, tagsConf, allConf)
}
func testAccCheckAlicloudOosStateConfigurationsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccStateConfiguration-%d"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_oos_state_configuration" "default" {
  template_name       = "ACS-ECS-InventoryDataCollection"
  configure_mode      = "ApplyOnly"
  description         = var.name
  schedule_type       = "rate"
  schedule_expression = "1 hour"
  resource_group_id   = data.alicloud_resource_manager_resource_groups.default.ids.0
  targets             = "{\"Filters\": [{\"Type\": \"All\", \"Parameters\": {\"InstanceChargeType\": \"PrePaid\"}}], \"ResourceType\": \"ALIYUN::ECS::Instance\"}"
  parameters          = "{\"policy\": {\"ACS:Application\": {\"Collection\": \"Enabled\"}}}"
  tags = {
    Created = "TF"
  }
}

data "alicloud_oos_state_configurations" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
