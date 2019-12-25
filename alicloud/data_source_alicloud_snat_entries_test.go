package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudSnatEntriesDataSourceBasic(t *testing.T) {
	rand := acctest.RandInt()
	snatIpConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSnatEntriesBasicConfig(rand, map[string]string{
			"snat_table_id": `alicloud_snat_entry.default.snat_table_id`,
			"snat_ip":       `alicloud_snat_entry.default.snat_ip`,
		}),
		fakeConfig: testAccCheckAlicloudSnatEntriesBasicConfig(rand, map[string]string{
			"snat_table_id": `alicloud_snat_entry.default.snat_table_id`,
			"snat_ip":       `"${alicloud_snat_entry.default.snat_ip}_fake"`,
		}),
	}

	sourceCidrConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSnatEntriesBasicConfig(rand, map[string]string{
			"snat_table_id": `alicloud_snat_entry.default.snat_table_id`,
			"snat_ip":       `alicloud_snat_entry.default.snat_ip`,
			"source_cidr":   `"172.16.0.0/21"`,
		}),
		fakeConfig: testAccCheckAlicloudSnatEntriesBasicConfig(rand, map[string]string{
			"snat_table_id": `alicloud_snat_entry.default.snat_table_id`,
			"snat_ip":       `alicloud_snat_entry.default.snat_ip`,
			"source_cidr":   `"172.16.0.0/20"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSnatEntriesBasicConfig(rand, map[string]string{
			"snat_table_id": `alicloud_snat_entry.default.snat_table_id`,
			"snat_ip":       `alicloud_snat_entry.default.snat_ip`,
			"source_cidr":   `"172.16.0.0/21"`,
		}),
		fakeConfig: testAccCheckAlicloudSnatEntriesBasicConfig(rand, map[string]string{
			"snat_table_id": `alicloud_snat_entry.default.snat_table_id`,
			"snat_ip":       `"${alicloud_snat_entry.default.snat_ip}_fake"`,
			"source_cidr":   `"172.16.0.0/21"`,
		}),
	}

	snatEntriesCheckInfo.dataSourceTestCheck(t, rand, snatIpConf, sourceCidrConf, allConf)

}

func testAccCheckAlicloudSnatEntriesBasicConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccForSnatEntriesDatasource%d"
}

data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
	name = var.name
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
	vpc_id = alicloud_vpc.default.id
	cidr_block = "172.16.0.0/21"
	availability_zone = data.alicloud_zones.default.zones.0.id
	name = var.name
}

resource "alicloud_nat_gateway" "default" {
	vpc_id = alicloud_vpc.default.id
	specification = "Small"
	name = var.name
}

resource "alicloud_eip" "default" {
	name = var.name
}

resource "alicloud_eip_association" "default" {
	allocation_id = alicloud_eip.default.id
	instance_id = alicloud_nat_gateway.default.id
}

resource "alicloud_snat_entry" "default" {
	snat_table_id = alicloud_nat_gateway.default.snat_table_ids
	source_vswitch_id = alicloud_vswitch.default.id
	snat_ip = alicloud_eip.default.ip_address
}

data "alicloud_snat_entries" "default" {
    %s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existSnatEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                 "1",
		"entries.#":             "1",
		"entries.0.id":          CHECKSET,
		"entries.0.snat_ip":     CHECKSET,
		"entries.0.status":      "Available",
		"entries.0.source_cidr": "172.16.0.0/21",
	}
}

var fakeSnatEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":     "0",
		"entries.#": "0",
	}
}

var snatEntriesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_snat_entries.default",
	existMapFunc: existSnatEntriesMapFunc,
	fakeMapFunc:  fakeSnatEntriesMapFunc,
}
