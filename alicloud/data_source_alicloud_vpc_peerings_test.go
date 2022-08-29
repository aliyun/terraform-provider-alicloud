package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVpcPeeringsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcPeeringsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_peering.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVpcPeeringsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_peering.default.id}_fake"]`,
		}),
	}
	peeringNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcPeeringsDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_vpc_peering.default.id}"]`,
			"peering_name": `"${alicloud_vpc_peering.default.peering_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcPeeringsDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_vpc_peering.default.id}"]`,
			"peering_name": `"${alicloud_vpc_peering.default.peering_name}_fake"`,
		}),
	}
	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcPeeringsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_peering.default.id}"]`,
			"vpc_id": `"${alicloud_vpc_peering.default.vpc_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcPeeringsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_peering.default.id}"]`,
			"vpc_id": `"${alicloud_vpc_peering.default.vpc_id}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcPeeringsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_vpc_peering.default.peering_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcPeeringsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_vpc_peering.default.peering_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcPeeringsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_peering.default.id}"]`,
			"status": `"${alicloud_vpc_peering.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcPeeringsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_peering.default.id}"]`,
			"status": `"Creating"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcPeeringsDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_vpc_peering.default.id}"]`,
			"name_regex":   `"${alicloud_vpc_peering.default.peering_name}"`,
			"peering_name": `"${alicloud_vpc_peering.default.peering_name}"`,
			"status":       `"${alicloud_vpc_peering.default.status}"`,
			"vpc_id":       `"${alicloud_vpc_peering.default.vpc_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcPeeringsDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_vpc_peering.default.id}_fake"]`,
			"name_regex":   `"${alicloud_vpc_peering.default.peering_name}_fake"`,
			"peering_name": `"${alicloud_vpc_peering.default.peering_name}_fake"`,
			"status":       `"Creating"`,
			"vpc_id":       `"${alicloud_vpc_peering.default.vpc_id}_fake"`,
		}),
	}
	var existAlicloudVpcPeeringsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"names.#":                        "1",
			"peerings.#":                     "1",
			"peerings.0.accepting_ali_uid":   CHECKSET,
			"peerings.0.accepting_region_id": defaultRegionToTest,
			"peerings.0.accepting_vpc_id":    CHECKSET,
			"peerings.0.bandwidth":           CHECKSET,
			"peerings.0.description":         fmt.Sprintf("tf-testAccPeering-%d", rand),
			"peerings.0.peering_name":        fmt.Sprintf("tf-testAccPeering-%d", rand),
			"peerings.0.vpc_id":              CHECKSET,
			"peerings.0.create_time":         CHECKSET,
			"peerings.0.id":                  CHECKSET,
			"peerings.0.peering_id":          CHECKSET,
			"peerings.0.status":              CHECKSET,
		}
	}
	var fakeAlicloudVpcPeeringsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudVpcPeeringsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_vpc_peerings.default",
		existMapFunc: existAlicloudVpcPeeringsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudVpcPeeringsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudVpcPeeringsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, peeringNameConf, vpcIdConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudVpcPeeringsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccPeering-%d"
}
data "alicloud_account" "default" {}

data "alicloud_vpcs" "default" {
	name_regex = "^default"
}

resource "alicloud_vpc_peering" "default" {
  peering_name        = var.name
  vpc_id              = data.alicloud_vpcs.default.ids.0
  accepting_ali_uid   = data.alicloud_account.default.id
  accepting_region_id = "%s"
  accepting_vpc_id    = data.alicloud_vpcs.default.ids.1
  description         = var.name
}

data "alicloud_vpc_peerings" "default" {
	%s	
}
`, rand, defaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
