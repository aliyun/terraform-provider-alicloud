package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudEIP_basic(t *testing.T) {
	var eip vpc.EipAddress

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_eip.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEIPConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(
						"alicloud_eip.foo", &eip),
					testAccCheckEIPAttributes(&eip),
				),
			},
			resource.TestStep{
				Config: testAccEIPConfigTwo,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(
						"alicloud_eip.foo", &eip),
					testAccCheckEIPAttributes(&eip),
					resource.TestCheckResourceAttr(
						"alicloud_eip.foo",
						"bandwidth",
						"10"),
				),
			},
		},
	})

}

func testAccCheckEIPExists(n string, eip *vpc.EipAddress) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EIP ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		d, err := client.DescribeEipAddress(rs.Primary.ID)

		log.Printf("[WARN] eip id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		if d.IpAddress == "" {
			return fmt.Errorf("EIP not found")
		}

		*eip = d
		return nil
	}
}

func testAccCheckEIPAttributes(eip *vpc.EipAddress) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if eip.IpAddress == "" {
			return fmt.Errorf("Empty Ip address")
		}

		return nil
	}
}

func testAccCheckEIPDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_eip" {
			continue
		}

		d, err := client.DescribeEipAddress(rs.Primary.ID)

		// Verify the error is what we want
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if d.AllocationId != "" {
			return fmt.Errorf("Error EIP still exist")
		}
	}

	return nil
}

const testAccEIPConfig = `
resource "alicloud_eip" "foo" {
}
`

const testAccEIPConfigTwo = `
resource "alicloud_eip" "foo" {
    bandwidth = "10"
    internet_charge_type = "PayByBandwidth"
}
`
