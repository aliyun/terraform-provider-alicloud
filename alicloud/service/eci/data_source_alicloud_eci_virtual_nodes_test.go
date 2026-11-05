package eci_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudEciVirtualNodesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_eci_virtual_node.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_eci_virtual_node.default.id}_fake"]`,
		}),
	}
	resourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_eci_virtual_node.default.id}"]`,
			"resource_group_id": `"${alicloud_eci_virtual_node.default.resource_group_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_eci_virtual_node.default.id}"]`,
			"resource_group_id": `"${alicloud_eci_virtual_node.default.resource_group_id}_fake"`,
		}),
	}
	securityGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_eci_virtual_node.default.id}"]`,
			"security_group_id": `"${alicloud_eci_virtual_node.default.security_group_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_eci_virtual_node.default.id}"]`,
			"security_group_id": `"${alicloud_eci_virtual_node.default.security_group_id}_fake"`,
		}),
	}
	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_eci_virtual_node.default.id}"]`,
			"tags": `{Created = "TF"}`,
		}),
		FakeConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_eci_virtual_node.default.id}"]`,
			"tags": `{Created = "TF_fake"}`,
		}),
	}
	vSwitchIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_eci_virtual_node.default.id}"]`,
			"vswitch_id": `"${alicloud_eci_virtual_node.default.vswitch_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_eci_virtual_node.default.id}"]`,
			"vswitch_id": `"${alicloud_eci_virtual_node.default.vswitch_id}_fake"`,
		}),
	}
	virtualNodeNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_eci_virtual_node.default.id}"]`,
			"virtual_node_name": `"${alicloud_eci_virtual_node.default.virtual_node_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_eci_virtual_node.default.id}"]`,
			"virtual_node_name": `"${alicloud_eci_virtual_node.default.virtual_node_name}_fake"`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_eci_virtual_node.default.virtual_node_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_eci_virtual_node.default.virtual_node_name}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_eci_virtual_node.default.id}"]`,
			"status": `"Ready"`,
		}),
		FakeConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_eci_virtual_node.default.id}"]`,
			"status": `"Pending"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_eci_virtual_node.default.id}"]`,
			"name_regex":        `"${alicloud_eci_virtual_node.default.virtual_node_name}"`,
			"resource_group_id": `"${alicloud_eci_virtual_node.default.resource_group_id}"`,
			"security_group_id": `"${alicloud_eci_virtual_node.default.security_group_id}"`,
			"status":            `"Ready"`,
			"tags":              `{Created = "TF"}`,
			"virtual_node_name": `"${alicloud_eci_virtual_node.default.virtual_node_name}"`,
			"vswitch_id":        `"${alicloud_eci_virtual_node.default.vswitch_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudEciVirtualNodesDataSourceName(rand, map[string]string{
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
			"nodes.0.cpu":               acctest.CHECKSET,
			"nodes.0.create_time":       acctest.CHECKSET,
			"nodes.0.events.#":          acctest.CHECKSET,
			"nodes.0.internet_ip":       acctest.CHECKSET,
			"nodes.0.intranet_ip":       acctest.CHECKSET,
			"nodes.0.memory":            acctest.CHECKSET,
			"nodes.0.resource_group_id": acctest.CHECKSET,
			"nodes.0.security_group_id": acctest.CHECKSET,
			"nodes.0.status":            acctest.CHECKSET,
			"nodes.0.tags.%":            "1",
			"nodes.0.tags.Created":      "TF",
			"nodes.0.vswitch_id":        acctest.CHECKSET,
			"nodes.0.id":                acctest.CHECKSET,
			"nodes.0.virtual_node_id":   acctest.CHECKSET,
			"nodes.0.virtual_node_name": acctest.CHECKSET,
			"nodes.0.vpc_id":            acctest.CHECKSET,
			"nodes.0.zone_id":           acctest.CHECKSET,
		}
	}
	var fakeAlicloudEciVirtualNodesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudEciVirtualNodesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_eci_virtual_nodes.default",
		ExistMapFunc: existAlicloudEciVirtualNodesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudEciVirtualNodesDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithEnvVariable(t, "KUBE_CONFIG")
	}
	alicloudEciVirtualNodesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, resourceGroupIdConf, securityGroupIdConf, tagsConf, virtualNodeNameConf, vSwitchIdConf, nameRegexConf, statusConf, allConf)
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
    name_regex = "^default-NODELETING$"
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
