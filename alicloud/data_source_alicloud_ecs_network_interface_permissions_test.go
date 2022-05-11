package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEcsNetworkInterfacePermissionsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsNetworkInterfacePermissionsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_network_interface_permission.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsNetworkInterfacePermissionsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_network_interface_permission.default.id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsNetworkInterfacePermissionsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecs_network_interface_permission.default.id}"]`,
			"status": `"Granted"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsNetworkInterfacePermissionsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecs_network_interface_permission.default.id}"]`,
			"status": `"Pending"`,
		}),
	}
	var existAlicloudEcsNetworkInterfacePermissionsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                         "1",
			"total_count":                                   "1",
			"permissions.#":                                 "1",
			"permissions.0.account_id":                      CHECKSET,
			"permissions.0.network_interface_id":            CHECKSET,
			"permissions.0.permission":                      "InstanceAttach",
			"permissions.0.status":                          "Granted",
			"permissions.0.service_name":                    "",
			"permissions.0.network_interface_permission_id": CHECKSET,
			"permissions.0.id":                              CHECKSET,
		}
	}
	var fakeAlicloudEcsNetworkInterfacePermissionsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudEcsNetworkInterfacePermissionsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecs_network_interface_permissions.default",
		existMapFunc: existAlicloudEcsNetworkInterfacePermissionsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcsNetworkInterfacePermissionsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudEcsNetworkInterfacePermissionsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf)
}
func testAccCheckAlicloudEcsNetworkInterfacePermissionsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccNetworkInterfacePermission-%d"
}

data "alicloud_zones" "default" {
    available_resource_creation= "VSwitch"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
    name = var.name
    vpc_id = data.alicloud_vpcs.default.ids.0
}
data "alicloud_resource_manager_resource_groups" "default"{}

resource "alicloud_ecs_network_interface" "default" {
    network_interface_name = var.name
    vswitch_id = data.alicloud_vswitches.default.ids.0
    security_group_ids = [alicloud_security_group.default.id]
	description = "Basic test"
	primary_ip_address = cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 16)
	tags = {
		Created = "TF",
		For =    "Test",
	}
	resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}
data "alicloud_account" "default" {}
resource "alicloud_ecs_network_interface_permission" "default" {
	account_id = data.alicloud_account.default.id
	network_interface_id = alicloud_ecs_network_interface.default.id
	permission = "InstanceAttach"
	force = true
}

data "alicloud_ecs_network_interface_permissions" "default" {	
	network_interface_id = alicloud_ecs_network_interface.default.id
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
