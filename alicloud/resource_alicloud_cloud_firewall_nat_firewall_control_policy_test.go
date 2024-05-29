package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var AlicloudCloudFirewallNatFirewallControlPolicyMap6280 = map[string]string{
	"create_time": CHECKSET,
	"acl_uuid":    CHECKSET,
}

func AlicloudCloudFirewallNatFirewallControlPolicyBasicDependence6280(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-hangzhou-i"
}

variable "dst_port_group_name" {
  default = "ALL"
}

variable "destination_group_name" {
  default = "tf-destination-group"
}

variable "region_id" {
  default = "cn-hangzhou"
}

variable "nat_gateway_id" {
  default = "ngw-wz90ep6rkkn4v24xyex3g"
}

variable "direction" {
  default = "out"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultLilUbB" {
  description = "TF-test-vpc"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "defaultvlFKz1" {
  vpc_id       = alicloud_vpc.defaultLilUbB.id
  cidr_block   = "172.16.1.0/24"
  vswitch_name = var.name
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_cloud_firewall_address_book" "port" {
  description      = "创建ALL的port类型地址簿"
  group_name       = format("%%s%%s", var.name, "port")
  group_type       = "port"
  address_list     = ["22/22", "23/23", "24/24"]
}

resource "alicloud_cloud_firewall_address_book" "ip" {
  description      = "创建ALL的ip类型地址簿"
  group_name       = format("%%s%%s", var.name, "ip")
  group_type       = "ip"
  address_list     = ["1.1.1.1/32", "2.2.2.2/32"]
}

resource "alicloud_cloud_firewall_address_book" "ip2" {
  description      = "创建ALL的ip类型地址簿2"
  group_name       = format("%%s%%s", var.name, "ip2")
  group_type       = "ip"
  address_list     = ["1.1.1.1/32", "2.2.2.2/32"]
}

resource "alicloud_nat_gateway" "defaultjJYdg2" {
  vpc_id           = alicloud_vpc.defaultLilUbB.id
  nat_gateway_name = var.name
  payment_type     = "PayAsYouGo"
  vswitch_id       = alicloud_vswitch.defaultvlFKz1.id
  nat_type         = "Enhanced"
}

`, name)
}

var AlicloudCloudFirewallNatFirewallControlPolicyMap5272 = map[string]string{
	"create_time": CHECKSET,
	"acl_uuid":    CHECKSET,
}

func AlicloudCloudFirewallNatFirewallControlPolicyBasicDependence5272(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-hangzhou-i"
}

variable "dst_port_group_name" {
  default = "ALL"
}

variable "destination_group_name" {
  default = "tf-destination-group"
}

variable "region_id" {
  default = "cn-hangzhou"
}

variable "nat_gateway_id" {
  default = "ngw-wz90ep6rkkn4v24xyex3g"
}

variable "direction" {
  default = "out"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultLilUbB" {
  description = "TF-test-vpc"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "defaultvlFKz1" {
  vpc_id       = alicloud_vpc.defaultLilUbB.id
  cidr_block   = "172.16.1.0/24"
  vswitch_name = var.name
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_cloud_firewall_address_book" "port" {
  description      = "创建ALL的port类型地址簿"
  group_name       = format("%%s%%s", var.name, "port")
  group_type       = "port"
  address_list     = ["22/22", "23/23", "24/24"]
}

resource "alicloud_cloud_firewall_address_book" "ip" {
  description      = "创建ALL的ip类型地址簿"
  group_name       = format("%%s%%s", var.name, "ip")
  group_type       = "ip"
  address_list     = ["1.1.1.1/32", "2.2.2.2/32"]
}

resource "alicloud_cloud_firewall_address_book" "ip2" {
  description      = "创建ALL的ip类型地址簿2"
  group_name       = format("%%s%%s", var.name, "ip2")
  group_type       = "ip"
  address_list     = ["1.1.1.1/32", "2.2.2.2/32"]
}

resource "alicloud_nat_gateway" "defaultjJYdg2" {
  vpc_id           = alicloud_vpc.defaultLilUbB.id
  nat_gateway_name = var.name
  payment_type     = "PayAsYouGo"
  vswitch_id       = alicloud_vswitch.defaultvlFKz1.id
  nat_type         = "Enhanced"
}


`, name)
}

