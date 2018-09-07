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
)

func init() {
	resource.AddTestSweepers("alicloud_vpn_customer_gateway", &resource.Sweeper{
		Name: "alicloud_vpn_customer_gateway",
		F:    testSweepVPNCustomerGateways,
	})
}

func testSweepVPNCustomerGateways(region string) error {
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

	var gws []vpc.CustomerGateway
	req := vpc.CreateDescribeCustomerGatewaysRequest()
	req.RegionId = conn.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		resp, err := conn.vpcconn.DescribeCustomerGateways(req)
		if err != nil {
			return fmt.Errorf("Error retrieving VPN Customer Gateways: %s", err)
		}
		if resp == nil || len(resp.CustomerGateways.CustomerGateway) < 1 {
			break
		}
		gws = append(gws, resp.CustomerGateways.CustomerGateway...)

		if len(resp.CustomerGateways.CustomerGateway) < PageSizeLarge {
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
		id := v.CustomerGatewayId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping VPN Customer Gateway: %s (%s)", name, id)
			continue
		}
		sweeped = true
		log.Printf("[INFO] Deleting VPN Customer Gateway: %s (%s)", name, id)
		req := vpc.CreateDeleteCustomerGatewayRequest()
		req.CustomerGatewayId = id
		if _, err := conn.vpcconn.DeleteCustomerGateway(req); err != nil {
			log.Printf("[ERROR] Failed to delete VPN Customer Gateway (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

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
						"alicloud_vpn_customer_gateway.foo", "name", "tf-testAccVpnCgwName_Create"),
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
						"alicloud_vpn_customer_gateway.foo", "name", "tf-testAccVpnCgwName_Create"),
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
						"alicloud_vpn_customer_gateway.foo", "name", "tf-testAccVpnCgwName_Update"),
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
			if NotFoundError(err) {
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
	name = "tf-testAccVpnCgwName_Create"
	ip_address = "43.104.22.228"
	description = "testAccVpnCgwDesc_Create"
}
`

const testAccVpnCustomerGatewayConfigUpdate = `
resource "alicloud_vpn_customer_gateway" "foo" {
	name = "tf-testAccVpnCgwName_Update"
	ip_address = "43.104.22.228"
	description = "testAccVpnCgwDesc_Update"
}
`
const testAccVpnCgwIpConfig = `
resource "alicloud_vpn_customer_gateway" "foo" {
	name = "tf-testAccVpnCgwName_Create"
	ip_address = "43.104.22.228"
	description = "testAccVpnCgwDesc_Create"
}
`
const testAccVpnCgwIpConfigUpdate = `
resource "alicloud_vpn_customer_gateway" "foo" {
	name = "tf-testAccVpnCgwName_Update"
	ip_address = "43.104.22.229"
	description = "testAccVpnCgwDesc_Update"
}
`
