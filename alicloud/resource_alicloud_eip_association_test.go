package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func testAccCheckEIPAssociationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_eip_association" {
			continue
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No EIP Association ID is set"))
		}

		// Try to find the EIP
		_, err := vpcService.DescribeEipAssociation(rs.Primary.ID)

		// Verify the error is what we want
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}

func TestAccAlicloudEIPAssociationBasic(t *testing.T) {
	var v vpc.EipAddress
	resourceId := "alicloud_eip_association.default"
	ra := resourceAttrInit(resourceId, testAccCheckEipAssociationBasicMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEIPAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPAssociationConfigBaisc(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAlicloudEIPAssociationMulti(t *testing.T) {
	var v vpc.EipAddress
	resourceId := "alicloud_eip_association.default.1"
	ra := resourceAttrInit(resourceId, testAccCheckEipAssociationBasicMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEIPAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPAssociationConfigMulti(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAlicloudEIPAssociationEni(t *testing.T) {
	var v vpc.EipAddress
	resourceId := "alicloud_eip_association.default"
	ra := resourceAttrInit(resourceId, testAccCheckEipAssociationBasicMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEIPAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPAssociationConfigEni(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":      "NetworkInterface",
						"private_ip_address": CHECKSET,
					}),
				),
			},
		},
	})
}

func testAccEIPAssociationConfigBaisc(rand int) string {
	return fmt.Sprintf(`
data "alicloud_zones" "default" {
}

data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

data "alicloud_images" "default" {
	name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
	owners = "system"
}

variable "name" {
	default = "tf-testAccEipAssociation%d"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_security_group" "default" {
  name = "${var.name}"
  description = "New security group"
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_instance" "default" {
  vswitch_id = local.vswitch_id
  image_id = "${data.alicloud_images.default.images.1.id}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  system_disk_category = "cloud_ssd"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"

  security_groups = ["${alicloud_security_group.default.id}"]
  instance_name = "${var.name}"
  tags = {
    Name = "TerraformTest-instance"
  }
  lifecycle {
    ignore_changes = [
      instance_type,
    ]
  }
}

resource "alicloud_eip_address" "default" {
	address_name = "${var.name}"
}

resource "alicloud_eip_association" "default" {
  allocation_id = "${alicloud_eip_address.default.id}"
  instance_id = "${alicloud_instance.default.id}"
  force = true
}
`, rand)
}

func testAccEIPAssociationConfigMulti(rand int) string {
	return fmt.Sprintf(`
data "alicloud_zones" "default" {
}

data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

data "alicloud_images" "default" {
    name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
	owners = "system"
}

variable "name" {
	default = "tf-testAccEipAssociation%d"
}

variable "number" {
		default = "2"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_security_group" "default" {
  name = "${var.name}"
  description = "New security group"
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_instance" "default" {
  count = "${var.number}"
  vswitch_id = local.vswitch_id
  image_id = "${data.alicloud_images.default.images.0.id}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  system_disk_category = "cloud_ssd"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"

  security_groups = ["${alicloud_security_group.default.id}"]
  instance_name = "${var.name}"
  tags = {
    Name = "TerraformTest-instance"
  }
  lifecycle {
    ignore_changes = [
      instance_type,
    ]
  }
}

resource "alicloud_eip_address" "default" {
	count = "${var.number}"
	address_name = "${var.name}"
}

resource "alicloud_eip_association" "default" {
  count = "${var.number}"
  allocation_id = "${element(alicloud_eip_address.default.*.id,count.index)}"
  instance_id = "${element(alicloud_instance.default.*.id,count.index)}"
}
`, rand)
}

func testAccEIPAssociationConfigEni(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccEipAssociation%d"
}

data "alicloud_zones" "default" {
    available_resource_creation= "VSwitch"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_network_interface" "default" {
	name = "${var.name}"
    vswitch_id = local.vswitch_id
	security_groups = [ "${alicloud_security_group.default.id}" ]
	private_ip = cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 1)
}

resource "alicloud_eip_address" "default" {
	address_name = "${var.name}"
}

resource "alicloud_eip_association" "default" {
  allocation_id = "${alicloud_eip_address.default.id}"
  instance_id = "${alicloud_network_interface.default.id}"
  instance_type = "NetworkInterface"
  private_ip_address = cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 1)
}
`, rand)
}

var testAccCheckEipAssociationBasicMap = map[string]string{
	"allocation_id": CHECKSET,
	"instance_id":   CHECKSET,
}
