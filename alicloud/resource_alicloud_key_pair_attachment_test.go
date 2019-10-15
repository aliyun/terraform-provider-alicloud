package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func testAccCheckKeyPairAttachmentDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_key_pair_attachment" {
			continue
		}
		// Try to find the Disk
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsService := EcsService{client}

		instanceIds := rs.Primary.Attributes["instance_ids"]

		for _, inst := range instanceIds {
			response, err := ecsService.DescribeInstance(string(inst))
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

func TestAccAlicloudKeyPairAttachmentBasic(t *testing.T) {
	var v ecs.KeyPair
	resourceId := "alicloud_key_pair_attachment.default"
	ra := resourceAttrInit(resourceId, testAccCheckKeyPairAttachmentBasicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKeyPairAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyPairAttachmentConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

const testAccKeyPairAttachmentConfigBasic = `
data "alicloud_zones" "default" {
	available_disk_category = "cloud_ssd"
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
	name_regex = "^ubuntu_18.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "tf-testAccKeyPairAttachmentConfig"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}
resource "alicloud_security_group" "default" {
  name = "${var.name}"
  description = "New security group"
  vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_instance" "default" {
  instance_name = "${var.name}-${count.index+1}"
  image_id = "${data.alicloud_images.default.images.0.id}"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  count = 2
  security_groups = ["${alicloud_security_group.default.id}"]
  vswitch_id = "${alicloud_vswitch.default.id}"

  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = 5
  password = "Yourpassword1234"

  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_ssd"
}

resource "alicloud_key_pair" "default" {
  key_name = "${var.name}"
}

resource "alicloud_key_pair_attachment" "default" {
  key_name = "${alicloud_key_pair.default.id}"
  instance_ids = "${alicloud_instance.default.*.id}"
}
`

var testAccCheckKeyPairAttachmentBasicMap = map[string]string{
	"key_name":       CHECKSET,
	"instance_ids.#": "2",
}
