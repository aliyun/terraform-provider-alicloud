// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test Wafv3 AddressBook. >>> Resource test cases, automatically generated.
// Case 测试地址簿 12719
func TestAccAliCloudWafv3AddressBook_basic12719(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_wafv3_address_book.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3AddressBookMap12719)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Wafv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3AddressBook")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3AddressBookBasicDependence12719)
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
					"description":       "test",
					"instance_id":       "${data.alicloud_wafv3_instances.default.ids.0}",
					"address_book_name": name,
					"address_list": []string{
						"100.100.100.100/32", "101.101.101.101/32", "102.102.102.102/32"},
					"address_book_type": "ip",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "test",
						"instance_id":       CHECKSET,
						"address_book_name": name,
						"address_list.#":    "3",
						"address_book_type": "ip",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "test1",
					"address_book_type": "ip",
					"address_list": []string{
						"101.101.101.101/32"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "test1",
						"address_book_type": "ip",
						"address_list.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_list": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_list.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test2",
					"address_list": []string{
						"102.102.102.102/32"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":    "test2",
						"address_list.#": "1",
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

var AlicloudWafv3AddressBookMap12719 = map[string]string{
	"id":              CHECKSET,
	"address_book_id": CHECKSET,
}

func AlicloudWafv3AddressBookBasicDependence12719(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_wafv3_instances" "default" {
}


`, name)
}

// Case 地址库 12702
func TestAccAliCloudWafv3AddressBook_basic12702(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_wafv3_address_book.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3AddressBookMap12702)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Wafv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3AddressBook")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccwafv3%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3AddressBookBasicDependence12702)
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
					"description":       "创建",
					"instance_id":       "${data.alicloud_wafv3_instances.default.ids.0}",
					"address_book_name": name,
					"address_list": []string{
						"100.100.100.100/32", "101.101.101.101/32", "102.102.102.102/32"},
					"address_book_type": "ip",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "创建",
						"instance_id":       CHECKSET,
						"address_book_name": name,
						"address_list.#":    "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "更新1",
					"address_book_name": name + "_update",
					"address_list": []string{
						"101.101.101.101/32"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "更新1",
						"address_book_name": name + "_update",
						"address_list.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "更新2",
					"address_book_name": name + "_update",
					"address_list":      []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "更新2",
						"address_book_name": name + "_update",
						"address_list.#":    "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "更新3",
					"address_book_name": name + "_update",
					"address_list": []string{
						"102.102.102.102/32"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "更新3",
						"address_book_name": name + "_update",
						"address_list.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "update4",
					"address_book_name": name + "_update",
					"address_list": []string{
						"100.100.100.100/26"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "update4",
						"address_book_name": name + "_update",
						"address_list.#":    "1",
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

var AlicloudWafv3AddressBookMap12702 = map[string]string{
	"id":              CHECKSET,
	"address_book_id": CHECKSET,
}

func AlicloudWafv3AddressBookBasicDependence12702(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_wafv3_instances" "default" {
}


`, name)
}

// Case test 12677
func TestAccAliCloudWafv3AddressBook_basic12677(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_wafv3_address_book.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3AddressBookMap12677)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Wafv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3AddressBook")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3AddressBookBasicDependence12677)
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
					"description":       "test",
					"instance_id":       "${data.alicloud_wafv3_instances.default.ids.0}",
					"address_book_name": name,
					"address_list": []string{
						"100.100.100.100/32", "101.101.101.101/32", "102.102.102.102/32"},
					"address_book_type": "ip",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "test",
						"instance_id":       CHECKSET,
						"address_book_name": name,
						"address_list.#":    "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "test1",
					"address_book_name": name + "_update",
					"address_list": []string{
						"101.101.101.101/32"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "test1",
						"address_book_name": name + "_update",
						"address_list.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_book_name": name + "_update",
					"address_list":      []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_book_name": name + "_update",
						"address_list.#":    "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "test2",
					"address_book_name": name + "_update",
					"address_list": []string{
						"102.102.102.102/32"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "test2",
						"address_book_name": name + "_update",
						"address_list.#":    "1",
					}),
				),
			},
			{
				// Overlap replace: keep 102, add 103 + 104. Exercises Update where
				// DeleteAddress + AddAddress fire in the same call.
				Config: testAccConfig(map[string]interface{}{
					"description":       "test3",
					"address_book_name": name + "_update",
					"address_list": []string{
						"102.102.102.102/32", "103.103.103.103/32", "104.104.104.104/32"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "test3",
						"address_book_name": name + "_update",
						"address_list.#":    "3",
					}),
				),
			},
			{
				// No-op re-apply: re-applying the prior config must produce an empty
				// plan. Guards against perpetual diff in Read.
				Config: testAccConfig(map[string]interface{}{
					"description":       "test3",
					"address_book_name": name + "_update",
					"address_list": []string{
						"102.102.102.102/32", "103.103.103.103/32", "104.104.104.104/32"},
				}),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "test3",
						"address_book_name": name + "_update",
						"address_list.#":    "3",
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

var AlicloudWafv3AddressBookMap12677 = map[string]string{
	"id":              CHECKSET,
	"address_book_id": CHECKSET,
}

func AlicloudWafv3AddressBookBasicDependence12677(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_wafv3_instances" "default" {
}


`, name)
}

// Test Wafv3 AddressBook. <<< Resource test cases, automatically generated.
