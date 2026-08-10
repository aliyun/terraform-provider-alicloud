// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ens NatGatewayForwardEntry. >>> Resource test cases, automatically generated.
// Case Dnat规则_20241218_TCP 9627
func TestAccAliCloudEnsNatGatewayForwardEntry_basic9627(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_nat_gateway_forward_entry.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsNatGatewayForwardEntryMap9627)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsNatGatewayForwardEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccens%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsNatGatewayForwardEntryBasicDependence9627)
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
					"external_port":      "100/200",
					"external_ip":        "${alicloud_ens_eip.defaultLQgQB6.ip_address}",
					"ip_protocol":        "TCP",
					"internal_port":      "100/200",
					"health_check_port":  "150",
					"nat_gateway_id":     "${alicloud_ens_nat_gateway.defaultlZ7YKl.id}",
					"forward_entry_name": "测试用例-dnat",
					"internal_ip":        "${alicloud_ens_instance.defaulth6OQ3p.private_ip_address}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"external_port":      "100/200",
						"external_ip":        CHECKSET,
						"ip_protocol":        "TCP",
						"internal_port":      "100/200",
						"health_check_port":  "150",
						"nat_gateway_id":     CHECKSET,
						"forward_entry_name": "测试用例-dnat",
						"internal_ip":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"external_port":      "151/151",
					"external_ip":        "${alicloud_ens_eip.eip2.ip_address}",
					"ip_protocol":        "UDP",
					"internal_port":      "151/151",
					"health_check_port":  "151",
					"forward_entry_name": "test2",
					"internal_ip":        "${alicloud_ens_instance.instance2.private_ip_address}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"external_port":      "151/151",
						"external_ip":        CHECKSET,
						"ip_protocol":        "UDP",
						"internal_port":      "151/151",
						"health_check_port":  "151",
						"forward_entry_name": "test2",
						"internal_ip":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_protocol": "Any",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_protocol": "Any",
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

var AlicloudEnsNatGatewayForwardEntryMap9627 = map[string]string{
	"status": CHECKSET,
}

func AlicloudEnsNatGatewayForwardEntryBasicDependence9627(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "ens_region_id" {
  default = "cn-hangzhou-44"
}

resource "alicloud_ens_network" "default6T9qR2" {
  network_name  = "测试用例_Dnat"
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "default5BAAN2" {
  cidr_block    = "10.0.6.0/24"
  vswitch_name  = "测试用例-dnat"
  ens_region_id = alicloud_ens_network.default6T9qR2.ens_region_id
  network_id    = alicloud_ens_network.default6T9qR2.id
}

resource "alicloud_ens_nat_gateway" "defaultlZ7YKl" {
  vswitch_id    = alicloud_ens_vswitch.default5BAAN2.id
  ens_region_id = alicloud_ens_vswitch.default5BAAN2.ens_region_id
  network_id    = alicloud_ens_vswitch.default5BAAN2.network_id
  instance_type = "enat.default"
  nat_name      = "测试用例-dnat"
}

resource "alicloud_ens_eip" "defaultLQgQB6" {
  bandwidth            = "5"
  payment_type         = "PayAsYouGo"
  ens_region_id        = var.ens_region_id
  eip_name             = "测试用例-dnat"
  internet_charge_type = "95BandwidthByMonth"
}

resource "alicloud_ens_eip_instance_attachment" "defaultc19VZl" {
  instance_id   = alicloud_ens_nat_gateway.defaultlZ7YKl.id
  allocation_id = alicloud_ens_eip.defaultLQgQB6.id
  instance_type = "Nat"
}

resource "alicloud_ens_instance" "defaulth6OQ3p" {
  auto_renew = false
  system_disk {
    size     = "20"
    category = "cloud_efficiency"
  }
  scheduling_strategy        = "Concentrate"
  schedule_area_level        = "Region"
  image_id                   = "centos_6_08_64_20G_alibase_20171208"
  payment_type               = "Subscription"
  instance_type              = "ens.sn1.stiny"
  password_inherit           = false
  password                   = "12345678abcABC"
  status                     = "Running"
  amount                     = "1"
  vswitch_id                 = alicloud_ens_vswitch.default5BAAN2.id
  internet_charge_type       = "95BandwidthByMonth"
  instance_name              = "测试用例-dnat"
  internet_max_bandwidth_out = "0"
  unique_suffix              = false
  auto_use_coupon            = "true"
  public_ip_identification   = false
  instance_charge_strategy   = "PriceHighPriority"
  ens_region_id              = var.ens_region_id
  period_unit                = "Month"
}

resource "alicloud_ens_eip" "eip2" {
  bandwidth            = "5"
  payment_type         = "PayAsYouGo"
  ens_region_id        = var.ens_region_id
  eip_name             = "测试用例-dnat2"
  internet_charge_type = "95BandwidthByMonth"
}

resource "alicloud_ens_eip_instance_attachment" "default4Ph8bE" {
  instance_id   = alicloud_ens_nat_gateway.defaultlZ7YKl.id
  allocation_id = alicloud_ens_eip.eip2.id
  instance_type = "Nat"
}

resource "alicloud_ens_instance" "instance2" {
  auto_renew = false
  system_disk {
    size     = "20"
    category = "cloud_efficiency"
  }
  scheduling_strategy        = "Concentrate"
  schedule_area_level        = "Region"
  image_id                   = "centos_6_08_64_20G_alibase_20171208"
  payment_type               = "Subscription"
  instance_type              = "ens.sn1.stiny"
  password_inherit           = false
  password                   = "12345678abcABC"
  status                     = "Running"
  amount                     = "1"
  vswitch_id                 = alicloud_ens_vswitch.default5BAAN2.id
  internet_charge_type       = "95BandwidthByMonth"
  instance_name              = "测试用例-dnat2"
  internet_max_bandwidth_out = "0"
  unique_suffix              = false
  auto_use_coupon            = "true"
  public_ip_identification   = false
  instance_charge_strategy   = "PriceHighPriority"
  ens_region_id              = var.ens_region_id
  period_unit                = "Month"
}


`, name)
}

// Test Ens NatGatewayForwardEntry. <<< Resource test cases, automatically generated.
