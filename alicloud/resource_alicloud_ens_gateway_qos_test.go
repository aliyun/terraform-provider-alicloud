// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ens GatewayQos. >>> Resource test cases, automatically generated.
// Case 网关限速_20250620_添加、移除实例 10964
func TestAccAliCloudEnsGatewayQos_basic10964(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_gateway_qos.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsGatewayQosMap10964)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsGatewayQos")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccens%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsGatewayQosBasicDependence10964)
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
					"gateway_qos_name": name,
					"bandwidth_in":     "10",
					"gateway_qos_type": "Nat",
					"network_id":       "${alicloud_ens_network.defaultC7YqlT.id}",
					"bandwidth_out":    "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_qos_name": name,
						"bandwidth_in":     "10",
						"gateway_qos_type": "Nat",
						"network_id":       CHECKSET,
						"bandwidth_out":    "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_qos_name": name + "_update",
					"bandwidth_in":     "20",
					"instances": []map[string]interface{}{
						{
							"status":        "Available",
							"instance_id":   "${alicloud_ens_instance.defaultvuczsY.id}",
							"instance_type": "ens",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_qos_name": name + "_update",
						"bandwidth_in":     "20",
						"instances.#":      "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instances": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instances.#": "0",
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

var AlicloudEnsGatewayQosMap10964 = map[string]string{
	"status":        CHECKSET,
	"creation_time": CHECKSET,
	"ens_region_id": CHECKSET,
}

func AlicloudEnsGatewayQosBasicDependence10964(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "ens_region_id" {
  default = "cn-hangzhou-63"
}

resource "alicloud_ens_network" "defaultC7YqlT" {
  network_name  = "镇元-网关限速测试使用"
  cidr_block    = "10.0.0.0/10"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "default5giQWR" {
  cidr_block    = "10.0.8.0/24"
  vswitch_name  = "镇元-网关限速测试"
  ens_region_id = var.ens_region_id
  network_id    = alicloud_ens_network.defaultC7YqlT.id
}

resource "alicloud_ens_instance" "defaultvuczsY" {
  amount      = "1"
  period_unit = "Month"
  auto_renew  = false
  system_disk {
    size = "20"
  }
  scheduling_strategy        = "Concentrate"
  schedule_area_level        = "Region"
  image_id                   = "centos_6_08_64_20G_alibase_20171208"
  payment_type               = "Subscription"
  instance_type              = "ens.sn1.small"
  password_inherit           = false
  password                   = "12345678abcABC"
  status                     = "Running"
  vswitch_id                 = alicloud_ens_vswitch.default5giQWR.id
  internet_charge_type       = "95BandwidthByMonth"
  instance_name              = "镇元-网关限速测试"
  internet_max_bandwidth_out = "0"
  unique_suffix              = false
  auto_use_coupon            = "true"
  public_ip_identification   = false
  instance_charge_strategy   = "PriceHighPriority"
  ens_region_id              = var.ens_region_id
}


`, name)
}

// Case 网关限速_无实例z_修改名称、带宽_20250619 10966
func TestAccAliCloudEnsGatewayQos_basic10966(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_gateway_qos.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsGatewayQosMap10966)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsGatewayQos")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccens%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsGatewayQosBasicDependence10966)
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
					"gateway_qos_name": name,
					"bandwidth_in":     "10",
					"gateway_qos_type": "Nat",
					"network_id":       "${alicloud_ens_network.defaultC7YqlT.id}",
					"bandwidth_out":    "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_qos_name": name,
						"bandwidth_in":     "10",
						"gateway_qos_type": "Nat",
						"network_id":       CHECKSET,
						"bandwidth_out":    "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_qos_name": name + "_update",
					"bandwidth_in":     "20",
					"bandwidth_out":    "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_qos_name": name + "_update",
						"bandwidth_in":     "20",
						"bandwidth_out":    "20",
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

var AlicloudEnsGatewayQosMap10966 = map[string]string{
	"status":        CHECKSET,
	"creation_time": CHECKSET,
	"ens_region_id": CHECKSET,
}

func AlicloudEnsGatewayQosBasicDependence10966(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "ens_region_id" {
  default = "cn-hangzhou-63"
}

resource "alicloud_ens_network" "defaultC7YqlT" {
  network_name  = "镇元-网关限速测试使用"
  cidr_block    = "10.0.0.0/10"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "default5giQWR" {
  cidr_block    = "10.0.8.0/24"
  vswitch_name  = "镇元-网关限速测试"
  ens_region_id = var.ens_region_id
  network_id    = alicloud_ens_network.defaultC7YqlT.id
}


`, name)
}

// Case 网关限速_创建带实例_20250619 10960
func TestAccAliCloudEnsGatewayQos_basic10960(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_gateway_qos.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsGatewayQosMap10960)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsGatewayQos")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccens%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsGatewayQosBasicDependence10960)
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
					"instances": []map[string]interface{}{
						{
							"status":        "Available",
							"instance_id":   "${alicloud_ens_instance.defaultvuczsY.id}",
							"instance_type": "ens",
						},
					},
					"gateway_qos_name": name,
					"bandwidth_in":     "10",
					"gateway_qos_type": "Nat",
					"network_id":       "${alicloud_ens_network.defaultC7YqlT.id}",
					"bandwidth_out":    "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instances.#":      "1",
						"gateway_qos_name": name,
						"bandwidth_in":     "10",
						"gateway_qos_type": "Nat",
						"network_id":       CHECKSET,
						"bandwidth_out":    "10",
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

var AlicloudEnsGatewayQosMap10960 = map[string]string{
	"status":        CHECKSET,
	"creation_time": CHECKSET,
	"ens_region_id": CHECKSET,
}

func AlicloudEnsGatewayQosBasicDependence10960(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "ens_region_id" {
  default = "cn-hangzhou-63"
}

resource "alicloud_ens_network" "defaultC7YqlT" {
  network_name  = "镇元-网关限速测试使用"
  cidr_block    = "10.0.0.0/10"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "default5giQWR" {
  cidr_block    = "10.0.8.0/24"
  vswitch_name  = "镇元-网关限速测试"
  ens_region_id = var.ens_region_id
  network_id    = alicloud_ens_network.defaultC7YqlT.id
}

resource "alicloud_ens_instance" "defaultvuczsY" {
  amount      = "1"
  period_unit = "Month"
  auto_renew  = false
  system_disk {
    size = "20"
  }
  scheduling_strategy        = "Concentrate"
  schedule_area_level        = "Region"
  image_id                   = "centos_6_08_64_20G_alibase_20171208"
  payment_type               = "Subscription"
  instance_type              = "ens.sn1.small"
  password_inherit           = false
  password                   = "12345678abcABC"
  status                     = "Running"
  vswitch_id                 = alicloud_ens_vswitch.default5giQWR.id
  internet_charge_type       = "95BandwidthByMonth"
  instance_name              = "镇元-网关限速测试"
  internet_max_bandwidth_out = "0"
  unique_suffix              = false
  auto_use_coupon            = "true"
  public_ip_identification   = false
  instance_charge_strategy   = "PriceHighPriority"
  ens_region_id              = var.ens_region_id
}


`, name)
}

// Test Ens GatewayQos. <<< Resource test cases, automatically generated.
