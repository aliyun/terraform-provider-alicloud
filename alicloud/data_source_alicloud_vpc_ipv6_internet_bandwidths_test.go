package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVPCIpv6InternetBandwidthsDataSource(t *testing.T) {
	resourceId := "data.alicloud_vpc_ipv6_internet_bandwidths.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-vpcipv6internetbandwidth-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceVpcIpv6InternetBandwidthsDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_vpc_ipv6_internet_bandwidth.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_vpc_ipv6_internet_bandwidth.default.id}-fake"},
		}),
	}
	ipv6InternetBandwidthIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ipv6_internet_bandwidth_id": "${alicloud_vpc_ipv6_internet_bandwidth.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ipv6_internet_bandwidth_id": "${alicloud_vpc_ipv6_internet_bandwidth.default.id}_fake",
		}),
	}
	ipv6AddressIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":             []string{"${alicloud_vpc_ipv6_internet_bandwidth.default.id}"},
			"ipv6_address_id": "${alicloud_vpc_ipv6_internet_bandwidth.default.ipv6_address_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":             []string{"${alicloud_vpc_ipv6_internet_bandwidth.default.id}"},
			"ipv6_address_id": "${alicloud_vpc_ipv6_internet_bandwidth.default.ipv6_address_id}_fake",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_vpc_ipv6_internet_bandwidth.default.id}"},
			"status": "Normal",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_vpc_ipv6_internet_bandwidth.default.id}"},
			"status": "FinacialLocked",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                        []string{"${alicloud_vpc_ipv6_internet_bandwidth.default.id}"},
			"ipv6_internet_bandwidth_id": "${alicloud_vpc_ipv6_internet_bandwidth.default.id}",
			"ipv6_address_id":            "${alicloud_vpc_ipv6_internet_bandwidth.default.ipv6_address_id}",
			"status":                     "Normal",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                        []string{"${alicloud_vpc_ipv6_internet_bandwidth.default.id}"},
			"ipv6_internet_bandwidth_id": "${alicloud_vpc_ipv6_internet_bandwidth.default.id}_fake",
			"ipv6_address_id":            "${alicloud_vpc_ipv6_internet_bandwidth.default.ipv6_address_id}_fake",
			"status":                     "FinacialLocked",
		}),
	}
	var existVpcIpv6InternetBandwidthMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"ids.0":                        CHECKSET,
			"bandwidths.#":                 "1",
			"bandwidths.0.id":              CHECKSET,
			"bandwidths.0.status":          "Normal",
			"bandwidths.0.ipv6_address_id": CHECKSET,
			"bandwidths.0.ipv6_gateway_id": CHECKSET,
			"bandwidths.0.ipv6_internet_bandwidth_id": CHECKSET,
			"bandwidths.0.internet_charge_type":       "PayByBandwidth",
			"bandwidths.0.bandwidth":                  "20",
			"bandwidths.0.payment_type":               "PayAsYouGo",
		}
	}

	var fakeVpcIpv6InternetBandwidthMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        "0",
			"bandwidths.#": "0",
		}
	}

	var VpcIpv6InternetBandwidthCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existVpcIpv6InternetBandwidthMapFunc,
		fakeMapFunc:  fakeVpcIpv6InternetBandwidthMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithEnvVariable(t, "ECS_WITH_IPV6_ADDRESS")
	}

	VpcIpv6InternetBandwidthCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, ipv6InternetBandwidthIdConf, ipv6AddressIdConf, statusConf, allConf)
}

func dataSourceVpcIpv6InternetBandwidthsDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_instances" "default" {
  name_regex = "no-deleteing-ipv6-address"
  status     = "Running"
}

data "alicloud_vpc_ipv6_addresses" "default" {
  associated_instance_id = data.alicloud_instances.default.instances.0.id
  status                 = "Available"
}

resource "alicloud_vpc_ipv6_internet_bandwidth" "default" {
  ipv6_address_id      = data.alicloud_vpc_ipv6_addresses.default.addresses.0.id
  ipv6_gateway_id      = data.alicloud_vpc_ipv6_addresses.default.addresses.0.ipv6_gateway_id
  internet_charge_type = "PayByBandwidth"
  bandwidth            = "20"
}`, name)
}
