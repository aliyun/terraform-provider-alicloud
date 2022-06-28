package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudOOSStateConfigurationsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.OOSSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosStateConfigurationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_oos_state_configuration.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudOosStateConfigurationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_oos_state_configuration.default.id}_fake"]`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosStateConfigurationsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_oos_state_configuration.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
		}),
		fakeConfig: "",
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosStateConfigurationsDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_oos_state_configuration.default.id}"]`,
			"tags": `{Created = "TF"}`,
		}),
		fakeConfig: testAccCheckAlicloudOosStateConfigurationsDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_oos_state_configuration.default.id}"]`,
			"tags": `{Created-Fake = "TF"}`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosStateConfigurationsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_oos_state_configuration.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
			"tags":              `{Created = "TF"}`,
		}),
		fakeConfig: testAccCheckAlicloudOosStateConfigurationsDataSourceName(rand, map[string]string{
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
			"configurations.0.create_time":            CHECKSET,
			"configurations.0.update_time":            CHECKSET,
			"configurations.0.description":            fmt.Sprintf("tf-testAccStateConfiguration-%d", rand),
			"configurations.0.parameters":             CHECKSET,
			"configurations.0.resource_group_id":      CHECKSET,
			"configurations.0.schedule_expression":    CHECKSET,
			"configurations.0.schedule_type":          "rate",
			"configurations.0.id":                     CHECKSET,
			"configurations.0.state_configuration_id": CHECKSET,
			"configurations.0.tags.%":                 "1",
			"configurations.0.tags.Created":           "TF",
			"configurations.0.targets":                CHECKSET,
			"configurations.0.template_id":            CHECKSET,
			"configurations.0.template_name":          "ACS-ECS-InventoryDataCollection",
			"configurations.0.template_version":       CHECKSET,
		}
	}
	var fakeAlicloudOosStateConfigurationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudOosStateConfigurationsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_oos_state_configurations.default",
		existMapFunc: existAlicloudOosStateConfigurationsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudOosStateConfigurationsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudOosStateConfigurationsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, resourceGroupIdConf, tagsConf, allConf)
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
