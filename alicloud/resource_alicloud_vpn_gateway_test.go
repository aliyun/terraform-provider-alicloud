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
	resource.AddTestSweepers("alicloud_vpn_gateway", &resource.Sweeper{
		Name: "alicloud_vpn_gateway",
		F:    testSweepVPNGateways,
		Dependencies: []string{
			"alicloud_ssl_vpn_server",
			"alicloud_ssl_vpn_client_cert",
		},
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
			log.Printf("[ERROR] Error retrieving VPN Gateways: %s", err)
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
		if !sweepAll() {
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
		time.Sleep(10 * time.Second)
	}
	return nil
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
			return WrapError(err)
		}

		if instance.VpnGatewayId != "" {
			return WrapError(Error("VPN %s still exist", instance.VpnGatewayId))
		}
	}

	return nil
}

// At present, some properties of this resource do not support modification, including: period, bandwidth, enable_ipsec,
// enable_ssl, ssl_connections etc.
func TestAccAliCloudVPNGatewayBasic(t *testing.T) {
	var v vpc.DescribeVpnGatewayResponse

	resourceId := "alicloud_vpn_gateway.default"
	ra := resourceAttrInit(resourceId, testAccVpnGatewayCheckMap)
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
		CheckDestroy:  testAccCheckVpnGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccVpnConfig%d", rand),
					}),
				),
			},
			{
				Config: testAccVpnConfig_name(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccVpnConfig%d_change", rand),
					}),
				),
			},
			{
				Config: testAccVpnConfig_description(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccVpnConfig%d_description", rand),
					}),
				),
			},
			{
				Config: testAccVpnConfig_all(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("tf-testAccVpnConfig%d", rand),
						"description": fmt.Sprintf("tf-testAccVpnConfig%d", rand),
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "auto_pay"},
			},
		},
	})

}

func testAccVpnConfigBasic(rand int) string {
	return fmt.Sprintf(`
variable "name" {
default =  "tf-testAccVpnConfig%d"
}

data "alicloud_zones" "default" {
available_resource_creation= "VSwitch"
}

data "alicloud_vpcs" "default" {
name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
vpc_id = data.alicloud_vpcs.default.ids.0
zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_vpn_gateway" "default" {
name = "${var.name}"
vpc_id = data.alicloud_vpcs.default.ids.0
bandwidth = "10"
enable_ssl = true
instance_charge_type = "PrePaid"
vswitch_id = local.vswitch_id
}
`, rand)
}

func testAccVpnConfig_name(rand int) string {
	return fmt.Sprintf(`
variable "name" {
default =  "tf-testAccVpnConfig%d"
}

data "alicloud_zones" "default" {
available_resource_creation= "VSwitch"
}

data "alicloud_vpcs" "default" {
name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
vpc_id = data.alicloud_vpcs.default.ids.0
zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_vpn_gateway" "default" {
name = "${var.name}_change"
vpc_id = data.alicloud_vpcs.default.ids.0
bandwidth = "10"
enable_ssl = true
instance_charge_type = "PrePaid"
vswitch_id = local.vswitch_id
}
`, rand)
}

func testAccVpnConfig_description(rand int) string {
	return fmt.Sprintf(`
variable "name" {
default =  "tf-testAccVpnConfig%d"
}

data "alicloud_zones" "default" {
available_resource_creation= "VSwitch"
}

data "alicloud_vpcs" "default" {
name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
vpc_id = data.alicloud_vpcs.default.ids.0
zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_vpn_gateway" "default" {
name = "${var.name}_change"
vpc_id = data.alicloud_vpcs.default.ids.0
bandwidth = "10"
enable_ssl = true
instance_charge_type = "PrePaid"
description = "${var.name}_description"
vswitch_id = local.vswitch_id
}
`, rand)
}

func testAccVpnConfig_all(rand int) string {
	return fmt.Sprintf(`
variable "name" {
default =  "tf-testAccVpnConfig%d"
}

data "alicloud_zones" "default" {
available_resource_creation= "VSwitch"
}

data "alicloud_vpcs" "default" {
name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
vpc_id = data.alicloud_vpcs.default.ids.0
zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_vpn_gateway" "default" {
name = "${var.name}"
vpc_id = data.alicloud_vpcs.default.ids.0
bandwidth = "10"
enable_ssl = true
instance_charge_type = "PrePaid"
description = "${var.name}"
vswitch_id = local.vswitch_id
}
`, rand)
}

var testAccVpnGatewayCheckMap = map[string]string{
	"vpc_id":       CHECKSET,
	"bandwidth":    "10",
	"enable_ssl":   "true",
	"enable_ipsec": "true",
	"description":  "",
	"vswitch_id":   CHECKSET,
}

func TestAccAliCloudVPNGateway_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_gateway.default"
	ra := resourceAttrInit(resourceId, AlicloudVpnGatewayMap3)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpnGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpnGatewayBasicDependence3)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":               "${data.alicloud_vpcs.default.ids.0}",
					"name":                 name,
					"vswitch_id":           "${data.alicloud_vswitches.default.vswitches.0.id}",
					"description":          name,
					"bandwidth":            "10",
					"enable_ssl":           "true",
					"auto_propagate":       "true",
					"period":               "1",
					"instance_charge_type": "PrePaid",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"network_type": "public",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":               CHECKSET,
						"name":                 name,
						"vswitch_id":           CHECKSET,
						"description":          name,
						"bandwidth":            "10",
						"enable_ssl":           "true",
						"auto_propagate":       "true",
						"period":               "1",
						"instance_charge_type": "PrePaid",
						"tags.%":               "2",
						"network_type":         "public",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "TestF-update",
						"From":    "update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_propagate": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_propagate": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "auto_pay"},
			},
		},
	})
}

