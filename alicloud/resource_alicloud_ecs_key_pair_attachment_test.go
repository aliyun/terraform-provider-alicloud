package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func testAccCheckKeyPairAttachmentDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ecs_key_pair_attachment" {
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

func TestAccAlicloudECSKeyPairAttachmentBasic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_key_pair_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsKeyPairAttachmentMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlicloudEcsKeyPairAttachment%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKeyPairAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: AlicloudEcsKeyPairAttachmentBasicDependence(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

var AlicloudEcsKeyPairAttachmentMap = map[string]string{
	"key_pair_name":  CHECKSET,
	"instance_ids.#": "1",
}

func AlicloudEcsKeyPairAttachmentBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu"
  most_recent = true
  owners      = "system"
}
resource "alicloud_vpc" "default" {
  vpc_name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name              = "${var.name}"
}
resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
}
resource "alicloud_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alicloud_security_group.default.id}"
  	cidr_ip = "172.16.0.0/24"
}

resource "alicloud_instance" "default" {
	vswitch_id = "${alicloud_vswitch.default.id}"
	image_id = "${data.alicloud_images.default.images.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	system_disk_category = "cloud_efficiency"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5
	security_groups = ["${alicloud_security_group.default.id}"]
	instance_name = "${var.name}"
}

resource "alicloud_ecs_key_pair" "default" {
  key_pair_name = "${var.name}"
}
resource "alicloud_ecs_key_pair_attachment" "default" {
  key_pair_name = "${alicloud_ecs_key_pair.default.id}"
  instance_ids = [alicloud_instance.default.id]
}
`, name)
}
