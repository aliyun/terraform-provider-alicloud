package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudVPCNatIp_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_nat_ip.default"
	ra := resourceAttrInit(resourceId, AlicloudVPCNatIpMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcNatIp")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcnatip%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPCNatIpBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"nat_ip":             "192.168.0.37",
					"nat_ip_cidr":        "${alicloud_vpc_nat_ip_cidr.default.nat_ip_cidr}",
					"nat_gateway_id":     "${alicloud_nat_gateway.default.id}",
					"nat_ip_description": "${var.name}",
					"nat_ip_name":        "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nat_ip":             "192.168.0.37",
						"nat_ip_cidr":        CHECKSET,
						"nat_gateway_id":     CHECKSET,
						"nat_ip_description": name,
						"nat_ip_name":        name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nat_ip_description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nat_ip_description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nat_ip_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nat_ip_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nat_ip_description": "${var.name}",
					"nat_ip_name":        "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nat_ip_description": name,
						"nat_ip_name":        name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

func TestAccAlicloudVPCNatIp_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_nat_ip.default"
	ra := resourceAttrInit(resourceId, AlicloudVPCNatIpMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcNatIp")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcnatip%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPCNatIpBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"nat_ip":             "192.168.0.37",
					"nat_ip_cidr":        "${alicloud_vpc_nat_ip_cidr.default.nat_ip_cidr}",
					"nat_gateway_id":     "${alicloud_nat_gateway.default.id}",
					"nat_ip_description": "${var.name}",
					"nat_ip_name":        "${var.name}",
					"dry_run":            "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nat_ip":             "192.168.0.37",
						"nat_ip_cidr":        CHECKSET,
						"nat_gateway_id":     CHECKSET,
						"nat_ip_description": name,
						"nat_ip_name":        name,
						"dry_run":            "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudVPCNatIpMap0 = map[string]string{
	"status":         CHECKSET,
	"nat_ip_cidr_id": NOSET,
}

func AlicloudVPCNatIpBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}

data "alicloud_vpcs" "default" {
 cidr_block = "172.16.0.0/12"
}

resource "alicloud_vpc" "default" {
    count = length(data.alicloud_vpcs.default.ids) > 0 ? 0 : 1
	vpc_name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

data "alicloud_vswitches" "default" {
  vpc_id = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
  zone_id = data.alicloud_zones.default.zones.0.id
}
resource "alicloud_vswitch" "default" {
  count         = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id        = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
  cidr_block    = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  zone_id     	= data.alicloud_zones.default.zones.0.id
  vswitch_name  = var.name
}

resource "alicloud_nat_gateway" "default" {
	vpc_id = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
	internet_charge_type = "PayByLcu"
	nat_gateway_name = var.name
    description = "${var.name}_description"
	nat_type = "Enhanced"
	vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : alicloud_vswitch.default[0].id
	network_type = "intranet"
}

resource "alicloud_vpc_nat_ip_cidr" "default" {
	nat_ip_cidr = "192.168.0.0/16"
	nat_gateway_id =  alicloud_nat_gateway.default.id
	nat_ip_cidr_description = var.name
	nat_ip_cidr_name = var.name
}
`, name)
}
