package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSlbAttachmentsDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbAttachmentsDataSourceConfig(rand, map[string]string{
			"load_balancer_id": `"${alicloud_slb_attachment.default.load_balancer_id}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbAttachmentsDataSourceConfig(rand, map[string]string{
			"load_balancer_id": `"${alicloud_slb_attachment.default.load_balancer_id}"`,
			"instance_ids":     `["${alicloud_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSlbAttachmentsDataSourceConfig(rand, map[string]string{
			"load_balancer_id": `"${alicloud_slb_attachment.default.load_balancer_id}"`,
			"instance_ids":     `["${alicloud_instance.default.id}_fake"]`,
		}),
	}

	var existDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"slb_attachments.#":             "1",
			"slb_attachments.0.instance_id": CHECKSET,
			"slb_attachments.0.weight":      "42",
		}
	}

	var fakeDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"slb_attachments.#": "0",
		}
	}

	var slbAttachmentCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_slb_attachments.default",
		existMapFunc: existDnsRecordsMapFunc,
		fakeMapFunc:  fakeDnsRecordsMapFunc,
	}

	slbAttachmentCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

func testAccCheckAlicloudSlbAttachmentsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCheckAlicloudSlbAttachmentsDataSourceBasic-%d"
}

data "alicloud_images" "default" {
  name_regex = "^ubuntu_18.*64"
  most_recent = true
  owners = "system"
}
data "alicloud_instance_types" "default" {
	cpu_core_count = 1
	memory_size = 2
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
  name = "${var.name}"
}

resource "alicloud_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_security_group" "default" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_instance" "default" {
  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
  image_id = "${data.alicloud_images.default.images.0.id}"

  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  internet_charge_type = "PayByTraffic"
  system_disk_category = "cloud_efficiency"

  security_groups = ["${alicloud_security_group.default.id}"]
  instance_name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb_attachment" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  instance_ids = ["${alicloud_instance.default.id}"]
  weight = 42
}

data "alicloud_slb_attachments" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
