package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudPvtzZoneAttachment_Basic(t *testing.T) {
	var zone pvtz.DescribeZoneInfoResponse
	var vpc vpc.DescribeVpcAttributeResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_pvtz_zone_attachment.zone-attachment",
		Providers:     testAccProviders,
		CheckDestroy:  testAccAlicloudPvtzZoneAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPvtzZoneAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneExists("alicloud_pvtz_zone.zone", &zone),
					testAccCheckVpcExists("alicloud_vpc.vpc", &vpc),
					testAccAlicloudPvtzZoneAttachmentExists("alicloud_pvtz_zone_attachment.zone-attachment", &zone, &vpc),
				),
			},
		},
	})

}

func TestAccAlicloudPvtzZoneAttachment_update(t *testing.T) {
	var zone pvtz.DescribeZoneInfoResponse
	var vpc vpc.DescribeVpcAttributeResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAlicloudPvtzZoneAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPvtzZoneAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneExists("alicloud_pvtz_zone.zone", &zone),
					testAccCheckVpcExists("alicloud_vpc.vpc", &vpc),
					testAccAlicloudPvtzZoneAttachmentExists("alicloud_pvtz_zone_attachment.zone-attachment", &zone, &vpc),
				),
			},
			resource.TestStep{
				Config: testAccPvtzZoneAttachmentConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneExists("alicloud_pvtz_zone.zone", &zone),
					testAccCheckVpcExists("alicloud_vpc.vpc1", &vpc),
					testAccAlicloudPvtzZoneAttachmentExists("alicloud_pvtz_zone_attachment.zone-attachment1", &zone, &vpc),
				),
			},
		},
	})

}

func TestAccAlicloudPvtzZoneAttachment_multi(t *testing.T) {
	var zone pvtz.DescribeZoneInfoResponse
	var vpc vpc.DescribeVpcAttributeResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_pvtz_zone_attachment.zone-attachment.0",
		Providers:     testAccProviders,
		CheckDestroy:  testAccAlicloudPvtzZoneAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPvtzZoneAttachmentConfigMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneExists("alicloud_pvtz_zone.zone", &zone),
					testAccCheckVpcExists("alicloud_vpc.vpcs.0", &vpc),
					testAccAlicloudPvtzZoneAttachmentExists("alicloud_pvtz_zone_attachment.zone-attachment.0", &zone, &vpc),
				),
			},
		},
	})
}

func testAccAlicloudPvtzZoneAttachmentExists(n string, zone *pvtz.DescribeZoneInfoResponse, vpc *vpc.DescribeVpcAttributeResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ZONE ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		instance, err := client.DescribePvtzZoneInfo(rs.Primary.ID)

		if err != nil {
			return err
		}

		if len(instance.BindVpcs.Vpc) == 0 {
			return fmt.Errorf("zone do not bind vpcs")
		}

		vpcId := vpc.VpcId
		vpcs := instance.BindVpcs.Vpc
		for i := 0; i < len(vpcs); i++ {
			if vpcId == vpcs[i].VpcId {
				*zone = instance
			}
		}

		return nil
	}
}

func testAccAlicloudPvtzZoneAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_pvtz_zone_attachment" {
			continue
		}

		instance, err := client.DescribePvtzZoneInfo(rs.Primary.ID)

		if err != nil && !NotFoundError(err) {
			return err
		}

		if len(instance.BindVpcs.Vpc) > 0 {
			return fmt.Errorf("zone %s still bind vpcs", rs.Primary.ID)
		}
	}

	return nil
}

const testAccPvtzZoneAttachmentConfig = `
resource "alicloud_pvtz_zone" "zone" {
	name = "foo.test.com"
}

resource "alicloud_vpc" "vpc" {
	name = "tf_test_foo"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_pvtz_zone_attachment" "zone-attachment" {
	zone_id = "${alicloud_pvtz_zone.zone.id}"
	vpc_ids = ["${alicloud_vpc.vpc.id}"]
}
`

const testAccPvtzZoneAttachmentConfigUpdate = `
resource "alicloud_pvtz_zone" "zone" {
	name = "foo.test.com"
}

resource "alicloud_vpc" "vpc1" {
	name = "vpc1"
	cidr_block = "192.168.0.0/16"
}

resource "alicloud_pvtz_zone_attachment" "zone-attachment1" {
	zone_id = "${alicloud_pvtz_zone.zone.id}"
	vpc_ids = ["${alicloud_vpc.vpc1.id}"]
}
`

const testAccPvtzZoneAttachmentConfigMulti = `
variable "count" {
  	default = "2"
}

resource "alicloud_vpc" "vpcs" {
	count = "${var.count}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_pvtz_zone" "zone" {
	name = "foo.test.com"
}

resource "alicloud_pvtz_zone_attachment" "zone-attachment" {
	count = "${var.count}"
	zone_id = "${alicloud_pvtz_zone.zone.id}"
	vpc_ids = ["${alicloud_vpc.vpcs.*.id}"]
}
`
