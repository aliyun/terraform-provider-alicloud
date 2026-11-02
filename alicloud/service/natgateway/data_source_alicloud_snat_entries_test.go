package natgateway_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudVPCSnatEntriesDataSourceBasic(t *testing.T) {
	rand := tfacctest.RandInt()
	snatIpConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSnatEntriesBasicConfig(rand, map[string]string{
			"snat_table_id": `"${alicloud_snat_entry.default.snat_table_id}"`,
			"snat_ip":       `"${alicloud_snat_entry.default.snat_ip}"`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSnatEntriesBasicConfig(rand, map[string]string{
			"snat_table_id": `"${alicloud_snat_entry.default.snat_table_id}"`,
			"name_regex":    `"${alicloud_snat_entry.default.snat_entry_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudSnatEntriesBasicConfig(rand, map[string]string{
			"snat_table_id": `"${alicloud_snat_entry.default.snat_table_id}"`,
			"name_regex":    `"${alicloud_snat_entry.default.snat_entry_name}_fakeii"`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSnatEntriesBasicConfig(rand, map[string]string{
			"snat_table_id": `"${alicloud_snat_entry.default.snat_table_id}"`,
			"status":        `"Available"`,
		}),
		FakeConfig: testAccCheckAlicloudSnatEntriesBasicConfig(rand, map[string]string{
			"snat_table_id": `"${alicloud_snat_entry.default.snat_table_id}"`,
			"status":        `"Deleting"`,
		}),
	}

	idsConfig := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSnatEntriesBasicConfig(rand, map[string]string{
			"snat_table_id": `"${alicloud_snat_entry.default.snat_table_id}"`,
			"ids":           `["${alicloud_snat_entry.default.snat_entry_id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudSnatEntriesBasicConfig(rand, map[string]string{
			"snat_table_id": `"${alicloud_snat_entry.default.snat_table_id}"`,
			"ids":           `["${alicloud_snat_entry.default.snat_entry_id}_fake"]`,
		}),
	}

	sourceCidrConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSnatEntriesBasicConfig(rand, map[string]string{
			"snat_table_id": `"${alicloud_snat_entry.default.snat_table_id}"`,
			"snat_ip":       `"${alicloud_snat_entry.default.snat_ip}"`,
			"source_cidr":   `"172.16.0.0/21"`,
		}),
		FakeConfig: testAccCheckAlicloudSnatEntriesBasicConfig(rand, map[string]string{
			"snat_table_id": `"${alicloud_snat_entry.default.snat_table_id}"`,
			"snat_ip":       `"${alicloud_snat_entry.default.snat_ip}"`,
			"source_cidr":   `"172.16.0.0/20"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSnatEntriesBasicConfig(rand, map[string]string{
			"snat_table_id": `"${alicloud_snat_entry.default.snat_table_id}"`,
			"snat_ip":       `"${alicloud_snat_entry.default.snat_ip}"`,
			"source_cidr":   `"172.16.0.0/21"`,
			"name_regex":    `"${alicloud_snat_entry.default.snat_entry_name}"`,
			"status":        `"Available"`,
			"ids":           `["${alicloud_snat_entry.default.snat_entry_id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudSnatEntriesBasicConfig(rand, map[string]string{
			"snat_table_id": `"${alicloud_snat_entry.default.snat_table_id}"`,
			"source_cidr":   `"172.16.0.0/21"`,
			"name_regex":    `"${alicloud_snat_entry.default.snat_entry_name}_fake"`,
			"status":        `"Deleting"`,
			"ids":           `["${alicloud_snat_entry.default.snat_entry_id}_fake"]`,
		}),
	}

	snatEntriesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, natGatewayDataSourcePreCheck(t), snatIpConf, nameRegexConf, statusConf, idsConfig, sourceCidrConf, allConf)

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
	vpc_name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	zone_id = "${data.alicloud_zones.default.zones.0.id}"
	vswitch_name = "${var.name}"
}

resource "alicloud_nat_gateway" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	nat_gateway_name = "${var.name}"
    vswitch_id    = alicloud_vswitch.default.id
    nat_type      = "Enhanced"
}

resource "alicloud_eip_address" "default" {
	address_name = "${var.name}"
}

resource "alicloud_eip_association" "default" {
	allocation_id = "${alicloud_eip_address.default.id}"
	instance_id = "${alicloud_nat_gateway.default.id}"
}

resource "alicloud_snat_entry" "default" {
	snat_table_id = "${alicloud_nat_gateway.default.snat_table_ids}"
	source_vswitch_id = "${alicloud_vswitch.default.id}"
	snat_ip = "${alicloud_eip_address.default.ip_address}"
   snat_entry_name = "${var.name}"
}

data "alicloud_snat_entries" "default" {
    %s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existSnatEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                       "1",
		"names.#":                     "1",
		"entries.#":                   "1",
		"entries.0.id":                acctest.CHECKSET,
		"entries.0.snat_ip":           acctest.CHECKSET,
		"entries.0.status":            "Available",
		"entries.0.source_cidr":       "172.16.0.0/21",
		"entries.0.snat_entry_id":     acctest.CHECKSET,
		"entries.0.snat_entry_name":   fmt.Sprintf("tf-testAccForSnatEntriesDatasource%d", rand),
		"entries.0.source_vswitch_id": acctest.CHECKSET,
	}
}

var fakeSnatEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":     "0",
		"entries.#": "0",
		"names.#":   "0",
	}
}

var snatEntriesCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_snat_entries.default",
	ExistMapFunc: existSnatEntriesMapFunc,
	FakeMapFunc:  fakeSnatEntriesMapFunc,
}
