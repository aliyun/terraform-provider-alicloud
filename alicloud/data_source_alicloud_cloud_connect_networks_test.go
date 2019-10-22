package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudCloudConnectNetworkDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudConnectNetworkDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cloud_connect_network.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudConnectNetworkDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cloud_connect_network.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudConnectNetworkDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cloud_connect_network.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudConnectNetworkDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cloud_connect_network.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudConnectNetworkDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_cloud_connect_network.default.id}"]`,
			"name_regex": `"${alicloud_cloud_connect_network.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudConnectNetworkDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_cloud_connect_network.default.id}_fake"]`,
			"name_regex": `"${alicloud_cloud_connect_network.default.name}"`,
		}),
	}

	var existDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#":             "1",
			"ids.#":                   "1",
			"names.#":                 "1",
			"instances.0.id":          CHECKSET,
			"instances.0.name":        fmt.Sprintf("tf-testAccCcnInstanceDataSourceBisic-%d", rand),
			"instances.0.description": "tf-testAccCcnInstanceDescription",
			"instances.0.cidr_block":  "192.168.0.0/24,192.168.1.0/24",
			"instances.0.is_default":  "true",
		}
	}

	var fakeDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#": "0",
			"ids.#":       "0",
			"names.#":     "0",
		}
	}

	var ccnInstancesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cloud_connect_networks.default",
		existMapFunc: existDnsRecordsMapFunc,
		fakeMapFunc:  fakeDnsRecordsMapFunc,
	}

	ccnInstancesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlicloudCloudConnectNetworkDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
		variable "name" {
			default = "tf-testAccCcnInstanceDataSourceBisic-%d"
		}
		resource "alicloud_cloud_connect_network" "default" {
			name = "${var.name}"
			description = "tf-testAccCcnInstanceDescription"
			cidr_block = "192.168.0.0/24,192.168.1.0/24"
			is_default = true
		}

		data "alicloud_cloud_connect_networks" "default" {
		  %s
		}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
