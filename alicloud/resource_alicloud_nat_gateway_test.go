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
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_nat_gateway", &resource.Sweeper{
		Name: "alicloud_nat_gateway",
		F:    testSweepNatGateways,
	})
}

func testSweepNatGateways(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
	}

	var gws []vpc.NatGateway
	req := vpc.CreateDescribeNatGatewaysRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeNatGateways(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Nat Gateways: %s", err)
		}
		resp, _ := raw.(*vpc.DescribeNatGatewaysResponse)
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
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteNatGateway(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Nat Gateway (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}
func TestAccAlicloudNatGateway_update(t *testing.T) {
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
			{
				Config: testAccNatGatewayConfigSpec,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists("alicloud_nat_gateway.foo", &nat),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "name", "tf-testAccNatGatewayConfigSpec"),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "specification", "Small"),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "description", ""),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "bandwidth_package_ids", ""),
					resource.TestCheckResourceAttrSet("alicloud_nat_gateway.foo", "forward_table_ids"),
					resource.TestCheckResourceAttrSet("alicloud_nat_gateway.foo", "snat_table_ids"),
				),
			},

			{
				Config: testAccNatGatewayConfigUpName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists("alicloud_nat_gateway.foo", &nat),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "name", "tf-testAccNatGatewayConfigUpName"),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "specification", "Small"),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "description", ""),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "bandwidth_package_ids", ""),
					resource.TestCheckResourceAttrSet("alicloud_nat_gateway.foo", "forward_table_ids"),
					resource.TestCheckResourceAttrSet("alicloud_nat_gateway.foo", "snat_table_ids"),
				),
			},
			{
				Config: testAccNatGatewayConfigUpDesc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists("alicloud_nat_gateway.foo", &nat),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "name", "tf-testAccNatGatewayConfigUpName"),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "specification", "Small"),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "description", "testAccNatGatewayConfig_description"),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "bandwidth_package_ids", ""),
					resource.TestCheckResourceAttrSet("alicloud_nat_gateway.foo", "forward_table_ids"),
					resource.TestCheckResourceAttrSet("alicloud_nat_gateway.foo", "snat_table_ids"),
				),
			},
			{
				Config: testAccNatGatewayConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists("alicloud_nat_gateway.foo", &nat),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "name", "tf-testAccNatGatewayConfigUpdate"),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "specification", "Small"),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "description", "testAccNatGatewayConfigUpdate_description"),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "bandwidth_package_ids", ""),
					resource.TestCheckResourceAttrSet("alicloud_nat_gateway.foo", "forward_table_ids"),
					resource.TestCheckResourceAttrSet("alicloud_nat_gateway.foo", "snat_table_ids"),
				),
			},
		},
	})
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
			{
				Config: testAccNatGatewayConfigSpec,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists("alicloud_nat_gateway.foo", &nat),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "name", "tf-testAccNatGatewayConfigSpec"),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "specification", "Small"),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "description", ""),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "bandwidth_package_ids", ""),
					resource.TestCheckResourceAttrSet("alicloud_nat_gateway.foo", "forward_table_ids"),
					resource.TestCheckResourceAttrSet("alicloud_nat_gateway.foo", "snat_table_ids"),
				),
			},

			{
				Config: testAccNatGatewayConfigSpecUpgrade,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists("alicloud_nat_gateway.foo", &nat),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "name", "tf-testAccNatGatewayConfigSpec"),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "specification", "Middle"),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "bandwidth_package_ids", ""),
					resource.TestCheckResourceAttr("alicloud_nat_gateway.foo", "description", ""),
					resource.TestCheckResourceAttrSet("alicloud_nat_gateway.foo", "forward_table_ids"),
					resource.TestCheckResourceAttrSet("alicloud_nat_gateway.foo", "snat_table_ids"),
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

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		vpcService := VpcService{client}
		instance, err := vpcService.DescribeNatGateway(rs.Primary.ID)

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
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_nat_gateway" {
			continue
		}

		// Try to find the Nat gateway
		if _, err := vpcService.DescribeNatGateway(rs.Primary.ID); err != nil {
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
const testAccNatGatewayConfigUpName = `
variable "name" {
	default = "tf-testAccNatGatewayConfigUpName"
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
const testAccNatGatewayConfigUpDesc = `
variable "name" {
	default = "tf-testAccNatGatewayConfigUpName"
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
	description="testAccNatGatewayConfig_description"
	name = "${var.name}"
}
`
const testAccNatGatewayConfigUpdate = `
variable "name" {
	default = "tf-testAccNatGatewayConfigUpdate"
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
	description="testAccNatGatewayConfigUpdate_description"
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