var AlicloudVpnGatewayMap3 = map[string]string{}

func AlicloudVpnGatewayBasicDependence3(name string) string {
	return fmt.Sprintf(`
variable "name" {
default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
}

`, name)
}

// Test VPNGateway VPNGateway. >>> Resource test cases, automatically generated.
// Case 3675
func TestAccAliCloudVPNGatewayVPNGateway_basic3675(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_gateway.default"
	ra := resourceAttrInit(resourceId, AlicloudVPNGatewayVPNGatewayMap3675)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VPNGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVPNGatewayVPNGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpngatewayvpngateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudVPNGatewayVPNGatewayBasicDependence3675)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":        "5",
					"vpc_id":           "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":       "${data.alicloud_vswitches.default.vswitches.0.id}",
					"vpn_gateway_name": name,
					"enable_ssl":       "true",
					"enable_ipsec":     "true",
					"vpn_type":         "Normal",
					"ssl_connections":  "5",
					"payment_type":     "Subscription",
					"period":           "1",
					"network_type":     "public",
					"auto_pay":         "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":        "5",
						"vpc_id":           CHECKSET,
						"vswitch_id":       CHECKSET,
						"vpn_gateway_name": name,
						"enable_ssl":       "true",
						"enable_ipsec":     "true",
						"vpn_type":         "Normal",
						"ssl_connections":  "5",
						"payment_type":     "Subscription",
						"period":           "1",
						"network_type":     "public",
						"auto_pay":         "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test_vpn",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test_vpn",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpn_gateway_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpn_gateway_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tes_vpn_new",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tes_vpn_new",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpn_gateway_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpn_gateway_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_propagate": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_propagate": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test_vpn",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test_vpn",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpn_gateway_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpn_gateway_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_connections":  "5",
					"enable_ipsec":     "true",
					"enable_ssl":       "true",
					"auto_pay":         "true",
					"vpn_type":         "Normal",
					"description":      "test_vpn",
					"vpn_gateway_name": name + "_update",
					"network_type":     "public",
					"bandwidth":        "5",
					"vswitch_id":       "${data.alicloud_vswitches.default.vswitches.0.id}",
					"vpc_id":           "${data.alicloud_vpcs.default.ids.0}",
					"period":           "1",
					"payment_type":     "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_connections":  "5",
						"enable_ipsec":     "true",
						"enable_ssl":       "true",
						"auto_pay":         "true",
						"vpn_type":         "Normal",
						"description":      "test_vpn",
						"vpn_gateway_name": name + "_update",
						"network_type":     "public",
						"bandwidth":        "5",
						"vswitch_id":       CHECKSET,
						"vpc_id":           CHECKSET,
						"period":           "1",
						"payment_type":     "Subscription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "enable_ipsec", "enable_ssl", "period", "ssl_connections"},
			},
		},
	})
}

var AlicloudVPNGatewayVPNGatewayMap3675 = map[string]string{
	"payment_type": CHECKSET,
	"status":       CHECKSET,
	"create_time":  CHECKSET,
}

func AliCloudVPNGatewayVPNGatewayBasicDependence3675(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
}


`, name)
}

// Case 3675  twin
func TestAccAliCloudVPNGatewayVPNGateway_basic3675_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_gateway.default"
	ra := resourceAttrInit(resourceId, AlicloudVPNGatewayVPNGatewayMap3675)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VPNGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVPNGatewayVPNGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpngatewayvpngateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudVPNGatewayVPNGatewayBasicDependence3675)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_connections":  "5",
					"enable_ipsec":     "true",
					"enable_ssl":       "true",
					"auto_pay":         "true",
					"vpn_type":         "Normal",
					"description":      "test_vpn",
					"vpn_gateway_name": name,
					"network_type":     "public",
					"bandwidth":        "5",
					"vswitch_id":       "${data.alicloud_vswitches.default.vswitches.0.id}",
					"vpc_id":           "${data.alicloud_vpcs.default.ids.0}",
					"period":           "1",
					"payment_type":     "Subscription",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_connections":  "5",
						"enable_ipsec":     "true",
						"enable_ssl":       "true",
						"auto_pay":         "true",
						"vpn_type":         "Normal",
						"description":      "test_vpn",
						"vpn_gateway_name": name,
						"network_type":     "public",
						"bandwidth":        "5",
						"vswitch_id":       CHECKSET,
						"vpc_id":           CHECKSET,
						"period":           "1",
						"payment_type":     "Subscription",
						"tags.%":           "2",
						"tags.Created":     "TF",
						"tags.For":         "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "enable_ipsec", "enable_ssl", "period", "ssl_connections"},
			},
		},
	})
}

// Test VPNGateway VPNGateway. <<< Resource test cases, automatically generated.
