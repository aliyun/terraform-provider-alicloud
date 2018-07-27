package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudOtsInstance_Basic(t *testing.T) {
	var instance ots.InstanceInfo
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ots_instance.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOtsInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist(
						"alicloud_ots_instance.basic", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_ots_instance.basic",
						"name", "tftestInstance"),
					resource.TestCheckResourceAttr(
						"alicloud_ots_instance.basic",
						"accessed_by", "Any"),
				),
			},
		},
	})

}

func TestAccAlicloudOtsInstance_Tags(t *testing.T) {
	var instance ots.InstanceInfo
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ots_instance.tags",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOtsInstanceTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist(
						"alicloud_ots_instance.tags", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_ots_instance.tags",
						"name", "tftestInstTag"),
					resource.TestCheckResourceAttr(
						"alicloud_ots_instance.tags",
						"instance_type", "HighPerformance"),
					resource.TestCheckResourceAttr(
						"alicloud_ots_instance.tags",
						"tags.Created", "TF"),
					resource.TestCheckResourceAttr(
						"alicloud_ots_instance.tags",
						"tags.For", "acceptance test"),
				),
			},
		},
	})

}

func testAccCheckOtsInstanceExist(n string, instance *ots.InstanceInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found OTS table: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no OTS table ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)

		response, err := client.DescribeOtsInstance(rs.Primary.ID)

		if err != nil {
			return err
		}
		instance = &response
		return nil
	}
}

func testAccCheckOtsInstanceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ots_instance" {
			continue
		}

		client := testAccProvider.Meta().(*AliyunClient)

		if _, err := client.DescribeOtsInstance(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		return fmt.Errorf("Ots instance %s still exists.", rs.Primary.ID)
	}

	return nil
}

const testAccOtsInstance = `
variable "name" {
  default = "tftestInstance"
}
resource "alicloud_ots_instance" "basic" {
  name = "${var.name}"
  description = "${var.name}"
}
`

const testAccOtsInstanceTags = `
variable "name" {
  default = "tftestInstTag"
}
resource "alicloud_ots_instance" "tags" {
  name = "${var.name}"
  description = "${var.name}"
  accessed_by = "Vpc"
  tags {
	Created = "TF"
	For = "acceptance test"
  }
}
`
