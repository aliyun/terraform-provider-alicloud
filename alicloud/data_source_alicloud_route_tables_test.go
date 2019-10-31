package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudRouteTablesDataSourceBasic(t *testing.T) {
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithRegions(t, false, connectivity.RouteTableNoSupportedRegions)
	}
	rand := acctest.RandInt()

	nameRegexConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"name_regex": `"${alicloud_route_table.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"name_regex": `"${alicloud_route_table.default.name}_fake"`,
		}),
	}

	vpcIdConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"name_regex": `"${alicloud_route_table.default.name}"`,
			"vpc_id":     `"${alicloud_vpc.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"name_regex": `"${alicloud_route_table.default.name}"`,
			"vpc_id":     `"${alicloud_vpc.default.id}_fake"`,
		}),
	}

	idsConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"ids": `[ "${alicloud_route_table.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlicloudRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"ids": `[ "${alicloud_route_table.default.id}_fake" ]`,
		}),
	}

	tagsConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"name_regex": `"${alicloud_route_table.default.name}"`,
			"tags": `{
							Created = "TF"
							For 	= "acceptance test"
					  }`,
		}),
		fakeConfig: testAccCheckAlicloudRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"name_regex": `"${alicloud_route_table.default.name}"`,
			"tags": `{
							Created = "TF-fake"
							For 	= "acceptance test-fake"
					  }`,
		}),
	}

	resourceGroupIdConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"name_regex": `"${alicloud_route_table.default.name}"`,
			// The resource route tables do not support resource_group_id, so it was set empty.
			"resource_group_id": `""`,
		}),
		fakeConfig: testAccCheckAlicloudRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"name_regex":        `"${alicloud_route_table.default.name}"`,
			"resource_group_id": fmt.Sprintf(`"%s_fake"`, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID")),
		}),
	}

	allConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"name_regex":        `"${alicloud_route_table.default.name}"`,
			"vpc_id":            `"${alicloud_vpc.default.id}"`,
			"ids":               `[ "${alicloud_route_table.default.id}" ]`,
			"resource_group_id": `""`,
		}),
		fakeConfig: testAccCheckAlicloudRouteTablesDataSourceConfigBaisc(rand, map[string]string{
			"name_regex":        `"${alicloud_route_table.default.name}_fake"`,
			"vpc_id":            `"${alicloud_vpc.default.id}"`,
			"ids":               `[ "${alicloud_route_table.default.id}" ]`,
			"resource_group_id": `""`,
		}),
	}

	routeTablesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConfig, vpcIdConfig, idsConfig, tagsConfig, resourceGroupIdConfig, allConfig)
}

func testAccCheckAlicloudRouteTablesDataSourceConfigBaisc(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccRouteTablesDatasource%d"
}

resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

resource "alicloud_route_table" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  name = "${var.name}"
  description = "${var.name}_description"
  tags 		= {
		Created = "TF"
		For 	= "acceptance test"
  }
}

data "alicloud_route_tables" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existRouteTablesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                     "1",
		"names.#":                   "1",
		"tables.#":                  "1",
		"tables.0.id":               CHECKSET,
		"tables.0.route_table_type": CHECKSET,
		"tables.0.creation_time":    CHECKSET,
		"tables.0.router_id":        CHECKSET,
		"tables.0.name":             fmt.Sprintf("tf-testAccRouteTablesDatasource%d", rand),
		"tables.0.description":      fmt.Sprintf("tf-testAccRouteTablesDatasource%d_description", rand),
	}
}

var fakeRouteTablesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":    "0",
		"names.#":  "0",
		"tables.#": "0",
	}
}

var routeTablesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_route_tables.default",
	existMapFunc: existRouteTablesMapFunc,
	fakeMapFunc:  fakeRouteTablesMapFunc,
}
