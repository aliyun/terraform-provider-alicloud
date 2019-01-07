package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudKeyPairAttachment_basic(t *testing.T) {
	var keypair ecs.KeyPair
	var instance ecs.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_key_pair_attachment.attach",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKeyPairAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyPairAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyPairExists(
						"alicloud_key_pair.key", &keypair),
					testAccCheckInstanceExists(
						"alicloud_instance.instance.0", &instance),
					testAccCheckKeyPairAttachmentExists(
						"alicloud_key_pair_attachment.attach", &instance, &keypair),
				),
			},
		},
	})

}

func testAccCheckKeyPairAttachmentExists(n string, instance *ecs.Instance, keypair *ecs.KeyPair) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Key Pair Attachment ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsService := EcsService{client}

		response, err := ecsService.DescribeInstanceById(instance.InstanceId)
		if err != nil {
			return fmt.Errorf("Error QueryInstancesById: %#v", err)
		}

		if response.KeyPairName == keypair.KeyPairName {
			keypair.KeyPairName = response.KeyPairName
			*instance = response
			return nil

		}
		return fmt.Errorf("Error KeyPairAttachment is not exist.")
	}
}

func testAccCheckKeyPairAttachmentDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_key_pair_attachment" {
			continue
		}
		// Try to find the Disk
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsService := EcsService{client}

		instance_ids := rs.Primary.Attributes["instance_ids"]

		for _, inst := range instance_ids {
			response, err := ecsService.DescribeInstanceById(string(inst))
			if err != nil {
				return err
			}

			if response.KeyPairName != "" {
				return fmt.Errorf("Error Key Pair Attachment still exist")
			}

		}
	}

	return nil
}

const testAccKeyPairAttachmentConfig = `
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_ssd"
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
	default = "tf-testAccKeyPairAttachmentConfig"
}

resource "alicloud_vpc" "main" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "main" {
  vpc_id = "${alicloud_vpc.main.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}
resource "alicloud_security_group" "group" {
  name = "${var.name}"
  description = "New security group"
  vpc_id = "${alicloud_vpc.main.id}"
}

resource "alicloud_instance" "instance" {
  instance_name = "${var.name}-${count.index+1}"
  image_id = "${data.alicloud_images.default.images.0.id}"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  count = 2
  security_groups = ["${alicloud_security_group.group.id}"]
  vswitch_id = "${alicloud_vswitch.main.id}"

  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = 5
  password = "Test12345"

  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_ssd"
}

resource "alicloud_key_pair" "key" {
  key_name = "${var.name}"
}

resource "alicloud_key_pair_attachment" "attach" {
  key_name = "${alicloud_key_pair.key.id}"
  instance_ids = ["${alicloud_instance.instance.*.id}"]
}
`
