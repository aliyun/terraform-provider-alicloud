package apigateway_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudApigatewayGroupsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_api_gateway_groups.default"
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccGroup_%d", rand),
		dataSourceApiGatewayGroupsConfigDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_api_gateway_group.default.name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_api_gateway_group.default.name}_fake",
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_api_gateway_group.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_api_gateway_group.default.id}_fake"},
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_api_gateway_group.default.name}",
			"ids":        []string{"${alicloud_api_gateway_group.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_api_gateway_group.default.name}_fake",
			"ids":        []string{"${alicloud_api_gateway_group.default.id}"},
		}),
	}

	var existApiGatewayGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"names.#":                 "1",
			"names.0":                 fmt.Sprintf("tf_testAccGroup_%d", rand),
			"groups.#":                "1",
			"groups.0.name":           fmt.Sprintf("tf_testAccGroup_%d", rand),
			"groups.0.description":    "tf_testAcc api gateway description",
			"groups.0.region_id":      acctest.CHECKSET,
			"groups.0.sub_domain":     acctest.CHECKSET,
			"groups.0.created_time":   acctest.CHECKSET,
			"groups.0.modified_time":  acctest.CHECKSET,
			"groups.0.traffic_limit":  acctest.CHECKSET,
			"groups.0.billing_status": acctest.CHECKSET,
			"groups.0.illegal_status": acctest.CHECKSET,
		}
	}
	var fakeApiGatewayGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"groups.#": "0",
		}
	}
	var apiGatewayGroupsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existApiGatewayGroupsMapFunc,
		FakeMapFunc:  fakeApiGatewayGroupsMapFunc,
	}

	apiGatewayGroupsCheckInfo.DataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}
func dataSourceApiGatewayGroupsConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}

	variable "description" {
	  default = "tf_testAcc api gateway description"
	}

	resource "alicloud_api_gateway_group" "default" {
	  name = "${var.name}"
	  description = "${var.description}"
	}
	`, name)
}
