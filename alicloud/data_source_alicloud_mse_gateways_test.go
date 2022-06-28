package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMSEGatewaysDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.MSEGatewaySupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mse_gateway.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mse_gateway.default.id}_fake"]`,
		}),
	}
	gatewayNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_mse_gateway.default.id}"]`,
			"gateway_name": `"${alicloud_mse_gateway.default.gateway_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_mse_gateway.default.id}"]`,
			"gateway_name": `"${alicloud_mse_gateway.default.gateway_name}_fake"`,
		}),
	}
	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_mse_gateway.default.id}"]`,
			"vpc_id": `"${alicloud_mse_gateway.default.vpc_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_mse_gateway.default.id}"]`,
			"vpc_id": `"${alicloud_mse_gateway.default.vpc_id}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_mse_gateway.default.gateway_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_mse_gateway.default.gateway_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_mse_gateway.default.id}"]`,
			"status": `"2"`,
		}),
		fakeConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_mse_gateway.default.id}"]`,
			"status": `"8"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
			"gateway_name": `"${alicloud_mse_gateway.default.gateway_name}"`,
			"ids":          `["${alicloud_mse_gateway.default.id}"]`,
			"name_regex":   `"${alicloud_mse_gateway.default.gateway_name}"`,
			"status":       `"2"`,
			"vpc_id":       `"${alicloud_mse_gateway.default.vpc_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
			"gateway_name": `"${alicloud_mse_gateway.default.gateway_name}_fake"`,
			"ids":          `["${alicloud_mse_gateway.default.id}_fake"]`,
			"name_regex":   `"${alicloud_mse_gateway.default.gateway_name}_fake"`,
			"status":       `"8"`,
			"vpc_id":       `"${alicloud_mse_gateway.default.vpc_id}_fake"`,
		}),
	}
	var existAlicloudMseGatewaysDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                    "1",
			"names.#":                                  "1",
			"gateways.#":                               "1",
			"gateways.0.id":                            CHECKSET,
			"gateways.0.gateway_name":                  fmt.Sprintf("tf-testAccGateway-%d", rand),
			"gateways.0.replica":                       "2",
			"gateways.0.spec":                          "MSE_GTW_2_4_200_c",
			"gateways.0.vswitch_id":                    CHECKSET,
			"gateways.0.backup_vswitch_id":             CHECKSET,
			"gateways.0.vpc_id":                        CHECKSET,
			"gateways.0.status":                        "2",
			"gateways.0.gateway_unique_id":             CHECKSET,
			"gateways.0.payment_type":                  CHECKSET,
			"gateways.0.slb_list.#":                    "1",
			"gateways.0.slb_list.0.associate_id":       CHECKSET,
			"gateways.0.slb_list.0.slb_id":             CHECKSET,
			"gateways.0.slb_list.0.slb_ip":             CHECKSET,
			"gateways.0.slb_list.0.slb_port":           CHECKSET,
			"gateways.0.slb_list.0.type":               CHECKSET,
			"gateways.0.slb_list.0.gmt_create":         CHECKSET,
			"gateways.0.slb_list.0.gateway_slb_mode":   CHECKSET,
			"gateways.0.slb_list.0.gateway_slb_status": CHECKSET,
		}
	}
	var fakeAlicloudMseGatewaysDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudMseGatewaysCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_mse_gateways.default",
		existMapFunc: existAlicloudMseGatewaysDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudMseGatewaysDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudMseGatewaysCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, gatewayNameConf, vpcIdConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudMseGatewaysDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccGateway-%d"
}

data "alicloud_zones" "default" {
    available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_mse_gateway" "default" {
	gateway_name = var.name
	replica = 2
	spec = "MSE_GTW_2_4_200_c"
	vswitch_id = data.alicloud_vswitches.default.ids.0
	backup_vswitch_id = data.alicloud_vswitches.default.ids.1
	vpc_id = data.alicloud_vpcs.default.ids.0
}

data "alicloud_mse_gateways" "default" {	
	enable_details = true
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
