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
			"vswitch_id": "${alicloud_vswitch.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ecs_network_interface.default.id}"},
			"vswitch_id": "${alicloud_vswitch.default.id}_fake",
		}),
	}
	privateIpConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ecs_network_interface.default.id}"},
			"private_ip": "${cidrhost(alicloud_vswitch.default.cidr_block, 100)}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ecs_network_interface.default.id}"},
			"private_ip": "${cidrhost(alicloud_vswitch.default.cidr_block, 101)}",
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
	pageNumberConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vswitch_id":  "${alicloud_vswitch.default.id}",
			"page_number": 1,
			"page_size":   10,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vswitch_id":  "${alicloud_vswitch.default.id}",
			"page_number": 99,
			"page_size":   10,
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

	EcsNetworkInterfacesInfo.dataSourceTestCheck(t, 0, nameRegexConf, idsConf, tagsConf, statusConf, vswitchIdConf, privateIpConf, securityGroupIdConf, resourceGroupIdConf, pageNumberConf)
}

func dataSourceEcsNetworkInterfacesDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name    = var.name
  cidr_block  = "192.168.0.0/16"
  enable_ipv6 = true
}
resource "alicloud_vswitch" "default" {
  vpc_id               = alicloud_vpc.default.id
  zone_id              = data.alicloud_zones.default.zones.0.id
  cidr_block           = "192.168.64.0/24"
  vswitch_name         = var.name
  ipv6_cidr_block_mask = 64
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}
data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_ecs_network_interface" "default" {
  network_interface_name = var.name
  vswitch_id             = alicloud_vswitch.default.id
  security_group_ids     = [alicloud_security_group.default.id]
  description            = "Basic test"
  primary_ip_address     = cidrhost(alicloud_vswitch.default.cidr_block, 100)
  ipv6_address_count     = 1
  tags = {
    Created = "TF",
    For     = "Test",
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}

`, name)
}
