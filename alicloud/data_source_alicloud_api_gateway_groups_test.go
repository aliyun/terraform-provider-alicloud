package alicloud

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudApigatewayGroupsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_api_gateway_groups.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccGroup_%d", rand),
		dataSourceApiGatewayGroupsConfigDependence)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_api_gateway_group.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_api_gateway_group.default.name}_fake",
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
			"groups.0.region_id":      CHECKSET,
			"groups.0.sub_domain":     CHECKSET,
			"groups.0.created_time":   CHECKSET,
			"groups.0.modified_time":  CHECKSET,
			"groups.0.traffic_limit":  CHECKSET,
			"groups.0.billing_status": CHECKSET,
			"groups.0.illegal_status": CHECKSET,
		}
	}
	var fakeApiGatewayGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"groups.#": "0",
		}
	}
	var apiGatewayGroupsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existApiGatewayGroupsMapFunc,
		fakeMapFunc:  fakeApiGatewayGroupsMapFunc,
	}

	apiGatewayGroupsCheckInfo.dataSourceTestCheck(t, rand, allConf)
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
