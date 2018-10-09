package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudHaVipAttachment_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: "alicloud_havip_attachment.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckHaVipAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccHaVipAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHaVipAttachmentExists("alicloud_havip_attachment.foo"),
					resource.TestCheckResourceAttrSet(
						"alicloud_havip_attachment.foo", "havip_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_havip_attachment.foo", "instance_id"),
				),
			},
		},
	})
}

func testAccCheckHaVipAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No HaVip ID is set")
		}
		client := testAccProvider.Meta().(*AliyunClient)
		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		if len(parts) != 2 {
			return fmt.Errorf("invalid resource id")
		}
		err := client.DescribeHaVipAttachment(parts[0], parts[1])
		if err != nil {
			return fmt.Errorf("Describe HaVip attachment error %#v", err)
		}
		return nil
	}
}

func testAccCheckHaVipAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_havip_attachment" {
			continue
		}
		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		if len(parts) != 2 {
			return fmt.Errorf("invalid resource id")
		}
		err := client.DescribeHaVipAttachment(parts[0], parts[1])
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Describe HaVip attachment error %#v", err)
		}
	}
	return nil
}

const testAccHaVipAttachmentConfig = `

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}

data "alicloud_images" "default" {
	name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}

variable "name" {
	default = "tf-testAccHaVipAttachment"
}

resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_havip" "foo" {
	vswitch_id = "${alicloud_vswitch.foo.id}"
	description = "${var.name}"
}

resource "alicloud_havip_attachment" "foo" {
	havip_id = "${alicloud_havip.foo.id}"
	instance_id = "${alicloud_instance.foo.id}"
}

resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "foo" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	vswitch_id = "${alicloud_vswitch.foo.id}"
	image_id = "${data.alicloud_images.default.images.0.id}"
	# series III
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	system_disk_category = "cloud_efficiency"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5
	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
	instance_name = "${var.name}"
	user_data = "echo 'net.ipv4.ip_forward=1'>> /etc/sysctl.conf"
}

`
