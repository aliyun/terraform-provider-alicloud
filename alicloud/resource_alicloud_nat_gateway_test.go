package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_nat_gateway", &resource.Sweeper{
		Name: "alicloud_nat_gateway",
		F:    testSweepNatGateways,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_cs_cluster",
		},
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

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}
	service := VpcService{client}
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
		// If a nat gateway name is not set successfully, it should be fetched by vpc name and deleted.
		if skip {
			if need, err := service.needSweepVpc(v.VpcId, ""); err == nil {
				skip = !need
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Nat Gateway: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting Nat Gateway: %s (%s)", name, id)
		if err := service.sweepNatGateway(id); err != nil {
			log.Printf("[ERROR] Failed to delete Nat Gateway (%s (%s)): %s", name, id, err)
		}
	}
	return nil
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

func TestAccAlicloudNatGatewayBasic(t *testing.T) {
	var v vpc.NatGateway
	resourceId := "alicloud_nat_gateway.default"
	ra := resourceAttrInit(resourceId, testAccCheckNatGatewayBasicMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatGatewayConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccNatGatewayConfig%d", rand),
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period"},
			},
			{
				Config: testAccNatGatewayConfig_type(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": "PostPaid",
					}),
				),
			},
			{
				Config: testAccNatGatewayConfig_name(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccNatGatewayConfig%d_change", rand),
					}),
				),
			},
			{
				Config: testAccNatGatewayConfig_description(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccNatGatewayConfig%d_description", rand),
					}),
				),
			},
			{
				Config: testAccNatGatewayConfig_specification(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"specification": "Middle",
					}),
				),
			},
			{
				Config: testAccNatGatewayConfig_all(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"specification": "Small",
						"name":          fmt.Sprintf("tf-testAccNatGatewayConfig%d_all", rand),
						"description":   fmt.Sprintf("tf-testAccNatGatewayConfig%d_description_all", rand),
					}),
				),
			},
		},
	})
}

func TestAccAlicloudNatGatewayEnhanced(t *testing.T) {
	var v vpc.NatGateway
	resourceId := "alicloud_nat_gateway.default"
	ra := resourceAttrInit(resourceId, testAccCheckNatGatewayBasicMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatGatewayConfig_natType(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nat_type": "Enhanced",
						"name":     fmt.Sprintf("tf-testAccNatGatewayConfig%d", rand),
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period"},
			},
		},
	})
}

func testAccNatGatewayConfigBasic(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccNatGatewayConfig%d"
}

data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_nat_gateway" "default" {
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
	name = "${var.name}"
}
`, rand)
}

func testAccNatGatewayConfig_type(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccNatGatewayConfig%d"
}

data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_nat_gateway" "default" {
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
	name = "${var.name}"
	instance_charge_type = "PostPaid"
}
`, rand)
}

func testAccNatGatewayConfig_name(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccNatGatewayConfig%d"
}

data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_nat_gateway" "default" {
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
	name = "${var.name}_change"
}
`, rand)
}

func testAccNatGatewayConfig_description(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccNatGatewayConfig%d"
}

data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_nat_gateway" "default" {
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
	name = "${var.name}_change"
	description = "${var.name}_description"
}
`, rand)
}

func testAccNatGatewayConfig_specification(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccNatGatewayConfig%d"
}

data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_nat_gateway" "default" {
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
	name = "${var.name}_change"
	description = "${var.name}_description"
	specification = "Middle"
}
`, rand)
}

func testAccNatGatewayConfig_natType(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccNatGatewayConfig%d"
}

resource "alicloud_vpc" "default" {
 name       = var.name
 cidr_block = "10.0.0.0/8"
}

data "alicloud_enhanced_nat_available_zones" "default"{
}

resource "alicloud_vswitch" "default" {
 name              = var.name
 availability_zone = data.alicloud_enhanced_nat_available_zones.default.zones.0.zone_id
 cidr_block        = "10.10.0.0/20"
 vpc_id            = alicloud_vpc.default.id
}

resource "alicloud_nat_gateway" "default" {
 depends_on           = [alicloud_vswitch.default]
 vpc_id               = alicloud_vpc.default.id
 specification        = "Small"
 name                 = var.name
 instance_charge_type = "PostPaid"
 vswitch_id           = alicloud_vswitch.default.id
 nat_type             = "Enhanced"
}
`, rand)
}

func testAccNatGatewayConfig_all(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccNatGatewayConfig%d"
}

data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_nat_gateway" "default" {
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
	name = "${var.name}_all"
	description = "${var.name}_description_all"
	specification = "Small"
}
`, rand)
}

var testAccCheckNatGatewayBasicMap = map[string]string{
	"name":                  "tf-testAccNatGatewayConfigSpec",
	"specification":         "Small",
	"description":           "",
	"bandwidth_package_ids": "",
	"forward_table_ids":     CHECKSET,
	"snat_table_ids":        CHECKSET,
}
