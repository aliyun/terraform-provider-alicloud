// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudFirewall PrivateDns. >>> Resource test cases, automatically generated.
// Case 私有dns资源测试_多参数 11848
func TestAccAliCloudCloudFirewallPrivateDns_basic11848(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_private_dns.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallPrivateDnsMap11848)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallPrivateDns")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallPrivateDnsBasicDependence11848)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"region_no":            "cn-hangzhou",
					"access_instance_name": name,
					"port":                 "53",
					"primary_vswitch_id":   "${alicloud_vswitch.vpcvsw1.id}",
					"standby_dns":          "4.4.4.4",
					"primary_dns":          "8.8.8.8",
					"vpc_id":               "${alicloud_vpc.vpc.id}",
					"private_dns_type":     "Custom",
					"firewall_type": []string{
						"internet"},
					"ip_protocol":        "UDP",
					"standby_vswitch_id": "${alicloud_vswitch.vpcvsw2.id}",
					"domain_name_list": []string{
						"www.baidu.com"},
					"primary_vswitch_ip": "172.16.3.1",
					"standby_vswitch_ip": "172.16.4.1",
					"member_uid":         "${data.alicloud_account.current.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"region_no":            "cn-hangzhou",
						"access_instance_name": name,
						"port":                 "53",
						"primary_vswitch_id":   CHECKSET,
						"standby_dns":          "4.4.4.4",
						"primary_dns":          "8.8.8.8",
						"vpc_id":               CHECKSET,
						"private_dns_type":     "Custom",
						"firewall_type.#":      "1",
						"ip_protocol":          "UDP",
						"standby_vswitch_id":   CHECKSET,
						"domain_name_list.#":   "1",
						"primary_vswitch_ip":   "172.16.3.1",
						"standby_vswitch_ip":   "172.16.4.1",
						"member_uid":           CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_instance_name": "深圳测试1",
					"standby_dns":          "2.2.2.2",
					"primary_dns":          "1.1.1.1",
					"domain_name_list": []string{
						"www.baidu.com", "www.163.com", "www.taobao.com", "www.sina.com"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_instance_name": "深圳测试1",
						"standby_dns":          "2.2.2.2",
						"primary_dns":          "1.1.1.1",
						"domain_name_list.#":   "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_instance_name": "深圳测试2",
					"private_dns_type":     "PrivateZone",
					"domain_name_list": []string{
						"www.sina.com"},
					"standby_dns": REMOVEKEY,
					"primary_dns": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_instance_name": "深圳测试2",
						"private_dns_type":     "PrivateZone",
						"domain_name_list.#":   "1",
						"standby_dns":          CHECKSET,
						"primary_dns":          CHECKSET,
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

var AlicloudCloudFirewallPrivateDnsMap11848 = map[string]string{
	"status":             CHECKSET,
	"access_instance_id": CHECKSET,
}

func AlicloudCloudFirewallPrivateDnsBasicDependence11848(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "current" {
}

resource "alicloud_vpc" "vpc" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "yqc-test-vpc"
}

resource "alicloud_vswitch" "vpcvsw1" {
  vpc_id     = alicloud_vpc.vpc.id
  zone_id    = "cn-hangzhou-i"
  cidr_block = "172.16.3.0/24"
}

resource "alicloud_vswitch" "vpcvsw2" {
  vpc_id     = alicloud_vpc.vpc.id
  zone_id    = "cn-hangzhou-j"
  cidr_block = "172.16.4.0/24"
}


`, name)
}

// Test CloudFirewall PrivateDns. <<< Resource test cases, automatically generated.