var AlicloudCloudFirewallNatFirewallControlPolicyMap5307 = map[string]string{
	"create_time": CHECKSET,
	"acl_uuid":    CHECKSET,
}

func AlicloudCloudFirewallNatFirewallControlPolicyBasicDependence5307(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "direction" {
  default = "out"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultDEiWfM" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "defaultFHDM3F" {
  vpc_id     = alicloud_vpc.defaultDEiWfM.id
  zone_id    = data.alicloud_zones.default.zones.0.id
  cidr_block = "172.16.2.0/24"
}

resource "alicloud_nat_gateway" "defaultMbS2Ts" {
  depends_on       = ["alicloud_cloud_firewall_address_book.port", "alicloud_cloud_firewall_address_book.ip"]
  vpc_id           = alicloud_vpc.defaultDEiWfM.id
  nat_gateway_name = var.name
  payment_type     = "PayAsYouGo"
  vswitch_id       = alicloud_vswitch.defaultFHDM3F.id
  nat_type         = "Enhanced"
}

resource "alicloud_cloud_firewall_address_book" "port" {
  description      = "创建ALL的port类型地址簿"
  group_name       = format("%%s%%s", var.name, "port")
  group_type       = "port"
  address_list     = ["22/22", "23/23", "24/24"]
}

resource "alicloud_cloud_firewall_address_book" "port-update" {
  description      = "创建ALL的port类型地址簿-update"
  group_name       = format("%%s%%s", var.name, "port-update")
  group_type       = "port"
  address_list     = ["22/22", "23/23", "24/24"]
}

resource "alicloud_cloud_firewall_address_book" "domain" {
  description      = "创建ALL的port类型地址簿"
  group_name       = format("%%s%%s", var.name, "domain")
  group_type       = "domain"
  address_list     = ["alibaba.com", "aliyun.com", "alicloud.com"]
}

resource "alicloud_cloud_firewall_address_book" "ip" {
  description      = "tf-destination-group"
  group_name       = var.name
  group_type       = "ip"
  address_list     = ["1.1.1.1/32", "2.2.2.2/32"]
}

resource "alicloud_cloud_firewall_nat_firewall_control_policy" "default0" {
  application_name_list = [
    "ANY"
  ]
  description = var.name
  release     = "false"
  ip_version  = "4"
  repeat_days = [
    "1",
    "2",
    "3"
  ]
  repeat_start_time   = "21:00"
  acl_action          = "log"
  dest_port_group     = alicloud_cloud_firewall_address_book.port.group_name
  repeat_type         = "Weekly"
  nat_gateway_id      = alicloud_nat_gateway.defaultMbS2Ts.id
  source              = "1.1.1.1/32"
  direction           = "out"
  repeat_end_time     = "21:30"
  start_time          = "1699156800"
  destination         = "1.1.1.1/32"
  end_time            = "1888545600"
  source_type         = "net"
  proto               = "TCP"
  new_order           = "1"
  destination_type    = "net"
  dest_port_type      = "group"
  domain_resolve_type = "0"

  lifecycle {
    ignore_changes = [new_order]
  }
}

`, name)
}

// Case 全生命周期_1.1 6280  twin
func TestAccAliCloudCloudFirewallNatFirewallControlPolicy_basic6280_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_nat_firewall_control_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallNatFirewallControlPolicyMap6280)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallNatFirewallControlPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallnatfirewallcontrolpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallNatFirewallControlPolicyBasicDependence6280)
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
					"destination":         "1.1.1.1/32",
					"description":         "测试nat防火墙规则",
					"dest_port":           "22/22",
					"acl_action":          "log",
					"destination_type":    "net",
					"direction":           "out",
					"source":              "2.2.2.2/32",
					"dest_port_type":      "port",
					"proto":               "TCP",
					"domain_resolve_type": "0",
					"repeat_type":         "Weekly",
					"nat_gateway_id":      "${alicloud_nat_gateway.defaultjJYdg2.id}",
					"new_order":           "1",
					"release":             "false",
					"source_type":         "net",
					"ip_version":          "4",
					"application_name_list": []string{
						"ANY"},
					"end_time":        "1888545600",
					"start_time":      "1699156800",
					"repeat_end_time": "21:30",
					"repeat_days": []string{
						"1", "2", "3", "4"},
					"repeat_start_time": "21:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":             "1.1.1.1/32",
						"description":             "测试nat防火墙规则",
						"dest_port":               "22/22",
						"acl_action":              "log",
						"destination_type":        "net",
						"direction":               "out",
						"source":                  "2.2.2.2/32",
						"dest_port_type":          "port",
						"proto":                   "TCP",
						"domain_resolve_type":     "0",
						"repeat_type":             "Weekly",
						"nat_gateway_id":          CHECKSET,
						"new_order":               "1",
						"release":                 "false",
						"source_type":             "net",
						"ip_version":              "4",
						"application_name_list.#": "1",
						"end_time":                "1888545600",
						"start_time":              "1699156800",
						"repeat_end_time":         "21:30",
						"repeat_days.#":           "4",
						"repeat_start_time":       "21:00",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ip_version"},
			},
		},
	})
}

// Case 全生命周期 5272  twin
func TestAccAliCloudCloudFirewallNatFirewallControlPolicy_basic5272_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_nat_firewall_control_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallNatFirewallControlPolicyMap5272)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallNatFirewallControlPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallnatfirewallcontrolpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallNatFirewallControlPolicyBasicDependence5272)
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
					"destination":         "1.1.1.1/32",
					"description":         "测试nat防火墙规则",
					"dest_port":           "22/22",
					"acl_action":          "log",
					"destination_type":    "net",
					"direction":           "out",
					"source":              "2.2.2.2/32",
					"dest_port_type":      "port",
					"proto":               "TCP",
					"domain_resolve_type": "0",
					"repeat_type":         "Weekly",
					"nat_gateway_id":      "${alicloud_nat_gateway.defaultjJYdg2.id}",
					"new_order":           "1",
					"release":             "false",
					"source_type":         "net",
					"ip_version":          "4",
					"application_name_list": []string{
						"ANY"},
					"end_time":        "1888545600",
					"start_time":      "1699156800",
					"repeat_end_time": "21:30",
					"repeat_days": []string{
						"1", "2", "3", "4"},
					"repeat_start_time": "21:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":             "1.1.1.1/32",
						"description":             "测试nat防火墙规则",
						"dest_port":               "22/22",
						"acl_action":              "log",
						"destination_type":        "net",
						"direction":               "out",
						"source":                  "2.2.2.2/32",
						"dest_port_type":          "port",
						"proto":                   "TCP",
						"domain_resolve_type":     "0",
						"repeat_type":             "Weekly",
						"nat_gateway_id":          CHECKSET,
						"new_order":               "1",
						"release":                 "false",
						"source_type":             "net",
						"ip_version":              "4",
						"application_name_list.#": "1",
						"end_time":                "1888545600",
						"start_time":              "1699156800",
						"repeat_end_time":         "21:30",
						"repeat_days.#":           "4",
						"repeat_start_time":       "21:00",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ip_version"},
			},
		},
	})
}

// Case 继续增加覆盖 5307  twin
func TestAccAliCloudCloudFirewallNatFirewallControlPolicy_basic5307_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_nat_firewall_control_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallNatFirewallControlPolicyMap5307)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallNatFirewallControlPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallnatfirewallcontrolpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallNatFirewallControlPolicyBasicDependence5307)
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
					"destination":         "1.1.1.1/32",
					"description":         "测试nat防火墙规则",
					"end_time":            "1888545600",
					"ip_version":          "4",
					"source_type":         "net",
					"start_time":          "1699156800",
					"acl_action":          "log",
					"destination_type":    "net",
					"direction":           "out",
					"source":              "1.1.1.1/32",
					"dest_port_type":      "group",
					"dest_port_group":     "${alicloud_cloud_firewall_address_book.port.group_name}",
					"proto":               "TCP",
					"repeat_end_time":     "21:30",
					"domain_resolve_type": "0",
					"repeat_type":         "Weekly",
					"repeat_days": []string{
						"1", "2", "3"},
					"repeat_start_time": "21:00",
					"nat_gateway_id":    "${alicloud_cloud_firewall_nat_firewall_control_policy.default0.nat_gateway_id}",
					"new_order":         "1",
					"release":           "false",
					"application_name_list": []string{
						"ANY"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":             "1.1.1.1/32",
						"description":             "测试nat防火墙规则",
						"end_time":                "1888545600",
						"ip_version":              "4",
						"source_type":             "net",
						"start_time":              "1699156800",
						"acl_action":              "log",
						"destination_type":        "net",
						"direction":               "out",
						"source":                  "1.1.1.1/32",
						"dest_port_type":          "group",
						"proto":                   "TCP",
						"repeat_end_time":         "21:30",
						"domain_resolve_type":     "0",
						"repeat_type":             "Weekly",
						"repeat_days.#":           "3",
						"repeat_start_time":       "21:00",
						"nat_gateway_id":          CHECKSET,
						"new_order":               "1",
						"release":                 "false",
						"dest_port_group":         CHECKSET,
						"application_name_list.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ip_version"},
			},
		},
	})
}

// Case 全生命周期_1.1 6280  raw
func TestAccAliCloudCloudFirewallNatFirewallControlPolicy_basic6280_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_nat_firewall_control_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallNatFirewallControlPolicyMap6280)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallNatFirewallControlPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallnatfirewallcontrolpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallNatFirewallControlPolicyBasicDependence6280)
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
					"destination":         "1.1.1.1/32",
					"description":         "测试nat防火墙规则",
					"dest_port":           "22/22",
					"acl_action":          "log",
					"destination_type":    "net",
					"direction":           "out",
					"source":              "2.2.2.2/32",
					"dest_port_type":      "port",
					"proto":               "TCP",
					"domain_resolve_type": "0",
					"repeat_type":         "Weekly",
					"nat_gateway_id":      "${alicloud_nat_gateway.defaultjJYdg2.id}",
					"new_order":           "1",
					"release":             "false",
					"source_type":         "net",
					"ip_version":          "4",
					"application_name_list": []string{
						"ANY"},
					"end_time":        "1888545600",
					"start_time":      "1699156800",
					"repeat_end_time": "21:30",
					"repeat_days": []string{
						"1", "2", "3", "4"},
					"repeat_start_time": "21:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":             "1.1.1.1/32",
						"description":             "测试nat防火墙规则",
						"dest_port":               "22/22",
						"acl_action":              "log",
						"destination_type":        "net",
						"direction":               "out",
						"source":                  "2.2.2.2/32",
						"dest_port_type":          "port",
						"proto":                   "TCP",
						"domain_resolve_type":     "0",
						"repeat_type":             "Weekly",
						"nat_gateway_id":          CHECKSET,
						"new_order":               "1",
						"release":                 "false",
						"source_type":             "net",
						"ip_version":              "4",
						"application_name_list.#": "1",
						"end_time":                "1888545600",
						"start_time":              "1699156800",
						"repeat_end_time":         "21:30",
						"repeat_days.#":           "4",
						"repeat_start_time":       "21:00",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":         "1.1.1.1/32",
					"description":         "测试nat防火墙规则",
					"dest_port":           "22/22",
					"acl_action":          "log",
					"destination_type":    "net",
					"direction":           "out",
					"source":              "2.2.2.2/32",
					"dest_port_type":      "port",
					"proto":               "TCP",
					"domain_resolve_type": "0",
					"repeat_type":         "Weekly",
					"nat_gateway_id":      "${alicloud_nat_gateway.defaultjJYdg2.id}",
					"new_order":           "1",
					"release":             "false",
					"source_type":         "net",
					"ip_version":          "4",
					"application_name_list": []string{
						"ANY"},
					"end_time":        "1888545600",
					"start_time":      "1699156800",
					"repeat_end_time": "21:30",
					"repeat_days": []string{
						"1", "2", "3", "4"},
					"repeat_start_time": "21:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":             "1.1.1.1/32",
						"description":             "测试nat防火墙规则",
						"dest_port":               "22/22",
						"acl_action":              "log",
						"destination_type":        "net",
						"direction":               "out",
						"source":                  "2.2.2.2/32",
						"dest_port_type":          "port",
						"proto":                   "TCP",
						"domain_resolve_type":     "0",
						"repeat_type":             "Weekly",
						"nat_gateway_id":          CHECKSET,
						"new_order":               "1",
						"release":                 "false",
						"source_type":             "net",
						"ip_version":              "4",
						"application_name_list.#": "1",
						"end_time":                "1888545600",
						"start_time":              "1699156800",
						"repeat_end_time":         "21:30",
						"repeat_days.#":           "4",
						"repeat_start_time":       "21:00",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":         "www.baidu.com",
					"description":         "测试1",
					"dest_port":           "80/80",
					"acl_action":          "drop",
					"destination_type":    "domain",
					"source":              "1.1.1.1/32",
					"domain_resolve_type": "1",
					"repeat_type":         "Permanent",
					"release":             "true",
					"application_name_list": []string{
						"HTTP", "HTTPS", "SMTP", "SMTPS", "SSL"},
					"repeat_days":       []string{},
					"start_time":        REMOVEKEY,
					"end_time":          REMOVEKEY,
					"repeat_start_time": REMOVEKEY,
					"repeat_end_time":   REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":             "www.baidu.com",
						"description":             "测试1",
						"dest_port":               "80/80",
						"acl_action":              "drop",
						"destination_type":        "domain",
						"source":                  "1.1.1.1/32",
						"domain_resolve_type":     "1",
						"repeat_type":             "Permanent",
						"release":                 "true",
						"application_name_list.#": "5",
						"repeat_days.#":           "0",
						"repeat_start_time":       "",
						"repeat_end_time":         "",
						"start_time":              "0",
						"end_time":                "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":         "2.2.2.2/32",
					"description":         "更新2",
					"dest_port":           "22/22",
					"destination_type":    "net",
					"domain_resolve_type": "0",
					"new_order":           "1",
					"release":             "false",
					"application_name_list": []string{
						"HTTP", "HTTPS", "SSL"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":             "2.2.2.2/32",
						"description":             "更新2",
						"dest_port":               "22/22",
						"destination_type":        "net",
						"domain_resolve_type":     "0",
						"new_order":               "1",
						"release":                 "false",
						"application_name_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":      "${alicloud_cloud_firewall_address_book.ip.group_name}",
					"description":      "测试修改",
					"acl_action":       "accept",
					"destination_type": "group",
					"source":           "${alicloud_cloud_firewall_address_book.ip.group_name}",
					"dest_port_type":   "group",
					"dest_port_group":  "${alicloud_cloud_firewall_address_book.port.group_name}",
					"dest_port":        REMOVEKEY,
					"proto":            "UDP",
					"repeat_type":      "Daily",
					"source_type":      "group",
					"application_name_list": []string{
						"ANY"},
					"end_time":          "1701424800",
					"start_time":        "1701421200",
					"repeat_end_time":   "22:30",
					"repeat_start_time": "22:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":             CHECKSET,
						"description":             "测试修改",
						"acl_action":              "accept",
						"destination_type":        "group",
						"source":                  CHECKSET,
						"dest_port_type":          "group",
						"proto":                   "UDP",
						"repeat_type":             "Daily",
						"source_type":             "group",
						"application_name_list.#": "1",
						"end_time":                "1701424800",
						"start_time":              "1701421200",
						"repeat_end_time":         "22:30",
						"repeat_start_time":       "22:00",
						"dest_port":               "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":      "2.2.2.2/32",
					"description":      "更新3",
					"destination_type": "net",
					"source":           "1.1.1.1/32",
					"dest_port_type":   "port",
					"proto":            "TCP",
					"repeat_type":      "Weekly",
					"source_type":      "net",
					"application_name_list": []string{
						"SSL", "HTTP", "HTTPS"},
					"repeat_days": []string{
						"1", "2", "3", "4", "5", "6"},
					"dest_port": "22/22",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":             "2.2.2.2/32",
						"description":             "更新3",
						"destination_type":        "net",
						"source":                  "1.1.1.1/32",
						"dest_port_type":          "port",
						"proto":                   "TCP",
						"repeat_type":             "Weekly",
						"source_type":             "net",
						"application_name_list.#": "3",
						"repeat_days.#":           "6",
						"dest_port":               "22/22",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "更新4",
					"application_name_list": []string{
						"SSL", "HTTP", "HTTPS"},
					"repeat_days": []string{
						"1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":             "更新4",
						"application_name_list.#": "3",
						"repeat_days.#":           "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ip_version"},
			},
		},
	})
}

// Case 全生命周期 5272  raw
func TestAccAliCloudCloudFirewallNatFirewallControlPolicy_basic5272_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_nat_firewall_control_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallNatFirewallControlPolicyMap5272)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallNatFirewallControlPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallnatfirewallcontrolpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallNatFirewallControlPolicyBasicDependence5272)
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
					"destination":         "1.1.1.1/32",
					"description":         "测试nat防火墙规则",
					"dest_port":           "22/22",
					"acl_action":          "log",
					"destination_type":    "net",
					"direction":           "out",
					"source":              "2.2.2.2/32",
					"dest_port_type":      "port",
					"proto":               "TCP",
					"domain_resolve_type": "0",
					"repeat_type":         "Weekly",
					"nat_gateway_id":      "${alicloud_nat_gateway.defaultjJYdg2.id}",
					"new_order":           "1",
					"release":             "false",
					"source_type":         "net",
					"ip_version":          "4",
					"application_name_list": []string{
						"ANY"},
					"end_time":        "1888545600",
					"start_time":      "1699156800",
					"repeat_end_time": "21:30",
					"repeat_days": []string{
						"1", "2", "3", "4"},
					"repeat_start_time": "21:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":             "1.1.1.1/32",
						"description":             "测试nat防火墙规则",
						"dest_port":               "22/22",
						"acl_action":              "log",
						"destination_type":        "net",
						"direction":               "out",
						"source":                  "2.2.2.2/32",
						"dest_port_type":          "port",
						"proto":                   "TCP",
						"domain_resolve_type":     "0",
						"repeat_type":             "Weekly",
						"nat_gateway_id":          CHECKSET,
						"new_order":               "1",
						"release":                 "false",
						"source_type":             "net",
						"ip_version":              "4",
						"application_name_list.#": "1",
						"end_time":                "1888545600",
						"start_time":              "1699156800",
						"repeat_end_time":         "21:30",
						"repeat_days.#":           "4",
						"repeat_start_time":       "21:00",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":         "www.baidu.com",
					"description":         "测试1",
					"dest_port":           "80/80",
					"acl_action":          "drop",
					"destination_type":    "domain",
					"source":              "1.1.1.1/32",
					"domain_resolve_type": "1",
					"repeat_type":         "Permanent",
					"release":             "true",
					"application_name_list": []string{
						"HTTP", "HTTPS", "SMTP", "SMTPS", "SSL"},
					"repeat_days": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":             "www.baidu.com",
						"description":             "测试1",
						"dest_port":               "80/80",
						"acl_action":              "drop",
						"destination_type":        "domain",
						"source":                  "1.1.1.1/32",
						"domain_resolve_type":     "1",
						"repeat_type":             "Permanent",
						"release":                 "true",
						"application_name_list.#": "5",
						"repeat_days.#":           "0",
						"repeat_start_time":       "",
						"repeat_end_time":         "",
						"start_time":              "0",
						"end_time":                "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":         "2.2.2.2/32",
					"description":         "更新2",
					"dest_port":           "22/22",
					"destination_type":    "net",
					"domain_resolve_type": "0",
					"release":             "false",
					"application_name_list": []string{
						"HTTP", "HTTPS", "SSL"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":             "2.2.2.2/32",
						"description":             "更新2",
						"dest_port":               "22/22",
						"destination_type":        "net",
						"domain_resolve_type":     "0",
						"release":                 "false",
						"application_name_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dest_port": "22/23",
					"application_name_list": []string{
						"HTTP", "HTTPS"},
					"release": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"release":                 "true",
						"dest_port":               "22/23",
						"application_name_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"release": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"release": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":      "${alicloud_cloud_firewall_address_book.ip2.group_name}",
					"description":      "测试修改",
					"acl_action":       "accept",
					"destination_type": "group",
					"source":           "${alicloud_cloud_firewall_address_book.ip.group_name}",
					"dest_port_type":   "group",
					"dest_port_group":  "${alicloud_cloud_firewall_address_book.port.group_name}",
					"proto":            "UDP",
					"repeat_type":      "Daily",
					"source_type":      "group",
					"application_name_list": []string{
						"ANY"},
					"end_time":          "1701424800",
					"start_time":        "1701421200",
					"repeat_end_time":   "22:30",
					"repeat_start_time": "22:00",
					"dest_port":         "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":             CHECKSET,
						"description":             "测试修改",
						"acl_action":              "accept",
						"destination_type":        "group",
						"source":                  CHECKSET,
						"dest_port_type":          "group",
						"proto":                   "UDP",
						"repeat_type":             "Daily",
						"source_type":             "group",
						"application_name_list.#": "1",
						"end_time":                "1701424800",
						"start_time":              "1701421200",
						"repeat_end_time":         "22:30",
						"repeat_start_time":       "22:00",
						"dest_port":               "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":      "2.2.2.2/32",
					"description":      "更新3",
					"destination_type": "net",
					"source":           "1.1.1.1/32",
					"dest_port_type":   "port",
					"proto":            "TCP",
					"repeat_type":      "Weekly",
					"source_type":      "net",
					"application_name_list": []string{
						"SSL", "HTTP", "HTTPS"},
					"repeat_days": []string{
						"1", "2", "3", "4", "5", "6"},
					"dest_port": "0/0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":             "2.2.2.2/32",
						"description":             "更新3",
						"destination_type":        "net",
						"source":                  "1.1.1.1/32",
						"dest_port_type":          "port",
						"proto":                   "TCP",
						"repeat_type":             "Weekly",
						"source_type":             "net",
						"application_name_list.#": "3",
						"repeat_days.#":           "6",
						"dest_port":               "0/0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "更新4",
					"application_name_list": []string{
						"SSL", "HTTP", "HTTPS"},
					"repeat_days": []string{
						"1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":             "更新4",
						"application_name_list.#": "3",
						"repeat_days.#":           "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ip_version"},
			},
		},
	})
}

// Case 继续增加覆盖 5307  raw
func TestAccAliCloudCloudFirewallNatFirewallControlPolicy_basic5307_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_nat_firewall_control_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallNatFirewallControlPolicyMap5307)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallNatFirewallControlPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallnatfirewallcontrolpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallNatFirewallControlPolicyBasicDependence5307)
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
					"destination":         "1.1.1.1/32",
					"description":         "nat防火墙规则",
					"end_time":            "1888545600",
					"ip_version":          "4",
					"source_type":         "net",
					"start_time":          "1699156800",
					"acl_action":          "log",
					"destination_type":    "net",
					"direction":           "out",
					"source":              "1.1.1.1/32",
					"dest_port_type":      "group",
					"proto":               "TCP",
					"repeat_end_time":     "21:30",
					"domain_resolve_type": "0",
					"repeat_type":         "Weekly",
					"repeat_days": []string{
						"1", "2", "3"},
					"repeat_start_time": "21:00",
					"nat_gateway_id":    "${alicloud_cloud_firewall_nat_firewall_control_policy.default0.nat_gateway_id}",
					"new_order":         "2",
					"release":           "false",
					"dest_port_group":   "${alicloud_cloud_firewall_address_book.port.group_name}",
					"application_name_list": []string{
						"ANY"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":             "1.1.1.1/32",
						"description":             "nat防火墙规则",
						"end_time":                "1888545600",
						"ip_version":              "4",
						"source_type":             "net",
						"start_time":              "1699156800",
						"acl_action":              "log",
						"destination_type":        "net",
						"direction":               "out",
						"source":                  "1.1.1.1/32",
						"dest_port_type":          "group",
						"proto":                   "TCP",
						"repeat_end_time":         "21:30",
						"domain_resolve_type":     "0",
						"repeat_type":             "Weekly",
						"repeat_days.#":           "3",
						"repeat_start_time":       "21:00",
						"nat_gateway_id":          CHECKSET,
						"new_order":               "2",
						"release":                 "false",
						"dest_port_group":         CHECKSET,
						"application_name_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":      "alicloud.com",
					"description":      "update1",
					"acl_action":       "accept",
					"destination_type": "domain",
					"repeat_type":      "Permanent",
					"release":          "true",
					//"dest_port_group":  "${alicloud_cloud_firewall_address_book.port-update.group_name}",
					"application_name_list": []string{
						"HTTPS", "HTTP", "SMTP", "SMTPS", "SSL"},
					"repeat_days":       REMOVEKEY,
					"start_time":        REMOVEKEY,
					"end_time":          REMOVEKEY,
					"repeat_start_time": REMOVEKEY,
					"repeat_end_time":   REMOVEKEY,
					"dest_port_type":    "port",
					"new_order":         "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":             "alicloud.com",
						"description":             "update1",
						"acl_action":              "accept",
						"destination_type":        "domain",
						"repeat_type":             "Permanent",
						"release":                 "true",
						"dest_port_group":         "",
						"dest_port_type":          "port",
						"application_name_list.#": "5",
						"repeat_start_time":       "",
						"repeat_end_time":         "",
						"start_time":              "0",
						"end_time":                "0",
						"repeat_days.#":           "0",
						"new_order":               "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ip_version"},
			},
		},
	})
}
