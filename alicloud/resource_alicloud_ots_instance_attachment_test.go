package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudOtsInstanceCapacityAttachment(t *testing.T) {
	var instance ots.InstanceInfo
	var vpcInfo ots.VpcInfo
	var vsw vpc.DescribeVSwitchAttributesResponse
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},

		// module name
		IDRefreshName: "alicloud_ots_instance_attachment.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsInstanceAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOtsInstanceAttachment(string(OtsCapacity), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					testAccCheckVSwitchExists("alicloud_vswitch.foo", &vsw),
					testAccCheckOtsInstanceAttachmentExist("alicloud_ots_instance_attachment.foo", &vpcInfo),
					resource.TestCheckResourceAttr("alicloud_ots_instance_attachment.foo", "instance_name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance_attachment.foo", "vpc_name", "test"),
					resource.TestCheckResourceAttrSet("alicloud_ots_instance_attachment.foo", "vswitch_id"),
					resource.TestCheckResourceAttrSet("alicloud_ots_instance_attachment.foo", "vpc_id"),
				),
			},
		},
	})

}

func TestAccAlicloudOtsInstanceHighPerformanceAttachment(t *testing.T) {
	var instance ots.InstanceInfo
	var vpcInfo ots.VpcInfo
	var vsw vpc.DescribeVSwitchAttributesResponse
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.OtsHighPerformanceNoSupportedRegions)
		},

		// module name
		IDRefreshName: "alicloud_ots_instance_attachment.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsInstanceAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOtsInstanceAttachment(string(OtsHighPerformance), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					testAccCheckVSwitchExists("alicloud_vswitch.foo", &vsw),
					testAccCheckOtsInstanceAttachmentExist("alicloud_ots_instance_attachment.foo", &vpcInfo),
					resource.TestCheckResourceAttr("alicloud_ots_instance_attachment.foo", "instance_name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance_attachment.foo", "vpc_name", "test"),
					resource.TestCheckResourceAttrSet("alicloud_ots_instance_attachment.foo", "vswitch_id"),
					resource.TestCheckResourceAttrSet("alicloud_ots_instance_attachment.foo", "vpc_id"),
				),
			},
		},
	})

}

func testAccCheckOtsInstanceAttachmentExist(n string, instance *ots.VpcInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found OTS table: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no OTS table ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		otsService := OtsService{client}

		response, err := otsService.DescribeOtsInstanceVpc(rs.Primary.ID)

		if err != nil {
			return err
		}
		instance = &response
		return nil
	}
}

func testAccCheckOtsInstanceAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ots_instance_attachment" {
			continue
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		otsService := OtsService{client}

		if _, err := otsService.DescribeOtsInstanceVpc(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		return fmt.Errorf("Ots instance attachment %s still exists.", rs.Primary.ID)
	}

	return nil
}

func testAccOtsInstanceAttachment(instanceType string, rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%d"
	}

	resource "alicloud_ots_instance" "foo" {
	  name = "${var.name}"
	  description = "${var.name}"
	  accessed_by = "Vpc"
	  instance_type = "%s"
	}

	data "alicloud_zones" "foo" {
	  available_resource_creation = "VSwitch"
	}
	resource "alicloud_vpc" "foo" {
	  cidr_block = "172.16.0.0/16"
	  name = "${var.name}"
	}

	resource "alicloud_vswitch" "foo" {
	  vpc_id = "${alicloud_vpc.foo.id}"
	  name = "${var.name}"
	  cidr_block = "172.16.1.0/24"
	  availability_zone = "${data.alicloud_zones.foo.zones.0.id}"
	}
	resource "alicloud_ots_instance_attachment" "foo" {
	  instance_name = "${alicloud_ots_instance.foo.name}"
	  vpc_name = "test"
	  vswitch_id = "${alicloud_vswitch.foo.id}"
	}
	`, rand, instanceType)
}
