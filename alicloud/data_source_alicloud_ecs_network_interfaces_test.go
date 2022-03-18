package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECSNetworkInterfacesDataSource(t *testing.T) {
	resourceId := "data.alicloud_ecs_network_interfaces.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEcsNetworkInterfacesDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ecs_network_interface.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ecs_network_interface.default.name}-fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_network_interface.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_network_interface.default.id}-fake"},
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_network_interface.default.id}"},
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "Test",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_network_interface.default.id}"},
			"tags": map[string]interface{}{
				"Created": "TF-fake",
				"For":     "Test-fake",
			},
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_ecs_network_interface.default.id}"},
			"status": "Available",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_ecs_network_interface.default.id}"},
			"status": "Deleting",
		}),
	}
	vswitchIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ecs_network_interface.default.id}"},
			"vswitch_id": "${data.alicloud_vswitches.default.ids.0}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ecs_network_interface.default.id}"},
			"vswitch_id": "${data.alicloud_vswitches.default.ids.0}_fake",
		}),
	}
	privateIpConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ecs_network_interface.default.id}"},
			"private_ip": "${cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 100)}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ecs_network_interface.default.id}"},
			"private_ip": "${cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 101)}",
		}),
	}
	securityGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_ecs_network_interface.default.id}"},
			"security_group_id": "${alicloud_security_group.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_ecs_network_interface.default.id}"},
			"security_group_id": "${alicloud_security_group.default.id}_fake",
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_ecs_network_interface.default.id}"},
			"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_ecs_network_interface.default.id}"},
			"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake",
		}),
	}
	var existEcsNetworkInterfacesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                       "1",
			"ids.0":                                       CHECKSET,
			"names.#":                                     "1",
			"names.0":                                     name,
			"interfaces.#":                                "1",
			"interfaces.0.description":                    CHECKSET,
			"interfaces.0.creation_time":                  CHECKSET,
			"interfaces.0.instance_id":                    "",
			"interfaces.0.mac":                            CHECKSET,
			"interfaces.0.id":                             CHECKSET,
			"interfaces.0.network_interface_id":           CHECKSET,
			"interfaces.0.network_interface_name":         name,
			"interfaces.0.name":                           name,
			"interfaces.0.primary_ip_address":             CHECKSET,
			"interfaces.0.private_ip":                     CHECKSET,
			"interfaces.0.private_ip_addresses.#":         "0",
			"interfaces.0.private_ips.#":                  "0",
			"interfaces.0.queue_number":                   CHECKSET,
			"interfaces.0.resource_group_id":              CHECKSET,
			"interfaces.0.security_group_ids.#":           "1",
			"interfaces.0.security_groups.#":              "1",
			"interfaces.0.status":                         CHECKSET,
			"interfaces.0.tags.%":                         "2",
			"interfaces.0.tags.Created":                   "TF",
			"interfaces.0.tags.For":                       "Test",
			"interfaces.0.type":                           CHECKSET,
			"interfaces.0.vswitch_id":                     CHECKSET,
			"interfaces.0.vpc_id":                         CHECKSET,
			"interfaces.0.owner_id":                       CHECKSET,
			"interfaces.0.network_interface_traffic_mode": CHECKSET,
		}
	}

	var fakeEcsNetworkInterfacesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        "0",
			"names.#":      "0",
			"interfaces.#": "0",
		}
	}

	var EcsNetworkInterfacesInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEcsNetworkInterfacesMapFunc,
		fakeMapFunc:  fakeEcsNetworkInterfacesMapFunc,
	}

	EcsNetworkInterfacesInfo.dataSourceTestCheck(t, 0, nameRegexConf, idsConf, tagsConf, statusConf, vswitchIdConf, privateIpConf, securityGroupIdConf, resourceGroupIdConf)
}

func dataSourceEcsNetworkInterfacesDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
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
data "alicloud_resource_manager_resource_groups" "default"{
	status = "OK"
}

resource "alicloud_ecs_network_interface" "default" {
    network_interface_name = var.name
    vswitch_id = data.alicloud_vswitches.default.ids.0
    security_group_ids = [alicloud_security_group.default.id]
	description = "Basic test"
	primary_ip_address = cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 100)
	tags = {
		Created = "TF",
		For =    "Test",
	}
	resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}

`, name)
}
