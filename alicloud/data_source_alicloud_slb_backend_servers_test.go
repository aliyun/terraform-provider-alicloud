package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudSlbBackendServersDataSource_basic(t *testing.T) {

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbBackendServersDataSourceConfig(map[string]string{
			"load_balancer_id": `"${alicloud_slb_backend_server.default.load_balancer_id}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbBackendServersDataSourceConfig(map[string]string{
			"load_balancer_id": `"${alicloud_slb_backend_server.default.load_balancer_id}"`,
			"ids":              `["${alicloud_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSlbBackendServersDataSourceConfig(map[string]string{
			"load_balancer_id": `"${alicloud_slb_backend_server.default.load_balancer_id}"`,
			"ids":              `["${alicloud_instance.default.id}_fake"]`,
		}),
	}

	var existDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"backend_servers.#":             "1",
			"backend_servers.0.id":          CHECKSET,
			"backend_servers.0.weight":      "100",
			"backend_servers.0.server_type": CHECKSET,
		}
	}

	var fakeDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"backend_servers.#": "0",
			"ids.#":             "0",
		}
	}

	var slbServerGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_slb_backend_servers.default",
		existMapFunc: existDnsRecordsMapFunc,
		fakeMapFunc:  fakeDnsRecordsMapFunc,
	}

	slbServerGroupsCheckInfo.dataSourceTestCheck(t, acctest.RandInt(), idsConf, allConf)
}

func testAccCheckAlicloudSlbBackendServersDataSourceConfig(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccslbbackendserversdatasourcebasic"
}

data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}
data "alicloud_images" "default" {
  name_regex = "^ubuntu_18.*_64"
  most_recent = true
  owners = "system"
}
data "alicloud_instance_types" "default" {
	cpu_core_count = 2
	memory_size = 4
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
}

resource "alicloud_security_group" "default" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_instance" "default" {
  image_id = "${data.alicloud_images.default.images.0.id}"

  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  internet_charge_type = "PayByTraffic"
  system_disk_category = "cloud_efficiency"

  security_groups = ["${alicloud_security_group.default.id}"]
  instance_name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb_backend_server" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"

  backend_servers {
    server_id = "${alicloud_instance.default.id}"
    weight     = 100
  }
}

data "alicloud_slb_backend_servers" "default" {
  %s
}
`, strings.Join(pairs, "\n  "))
	return config
}
