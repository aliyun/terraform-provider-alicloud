package alicloud

import (
	"fmt"
	"log"

	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_vpn_router_entry", &resource.Sweeper{
		Name: "alicloud_vpn_router_entry",
		F:    testSweepVPNRouterEntry,
	})
}

func testSweepVPNRouterEntry(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	var gws []vpc.VpnRouteEntry
	request := vpc.CreateDescribeVpnRouteEntriesRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVpnRouteEntries(request)
		})
		if err != nil {
			return WrapError(err)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*vpc.DescribeVpnRouteEntriesResponse)
		if len(response.VpnRouteEntries.VpnRouteEntry) < 1 {
			break
		}
		gws = response.VpnRouteEntries.VpnRouteEntry

		if len(response.VpnRouteEntries.VpnRouteEntry) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	sweeped := false
	for _, v := range gws {
		id := v.VpnInstanceId

		sweeped = true
		log.Printf("[INFO] Deleting VPN route entry: (%s)", id)
		req := vpc.CreateDeleteVpnRouteEntryRequest()
		req.VpnGatewayId = id
		req.RouteDest = v.RouteDest
		req.NextHop = v.NextHop
		req.Weight = requests.NewInteger(v.Weight)

		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteVpnRouteEntry(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete VPN route entry(%s): %s", id, err)
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil

}

func testAccCheckVpnRouteEntryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	vpnRouteEntryService := VpnRouteEntryService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_vpn_route_entry" {
			continue
		}

		fmt.Printf("rs: %v \t %[1]T\n\n\n", rs)
		_, err := vpnRouteEntryService.DescribeVpnRouteEntry(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}

func TestAccAlicloudVpnRouteEntryMulti(t *testing.T) {
	var v VpnState

	resourceId := "alicloud_vpn_route_entry.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &VpnRouteEntryService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 200)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpnRouteEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnRouteEntryConfigAll(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"publish_vpc": fmt.Sprintf("%t", true),
						"description": fmt.Sprintf("tf-testAccVpnRouteEntryConfig%d", rand),
						"weight":      fmt.Sprintf("%d", 100),
					}),
				),
			},
		},
	})
}

func testAccVpnRouteEntryConfigAll(rand int) string {
	return fmt.Sprintf(`
resource "alicloud_vpc" "default" {
  name       = "tf-testAccVpnRouteEntryConfig"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.1.0/24"
  availability_zone = "cn-hangzhou-b"
}

resource "alicloud_vpn_gateway" "default" {
  name       = "tf-testAccVpnRouteEntryConfig"
  vpc_id =  "${alicloud_vpc.default.id}"
  bandwidth = 10
  instance_charge_type = "PostPaid"
  enable_ssl = false
}

resource "alicloud_vpn_connection" "default" {
  name = "tf-testAccVpnRouteEntryConfig"
  customer_gateway_id = "${alicloud_vpn_customer_gateway.default.id}"
  vpn_gateway_id = "${alicloud_vpn_gateway.default.id}"
  local_subnet = ["192.168.2.0/24"]
  remote_subnet = ["192.168.3.0/24"]
}

resource "alicloud_vpn_customer_gateway" "default" {
  name = "tf-testAccVpnRouteEntryConfig"
  ip_address = "192.168.1.1"
}

resource "alicloud_vpn_route_entry" "default" {
  description       = "tf-testAccVpnRouteEntryConfig%d"
  vpn_gateway_id = "${alicloud_vpn_gateway.default.id}"
  route_dest = "12.0.0.2/32"
  next_hop = "${alicloud_vpn_connection.default.id}"
  weight =100
  publish_vpc = true
}
`, rand)
}
