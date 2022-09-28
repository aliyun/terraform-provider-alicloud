package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/alibabacloud-go/tea-rpc/client"
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

func TestUnitAlicloudAlidnsAccessStrategy(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_alidns_access_strategy"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_alidns_access_strategy"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"default_addr_pools": []interface{}{
			map[string]interface{}{
				"lba_weight":   1,
				"addr_pool_id": "AddDnsGtmAccessStrategyValue",
			},
		},
		"default_addr_pool_type":         "AddDnsGtmAccessStrategyValue",
		"default_latency_optimization":   "AddDnsGtmAccessStrategyValue",
		"default_lba_strategy":           "AddDnsGtmAccessStrategyValue",
		"default_max_return_addr_num":    10,
		"default_min_available_addr_num": 10,
		"lines": []interface{}{
			map[string]interface{}{
				"line_code": "AddDnsGtmAccessStrategyValue",
			},
		},
		"failover_addr_pools": []interface{}{
			map[string]interface{}{
				"lba_weight":   1,
				"addr_pool_id": "AddDnsGtmAccessStrategyValue",
			},
		},
		"failover_addr_pool_type":         "AddDnsGtmAccessStrategyValue",
		"failover_latency_optimization":   "AddDnsGtmAccessStrategyValue",
		"failover_lba_strategy":           "AddDnsGtmAccessStrategyValue",
		"failover_max_return_addr_num":    10,
		"failover_min_available_addr_num": 10,
		"instance_id":                     "AddDnsGtmAccessStrategyValue",
		"lang":                            "AddDnsGtmAccessStrategyValue",
		"strategy_mode":                   "AddDnsGtmAccessStrategyValue",
		"strategy_name":                   "AddDnsGtmAccessStrategyValue",
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
		if err != nil {
			log.Printf("[ERROR] the field %s setting error", key)
		}
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// DescribeDnsGtmAccessStrategy
		"DefaultAddrPoolType": "AddDnsGtmAccessStrategyValue",
		"DefaultAddrPools": map[string]interface{}{
			"DefaultAddrPool": []interface{}{
				map[string]interface{}{
					"Id":        "AddDnsGtmAccessStrategyValue",
					"LbaWeight": 1,
				},
			},
		},
		"DefaultLatencyOptimization": "AddDnsGtmAccessStrategyValue",
		"DefaultLbaStrategy":         "AddDnsGtmAccessStrategyValue",
		"DefaultMaxReturnAddrNum":    10,
		"DefaultMinAvailableAddrNum": 10,
		"FailoverAddrPoolType":       "AddDnsGtmAccessStrategyValue",
		"FailoverAddrPools": map[string]interface{}{
			"FailoverAddrPool": []interface{}{
				map[string]interface{}{
					"Id":        "AddDnsGtmAccessStrategyValue",
					"LbaWeight": 1,
				},
			},
		},
		"Lines": map[string]interface{}{
			"Line": []interface{}{
				map[string]interface{}{
					"LineCode": "AddDnsGtmAccessStrategyValue",
				},
			},
		},
		"FailoverLatencyOptimization": "AddDnsGtmAccessStrategyValue",
		"FailoverLbaStrategy":         "AddDnsGtmAccessStrategyValue",
		"FailoverMaxReturnAddrNum":    10,
		"FailoverMinAvailableAddrNum": 10,
		"InstanceId":                  "AddDnsGtmAccessStrategyValue",
		"StrategyMode":                "AddDnsGtmAccessStrategyValue",
		"StrategyName":                "AddDnsGtmAccessStrategyValue",
		"AccessMode":                  "AddDnsGtmAccessStrategyValue",
		"StrategyId":                  "AddDnsGtmAccessStrategyValue",
	}
	CreateMockResponse := map[string]interface{}{
		// AddDnsGtmAccessStrategy
		"StrategyId": "AddDnsGtmAccessStrategyValue",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_alidns_access_strategy", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAlidnsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudAlidnsAccessStrategyCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeDnsGtmAccessStrategy Response
		"StrategyId": "AddDnsGtmAccessStrategyValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "AddDnsGtmAccessStrategy" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudAlidnsAccessStrategyCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_alidns_access_strategy"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAlidnsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudAlidnsAccessStrategyUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)

	// UpdateDnsGtmAccessStrategy
	attributesDiff := map[string]interface{}{
		"access_mode":            "UpdateDnsGtmAccessStrategyValue",
		"default_addr_pool_type": "UpdateDnsGtmAccessStrategyValue",
		"lines": []interface{}{
			map[string]interface{}{
				"line_code": "UpdateDnsGtmAccessStrategyValue",
			},
		},
		"default_addr_pools": []interface{}{
			map[string]interface{}{
				"lba_weight":   2,
				"addr_pool_id": "UpdateDnsGtmAccessStrategyValue",
			},
		},
		"failover_addr_pools": []map[string]interface{}{
			{
				"lba_weight":   2,
				"addr_pool_id": "UpdateDnsGtmAccessStrategyValue",
			},
		},
		"default_min_available_addr_num":  15,
		"strategy_name":                   "UpdateDnsGtmAccessStrategyValue",
		"default_latency_optimization":    "UpdateDnsGtmAccessStrategyValue",
		"default_lba_strategy":            "UpdateDnsGtmAccessStrategyValue",
		"failover_max_return_addr_num":    15,
		"failover_addr_pool_type":         "UpdateDnsGtmAccessStrategyValue",
		"failover_latency_optimization":   "UpdateDnsGtmAccessStrategyValue",
		"failover_lba_strategy":           "UpdateDnsGtmAccessStrategyValue",
		"failover_min_available_addr_num": 15,
		"default_max_return_addr_num":     15,
		"lang":                            "UpdateDnsGtmAccessStrategyValue",
	}
	diff, err := newInstanceDiff("alicloud_alidns_access_strategy", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_alidns_access_strategy"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDnsGtmAccessStrategy Response
		"DefaultAddrPoolType": "UpdateDnsGtmAccessStrategyValue",
		"DefaultAddrPools": map[string]interface{}{
			"DefaultAddrPool": []interface{}{
				map[string]interface{}{
					"Id":        "UpdateDnsGtmAccessStrategyValue",
					"LbaWeight": 2,
				},
			},
		},
		"DefaultLatencyOptimization": "UpdateDnsGtmAccessStrategyValue",
		"DefaultLbaStrategy":         "UpdateDnsGtmAccessStrategyValue",
		"DefaultMaxReturnAddrNum":    15,
		"DefaultMinAvailableAddrNum": 15,
		"FailoverAddrPoolType":       "UpdateDnsGtmAccessStrategyValue",
		"FailoverAddrPools": map[string]interface{}{
			"FailoverAddrPool": []interface{}{
				map[string]interface{}{
					"Id":        "UpdateDnsGtmAccessStrategyValue",
					"LbaWeight": 2,
				},
			},
		},
		"Lines": map[string]interface{}{
			"Line": []interface{}{
				map[string]interface{}{
					"LineCode": "UpdateDnsGtmAccessStrategyValue",
				},
			},
		},
		"FailoverLatencyOptimization": "UpdateDnsGtmAccessStrategyValue",
		"FailoverLbaStrategy":         "UpdateDnsGtmAccessStrategyValue",
		"FailoverMaxReturnAddrNum":    15,
		"FailoverMinAvailableAddrNum": 15,
		"StrategyName":                "UpdateDnsGtmAccessStrategyValue",
		"AccessMode":                  "UpdateDnsGtmAccessStrategyValue",
		"Lang":                        "UpdateDnsGtmAccessStrategyValue",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateDnsGtmAccessStrategy" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudAlidnsAccessStrategyUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_alidns_access_strategy"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeDnsGtmAccessStrategy" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudAlidnsAccessStrategyRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAlidnsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudAlidnsAccessStrategyDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteDnsGtmAccessStrategy" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudAlidnsAccessStrategyDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		}
	}
}
