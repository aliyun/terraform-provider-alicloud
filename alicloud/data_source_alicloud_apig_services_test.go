package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudApigServicesDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApigServicesDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_service.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudApigServicesDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_service.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApigServicesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_apig_service.default.service_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudApigServicesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_apig_service.default.service_name}_fake"`,
		}),
	}

	gatewayIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApigServicesDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_apig_service.default.id}"]`,
			"gateway_id": `"${alicloud_apig_service.default.gateway_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudApigServicesDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_apig_service.default.id}"]`,
			"gateway_id": `"${alicloud_apig_service.default.gateway_id}_fake"`,
		}),
	}

	sourceTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApigServicesDataSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_apig_service.default.id}"]`,
			"gateway_id":  `"${alicloud_apig_service.default.gateway_id}"`,
			"source_type": `"${alicloud_apig_service.default.source_type}"`,
		}),
		fakeConfig: testAccCheckAliCloudApigServicesDataSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_apig_service.default.id}"]`,
			"gateway_id":  `"${alicloud_apig_service.default.gateway_id}"`,
			"source_type": `"K8S"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApigServicesDataSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_apig_service.default.id}"]`,
			"name_regex":  `"${alicloud_apig_service.default.service_name}"`,
			"gateway_id":  `"${alicloud_apig_service.default.gateway_id}"`,
			"source_type": `"${alicloud_apig_service.default.source_type}"`,
		}),
		fakeConfig: testAccCheckAliCloudApigServicesDataSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_apig_service.default.id}_fake"]`,
			"name_regex":  `"${alicloud_apig_service.default.service_name}_fake"`,
			"gateway_id":  `"${alicloud_apig_service.default.gateway_id}"`,
			"source_type": `"${alicloud_apig_service.default.source_type}"`,
		}),
	}

	AliCloudApigServicesCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, gatewayIdConf, sourceTypeConf, allConf)
}

var existAliCloudApigServicesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"services.#":              "1",
		"services.0.id":           CHECKSET,
		"services.0.service_id":   CHECKSET,
		"services.0.service_name": CHECKSET,
		"services.0.source_type":  "VIP",
		"services.0.gateway_id":   CHECKSET,
		"services.0.addresses.#":  CHECKSET,
		"ids.#":                   "1",
		"names.#":                 "1",
	}
}

var fakeAliCloudApigServicesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"services.#": "0",
		"ids.#":      "0",
		"names.#":    "0",
	}
}

var AliCloudApigServicesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_apig_services.default",
	existMapFunc: existAliCloudApigServicesMapFunc,
	fakeMapFunc:  fakeAliCloudApigServicesMapFunc,
}

func testAccCheckAliCloudApigServicesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tfaccapigsvcds%d"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "default" {
  network_access_config {
    type = "Intranet"
  }
  vswitch {
    vswitch_id = data.alicloud_vswitches.default.ids.0
  }
  zone_config {
    select_option = "Auto"
  }
  vpc {
    vpc_id = data.alicloud_vpcs.default.ids.0
  }
  payment_type = "PayAsYouGo"
  gateway_name = var.name
  spec         = "apigw.small.x1"
  log_config {
    sls {
      enable = false
    }
  }
}

resource "alicloud_apig_service" "default" {
  addresses    = ["127.0.0.1:8080"]
  service_name = var.name
  source_type  = "VIP"
  gateway_id   = alicloud_apig_gateway.default.id
  namespace    = "default"
}

data "alicloud_apig_services" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
