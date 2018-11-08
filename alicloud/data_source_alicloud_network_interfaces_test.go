package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudNetworkInterfacesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckAlicloudNetworkInterfacesDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_network_interfaces.enis"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_network_interfaces.enis", "interfaces.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.name", "tf-testAcc-eni-xy"),
					resource.TestCheckResourceAttrSet("data.alicloud_network_interfaces.enis", "interfaces.0.status"),
					resource.TestCheckResourceAttrSet("data.alicloud_network_interfaces.enis", "interfaces.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_network_interfaces.enis", "interfaces.0.vswitch_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_network_interfaces.enis", "interfaces.0.zone_id"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.public_ip", ""),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.private_ip", "192.168.0.2"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.private_ips.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.security_groups.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.description", "Basic test"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.instance_id", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_network_interfaces.enis", "interfaces.0.creation_time"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.tags.%", "1"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.tags.TF-VER", "0.11.3"),
				),
			},
		},
	})
}

func TestAccAlicloudNetworkInterfacesDataSourceWithId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckAlicloudNetworkInterfacesDataSourceWithId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_network_interfaces.enis"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_network_interfaces.enis", "interfaces.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.name", "tf-testAcc-eni-xy"),
					resource.TestCheckResourceAttrSet("data.alicloud_network_interfaces.enis", "interfaces.0.status"),
					resource.TestCheckResourceAttrSet("data.alicloud_network_interfaces.enis", "interfaces.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_network_interfaces.enis", "interfaces.0.vswitch_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_network_interfaces.enis", "interfaces.0.zone_id"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.public_ip", ""),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.private_ip", "192.168.0.2"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.private_ips.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.security_groups.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.description", "Basic test"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.instance_id", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_network_interfaces.enis", "interfaces.0.creation_time"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.tags.%", "1"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.tags.TF-VER", "0.11.3"),
				),
			},
		},
	})
}

func TestAccAlicloudNetworkInterfacesDataSourceWithAllFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckAlicloudNetworkInterfacesDataSourceWithAllFields,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_network_interfaces.enis"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_network_interfaces.enis", "interfaces.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.name", "tf-testAcc-eni-xy"),
					resource.TestCheckResourceAttrSet("data.alicloud_network_interfaces.enis", "interfaces.0.status"),
					resource.TestCheckResourceAttrSet("data.alicloud_network_interfaces.enis", "interfaces.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_network_interfaces.enis", "interfaces.0.vswitch_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_network_interfaces.enis", "interfaces.0.zone_id"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.public_ip", ""),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.private_ip", "192.168.0.2"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.private_ips.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.security_groups.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.description", "Basic test"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.instance_id", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_network_interfaces.enis", "interfaces.0.creation_time"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.tags.%", "1"),
					resource.TestCheckResourceAttr("data.alicloud_network_interfaces.enis", "interfaces.0.tags.TF-VER", "0.11.3"),
				),
			},
		},
	})
}

const testAccCheckAlicloudNetworkInterfacesDataSourceBasic = `
resource "alicloud_vpc" "vpc" {
    name = "tf-testAcc-vpc-xy"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    "available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "vswitch" {
    name = "tf-testAcc-vswitch-xy"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group" "sg" {
    name = "tf-testAcc-sg-xy"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_network_interface" "eni" {
    name = "tf-testAcc-eni-xy"
    vswitch_id = "${alicloud_vswitch.vswitch.id}"
    security_groups = [ "${alicloud_security_group.sg.id}" ]
	description = "Basic test"
	private_ip = "192.168.0.2"
	tags = {
		TF-VER = "0.11.3"
	}
}

data "alicloud_network_interfaces" "enis"  {
	name_regex = "${alicloud_network_interface.eni.name}"
}
`

const testAccCheckAlicloudNetworkInterfacesDataSourceWithId = `
resource "alicloud_vpc" "vpc" {
    name = "tf-testAcc-vpc-xy"
    cidr_block = "192.168.0.0/24"
}

data "alicloud_zones" "default" {
    "available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "vswitch" {
    name = "tf-testAcc-vswitch-xy"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group" "sg" {
    name = "tf-testAcc-sg-xy"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_network_interface" "eni" {
    name = "tf-testAcc-eni-xy"
    vswitch_id = "${alicloud_vswitch.vswitch.id}"
    security_groups = [ "${alicloud_security_group.sg.id}" ]
	description = "Basic test"
	private_ip = "192.168.0.2"
	tags = {
		TF-VER = "0.11.3"
	}
}


data "alicloud_network_interfaces" "enis"  {
	ids = ["${alicloud_network_interface.eni.id}"]
}
`

const testAccCheckAlicloudNetworkInterfacesDataSourceWithAllFields = `
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

resource "alicloud_network_interface" "eni" {
    name = "tf-testAcc-eni-xy"
    vswitch_id = "${alicloud_vswitch.vswitch.id}"
    security_groups = [ "${alicloud_security_group.sg.id}" ]
	description = "Basic test"
	private_ip = "192.168.0.2"
	tags = {
		TF-VER = "0.11.3"
	}
}

data "alicloud_network_interfaces" "enis"  {
	ids = ["${alicloud_network_interface.eni.id}"]
	name_regex = "${alicloud_network_interface.eni.name}"
	vpc_id = "${alicloud_vpc.vpc.id}"
	vswitch_id = "${alicloud_vswitch.vswitch.id}"
	security_group_id = "${alicloud_security_group.sg.id}"
	private_ip = "192.168.0.2"
	tags = {
		TF-VER = "0.11.3"
	}
}
`
