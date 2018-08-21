package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudVpnCustomerGateway_basic(t *testing.T) {
	var vpnCgw vpc.DescribeCustomerGatewayResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_vpn_customer_gateway.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpnCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpnCustomerGatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnCustomerGatewayExists("alicloud_vpn_customer_gateway.foo", &vpnCgw),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_customer_gateway.foo", "name", "testAccVpnCgwName_Create"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_customer_gateway.foo", "description", "testAccVpnCgwDesc_Create"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_customer_gateway.foo", "ip_address", "43.104.22.228"),
				),
			},
		},
	})
}

func TestAccAlicloudVpnCustomerGateway_update(t *testing.T) {
	var vpnCgw vpc.DescribeCustomerGatewayResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpnCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpnCustomerGatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnCustomerGatewayExists("alicloud_vpn_customer_gateway.foo", &vpnCgw),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_customer_gateway.foo", "name", "testAccVpnCgwName_Create"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_customer_gateway.foo", "description", "testAccVpnCgwDesc_Create"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_customer_gateway.foo", "ip_address", "43.104.22.228"),
				),
			},
			resource.TestStep{
				Config: testAccVpnCustomerGatewayConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnCustomerGatewayExists("alicloud_vpn_customer_gateway.foo", &vpnCgw),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_customer_gateway.foo", "name", "testAccVpnCgwName_Update"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_customer_gateway.foo", "description", "testAccVpnCgwDesc_Update"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_customer_gateway.foo", "ip_address", "43.104.22.228"),
				),
			},
		},
	})
}

func TestAccAlicloudVpnCustomerGateway_updateIp(t *testing.T) {
	var vpnCgw vpc.DescribeCustomerGatewayResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpnCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpnCgwIpConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnCustomerGatewayExists("alicloud_vpn_customer_gateway.foo", &vpnCgw),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_customer_gateway.foo", "ip_address", "43.104.22.228"),
				),
			},
			resource.TestStep{
				Config: testAccVpnCgwIpConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnCustomerGatewayExists("alicloud_vpn_customer_gateway.foo", &vpnCgw),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_customer_gateway.foo", "ip_address", "43.104.22.229"),
				),
			},
		},
	})
}

func testAccCheckVpnCustomerGatewayExists(n string, vpn *vpc.DescribeCustomerGatewayResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		instance, err := client.DescribeCustomerGateway(rs.Primary.ID)

		if err != nil {
			return err
		}

		*vpn = instance
		return nil
	}
}

func testAccCheckVpnCustomerGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_vpn_customer_gateway" {
			continue
		}

		instance, err := client.DescribeCustomerGateway(rs.Primary.ID)

		if err != nil {
			if IsExceptedError(err, CgwNotFound) {
				continue
			}
			return err
		}

		if instance.CustomerGatewayId != "" {
			return fmt.Errorf("VPN Customer Gateway %s still exist", instance.CustomerGatewayId)
		}
	}

	return nil
}

const testAccVpnCustomerGatewayConfig = `
resource "alicloud_vpn_customer_gateway" "foo" {
  name = "testAccVpnCgwName_Create"
  ip_address = "43.104.22.228"
  description = "testAccVpnCgwDesc_Create"
}
`

const testAccVpnCustomerGatewayConfigUpdate = `
resource "alicloud_vpn_customer_gateway" "foo" {
  name = "testAccVpnCgwName_Update"
  ip_address = "43.104.22.228"
  description = "testAccVpnCgwDesc_Update"
}
`
const testAccVpnCgwIpConfig = `
resource "alicloud_vpn_customer_gateway" "foo" {
  name = "testAccVpnCgwName_Create"
  ip_address = "43.104.22.228"
  description = "testAccVpnCgwDesc_Create"
}
`
const testAccVpnCgwIpConfigUpdate = `
resource "alicloud_vpn_customer_gateway" "foo" {
  name = "testAccVpnCgwName_Update"
  ip_address = "43.104.22.229"
  description = "testAccVpnCgwDesc_Update"
}
`
