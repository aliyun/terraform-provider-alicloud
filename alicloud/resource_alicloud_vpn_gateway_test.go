package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_vpn_gateway", &resource.Sweeper{
		Name: "alicloud_vpn_gateway",
		F:    testSweepVPNGateways,
	})
}

func testSweepVPNGateways(region string) error {
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
		"testAcc",
	}

	var gws []vpc.VpnGateway
	req := vpc.CreateDescribeVpnGatewaysRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVpnGateways(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving VPN Gateways: %s", err)
		}
		resp, _ := raw.(*vpc.DescribeVpnGatewaysResponse)
		if resp == nil || len(resp.VpnGateways.VpnGateway) < 1 {
			break
		}
		gws = append(gws, resp.VpnGateways.VpnGateway...)

		if len(resp.VpnGateways.VpnGateway) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	sweeped := false
	for _, v := range gws {
		name := v.Name
		id := v.VpnGatewayId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping VPN Gateway: %s (%s)", name, id)
			continue
		}
		sweeped = true
		log.Printf("[INFO] Deleting VPN Gateway: %s (%s)", name, id)
		req := vpc.CreateDeleteVpnGatewayRequest()
		req.VpnGatewayId = id
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteVpnGateway(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete VPN Gateway (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudVpnGateway_basic(t *testing.T) {
	var vpn vpc.DescribeVpnGatewayResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_vpn_gateway.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpnGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnGatewayExists("alicloud_vpn_gateway.foo", &vpn),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "name", "tf-testAccVpnConfig_create"),
					resource.TestCheckResourceAttrSet(
						"alicloud_vpn_gateway.foo", "vpc_id"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "enable_ssl", "false"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "enable_ipsec", "true"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "description", "test_create_description"),
				),
			},
		},
	})

}

func TestAccAlicloudVpnGateway_update(t *testing.T) {
	var vpn vpc.DescribeVpnGatewayResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpnGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnGatewayExists("alicloud_vpn_gateway.foo", &vpn),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "name", "tf-testAccVpnConfig_create"),
					resource.TestCheckResourceAttrSet(
						"alicloud_vpn_gateway.foo", "vpc_id"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "enable_ssl", "false"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "enable_ipsec", "true"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "description", "test_create_description"),
				),
			},
			{
				Config: testAccVpnConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnGatewayExists("alicloud_vpn_gateway.foo", &vpn),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "name", "tf-testAccVpnConfig_update"),
					resource.TestCheckResourceAttrSet(
						"alicloud_vpn_gateway.foo", "vpc_id"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "enable_ssl", "false"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "enable_ipsec", "true"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "description", "test_update_description"),
				),
			},
		},
	})
}

func testAccCheckVpnGatewayExists(n string, vpn *vpc.DescribeVpnGatewayResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		vpnGatewayService := VpnGatewayService{client}
		instance, err := vpnGatewayService.DescribeVpnGateway(rs.Primary.ID)

		if err != nil {
			return err
		}

		*vpn = instance
		return nil
	}
}

func testAccCheckVpnGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_vpn" {
			continue
		}

		instance, err := vpnGatewayService.DescribeVpnGateway(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Describe Vpn error %#v", err)
		}

		if instance.VpnGatewayId != "" {
			return fmt.Errorf("VPN %s still exist", instance.VpnGatewayId)
		}
	}

	return nil
}

const testAccVpnConfig = `
variable "name" {
	default =  "tf-testAccVpnConfig_create"
}
resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_vpn_gateway" "foo" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.foo.id}"
	bandwidth = "10"
	enable_ssl = false
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}
`

const testAccVpnConfigUpdate = `
variable "name" {
	default =  "tf-testAccVpnConfig_update"
}
resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_vpn_gateway" "foo" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.foo.id}"
	bandwidth = "10"
	enable_ssl = false
	instance_charge_type = "PostPaid"
	description = "test_update_description"
}
`
