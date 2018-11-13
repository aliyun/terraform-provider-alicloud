package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudNetworkInterfaceAttachment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_network_interface_attachment.att",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNetworkInterfaceAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkInterfaceAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInterfaceAttachmentExists("alicloud_network_interface.eni", "alicloud_instance.instance"),
					resource.TestCheckResourceAttrSet("alicloud_network_interface_attachment.att", "network_interface_id"),
					resource.TestCheckResourceAttrSet("alicloud_network_interface_attachment.att", "instance_id"),
				),
			},
		},
	})
}

func TestAccAlicloudNetworkInterfaceAttachmentWithMultiEni(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_network_interface_attachment.att1",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNetworkInterfaceAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkInterfaceAttachmentConfigWithMultiEni,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkInterfaceAttachmentExists("alicloud_network_interface.eni0", "alicloud_instance.instance"),
					testAccCheckNetworkInterfaceAttachmentExists("alicloud_network_interface.eni1", "alicloud_instance.instance"),
					resource.TestCheckResourceAttrSet("alicloud_network_interface_attachment.att0", "network_interface_id"),
					resource.TestCheckResourceAttrSet("alicloud_network_interface_attachment.att0", "instance_id"),
					resource.TestCheckResourceAttrSet("alicloud_network_interface_attachment.att1", "network_interface_id"),
					resource.TestCheckResourceAttrSet("alicloud_network_interface_attachment.att1", "instance_id"),
				),
			},
		},
	})
}

func testAccCheckNetworkInterfaceAttachmentExists(eniName string, instanceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[eniName]
		if !ok {
			return fmt.Errorf("Not found: %s", eniName)
		}
		eniId := rs.Primary.ID

		rs, ok = s.RootModule().Resources[instanceName]
		if !ok {
			return fmt.Errorf("Not found: %s", instanceName)
		}
		instanceId := rs.Primary.ID

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsService := EcsService{client}
		_, err := ecsService.DescribeNetworkInterfaceById(instanceId, eniId)
		if err != nil {
			if NotFoundError(err) {
				return fmt.Errorf("Attach NetworkInterface (%s) to Instance (%s) failed", eniId, instanceId)
			}
			return fmt.Errorf("Describe NetworkInterface (%s) failed when checking attachment, %s", rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccCheckNetworkInterfaceAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_network_interface" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No NetworkInterface ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsService := EcsService{client}
		_, err := ecsService.DescribeNetworkInterfaceById("", rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
	}

	return nil
}

const testAccNetworkInterfaceAttachmentConfig = `
resource "alicloud_vpc" "vpc" {
    name = "tf-testAcc-vpc"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    "available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "vswitch" {
    name = "tf-testAcc-vswitch"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group" "sg" {
    name = "tf-testAcc-sg"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

data "alicloud_instance_types" "default" {
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    cpu_core_count = 2
    memory_size = 4
}

resource "alicloud_instance" "instance" {
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    security_groups = ["${alicloud_security_group.sg.id}"]

    instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
    system_disk_category = "cloud_efficiency"
    image_id             = "centos_7_04_64_20G_alibase_201701015.vhd"
    instance_name        = "tf-testAcc-instance"
    vswitch_id = "${alicloud_vswitch.vswitch.id}"
    internet_max_bandwidth_out = 10
}

resource "alicloud_network_interface" "eni" {
    name = "tf-testAcc-eni"
    vswitch_id = "${alicloud_vswitch.vswitch.id}"
    security_groups = [ "${alicloud_security_group.sg.id}" ]
}

resource "alicloud_network_interface_attachment" "att" {
    instance_id = "${alicloud_instance.instance.id}"
    network_interface_id = "${alicloud_network_interface.eni.id}"
}
`

const testAccNetworkInterfaceAttachmentConfigWithMultiEni = `
resource "alicloud_vpc" "vpc" {
    name = "tf-testAcc-vpc"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    "available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "vswitch" {
    name = "tf-testAcc-vswitch"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group" "sg" {
    name = "tf-testAcc-sg"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

data "alicloud_instance_types" "default" {
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    cpu_core_count = 2
    memory_size = 4
}

resource "alicloud_instance" "instance" {
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    security_groups = ["${alicloud_security_group.sg.id}"]

    instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
    system_disk_category = "cloud_efficiency"
    image_id             = "centos_7_04_64_20G_alibase_201701015.vhd"
    instance_name        = "tf-testAcc-instance"
    vswitch_id = "${alicloud_vswitch.vswitch.id}"
    internet_max_bandwidth_out = 10
}

resource "alicloud_network_interface" "eni0" {
    name = "tf-testAcc-eni0"
    vswitch_id = "${alicloud_vswitch.vswitch.id}"
    security_groups = [ "${alicloud_security_group.sg.id}" ]
}

resource "alicloud_network_interface" "eni1" {
    name = "tf-testAcc-eni1"
    vswitch_id = "${alicloud_vswitch.vswitch.id}"
    security_groups = [ "${alicloud_security_group.sg.id}" ]
}

resource "alicloud_network_interface_attachment" "att0" {
    instance_id = "${alicloud_instance.instance.id}"
    network_interface_id = "${alicloud_network_interface.eni0.id}"
}

resource "alicloud_network_interface_attachment" "att1" {
    instance_id = "${alicloud_instance.instance.id}"
    network_interface_id = "${alicloud_network_interface.eni1.id}"
}
`
