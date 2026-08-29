// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudEnsNatGatewaySnatEntryDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsNatGatewaySnatEntrySourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ens_nat_gateway_snat_entry.default.id}"]`,
			"nat_gateway_id": `"${alicloud_ens_nat_gateway.default2Kn0nu.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsNatGatewaySnatEntrySourceConfig(rand, map[string]string{
			"ids":            `["snat-fake0000"]`,
			"nat_gateway_id": `"${alicloud_ens_nat_gateway.default2Kn0nu.id}"`,
		}),
	}

	SnatEntryNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsNatGatewaySnatEntrySourceConfig(rand, map[string]string{
			"ids":             `["${alicloud_ens_nat_gateway_snat_entry.default.id}"]`,
			"nat_gateway_id":  `"${alicloud_ens_nat_gateway.default2Kn0nu.id}"`,
			"snat_entry_name": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsNatGatewaySnatEntrySourceConfig(rand, map[string]string{
			"ids":             `["${alicloud_ens_nat_gateway_snat_entry.default.id}"]`,
			"nat_gateway_id":  `"${alicloud_ens_nat_gateway.default2Kn0nu.id}"`,
			"snat_entry_name": `"${var.name}_fake"`,
		}),
	}
	SourceCidrConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsNatGatewaySnatEntrySourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ens_nat_gateway_snat_entry.default.id}"]`,
			"nat_gateway_id": `"${alicloud_ens_nat_gateway.default2Kn0nu.id}"`,
			"source_cidr":    `"10.0.0.0/8"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsNatGatewaySnatEntrySourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ens_nat_gateway_snat_entry.default.id}"]`,
			"nat_gateway_id": `"${alicloud_ens_nat_gateway.default2Kn0nu.id}"`,
			"source_cidr":    `"192.0.2.0/24"`,
		}),
	}
	SnatIpConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsNatGatewaySnatEntrySourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ens_nat_gateway_snat_entry.default.id}"]`,
			"nat_gateway_id": `"${alicloud_ens_nat_gateway.default2Kn0nu.id}"`,
			"snat_ip":        `"${alicloud_ens_eip.defaultiUbwh0.ip_address}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsNatGatewaySnatEntrySourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ens_nat_gateway_snat_entry.default.id}"]`,
			"nat_gateway_id": `"${alicloud_ens_nat_gateway.default2Kn0nu.id}"`,
			"snat_ip":        `"192.0.2.254"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsNatGatewaySnatEntrySourceConfig(rand, map[string]string{
			"ids":             `["${alicloud_ens_nat_gateway_snat_entry.default.id}"]`,
			"snat_entry_name": `"${var.name}"`,

			"nat_gateway_id": `"${alicloud_ens_nat_gateway.default2Kn0nu.id}"`,

			"source_cidr": `"10.0.0.0/8"`,

			"snat_ip": `"${alicloud_ens_eip.defaultiUbwh0.ip_address}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsNatGatewaySnatEntrySourceConfig(rand, map[string]string{
			"ids":             `["snat-fake0000"]`,
			"snat_entry_name": `"${var.name}_fake"`,

			"nat_gateway_id": `"${alicloud_ens_nat_gateway.default2Kn0nu.id}"`,

			"source_cidr": `"192.0.2.0/24"`,

			"snat_ip": `"192.0.2.254"`,
		}),
	}

	EnsNatGatewaySnatEntryCheckInfo.dataSourceTestCheck(t, rand, idsConf, SnatEntryNameConf, SourceCidrConf, SnatIpConf, allConf)
}

var existEnsNatGatewaySnatEntryMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"entries.#":                 "1",
		"entries.0.status":          CHECKSET,
		"entries.0.source_cidr":     CHECKSET,
		"entries.0.idle_timeout":    CHECKSET,
		"entries.0.snat_ip":         CHECKSET,
		"entries.0.eip_affinity":    CHECKSET,
		"entries.0.snat_entry_name": CHECKSET,
		"entries.0.isp_affinity":    CHECKSET,
		"entries.0.snat_entry_id":   CHECKSET,
		"entries.0.nat_gateway_id":  CHECKSET,
	}
}

var fakeEnsNatGatewaySnatEntryMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"entries.#": "0",
	}
}

var EnsNatGatewaySnatEntryCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ens_nat_gateway_snat_entries.default",
	existMapFunc: existEnsNatGatewaySnatEntryMapFunc,
	fakeMapFunc:  fakeEnsNatGatewaySnatEntryMapFunc,
}

func testAccCheckAlicloudEnsNatGatewaySnatEntrySourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccEnsNatGatewaySnatEntry%d"
}
variable "ens_region_id" {
    default = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_network" "defaultXqhlfk" {
        network_name = "${var.name}"
        cidr_block = "10.0.0.0/8"
        ens_region_id = "${var.ens_region_id}"
}

resource "alicloud_ens_vswitch" "defaultzkXvut" {
        cidr_block = "10.0.0.0/24"
        vswitch_name = "${var.name}"
        ens_region_id = "${alicloud_ens_network.defaultXqhlfk.ens_region_id}"
        network_id = "${alicloud_ens_network.defaultXqhlfk.id}"
}

resource "alicloud_ens_eip" "defaultiUbwh0" {
        depends_on = [alicloud_ens_nat_gateway.default2Kn0nu]
        bandwidth = "5"
        payment_type = "PayAsYouGo"
        ens_region_id = "${alicloud_ens_vswitch.defaultzkXvut.ens_region_id}"
        eip_name = "${var.name}"
        internet_charge_type = "95BandwidthByMonth"
}

resource "alicloud_ens_nat_gateway" "default2Kn0nu" {
        vswitch_id = "${alicloud_ens_vswitch.defaultzkXvut.id}"
        ens_region_id = "${alicloud_ens_vswitch.defaultzkXvut.ens_region_id}"
        network_id = "${alicloud_ens_vswitch.defaultzkXvut.network_id}"
        instance_type = "enat.default"
        nat_name = "${var.name}"
}

resource "alicloud_ens_eip_instance_attachment" "defaultlI0M0t" {
        instance_id = "${alicloud_ens_nat_gateway.default2Kn0nu.id}"
        allocation_id = "${alicloud_ens_eip.defaultiUbwh0.id}"
        instance_type = "Nat"
        standby = false
}



resource "alicloud_ens_nat_gateway_snat_entry" "default" {
        depends_on = [alicloud_ens_eip_instance_attachment.defaultlI0M0t]
        snat_entry_name = "${var.name}"
        source_cidr = "10.0.0.0/8"
        snat_ip = "${alicloud_ens_eip.defaultiUbwh0.ip_address}"
        nat_gateway_id = "${alicloud_ens_nat_gateway.default2Kn0nu.id}"
        idle_timeout = "50"
        isp_affinity = false
        eip_affinity = false
}

data "alicloud_ens_nat_gateway_snat_entries" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
