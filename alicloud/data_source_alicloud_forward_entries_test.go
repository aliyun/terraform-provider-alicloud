package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVPCForwardEntriesDataSourceBasic(t *testing.T) {
	rand := acctest.RandInt()
	forwardTableIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudForwardEntriesDataSourceConfigBasic(rand, map[string]string{
			"forward_table_id": `"${alicloud_forward_entry.default.forward_table_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudForwardEntriesDataSourceConfigBasic(rand, map[string]string{
			"forward_table_id": `"${alicloud_forward_entry.default.forward_table_id}_fake"`,
		}),
	}

	externalIpConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudForwardEntriesDataSourceConfigBasic(rand, map[string]string{
			"forward_table_id": `"${alicloud_forward_entry.default.forward_table_id}"`,
			"external_ip":      `"${alicloud_forward_entry.default.external_ip}"`,
		}),
		fakeConfig: testAccCheckAlicloudForwardEntriesDataSourceConfigBasic(rand, map[string]string{
			"forward_table_id": `"${alicloud_forward_entry.default.forward_table_id}"`,
			"external_ip":      ` "${alicloud_forward_entry.default.external_ip}_fake" `,
		}),
	}

	internalIpConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudForwardEntriesDataSourceConfigBasic(rand, map[string]string{
			"internal_ip":      `"${alicloud_forward_entry.default.internal_ip}"`,
			"forward_table_id": `"${alicloud_forward_entry.default.forward_table_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudForwardEntriesDataSourceConfigBasic(rand, map[string]string{
			"internal_ip":      `"${alicloud_forward_entry.default.internal_ip}_fake"`,
			"forward_table_id": `"${alicloud_forward_entry.default.forward_table_id}"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudForwardEntriesDataSourceConfigBasic(rand, map[string]string{
			"forward_table_id": `"${alicloud_forward_entry.default.forward_table_id}"`,
			"ids":              `[ "${alicloud_forward_entry.default.forward_entry_id}" ]`,
		}),
		fakeConfig: testAccCheckAlicloudForwardEntriesDataSourceConfigBasic(rand, map[string]string{
			"forward_table_id": `"${alicloud_forward_entry.default.forward_table_id}"`,
			"ids":              `[ "${alicloud_forward_entry.default.forward_entry_id}_fake" ]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudForwardEntriesDataSourceConfigBasic(rand, map[string]string{
			"forward_table_id": `"${alicloud_forward_entry.default.forward_table_id}"`,
			"name_regex":       `"${alicloud_forward_entry.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudForwardEntriesDataSourceConfigBasic(rand, map[string]string{
			"forward_table_id": `"${alicloud_forward_entry.default.forward_table_id}"`,
			"name_regex":       `"${alicloud_forward_entry.default.name}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudForwardEntriesDataSourceConfigBasic(rand, map[string]string{
			"forward_table_id": `"${alicloud_forward_entry.default.forward_table_id}"`,
			"status":           `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudForwardEntriesDataSourceConfigBasic(rand, map[string]string{
			"forward_table_id": `"${alicloud_forward_entry.default.forward_table_id}"`,
			"status":           `"Deleting"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudForwardEntriesDataSourceConfigBasic(rand, map[string]string{
			"forward_table_id": `"${alicloud_forward_entry.default.forward_table_id}"`,
			"external_ip":      `"${alicloud_forward_entry.default.external_ip}"`,
			"internal_ip":      `"${alicloud_forward_entry.default.internal_ip}"`,
			"ids":              `[ "${alicloud_forward_entry.default.forward_entry_id}" ]`,
			"name_regex":       `"${alicloud_forward_entry.default.name}"`,
			"status":           `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudForwardEntriesDataSourceConfigBasic(rand, map[string]string{
			"forward_table_id": `"${alicloud_forward_entry.default.forward_table_id}"`,
			"external_ip":      `"${alicloud_forward_entry.default.external_ip}"`,
			"internal_ip":      `"${alicloud_forward_entry.default.internal_ip}"`,
			"ids":              `[ "${alicloud_forward_entry.default.forward_entry_id}_fake" ]`,
			"name_regex":       `"${alicloud_forward_entry.default.name}"`,
			"status":           `"Deleting"`,
		}),
	}
	forwardEntriesCheckInfo.dataSourceTestCheck(t, rand, forwardTableIdConf, externalIpConf, internalIpConf, idsConf, nameRegexConf, statusConf, allConf)

}

func testAccCheckAlicloudForwardEntriesDataSourceConfigBasic(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccForwardEntryConfig%d"
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
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
    internet_charge_type = "PayByLcu"
	nat_gateway_name = "${var.name}"
    nat_type = "Enhanced"
	vswitch_id = alicloud_vswitch.default.id
}

resource "alicloud_eip_address" "default" {
	address_name = "${var.name}"
}

resource "alicloud_eip_association" "default" {
	allocation_id = "${alicloud_eip_address.default.id}"
	instance_id = "${alicloud_nat_gateway.default.id}"
}

resource "alicloud_forward_entry" "default"{
	forward_table_id = "${alicloud_nat_gateway.default.forward_table_ids}"
	external_ip = "${alicloud_eip_address.default.ip_address}"
	external_port = "80"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.3"
	internal_port = "8080"
}

data "alicloud_forward_entries" "default" {
	%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existForwardEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                   "1",
		"names.#":                 "1",
		"entries.#":               "1",
		"entries.0.id":            CHECKSET,
		"entries.0.external_ip":   CHECKSET,
		"entries.0.external_port": "80",
		"entries.0.internal_ip":   "172.16.0.3",
		"entries.0.internal_port": "8080",
		"entries.0.ip_protocol":   "tcp",
		"entries.0.status":        "Available",
	}
}

var fakeForwardEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":     "0",
		"names.#":   "0",
		"entries.#": "0",
	}
}

var forwardEntriesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_forward_entries.default",
	existMapFunc: existForwardEntriesMapFunc,
	fakeMapFunc:  fakeForwardEntriesMapFunc,
}
