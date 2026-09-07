// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Alidns CloudGtmAddress. >>> Resource test cases, automatically generated.
// Case resourceCase_20260320_c1u6VV 12687
func TestAccAliCloudAlidnsCloudGtmAddress_basic12687(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alidns_cloud_gtm_address.default"
	ra := resourceAttrInit(resourceId, AlicloudAlidnsCloudGtmAddressMap12687)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlidnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlidnsCloudGtmAddress")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalidns%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlidnsCloudGtmAddressBasicDependence12687)
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
					"type":             "IPv4",
					"health_judgement": "all_ok",
					"health_tasks": []map[string]interface{}{
						{
							"template_id": "${alicloud_alidns_cloud_gtm_monitor_template.create_ping_template.id}",
						},
						{
							"port":        "53",
							"template_id": "${alicloud_alidns_cloud_gtm_monitor_template.create_tcp_template.id}",
						},
						{
							"port":        "443",
							"template_id": "${alicloud_alidns_cloud_gtm_monitor_template.create_https_template.id}",
						},
					},
					"address":                 "1.1.1.1",
					"enable_status":           "enable",
					"available_mode":          "manual",
					"manual_available_status": "available",
					"name":                    name,
					"remark":                  "remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":                    "IPv4",
						"health_judgement":        "all_ok",
						"health_tasks.#":          "3",
						"address":                 "1.1.1.1",
						"enable_status":           "enable",
						"available_mode":          "manual",
						"manual_available_status": "available",
						"name":                    name,
						"remark":                  "remark",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_judgement": "any_ok",
					"health_tasks": []map[string]interface{}{
						{
							"template_id": "${alicloud_alidns_cloud_gtm_monitor_template.create_ping_template.id}",
						},
						{
							"port":        "443",
							"template_id": "${alicloud_alidns_cloud_gtm_monitor_template.create_https_template.id}",
						},
					},
					"address":       "2.2.2.2",
					"enable_status": "disable",
					"name":          name + "_update",
					"remark":        "add-test-modify",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_judgement": "any_ok",
						"health_tasks.#":   "2",
						"address":          "2.2.2.2",
						"enable_status":    "disable",
						"name":             name + "_update",
						"remark":           "add-test-modify",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_tasks":   REMOVEKEY,
					"available_mode": "auto",
					"name":           name + "_update",
					"remark":         "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_tasks.#": "0",
						"available_mode": "auto",
						"name":           name + "_update",
						"remark":         "test",
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

var AlicloudAlidnsCloudGtmAddressMap12687 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudAlidnsCloudGtmAddressBasicDependence12687(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_alidns_cloud_gtm_monitor_template" "create_tcp_template" {
  ip_version = "IPv4"
  timeout    = "3000"
  isp_city_nodes {
    city_code = "357"
    isp_code  = "465"
  }
  isp_city_nodes {
    city_code = "738"
    isp_code  = "465"
  }
  evaluation_count = "1"
  protocol         = "tcp"
  failure_rate     = "50"
  extend_info      = "{}"
  name             = "test-case-2"
  interval         = "60"
}

resource "alicloud_alidns_cloud_gtm_monitor_template" "create_https_template" {
  ip_version = "IPv4"
  timeout    = "2000"
  isp_city_nodes {
    city_code = "357"
    isp_code  = "465"
  }
  isp_city_nodes {
    city_code = "738"
    isp_code  = "465"
  }
  evaluation_count = "1"
  protocol         = "https"
  failure_rate     = "50"
  extend_info      = "{\"code\":400,\"followRedirect\":true,\"path\":\"/\",\"sni\":false}"
  name             = "test-case-3"
  interval         = "60"
}

resource "alicloud_alidns_cloud_gtm_monitor_template" "create_ping_template" {
  ip_version = "IPv4"
  timeout    = "3000"
  isp_city_nodes {
    city_code = "357"
    isp_code  = "465"
  }
  isp_city_nodes {
    city_code = "738"
    isp_code  = "465"
  }
  evaluation_count = "1"
  protocol         = "ping"
  failure_rate     = "50"
  extend_info      = "{\"packetLossRate\":10,\"packetNum\":20}"
  name             = "test-case-1"
  interval         = "60"
}


`, name)
}

// Case resourceCase_20260323_h98fTA 12680
func TestAccAliCloudAlidnsCloudGtmAddress_basic12680(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alidns_cloud_gtm_address.default"
	ra := resourceAttrInit(resourceId, AlicloudAlidnsCloudGtmAddressMap12680)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlidnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlidnsCloudGtmAddress")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalidns%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlidnsCloudGtmAddressBasicDependence12680)
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
					"type":             "domain",
					"health_judgement": "all_ok",
					"address":          "www.tianxuan.top",
					"enable_status":    "enable",
					"available_mode":   "auto",
					"name":             name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":             "domain",
						"health_judgement": "all_ok",
						"address":          "www.tianxuan.top",
						"enable_status":    "enable",
						"available_mode":   "auto",
						"name":             name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"available_mode":          "manual",
					"manual_available_status": "available",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"available_mode":          "manual",
						"manual_available_status": "available",
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

var AlicloudAlidnsCloudGtmAddressMap12680 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudAlidnsCloudGtmAddressBasicDependence12680(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resourceCase_20260323_8FMXi4 12683
func TestAccAliCloudAlidnsCloudGtmAddress_basic12683(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alidns_cloud_gtm_address.default"
	ra := resourceAttrInit(resourceId, AlicloudAlidnsCloudGtmAddressMap12683)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlidnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlidnsCloudGtmAddress")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalidns%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlidnsCloudGtmAddressBasicDependence12683)
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
					"type":             "IPv6",
					"health_judgement": "any_ok",
					"address":          "2400:3200:baba:0:0:0:0:1",
					"enable_status":    "enable",
					"available_mode":   "auto",
					"name":             name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":             "IPv6",
						"health_judgement": "any_ok",
						"address":          "2400:3200:baba:0:0:0:0:1",
						"enable_status":    "enable",
						"available_mode":   "auto",
						"name":             name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_judgement": "all_ok",
					"enable_status":    "disable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_judgement": "all_ok",
						"enable_status":    "disable",
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

var AlicloudAlidnsCloudGtmAddressMap12683 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudAlidnsCloudGtmAddressBasicDependence12683(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Alidns CloudGtmAddress. <<< Resource test cases, automatically generated.
