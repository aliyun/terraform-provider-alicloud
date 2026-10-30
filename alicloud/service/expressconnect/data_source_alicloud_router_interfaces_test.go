package expressconnect_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudVPCRouterInterfacesDataSourceBasic(t *testing.T) {
	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
	}
	rand := tfacctest.RandIntRange(1000, 9999)

	oppositeInterfaceIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRouterInterfacesDataSourceConfig(rand, map[string]string{
			"opposite_interface_id": `"${alicloud_router_interface_connection.foo.opposite_interface_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudRouterInterfacesDataSourceConfig(rand, map[string]string{
			"opposite_interface_id": `"${alicloud_router_interface_connection.foo.opposite_interface_id}_fake"`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRouterInterfacesDataSourceConfig(rand, map[string]string{
			"opposite_interface_id": `"${alicloud_router_interface_connection.foo.opposite_interface_id}"`,
			"status":                `"Active"`,
		}),
		FakeConfig: testAccCheckAlicloudRouterInterfacesDataSourceConfig(rand, map[string]string{
			"opposite_interface_id": `"${alicloud_router_interface_connection.foo.opposite_interface_id}"`,
			"status":                `"Inactive"`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRouterInterfacesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}_initiating"`,
		}),
		FakeConfig: testAccCheckAlicloudRouterInterfacesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}_fake"`,
		}),
	}

	specificationConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRouterInterfacesDataSourceConfig(rand, map[string]string{
			"opposite_interface_id": `"${alicloud_router_interface_connection.foo.opposite_interface_id}"`,
			"specification":         `"Large.2"`,
		}),
		FakeConfig: testAccCheckAlicloudRouterInterfacesDataSourceConfig(rand, map[string]string{
			"opposite_interface_id": `"${alicloud_router_interface_connection.foo.opposite_interface_id}"`,
			"specification":         `"Large.1"`,
		}),
	}

	routerIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRouterInterfacesDataSourceConfig(rand, map[string]string{
			"router_id": `"${alicloud_vpc.default.0.router_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudRouterInterfacesDataSourceConfig(rand, map[string]string{
			"router_id": `"${alicloud_vpc.default.0.router_id}_fake"`,
		}),
	}

	routerTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRouterInterfacesDataSourceConfig(rand, map[string]string{
			"opposite_interface_id": `"${alicloud_router_interface_connection.foo.opposite_interface_id}"`,
			"router_type":           `"VRouter"`,
		}),
		FakeConfig: testAccCheckAlicloudRouterInterfacesDataSourceConfig(rand, map[string]string{
			"opposite_interface_id": `"${alicloud_router_interface_connection.foo.opposite_interface_id}"`,
			"router_type":           `"VBR"`,
		}),
	}

	roleConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRouterInterfacesDataSourceConfig(rand, map[string]string{
			"opposite_interface_id": `"${alicloud_router_interface_connection.foo.opposite_interface_id}"`,
			"role":                  `"InitiatingSide"`,
		}),
		FakeConfig: testAccCheckAlicloudRouterInterfacesDataSourceConfig(rand, map[string]string{
			"opposite_interface_id": `"${alicloud_router_interface_connection.foo.opposite_interface_id}"`,
			"role":                  `"AcceptingSide"`,
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRouterInterfacesDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alicloud_router_interface.initiating.id}" ]`,
		}),
		FakeConfig: testAccCheckAlicloudRouterInterfacesDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alicloud_router_interface.initiating.id}_fake" ]`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRouterInterfacesDataSourceConfig(rand, map[string]string{
			"opposite_interface_id": `"${alicloud_router_interface_connection.foo.opposite_interface_id}"`,
			"role":                  `"InitiatingSide"`,
			"name_regex":            `"${var.name}_initiating"`,
			"specification":         `"Large.2"`,
			"router_id":             `"${alicloud_vpc.default.0.router_id}"`,
			"router_type":           `"VRouter"`,
			"ids":                   `[ "${alicloud_router_interface.initiating.id}" ]`,
		}),
		FakeConfig: testAccCheckAlicloudRouterInterfacesDataSourceConfig(rand, map[string]string{
			"opposite_interface_id": `"${alicloud_router_interface_connection.foo.opposite_interface_id}"`,
			"role":                  `"AcceptingSide"`,
			"name_regex":            `"${var.name}_initiating"`,
			"specification":         `"Large.2"`,
			"router_id":             `"${alicloud_vpc.default.0.router_id}"`,
			"router_type":           `"VRouter"`,
			"ids":                   `[ "${alicloud_router_interface.initiating.id}" ]`,
		}),
	}

	routerInterfacesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, oppositeInterfaceIdConf, statusConf, nameRegexConf, specificationConf, routerIdConf,
		routerTypeConf, roleConf, idsConf, allConf)

}

func testAccCheckAlicloudRouterInterfacesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCheckAlicloudRouterInterfacesDataSourceConfig%d"
}

variable cidr_block_list {
	type = list(string)
	default = [ "172.16.0.0/12", "192.168.0.0/16" ]
}

resource "alicloud_vpc" "default" {
  count = 2
  name = "${var.name}"
  cidr_block = "${element(var.cidr_block_list,count.index)}"
}

data "alicloud_regions" "current_regions" {
  current = true
}
resource "alicloud_router_interface" "initiating" {
  opposite_region = "${data.alicloud_regions.current_regions.regions.0.id}"
  router_type = "VRouter"
  router_id = "${alicloud_vpc.default.0.router_id}"
  role = "InitiatingSide"
  specification = "Large.2"
  name = "${var.name}_initiating"
  description = "${var.name}_decription"
  depends_on = [ "alicloud_vpc.default" ]
}
resource "alicloud_router_interface" "accepting" {
  opposite_region = "${data.alicloud_regions.current_regions.regions.0.id}"
  router_type = "VRouter"
  router_id = "${alicloud_vpc.default.1.router_id}"
  role = "AcceptingSide"
  specification = "Negative"
  name = "${var.name}_accepting"
  description = "${var.name}_decription"
  depends_on = [ "alicloud_vpc.default" ]
}
resource "alicloud_router_interface_connection" "bar" {
  interface_id = "${alicloud_router_interface.accepting.id}"
  opposite_interface_id = "${alicloud_router_interface.initiating.id}"
}
resource "alicloud_router_interface_connection" "foo" {
  interface_id = "${alicloud_router_interface.initiating.id}"
  opposite_interface_id = "${alicloud_router_interface.accepting.id}"
  depends_on = [
    "alicloud_router_interface_connection.bar"
  ]
}

data "alicloud_router_interfaces" "default" {
  %s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existRouterInterfacesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                    "1",
		"names.#":                                  "1",
		"interfaces.#":                             "1",
		"interfaces.0.id":                          acctest.CHECKSET,
		"interfaces.0.status":                      "Active",
		"interfaces.0.name":                        fmt.Sprintf("tf-testAccCheckAlicloudRouterInterfacesDataSourceConfig%d_initiating", rand),
		"interfaces.0.description":                 fmt.Sprintf("tf-testAccCheckAlicloudRouterInterfacesDataSourceConfig%d_decription", rand),
		"interfaces.0.role":                        "InitiatingSide",
		"interfaces.0.specification":               "Large.2",
		"interfaces.0.router_id":                   acctest.CHECKSET,
		"interfaces.0.router_type":                 "VRouter",
		"interfaces.0.vpc_id":                      acctest.CHECKSET,
		"interfaces.0.access_point_id":             "",
		"interfaces.0.creation_time":               acctest.CHECKSET,
		"interfaces.0.opposite_region_id":          acctest.CHECKSET,
		"interfaces.0.opposite_interface_id":       acctest.CHECKSET,
		"interfaces.0.opposite_router_id":          acctest.CHECKSET,
		"interfaces.0.opposite_router_type":        "VRouter",
		"interfaces.0.opposite_interface_owner_id": acctest.CHECKSET,
		"interfaces.0.health_check_source_ip":      "",
		"interfaces.0.health_check_target_ip":      "",
	}
}

var fakeRouterInterfacesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":        "0",
		"names.#":      "0",
		"interfaces.#": "0",
	}
}

var routerInterfacesCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_router_interfaces.default",
	ExistMapFunc: existRouterInterfacesMapFunc,
	FakeMapFunc:  fakeRouterInterfacesMapFunc,
}
