package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudNetworkInterfaceAttachmentBasic(t *testing.T) {
	var v ecs.NetworkInterfaceSet
	resourceId := "alicloud_network_interface_attachment.default"
	ra := resourceAttrInit(resourceId, testAccCheckNetworkInterfaceAttachmentCheckMap)
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

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNetworkInterfaceAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkInterfaceAttachmentConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlicloudNetworkInterfaceAttachmentMulti(t *testing.T) {
	var v ecs.NetworkInterfaceSet
	resourceId := "alicloud_network_interface_attachment.default.1"
	ra := resourceAttrInit(resourceId, testAccCheckNetworkInterfaceAttachmentCheckMap)
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

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNetworkInterfaceAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkInterfaceAttachmentConfigMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccCheckNetworkInterfaceAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_network_interface_Attachment" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No NetworkInterface ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsService := EcsService{client}
		_, err := ecsService.DescribeNetworkInterfaceAttachment(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
	}

	return nil
}

const testAccNetworkInterfaceAttachmentConfigBasic = `
variable "name" {
  default = "tf-testAccNetworkInterfaceAttachment"
}

resource "alicloud_vpc" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    "available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
}

data "alicloud_instance_types" "default" {
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    eni_amount = 2
}

data "alicloud_images" "default" {
	name_regex  = "^ubuntu_14.*_64"
  	most_recent = true
	owners = "system"
}

resource "alicloud_instance" "default" {
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    security_groups = ["${alicloud_security_group.default.id}"]

    instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
    system_disk_category = "cloud_efficiency"
    image_id             = "${data.alicloud_images.default.images.0.id}"
    instance_name        = "${var.name}"
    vswitch_id = "${alicloud_vswitch.default.id}"
    internet_max_bandwidth_out = 10
}

resource "alicloud_network_interface" "default" {
    name = "${var.name}"
    vswitch_id = "${alicloud_vswitch.default.id}"
    security_groups = [ "${alicloud_security_group.default.id}" ]
}

resource "alicloud_network_interface_attachment" "default" {
    instance_id = "${alicloud_instance.default.id}"
    network_interface_id = "${alicloud_network_interface.default.id}"
}
`

const testAccNetworkInterfaceAttachmentConfigMulti = `
variable "name" {
  default = "tf-testAccNetworkInterfaceAttachment"
}

variable "count" {
		default = "2"
	}

resource "alicloud_vpc" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    "available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
}

data "alicloud_instance_types" "default" {
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    eni_amount = 2
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_14.*_64"
  most_recent = true
  owners      = "system"
}

resource "alicloud_instance" "default" {
	count = "${var.count}"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    security_groups = ["${alicloud_security_group.default.id}"]

    instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
    system_disk_category = "cloud_efficiency"
    image_id             = "${data.alicloud_images.default.images.0.id}"
    instance_name        = "${var.name}"
    vswitch_id = "${alicloud_vswitch.default.id}"
    internet_max_bandwidth_out = 10
}

resource "alicloud_network_interface" "default" {
    count = "${var.count}"
    name = "${var.name}"
    vswitch_id = "${alicloud_vswitch.default.id}"
    security_groups = [ "${alicloud_security_group.default.id}" ]
}

resource "alicloud_network_interface_attachment" "default" {
	count = "${var.count}"
    instance_id = "${element(alicloud_instance.default.*.id, count.index)}"
    network_interface_id = "${element(alicloud_network_interface.default.*.id, count.index)}"
}
`

var testAccCheckNetworkInterfaceAttachmentCheckMap = map[string]string{
	"instance_id":          CHECKSET,
	"network_interface_id": CHECKSET,
}
