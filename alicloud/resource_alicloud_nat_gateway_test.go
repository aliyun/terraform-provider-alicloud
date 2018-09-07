package alicloud

import (
	"fmt"
	"testing"

	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_nat_gateway", &resource.Sweeper{
		Name: "alicloud_nat_gateway",
		F:    testSweepNatGateways,
	})
}

func testSweepNatGateways(region string) error {
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
	}

	var gws []vpc.NatGateway
	req := vpc.CreateDescribeNatGatewaysRequest()
	req.RegionId = conn.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		resp, err := conn.vpcconn.DescribeNatGateways(req)
		if err != nil {
			return fmt.Errorf("Error retrieving Nat Gateways: %s", err)
		}
		if resp == nil || len(resp.NatGateways.NatGateway) < 1 {
			break
		}
		gws = append(gws, resp.NatGateways.NatGateway...)

		if len(resp.NatGateways.NatGateway) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	for _, v := range gws {
		name := v.Name
		id := v.NatGatewayId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Nat Gateway: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting Nat Gateway: %s (%s)", name, id)
		req := vpc.CreateDeleteNatGatewayRequest()
		req.NatGatewayId = id
		req.Force = requests.NewBoolean(true)
		if _, err := conn.vpcconn.DeleteNatGateway(req); err != nil {
			log.Printf("[ERROR] Failed to delete Nat Gateway (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccAlicloudNatGateway_specification(t *testing.T) {
	var nat vpc.NatGateway

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_nat_gateway.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNatGatewayDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNatGatewayConfigSpec,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists(
						"alicloud_nat_gateway.foo", &nat),
					resource.TestCheckResourceAttr(
						"alicloud_nat_gateway.foo",
						"name",
						"tf-testAccNatGatewayConfigSpec"),
					resource.TestCheckResourceAttr(
						"alicloud_nat_gateway.foo",
						"specification",
						"Small"),
				),
			},

			resource.TestStep{
				Config: testAccNatGatewayConfigSpecUpgrade,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists(
						"alicloud_nat_gateway.foo", &nat),
					resource.TestCheckResourceAttr(
						"alicloud_nat_gateway.foo",
						"specification",
						"Middle"),
				),
			},
		},
	})
}

func testAccCheckNatGatewayExists(n string, nat *vpc.NatGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Gateway ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		instance, err := client.DescribeNatGateway(rs.Primary.ID)

		if err != nil {
			return err
		}
		if instance.NatGatewayId != rs.Primary.ID {
			return fmt.Errorf("Nat gateway not found")
		}

		*nat = instance
		return nil
	}
}

func testAccCheckNatGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_nat_gateway" {
			continue
		}

		// Try to find the Nat gateway
		if _, err := client.DescribeNatGateway(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		return fmt.Errorf("Nat gateway %s still exist", rs.Primary.ID)
	}

	return nil
}

const testAccNatGatewayConfigSpec = `
variable "name" {
	default = "tf-testAccNatGatewayConfigSpec"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	specification = "Small"
	name = "${var.name}"
}
`

const testAccNatGatewayConfigSpecUpgrade = `
variable "name" {
	default = "tf-testAccNatGatewayConfigSpec"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	specification = "Middle"
	name = "${var.name}"
}
`
