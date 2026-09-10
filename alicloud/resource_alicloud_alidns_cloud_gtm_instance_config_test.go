// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case resourceCase_20260320_bH90dh 12689
func TestAccAliCloudAlidnsCloudGtmInstanceConfig_basic12689(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alidns_cloud_gtm_instance_config.default"
	ra := resourceAttrInit(resourceId, AlicloudAlidnsCloudGtmInstanceConfigMap12689)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlidnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlidnsCloudGtmInstanceConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalidns%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlidnsCloudGtmInstanceConfigBasicDependence12689)
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
					"address_pool_lb_strategy":  "sequence",
					"schedule_rr_type":          "A",
					"schedule_zone_name":        "${alicloud_alidns_domain.default.domain_name}",
					"enable_status":             "disable",
					"schedule_host_name":        "example",
					"schedule_zone_mode":        "custom",
					"ttl":                       "600",
					"sequence_lb_strategy_mode": "preemptive",
					"remark":                    "remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_pool_lb_strategy":  "sequence",
						"schedule_rr_type":          "A",
						"schedule_zone_name":        CHECKSET,
						"enable_status":             "disable",
						"schedule_host_name":        "example",
						"schedule_zone_mode":        "custom",
						"ttl":                       "600",
						"sequence_lb_strategy_mode": "preemptive",
						"remark":                    "remark",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_zone_name":        "${alicloud_alidns_domain.default2.domain_name}",
					"enable_status":             "enable",
					"schedule_host_name":        "example-2",
					"ttl":                       "60",
					"sequence_lb_strategy_mode": "nonPreemptive",
					"remark":                    "modify-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"schedule_zone_name":        CHECKSET,
						"enable_status":             "enable",
						"schedule_host_name":        "example-2",
						"ttl":                       "60",
						"sequence_lb_strategy_mode": "nonPreemptive",
						"remark":                    "modify-test",
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

var AlicloudAlidnsCloudGtmInstanceConfigMap12689 = map[string]string{
	"config_id":   CHECKSET,
	"instance_id": CHECKSET,
}

func AlicloudAlidnsCloudGtmInstanceConfigBasicDependence12689(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_alidns_domain" "default" {
    domain_name = "%s.abc"
}

resource "alicloud_alidns_domain" "default2" {
    domain_name = "%s-2.abc"
}


`, name, name, name)
}

// Case resourceCase_20260323_w4vWin 12679
func TestAccAliCloudAlidnsCloudGtmInstanceConfig_basic12679(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alidns_cloud_gtm_instance_config.default"
	ra := resourceAttrInit(resourceId, AlicloudAlidnsCloudGtmInstanceConfigMap12679)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlidnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlidnsCloudGtmInstanceConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalidns%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlidnsCloudGtmInstanceConfigBasicDependence12679)
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
					"address_pool_lb_strategy": "round_robin",
					"schedule_rr_type":         "CNAME",
					"schedule_zone_name":       "${alicloud_alidns_domain.default.domain_name}",
					"enable_status":            "disable",
					"schedule_host_name":       "www",
					"schedule_zone_mode":       "custom",
					"ttl":                      "600",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_pool_lb_strategy": "round_robin",
						"schedule_rr_type":         "CNAME",
						"schedule_zone_name":       CHECKSET,
						"enable_status":            "disable",
						"schedule_host_name":       "www",
						"schedule_zone_mode":       "custom",
						"ttl":                      "600",
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
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudAlidnsCloudGtmInstanceConfigMap12679 = map[string]string{
	"config_id":   CHECKSET,
	"instance_id": CHECKSET,
}

func AlicloudAlidnsCloudGtmInstanceConfigBasicDependence12679(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_alidns_domain" "default" {
    domain_name = "%s.abc"
}


`, name, name)
}

// Case resourceCase_20260323_zV1Ijx 12682
func TestAccAliCloudAlidnsCloudGtmInstanceConfig_basic12682(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alidns_cloud_gtm_instance_config.default"
	ra := resourceAttrInit(resourceId, AlicloudAlidnsCloudGtmInstanceConfigMap12682)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlidnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlidnsCloudGtmInstanceConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalidns%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlidnsCloudGtmInstanceConfigBasicDependence12682)
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
					"address_pool_lb_strategy": "round_robin",
					"schedule_rr_type":         "AAAA",
					"schedule_zone_name":       "${alicloud_alidns_domain.default.domain_name}",
					"enable_status":            "disable",
					"schedule_host_name":       "www",
					"schedule_zone_mode":       "custom",
					"ttl":                      "600",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_pool_lb_strategy": "round_robin",
						"schedule_rr_type":         "AAAA",
						"schedule_zone_name":       CHECKSET,
						"enable_status":            "disable",
						"schedule_host_name":       "www",
						"schedule_zone_mode":       "custom",
						"ttl":                      "600",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_pool_lb_strategy": "weight",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_pool_lb_strategy": "weight",
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

var AlicloudAlidnsCloudGtmInstanceConfigMap12682 = map[string]string{
	"config_id":   CHECKSET,
	"instance_id": CHECKSET,
}

func AlicloudAlidnsCloudGtmInstanceConfigBasicDependence12682(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_alidns_domain" "default" {
    domain_name = "%s.abc"
}


`, name, name)
}
