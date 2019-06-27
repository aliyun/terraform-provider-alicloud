package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudPvtzZoneAttachment_update(t *testing.T) {
	var zone pvtz.DescribeZoneInfoResponse
	var vpc vpc.DescribeVpcAttributeResponse
	rand := acctest.RandIntRange(10000, 999999)
	resourceId := "alicloud_pvtz_zone_attachment.zone-attachment"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccAlicloudPvtzZoneAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPvtzZoneAttachmentConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneExists("alicloud_pvtz_zone.zone", &zone),
					testAccCheckVpcExists("alicloud_vpc.vpc", &vpc),
					testAccAlicloudPvtzZoneAttachmentExists("alicloud_pvtz_zone_attachment.zone-attachment", &zone, &vpc),
					resource.TestCheckResourceAttrSet("alicloud_pvtz_zone_attachment.zone-attachment", "zone_id"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_attachment.zone-attachment", "vpc_ids.#", "1"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPvtzZoneAttachmentConfigUpdate(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneExists("alicloud_pvtz_zone.zone", &zone),
					testAccCheckVpcExists("alicloud_vpc.vpc1", &vpc),
					testAccAlicloudPvtzZoneAttachmentExists("alicloud_pvtz_zone_attachment.zone-attachment", &zone, &vpc),
					resource.TestCheckResourceAttrSet("alicloud_pvtz_zone_attachment.zone-attachment", "zone_id"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_attachment.zone-attachment", "vpc_ids.#", "1"),
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
		IDRefreshName: "alicloud_pvtz_zone_attachment.zone-attachment",
		Providers:     testAccProviders,
		CheckDestroy:  testAccAlicloudPvtzZoneAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPvtzZoneAttachmentConfigMulti(acctest.RandIntRange(10000, 999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneExists("alicloud_pvtz_zone.zone", &zone),
					testAccCheckVpcExists("alicloud_vpc.vpcs.0", &vpc),
					testAccCheckVpcExists("alicloud_vpc.vpcs.1", &vpc),
					testAccAlicloudPvtzZoneAttachmentExists("alicloud_pvtz_zone_attachment.zone-attachment", &zone, &vpc),
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

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		pvtzService := PvtzService{client}
		instance, err := pvtzService.DescribePvtzZoneAttachment(rs.Primary.ID)

		if err != nil {
			return WrapError(err)
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
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_pvtz_zone_attachment" {
			continue
		}

		if _, err := pvtzService.DescribePvtzZoneAttachment(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		return WrapError(fmt.Errorf("zone %s still bind vpcs", rs.Primary.ID))
	}

	return nil
}

func testAccPvtzZoneAttachmentConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_pvtz_zone" "zone" {
		name = "tf-testacc%d.test.com"
	}

	resource "alicloud_vpc" "vpc" {
		name = "tf-testaccPvtzZoneAttachmentConfig"
		cidr_block = "172.16.0.0/12"
	}

	resource "alicloud_pvtz_zone_attachment" "zone-attachment" {
		zone_id = "${alicloud_pvtz_zone.zone.id}"
		vpc_ids = ["${alicloud_vpc.vpc.id}"]
	}
	`, rand)
}

func testAccPvtzZoneAttachmentConfigUpdate(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_pvtz_zone" "zone" {
		name = "tf-testacc%d.test.com"
	}

	resource "alicloud_vpc" "vpc1" {
		name = "tf-testaccPvtzZoneAttachmentConfigUpdate"
		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_pvtz_zone_attachment" "zone-attachment" {
		zone_id = "${alicloud_pvtz_zone.zone.id}"
		vpc_ids = ["${alicloud_vpc.vpc1.id}"]
	}
	`, rand)
}

func testAccPvtzZoneAttachmentConfigMulti(rand int) string {
	return fmt.Sprintf(`
	variable "number" {
		  default = "2"
	}

	resource "alicloud_vpc" "vpcs" {
		count = "${var.number}"
		cidr_block = "172.16.0.0/12"
		name = "tf-testaccPvtzZoneAttachmentConfigMulti"
	}

	resource "alicloud_pvtz_zone" "zone" {
		name = "tf-testacc%d.test.com"
	}

	resource "alicloud_pvtz_zone_attachment" "zone-attachment" {
		zone_id = "${alicloud_pvtz_zone.zone.id}"
		vpc_ids = "${alicloud_vpc.vpcs.*.id}"
	}
	`, rand)
}
