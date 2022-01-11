package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudAlidnsAccessStrategy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alidns_access_strategy.default"
	checkoutSupportedRegions(t, true, connectivity.AlidnsSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudAlidnsAccessStrategyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlidnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlidnsAccessStrategy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlidnsAccessStrategyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
			testAccPreCheckWithEnvVariable(t, "ALICLOUD_ICP_DOMAIN_NAME")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"default_addr_pool_type":         "IPV4",
					"default_lba_strategy":           "ALL_RR",
					"default_min_available_addr_num": "1",
					"default_addr_pools": []map[string]interface{}{
						{
							"addr_pool_id": "${alicloud_alidns_address_pool.ipv4.0.id}",
						},
					},
					"strategy_mode": "GEO",
					"lines": []map[string]interface{}{
						{
							"line_code": "default",
						},
					},
					"instance_id":   "${local.gtm_instance_id}",
					"strategy_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_addr_pool_type":         "IPV4",
						"strategy_mode":                  "GEO",
						"default_lba_strategy":           "ALL_RR",
						"instance_id":                    CHECKSET,
						"default_addr_pools.#":           "1",
						"lines.#":                        "1",
						"default_min_available_addr_num": "1",
						"strategy_name":                  name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"strategy_name": "${var.name}_1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"strategy_name": name + "_1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"failover_addr_pool_type":         "IPV4",
					"failover_lba_strategy":           "ALL_RR",
					"failover_min_available_addr_num": "1",
					"failover_addr_pools": []map[string]interface{}{
						{
							"addr_pool_id": "${alicloud_alidns_address_pool.ipv4.1.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"failover_min_available_addr_num": "1",
						"failover_lba_strategy":           "ALL_RR",
						"failover_addr_pool_type":         "IPV4",
						"failover_addr_pools.#":           "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"default_min_available_addr_num": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_min_available_addr_num": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"failover_min_available_addr_num": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"failover_min_available_addr_num": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_addr_pool_type": "IPV6",
					"default_addr_pools": []map[string]interface{}{
						{
							"addr_pool_id": "${alicloud_alidns_address_pool.ipv6.0.id}",
						},
					},
					"failover_addr_pool_type": "IPV6",
					"failover_addr_pools": []map[string]interface{}{
						{
							"addr_pool_id": "${alicloud_alidns_address_pool.ipv6.1.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_addr_pool_type":  "IPV6",
						"default_addr_pools.#":    "1",
						"failover_addr_pool_type": "IPV6",
						"failover_addr_pools.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lines": []map[string]interface{}{
						{
							"line_code": "cn_mobile_beijing",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lines.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_mode": "DEFAULT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_mode": "DEFAULT",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_lba_strategy": "RATIO",
					"default_addr_pools": []map[string]interface{}{
						{
							"lba_weight":   "1",
							"addr_pool_id": "${alicloud_alidns_address_pool.ipv6.0.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_lba_strategy": "RATIO",
						"default_addr_pools.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"failover_lba_strategy": "RATIO",
					"failover_addr_pools": []map[string]interface{}{
						{
							"lba_weight":   "1",
							"addr_pool_id": "${alicloud_alidns_address_pool.ipv6.1.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"failover_lba_strategy": "RATIO",
						"failover_addr_pools.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}
func TestAccAlicloudAlidnsAccessStrategy_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alidns_access_strategy.default"
	checkoutSupportedRegions(t, true, connectivity.AlidnsSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudAlidnsAccessStrategyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlidnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlidnsAccessStrategy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlidnsAccessStrategyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
			testAccPreCheckWithEnvVariable(t, "ALICLOUD_ICP_DOMAIN_NAME")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"access_mode":                    "AUTO",
					"default_addr_pool_type":         "IPV4",
					"default_lba_strategy":           "RATIO",
					"default_min_available_addr_num": "1",
					"default_addr_pools": []map[string]interface{}{
						{
							"lba_weight":   "1",
							"addr_pool_id": "${alicloud_alidns_address_pool.ipv4.0.id}",
						},
					},
					"failover_addr_pool_type":         "IPV4",
					"failover_lba_strategy":           "RATIO",
					"failover_min_available_addr_num": "1",
					"failover_addr_pools": []map[string]interface{}{
						{
							"lba_weight":   "1",
							"addr_pool_id": "${alicloud_alidns_address_pool.ipv4.1.id}",
						},
					},
					"strategy_mode": "GEO",
					"lines": []map[string]interface{}{
						{
							"line_code": "default",
						},
					},
					"instance_id":   "${local.gtm_instance_id}",
					"strategy_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"failover_min_available_addr_num": "1",
						"default_addr_pool_type":          "IPV4",
						"failover_lba_strategy":           "RATIO",
						"strategy_mode":                   "GEO",
						"default_lba_strategy":            "RATIO",
						"failover_addr_pool_type":         "IPV4",
						"failover_addr_pools.#":           "1",
						"instance_id":                     CHECKSET,
						"default_addr_pools.#":            "1",
						"lines.#":                         "1",
						"default_min_available_addr_num":  "1",
						"strategy_name":                   name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}
func TestAccAlicloudAlidnsAccessStrategy_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alidns_access_strategy.default"
	checkoutSupportedRegions(t, true, connectivity.AlidnsSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudAlidnsAccessStrategyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlidnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlidnsAccessStrategy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlidnsAccessStrategyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
			testAccPreCheckWithEnvVariable(t, "ALICLOUD_ICP_DOMAIN_NAME")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"access_mode":                  "AUTO",
					"strategy_name":                "${var.name}",
					"instance_id":                  "${local.gtm_instance_id}",
					"strategy_mode":                "LATENCY",
					"default_addr_pool_type":       "IPV4",
					"default_max_return_addr_num":  "2",
					"default_latency_optimization": "OPEN",
					"default_addr_pools": []map[string]interface{}{
						{
							"addr_pool_id": "${alicloud_alidns_address_pool.ipv4.0.id}",
						},
					},
					"default_min_available_addr_num": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_addr_pool_type":         "IPV4",
						"default_max_return_addr_num":    "2",
						"strategy_mode":                  "LATENCY",
						"instance_id":                    CHECKSET,
						"default_latency_optimization":   "OPEN",
						"default_addr_pools.#":           "1",
						"default_min_available_addr_num": "1",
						"strategy_name":                  name,
						"access_mode":                    "AUTO",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"failover_min_available_addr_num": "1",
					"failover_addr_pool_type":         "IPV4",
					"failover_addr_pools": []map[string]interface{}{
						{
							"addr_pool_id": "${alicloud_alidns_address_pool.ipv4.1.id}",
						},
					},
					"failover_latency_optimization": "OPEN",
					"failover_max_return_addr_num":  "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"failover_min_available_addr_num": "1",
						"failover_addr_pool_type":         "IPV4",
						"failover_addr_pools.#":           "1",
						"failover_latency_optimization":   "OPEN",
						"failover_max_return_addr_num":    "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"default_max_return_addr_num": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_max_return_addr_num": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"failover_max_return_addr_num": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"failover_max_return_addr_num": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"default_latency_optimization": "CLOSE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_latency_optimization": "CLOSE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"failover_latency_optimization": "CLOSE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"failover_latency_optimization": "CLOSE",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}
func TestAccAlicloudAlidnsAccessStrategy_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alidns_access_strategy.default"
	checkoutSupportedRegions(t, true, connectivity.AlidnsSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudAlidnsAccessStrategyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlidnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlidnsAccessStrategy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlidnsAccessStrategyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
			testAccPreCheckWithEnvVariable(t, "ALICLOUD_ICP_DOMAIN_NAME")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"access_mode":                  "AUTO",
					"strategy_name":                "${var.name}",
					"instance_id":                  "${local.gtm_instance_id}",
					"strategy_mode":                "LATENCY",
					"default_addr_pool_type":       "IPV4",
					"default_max_return_addr_num":  "2",
					"default_latency_optimization": "OPEN",
					"default_addr_pools": []map[string]interface{}{
						{
							"addr_pool_id": "${alicloud_alidns_address_pool.ipv4.0.id}",
						},
					},
					"default_min_available_addr_num":  "1",
					"failover_min_available_addr_num": "1",
					"failover_addr_pool_type":         "IPV4",
					"failover_addr_pools": []map[string]interface{}{
						{
							"addr_pool_id": "${alicloud_alidns_address_pool.ipv4.1.id}",
						},
					},
					"failover_latency_optimization": "OPEN",
					"failover_max_return_addr_num":  "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"failover_min_available_addr_num": "1",
						"default_addr_pool_type":          "IPV4",
						"default_max_return_addr_num":     "2",
						"strategy_mode":                   "LATENCY",
						"failover_addr_pool_type":         "IPV4",
						"failover_addr_pools.#":           "1",
						"instance_id":                     CHECKSET,
						"default_latency_optimization":    "OPEN",
						"default_addr_pools.#":            "1",
						"default_min_available_addr_num":  "1",
						"failover_latency_optimization":   "OPEN",
						"strategy_name":                   name,
						"failover_max_return_addr_num":    "2",
						"access_mode":                     "AUTO",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

var AlicloudAlidnsAccessStrategyMap0 = map[string]string{
	"default_addr_pools.#": CHECKSET,
	"lang":                 NOSET,
}

func AlicloudAlidnsAccessStrategyBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

variable "domain_name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_cms_alarm_contact_group" "default" {
  alarm_contact_group_name = var.name
}

data "alicloud_alidns_gtm_instances" "default" {}

resource "alicloud_alidns_gtm_instance" "default" {
  count                   = length(data.alicloud_alidns_gtm_instances.default.ids) > 0 ? 0 : 1
  instance_name           = var.name
  payment_type            = "Subscription"
  period                  = 1
  renewal_status          = "ManualRenewal"
  package_edition         = "ultimate"
  health_check_task_count = 100
  sms_notification_count  = 1000
  public_cname_mode       = "SYSTEM_ASSIGN"
  ttl                     = 60
  cname_type              = "PUBLIC"
  resource_group_id       = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  alert_group             = [alicloud_cms_alarm_contact_group.default.alarm_contact_group_name]
  public_user_domain_name = var.domain_name
  alert_config {
    sms_notice      = true
    notice_type     = "ADDR_ALERT"
    email_notice    = true
    dingtalk_notice = true
  }
}

locals {
  gtm_instance_id = length(data.alicloud_alidns_gtm_instances.default.ids) > 0 ? data.alicloud_alidns_gtm_instances.default.ids[0] : concat(alicloud_alidns_gtm_instance.default.*.id, [""])[0]
}

resource "alicloud_alidns_address_pool" "ipv4" {
  count             = 2
  address_pool_name = var.name
  instance_id       = local.gtm_instance_id
  lba_strategy      = "RATIO"
  type              = "IPV4"
  address {
    attribute_info = "{\"lineCodeRectifyType\":\"RECTIFIED\",\"lineCodes\":[\"os_namerica_us\"]}"
    remark         = "address_remark"
    address        = "1.1.1.1"
    mode           = "SMART"
    lba_weight     = 1
  }
}

resource "alicloud_alidns_address_pool" "ipv6" {
  count             = 2
  address_pool_name = var.name
  instance_id       = local.gtm_instance_id
  lba_strategy      = "RATIO"
  type              = "IPV6"
  address {
    attribute_info = "{\"lineCodeRectifyType\":\"RECTIFIED\",\"lineCodes\":[\"os_namerica_us\"]}"
    remark         = "address_remark"
    address        = "1:1:1:1:1:1:1:1"
    mode           = "SMART"
    lba_weight     = 1
  }
}

`, name, os.Getenv("ALICLOUD_ICP_DOMAIN_NAME"))
}
