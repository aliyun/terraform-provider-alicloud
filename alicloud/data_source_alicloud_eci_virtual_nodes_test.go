package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEciVirtualNodesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_eci_virtual_node.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_eci_virtual_node.default.id}_fake"]`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_eci_virtual_node.default.id}"]`,
			"resource_group_id": `"${alicloud_eci_virtual_node.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_eci_virtual_node.default.id}"]`,
			"resource_group_id": `"${alicloud_eci_virtual_node.default.resource_group_id}_fake"`,
		}),
	}
	securityGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_eci_virtual_node.default.id}"]`,
			"security_group_id": `"${alicloud_eci_virtual_node.default.security_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_eci_virtual_node.default.id}"]`,
			"security_group_id": `"${alicloud_eci_virtual_node.default.security_group_id}_fake"`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_eci_virtual_node.default.id}"]`,
			"tags": `{Created = "TF"}`,
		}),
		fakeConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_eci_virtual_node.default.id}"]`,
			"tags": `{Created = "TF_fake"}`,
		}),
	}
	vSwitchIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_eci_virtual_node.default.id}"]`,
			"vswitch_id": `"${alicloud_eci_virtual_node.default.vswitch_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_eci_virtual_node.default.id}"]`,
			"vswitch_id": `"${alicloud_eci_virtual_node.default.vswitch_id}_fake"`,
		}),
	}
	virtualNodeNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_eci_virtual_node.default.id}"]`,
			"virtual_node_name": `"${alicloud_eci_virtual_node.default.virtual_node_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_eci_virtual_node.default.id}"]`,
			"virtual_node_name": `"${alicloud_eci_virtual_node.default.virtual_node_name}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_eci_virtual_node.default.virtual_node_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_eci_virtual_node.default.virtual_node_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_eci_virtual_node.default.id}"]`,
			"status": `"Ready"`,
		}),
		fakeConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_eci_virtual_node.default.id}"]`,
			"status": `"Pending"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_eci_virtual_node.default.id}"]`,
			"name_regex":        `"${alicloud_eci_virtual_node.default.virtual_node_name}"`,
			"resource_group_id": `"${alicloud_eci_virtual_node.default.resource_group_id}"`,
			"security_group_id": `"${alicloud_eci_virtual_node.default.security_group_id}"`,
			"status":            `"Ready"`,
			"tags":              `{Created = "TF"}`,
			"virtual_node_name": `"${alicloud_eci_virtual_node.default.virtual_node_name}"`,
			"vswitch_id":        `"${alicloud_eci_virtual_node.default.vswitch_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_eci_virtual_node.default.id}_fake"]`,
			"name_regex":        `"${alicloud_eci_virtual_node.default.virtual_node_name}_fake"`,
			"resource_group_id": `"${alicloud_eci_virtual_node.default.resource_group_id}_fake"`,
			"security_group_id": `"${alicloud_eci_virtual_node.default.security_group_id}_fake"`,
			"status":            `"Pending"`,
			"tags":              `{Created = "TF_fake"}`,
			"virtual_node_name": `"${alicloud_eci_virtual_node.default.virtual_node_name}_fake"`,
			"vswitch_id":        `"${alicloud_eci_virtual_node.default.vswitch_id}_fake"`,
		}),
	}
	var existAlicloudEciVirtualNodesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"names.#":                   "1",
			"nodes.#":                   "1",
			"nodes.0.cpu":               CHECKSET,
			"nodes.0.create_time":       CHECKSET,
			"nodes.0.events.#":          CHECKSET,
			"nodes.0.internet_ip":       CHECKSET,
			"nodes.0.intranet_ip":       CHECKSET,
			"nodes.0.memory":            CHECKSET,
			"nodes.0.resource_group_id": CHECKSET,
			"nodes.0.security_group_id": CHECKSET,
			"nodes.0.status":            CHECKSET,
			"nodes.0.tags.%":            "1",
			"nodes.0.tags.Created":      "TF",
			"nodes.0.vswitch_id":        CHECKSET,
			"nodes.0.id":                CHECKSET,
			"nodes.0.virtual_node_id":   CHECKSET,
			"nodes.0.virtual_node_name": CHECKSET,
			"nodes.0.vpc_id":            CHECKSET,
			"nodes.0.zone_id":           CHECKSET,
		}
	}
	var fakeAlicloudEciVirtualNodesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudEciVirtualNodesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_eci_virtual_nodes.default",
		existMapFunc: existAlicloudEciVirtualNodesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEciVirtualNodesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithEnvVariable(t, "KUBE_CONFIG")
	}
	alicloudEciVirtualNodesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, resourceGroupIdConf, securityGroupIdConf, tagsConf, virtualNodeNameConf, vSwitchIdConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudEciVirtualNodesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testaccvirtualnode-%d"
}

variable "kube_config" {
  default = "%s"
}

data "alicloud_eci_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_eci_zones.default.zones.0.zone_ids.1
}

resource "alicloud_security_group" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  name   = var.name
}

resource "alicloud_eip_address" "default" {
  address_name = var.name
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_eci_virtual_node" "default" {
  security_group_id = alicloud_security_group.default.id
  virtual_node_name = var.name
  vswitch_id        = data.alicloud_vswitches.default.ids.1
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  enable_public_network = false
  eip_instance_id       = alicloud_eip_address.default.id
  kube_config           = var.kube_config
  tags = {
    Created = "TF"
  }
}

data "alicloud_eci_virtual_nodes" "default" {	
	%s
}
`, rand, os.Getenv("KUBE_CONFIG"), strings.Join(pairs, " \n "))
	return config
}
