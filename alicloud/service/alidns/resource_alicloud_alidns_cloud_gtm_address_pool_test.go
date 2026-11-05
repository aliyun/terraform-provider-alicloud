// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alidns_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/alidns"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Case resourceCase_20260318_o7cRMj 12688
func TestAccAliCloudAlidnsCloudGtmAddressPool_basic12688(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alidns_cloud_gtm_address_pool.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAlidnsCloudGtmAddressPoolMap12688)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alidns.NewAlidnsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAlidnsCloudGtmAddressPool")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalidns%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAlidnsCloudGtmAddressPoolBasicDependence12688)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"address_pool_name":         "pool-test-1",
					"health_judgement":          "all_ok",
					"address_pool_type":         "IPv4",
					"enable_status":             "enable",
					"address_lb_strategy":       "sequence",
					"sequence_lb_strategy_mode": "preemptive",
					"remark":                    "remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_pool_name":         "pool-test-1",
						"health_judgement":          "all_ok",
						"address_pool_type":         "IPv4",
						"enable_status":             "enable",
						"address_lb_strategy":       "sequence",
						"sequence_lb_strategy_mode": "preemptive",
						"remark":                    "remark",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_pool_name":         "pool-modify",
					"health_judgement":          "any_ok",
					"enable_status":             "disable",
					"sequence_lb_strategy_mode": "nonPreemptive",
					"remark":                    "add-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_pool_name":         "pool-modify",
						"health_judgement":          "any_ok",
						"enable_status":             "disable",
						"sequence_lb_strategy_mode": "nonPreemptive",
						"remark":                    "add-test",
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

var AlicloudAlidnsCloudGtmAddressPoolMap12688 = map[string]string{
	"create_time": acctest.CHECKSET,
}

func AlicloudAlidnsCloudGtmAddressPoolBasicDependence12688(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resourceCase_20260323_nenn2G 12681
func TestAccAliCloudAlidnsCloudGtmAddressPool_basic12681(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alidns_cloud_gtm_address_pool.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAlidnsCloudGtmAddressPoolMap12681)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alidns.NewAlidnsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAlidnsCloudGtmAddressPool")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalidns%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAlidnsCloudGtmAddressPoolBasicDependence12681)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"address_pool_name":   "resoure-test-pool-3",
					"health_judgement":    "p30_ok",
					"address_pool_type":   "domain",
					"enable_status":       "enable",
					"address_lb_strategy": "round_robin",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_pool_name":   "resoure-test-pool-3",
						"health_judgement":    "p30_ok",
						"address_pool_type":   "domain",
						"enable_status":       "enable",
						"address_lb_strategy": "round_robin",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_judgement": "p50_ok",
					"enable_status":    "disable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_judgement": "p50_ok",
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

var AlicloudAlidnsCloudGtmAddressPoolMap12681 = map[string]string{
	"create_time": acctest.CHECKSET,
}

func AlicloudAlidnsCloudGtmAddressPoolBasicDependence12681(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resourceCase_20260323_ZeTYRG 12686
func TestAccAliCloudAlidnsCloudGtmAddressPool_basic12686(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alidns_cloud_gtm_address_pool.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAlidnsCloudGtmAddressPoolMap12686)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alidns.NewAlidnsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAlidnsCloudGtmAddressPool")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalidns%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAlidnsCloudGtmAddressPoolBasicDependence12686)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"address_pool_name":   "resource-test-pool-2",
					"health_judgement":    "any_ok",
					"address_pool_type":   "IPv6",
					"enable_status":       "enable",
					"address_lb_strategy": "round_robin",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_pool_name":   "resource-test-pool-2",
						"health_judgement":    "any_ok",
						"address_pool_type":   "IPv6",
						"enable_status":       "enable",
						"address_lb_strategy": "round_robin",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_judgement":    "p30_ok",
					"remark":              "test",
					"address_lb_strategy": "weight",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_judgement":    "p30_ok",
						"remark":              "test",
						"address_lb_strategy": "weight",
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

var AlicloudAlidnsCloudGtmAddressPoolMap12686 = map[string]string{
	"create_time": acctest.CHECKSET,
}

func AlicloudAlidnsCloudGtmAddressPoolBasicDependence12686(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}
