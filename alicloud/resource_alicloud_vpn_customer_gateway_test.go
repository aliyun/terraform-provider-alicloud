package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_vpn_customer_gateway", &resource.Sweeper{
		Name: "alicloud_vpn_customer_gateway",
		F:    testSweepVPNCustomerGateways,
	})
}

func testSweepVPNCustomerGateways(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
	}

	var gws []vpc.CustomerGateway
	request := vpc.CreateDescribeCustomerGatewaysRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeCustomerGateways(request)
		})
		if err != nil {
			return WrapError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeCustomerGatewaysResponse)
		if len(response.CustomerGateways.CustomerGateway) < 1 {
			break
		}
		gws = append(gws, response.CustomerGateways.CustomerGateway...)

		if len(response.CustomerGateways.CustomerGateway) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
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
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteCustomerGateway(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete VPN Customer Gateway (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func testAccCheckVpnCustomerGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_vpn_customer_gateway" {
			continue
		}

		_, err := vpnGatewayService.DescribeVpnCustomerGateway(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}

func TestAccAlicloudVPNCustomerGatewayBasic(t *testing.T) {
	var v vpc.DescribeCustomerGatewayResponse

	resourceId := "alicloud_vpn_customer_gateway.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &VpnGatewayService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpnCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnCustomerGatewayConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("tf-testAccVpnCgwName%d", rand),
						"description": "",
						"ip_address":  "43.104.22.228",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccVpnCustomerGatewayConfig_name(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccVpnCgwName%d_change", rand),
					}),
				),
			},
			{
				Config: testAccVpnCustomerGatewayConfig_description(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccVpnCgwName%d_change", rand),
					}),
				),
			},
			{
				Config: testAccVpnCustomerGatewayConfig_all(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("tf-testAccVpnCgwName%d", rand),
						"description": fmt.Sprintf("tf-testAccVpnCgwName%d", rand),
					}),
				),
			},
		},
	})
}

func TestAccAlicloudVPNCustomerGatewayMulti(t *testing.T) {
	var v vpc.DescribeCustomerGatewayResponse

	resourceId := "alicloud_vpn_customer_gateway.default.4"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &VpnGatewayService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpnCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnCustomerGatewayConfig_multi(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("tf-testAccVpnCgwName%d", rand),
						"description": "",
						"ip_address":  "43.104.22.225",
					}),
				),
			},
		},
	})
}

func testAccVpnCustomerGatewayConfigBasic(rand int) string {
	return fmt.Sprintf(`
resource "alicloud_vpn_customer_gateway" "default" {
	name = "tf-testAccVpnCgwName%d"
	ip_address = "43.104.22.228"
}
`, rand)
}

func testAccVpnCustomerGatewayConfig_name(rand int) string {
	return fmt.Sprintf(`
resource "alicloud_vpn_customer_gateway" "default" {
	name = "tf-testAccVpnCgwName%d_change"
	ip_address = "43.104.22.228"
}
`, rand)
}

func testAccVpnCustomerGatewayConfig_description(rand int) string {
	return fmt.Sprintf(`
resource "alicloud_vpn_customer_gateway" "default" {
	name = "tf-testAccVpnCgwName%d_change"
	ip_address = "43.104.22.228"
	description = "tf-testAccVpnCgwName%d_change"
}
`, rand, rand)
}

func testAccVpnCustomerGatewayConfig_all(rand int) string {
	return fmt.Sprintf(`
resource "alicloud_vpn_customer_gateway" "default" {
	name = "tf-testAccVpnCgwName%d"
	ip_address = "43.104.22.228"
	description = "tf-testAccVpnCgwName%d"
}
`, rand, rand)
}

func testAccVpnCustomerGatewayConfig_multi(rand int) string {
	return fmt.Sprintf(`
resource "alicloud_vpn_customer_gateway" "default" {
	count = 5
	name = "tf-testAccVpnCgwName%d"
	ip_address = "43.104.22.${ 221 + count.index }"
}
`, rand)
}

func TestAccAlicloudVPNCustomerGateway_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_customer_gateway.default"
	ra := resourceAttrInit(resourceId, AlicloudVpnCustomerGatewayMap3)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpnCustomerGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpnCustomerGatewayBasicDependence3)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":       name,
					"ip_address": "192.104.22.228",
					"asn":        "12345",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":       name,
						"ip_address": "192.104.22.228",
						"asn":        "12345",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudVpnCustomerGatewayMap3 = map[string]string{}

func AlicloudVpnCustomerGatewayBasicDependence3(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_vpcs" "default"	{
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
}

`, name)
}
