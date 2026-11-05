package apig_test

import (
	"fmt"
	"strings"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudApigServiceDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudApigServicesDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_service.default.id}"]`,
		}),
		FakeConfig: testAccCheckAliCloudApigServicesDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_service.default.id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudApigServicesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_apig_service.default.service_name}"`,
		}),
		FakeConfig: testAccCheckAliCloudApigServicesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_apig_service.default.service_name}_fake"`,
		}),
	}

	gatewayIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudApigServicesDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_apig_service.default.id}"]`,
			"gateway_id": `"${alicloud_apig_service.default.gateway_id}"`,
		}),
		FakeConfig: testAccCheckAliCloudApigServicesDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_apig_service.default.id}_fake"]`,
			"gateway_id": `"${alicloud_apig_service.default.gateway_id}"`,
		}),
	}

	sourceTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudApigServicesDataSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_apig_service.default.id}"]`,
			"source_type": `"${alicloud_apig_service.default.source_type}"`,
		}),
		FakeConfig: testAccCheckAliCloudApigServicesDataSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_apig_service.default.id}_fake"]`,
			"source_type": `"${alicloud_apig_service.default.source_type}"`,
		}),
	}

	resourceGroupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudApigServicesDataSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_apig_service.default.id}"]`,
			"resource_group_id": `"${alicloud_apig_service.default.resource_group_id}"`,
		}),
		FakeConfig: testAccCheckAliCloudApigServicesDataSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_apig_service.default.id}_fake"]`,
			"resource_group_id": `"${alicloud_apig_service.default.resource_group_id}"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudApigServicesDataSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_apig_service.default.id}"]`,
			"name_regex":        `"${alicloud_apig_service.default.service_name}"`,
			"gateway_id":        `"${alicloud_apig_service.default.gateway_id}"`,
			"source_type":       `"${alicloud_apig_service.default.source_type}"`,
			"resource_group_id": `"${alicloud_apig_service.default.resource_group_id}"`,
		}),
		FakeConfig: testAccCheckAliCloudApigServicesDataSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_apig_service.default.id}_fake"]`,
			"name_regex":        `"${alicloud_apig_service.default.service_name}_fake"`,
			"gateway_id":        `"${alicloud_apig_service.default.gateway_id}"`,
			"source_type":       `"${alicloud_apig_service.default.source_type}"`,
			"resource_group_id": `"${alicloud_apig_service.default.resource_group_id}"`,
		}),
	}

	AliCloudApigServicesCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, gatewayIdConf, sourceTypeConf, resourceGroupConf, allConf)
}

var existAliCloudApigServicesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"services.#":                   "1",
		"services.0.id":                acctest.CHECKSET,
		"services.0.service_id":        acctest.CHECKSET,
		"services.0.service_name":      acctest.CHECKSET,
		"services.0.source_type":       "VIP",
		"services.0.gateway_id":        acctest.CHECKSET,
		"services.0.resource_group_id": acctest.CHECKSET,
		"services.0.create_timestamp":  acctest.CHECKSET,
		"ids.#":                        "1",
		"names.#":                      "1",
	}
}

var fakeAliCloudApigServicesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"services.#": "0",
		"ids.#":      "0",
		"names.#":    "0",
	}
}

var AliCloudApigServicesCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_apig_services.default",
	ExistMapFunc: existAliCloudApigServicesMapFunc,
	FakeMapFunc:  fakeAliCloudApigServicesMapFunc,
}

func testAccCheckAliCloudApigServicesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testacc-apig-service-ds-%d"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "default" {
  gateway_name = var.name
  spec         = "apigw.small.x1"
  payment_type = "PayAsYouGo"
  vpc {
    vpc_id = data.alicloud_vpcs.default.ids.0
  }
  vswitch {
    vswitch_id = data.alicloud_vswitches.default.ids.0
  }
  zone_config {
    select_option = "Auto"
  }
  network_access_config {
    type = "Intranet"
  }
  log_config {
    sls {
      enable = "false"
    }
  }
}

resource "alicloud_apig_service" "default" {
  service_name = var.name
  source_type  = "VIP"
  gateway_id   = alicloud_apig_gateway.default.id
  addresses    = ["127.0.0.1:8080"]
}

data "alicloud_apig_services" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
