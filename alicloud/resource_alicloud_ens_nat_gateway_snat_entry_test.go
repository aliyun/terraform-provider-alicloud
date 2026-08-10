// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ens NatGatewaySnatEntry. >>> Resource test cases, automatically generated.
// Case Snat_20241218 9626
func TestAccAliCloudEnsNatGatewaySnatEntry_basic9626(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_nat_gateway_snat_entry.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsNatGatewaySnatEntryMap9626)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsNatGatewaySnatEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccens%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsNatGatewaySnatEntryBasicDependence9626)
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
					"snat_entry_name": "测试用例-snat",
					"source_cidr":     "10.0.0.0/8",
					"snat_ip":         "${alicloud_ens_eip.defaultiUbwh0.ip_address}",
					"nat_gateway_id":  "${alicloud_ens_nat_gateway.default2Kn0nu.id}",
					"idle_timeout":    "50",
					"isp_affinity":    "false",
					"eip_affinity":    "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_entry_name": "测试用例-snat",
						"source_cidr":     "10.0.0.0/8",
						"snat_ip":         CHECKSET,
						"nat_gateway_id":  CHECKSET,
						"idle_timeout":    "50",
						"isp_affinity":    "false",
						"eip_affinity":    "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"snat_entry_name": "测试用例-snat2",
					"snat_ip":         "${alicloud_ens_eip.eip2.ip_address}",
					"isp_affinity":    "true",
					"eip_affinity":    "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_entry_name": "测试用例-snat2",
						"snat_ip":         CHECKSET,
						"isp_affinity":    "true",
						"eip_affinity":    "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudEnsNatGatewaySnatEntryMap9626 = map[string]string{
	"status": CHECKSET,
}

func AlicloudEnsNatGatewaySnatEntryBasicDependence9626(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "ens_region_id" {
  default = "cn-hangzhou-44"
}

resource "alicloud_ens_network" "defaultXqhlfk" {
  network_name  = "测试用例-snat"
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "defaultzkXvut" {
  cidr_block    = "10.0.0.0/24"
  vswitch_name  = "测试用例-snat"
  ens_region_id = alicloud_ens_network.defaultXqhlfk.ens_region_id
  network_id    = alicloud_ens_network.defaultXqhlfk.id
}

resource "alicloud_ens_eip" "defaultiUbwh0" {
  bandwidth            = "5"
  payment_type         = "PayAsYouGo"
  ens_region_id        = alicloud_ens_vswitch.defaultzkXvut.ens_region_id
  eip_name             = "测试用例-snat"
  internet_charge_type = "95BandwidthByMonth"
}

resource "alicloud_ens_nat_gateway" "default2Kn0nu" {
  vswitch_id    = alicloud_ens_vswitch.defaultzkXvut.id
  ens_region_id = alicloud_ens_vswitch.defaultzkXvut.ens_region_id
  network_id    = alicloud_ens_vswitch.defaultzkXvut.network_id
  instance_type = "enat.default"
  nat_name      = "测试用例-snat"
}

resource "alicloud_ens_eip_instance_attachment" "defaultlI0M0t" {
  instance_id   = alicloud_ens_nat_gateway.default2Kn0nu.id
  allocation_id = alicloud_ens_eip.defaultiUbwh0.id
  instance_type = "Nat"
  standby       = false
}

resource "alicloud_ens_eip" "eip2" {
  bandwidth            = "5"
  payment_type         = "PayAsYouGo"
  ens_region_id        = var.ens_region_id
  eip_name             = "测试用例-snat2"
  internet_charge_type = "95BandwidthByMonth"
}

resource "alicloud_ens_eip_instance_attachment" "defaultbMMEpj" {
  instance_id   = alicloud_ens_nat_gateway.default2Kn0nu.id
  allocation_id = alicloud_ens_eip.eip2.id
  instance_type = "Nat"
  standby       = false
}


`, name)
}

// Test Ens NatGatewaySnatEntry. <<< Resource test cases, automatically generated.
