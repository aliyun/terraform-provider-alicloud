package alicloud

import (
	"fmt"
	"log"
	"testing"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_eip", &resource.Sweeper{
		Name: "alicloud_eip",
		F:    testSweepEips,
	})
}

func testSweepEips(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	conn := client.(*AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
	}

	var eips []vpc.EipAddress
	req := vpc.CreateDescribeEipAddressesRequest()
	req.RegionId = conn.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		resp, err := conn.vpcconn.DescribeEipAddresses(req)
		if err != nil {
			return fmt.Errorf("Error retrieving EIPs: %s", err)
		}
		if resp == nil || len(resp.EipAddresses.EipAddress) < 1 {
			break
		}
		eips = append(eips, resp.EipAddresses.EipAddress...)

		if len(resp.EipAddresses.EipAddress) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	for _, v := range eips {
		name := v.Name
		id := v.AllocationId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping EIP: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting EIP: %s (%s)", name, id)
		req := vpc.CreateReleaseEipAddressRequest()
		req.AllocationId = id
		if _, err := conn.vpcconn.ReleaseEipAddress(req); err != nil {
			log.Printf("[ERROR] Failed to delete EIP (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

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
					resource.TestCheckResourceAttr("alicloud_eip.foo", "name", "tf-testAccEIPConfig"),
				),
			},
			resource.TestStep{
				Config: testAccEIPConfigTwo,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(
						"alicloud_eip.foo", &eip),
					testAccCheckEIPAttributes(&eip),
					resource.TestCheckResourceAttr("alicloud_eip.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("alicloud_eip.foo", "name", "tf-testAccEIPConfigTwo"),
					resource.TestCheckResourceAttr("alicloud_eip.foo", "description", "testAccEIPConfigTwo"),
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
	name = "tf-testAccEIPConfig"
}
`

const testAccEIPConfigTwo = `
resource "alicloud_eip" "foo" {
    bandwidth = "10"
    internet_charge_type = "PayByTraffic"
    name = "tf-testAccEIPConfigTwo"
    description = "testAccEIPConfigTwo"
}
`
