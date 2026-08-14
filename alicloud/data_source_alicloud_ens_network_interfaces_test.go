package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEnsNetworkInterfacesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsNetworkInterfacesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ens_network_interface.default.network_interface_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsNetworkInterfacesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ens_network_interface.default.network_interface_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsNetworkInterfacesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ens_network_interface.default.network_interface_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsNetworkInterfacesSourceConfig(rand, map[string]string{
			"name_regex": `"TestAccAlicloudEnsNetworkInterfacesDataSource_fake"`,
		}),
	}

	EnsNetworkInterfaceCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, allConf)
}

func testAccCheckAlicloudEnsNetworkInterfacesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccEnsNetworkInterfacesDataSource%d"
}

variable "ens_region_id" {
  default = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_security_group" "default" {
  security_group_name = "tf-testAccEnsNetworkInterfacesDataSource%d"
}

resource "alicloud_ens_network" "default" {
  network_name  = "tf-testAccEnsNetworkInterfacesDataSource%d"
  cidr_block    = "192.168.2.0/24"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "default" {
  cidr_block    = "192.168.2.0/24"
  vswitch_name  = "tf-testAccEnsNetworkInterfacesDataSource%d"
  ens_region_id = var.ens_region_id
  network_id    = alicloud_ens_network.default.id
}

resource "alicloud_ens_network_interface" "default" {
  network_interface_name = "tf-testAccEnsNetworkInterfacesDataSource%d"
  security_group_ids     = [alicloud_ens_security_group.default.id]
  vswitch_id             = alicloud_ens_vswitch.default.id
}

data "alicloud_ens_network_interfaces" "default" {
  enable_details = true
  %s
}
`, rand, rand, rand, rand, rand, strings.Join(pairs, "\n  "))
	return config
}

var existEnsNetworkInterfaceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"interfaces.#":                        "1",
		"interfaces.0.network_interface_name": fmt.Sprintf("tf-testAccEnsNetworkInterfacesDataSource%d", rand),
	}
}

var fakeEnsNetworkInterfaceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"interfaces.#": "0",
	}
}

var EnsNetworkInterfaceCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ens_network_interfaces.default",
	existMapFunc: existEnsNetworkInterfaceMapFunc,
	fakeMapFunc:  fakeEnsNetworkInterfaceMapFunc,
}
