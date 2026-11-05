package mse_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudMSEGatewaysDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.MSEGatewaySupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mse_gateway.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mse_gateway.default.id}_fake"]`,
		}),
	}
	gatewayNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_mse_gateway.default.id}"]`,
			"gateway_name": `"${alicloud_mse_gateway.default.gateway_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_mse_gateway.default.id}"]`,
			"gateway_name": `"${alicloud_mse_gateway.default.gateway_name}_fake"`,
		}),
	}
	vpcIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_mse_gateway.default.id}"]`,
			"vpc_id": `"${alicloud_mse_gateway.default.vpc_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_mse_gateway.default.id}"]`,
			"vpc_id": `"${alicloud_mse_gateway.default.vpc_id}_fake"`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_mse_gateway.default.gateway_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_mse_gateway.default.gateway_name}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_mse_gateway.default.id}"]`,
			"status": `"2"`,
		}),
		FakeConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_mse_gateway.default.id}"]`,
			"status": `"8"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
			"gateway_name": `"${alicloud_mse_gateway.default.gateway_name}"`,
			"ids":          `["${alicloud_mse_gateway.default.id}"]`,
			"name_regex":   `"${alicloud_mse_gateway.default.gateway_name}"`,
			"status":       `"2"`,
			"vpc_id":       `"${alicloud_mse_gateway.default.vpc_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudMseGatewaysDataSourceName(rand, map[string]string{
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
			"gateways.0.id":                            acctest.CHECKSET,
			"gateways.0.gateway_name":                  fmt.Sprintf("tf-testAccGateway-%d", rand),
			"gateways.0.replica":                       "2",
			"gateways.0.spec":                          "MSE_GTW_2_4_200_c",
			"gateways.0.vswitch_id":                    acctest.CHECKSET,
			"gateways.0.backup_vswitch_id":             acctest.CHECKSET,
			"gateways.0.vpc_id":                        acctest.CHECKSET,
			"gateways.0.status":                        "2",
			"gateways.0.gateway_unique_id":             acctest.CHECKSET,
			"gateways.0.payment_type":                  acctest.CHECKSET,
			"gateways.0.slb_list.#":                    "1",
			"gateways.0.slb_list.0.associate_id":       acctest.CHECKSET,
			"gateways.0.slb_list.0.slb_id":             acctest.CHECKSET,
			"gateways.0.slb_list.0.slb_ip":             acctest.CHECKSET,
			"gateways.0.slb_list.0.slb_port":           acctest.CHECKSET,
			"gateways.0.slb_list.0.type":               acctest.CHECKSET,
			"gateways.0.slb_list.0.gmt_create":         acctest.CHECKSET,
			"gateways.0.slb_list.0.gateway_slb_mode":   acctest.CHECKSET,
			"gateways.0.slb_list.0.gateway_slb_status": acctest.CHECKSET,
		}
	}
	var fakeAlicloudMseGatewaysDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudMseGatewaysCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_mse_gateways.default",
		ExistMapFunc: existAlicloudMseGatewaysDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudMseGatewaysDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudMseGatewaysCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, gatewayNameConf, vpcIdConf, nameRegexConf, statusConf, allConf)
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
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.17.0.0/16"
}
resource "alicloud_vswitch" "default" {
	vswitch_name = var.name
	cidr_block   = "172.17.3.0/24"
	vpc_id       = alicloud_vpc.default.id
	zone_id      = data.alicloud_zones.default.zones.0.id
}
resource "alicloud_vswitch" "backup" {
	vswitch_name = "${var.name}-backup"
	cidr_block   = "172.17.4.0/24"
	vpc_id       = alicloud_vpc.default.id
	zone_id      = data.alicloud_zones.default.zones.1.id
}

resource "alicloud_mse_gateway" "default" {
	gateway_name = var.name
	replica = 2
	spec = "MSE_GTW_2_4_200_c"
	vswitch_id = alicloud_vswitch.default.id
	backup_vswitch_id = alicloud_vswitch.backup.id
	vpc_id = alicloud_vpc.default.id
	enterprise_security_group = true
	slb_spec = "slb.s2.small"
	delete_slb = true
}

data "alicloud_mse_gateways" "default" {
	enable_details = true
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
