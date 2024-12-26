package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Apig HttpApi. >>> Resource test cases, automatically generated.
// Case HttpApi测试_CCApi_2 9288
func TestAccAliCloudApigHttpApi_basic9288(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_http_api.default"
	ra := resourceAttrInit(resourceId, AlicloudApigHttpApiMap9288)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigHttpApi")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapighttpapi%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigHttpApiBasicDependence9288)
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
					"http_api_name": name,
					"protocols": []string{
						"${var.protocol}"},
					"base_path":         "/v1",
					"description":       "zhiwei_pop_testcase",
					"type":              "Rest",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http_api_name":     name,
						"protocols.#":       "1",
						"base_path":         "/v1",
						"description":       "zhiwei_pop_testcase",
						"type":              "Rest",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocols": []string{
						"${var.protocol_https}", "${var.protocol}"},
					"base_path":         "/v2",
					"description":       "1735184737",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocols.#":       "2",
						"base_path":         "/v2",
						"description":       CHECKSET,
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocols": []string{
						"${var.protocol}"},
					"base_path": "/v1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocols.#": "1",
						"base_path":   "/v1",
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

var AlicloudApigHttpApiMap9288 = map[string]string{}

func AlicloudApigHttpApiBasicDependence9288(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "protocol" {
  default = "HTTP"
}

variable "protocol_https" {
  default = "HTTPS"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case HttpApi测试_CCApi 9287
func TestAccAliCloudApigHttpApi_basic9287(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_http_api.default"
	ra := resourceAttrInit(resourceId, AlicloudApigHttpApiMap9287)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigHttpApi")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapighttpapi%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigHttpApiBasicDependence9287)
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
					"http_api_name": name,
					"protocols": []string{
						"${var.protocol}"},
					"base_path":         "/v1",
					"description":       "zhiwei_pop_testcase",
					"type":              "Rest",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http_api_name":     name,
						"protocols.#":       "1",
						"base_path":         "/v1",
						"description":       "zhiwei_pop_testcase",
						"type":              "Rest",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocols": []string{
						"${var.protocol_https}", "${var.protocol}"},
					"base_path":         "/v2",
					"description":       "1735184737",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocols.#":       "2",
						"base_path":         "/v2",
						"description":       CHECKSET,
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocols": []string{
						"${var.protocol}"},
					"base_path": "/v1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocols.#": "1",
						"base_path":   "/v1",
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

var AlicloudApigHttpApiMap9287 = map[string]string{}

func AlicloudApigHttpApiBasicDependence9287(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "protocol" {
  default = "HTTP"
}

variable "protocol_https" {
  default = "HTTPS"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case HttpApi测试_副本1732026511257 9021
func TestAccAliCloudApigHttpApi_basic9021(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_http_api.default"
	ra := resourceAttrInit(resourceId, AlicloudApigHttpApiMap9021)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigHttpApi")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapighttpapi%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigHttpApiBasicDependence9021)
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
					"http_api_name": name,
					"protocols": []string{
						"${var.protocol}"},
					"base_path":         "/v1",
					"description":       "zhiwei_pop_testcase",
					"type":              "Rest",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http_api_name":     name,
						"protocols.#":       "1",
						"base_path":         "/v1",
						"description":       "zhiwei_pop_testcase",
						"type":              "Rest",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocols": []string{
						"${var.protocol_https}", "${var.protocol}"},
					"base_path":         "/v2",
					"description":       "1735184737",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocols.#":       "2",
						"base_path":         "/v2",
						"description":       CHECKSET,
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocols": []string{
						"${var.protocol}"},
					"base_path": "/v1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocols.#": "1",
						"base_path":   "/v1",
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

var AlicloudApigHttpApiMap9021 = map[string]string{}

func AlicloudApigHttpApiBasicDependence9021(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "protocol" {
  default = "HTTP"
}

variable "protocol_https" {
  default = "HTTPS"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case HttpApi测试 6880
func TestAccAliCloudApigHttpApi_basic6880(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_http_api.default"
	ra := resourceAttrInit(resourceId, AlicloudApigHttpApiMap6880)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigHttpApi")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sapighttpapi%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigHttpApiBasicDependence6880)
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
					"http_api_name": name,
					"protocols": []string{
						"${var.protocol}"},
					"base_path":   "/v1",
					"description": "zhiwei_pop_testcase",
					"type":        "Rest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http_api_name": name,
						"protocols.#":   "1",
						"base_path":     "/v1",
						"description":   "zhiwei_pop_testcase",
						"type":          "Rest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocols": []string{
						"${var.protocol_https}", "${var.protocol}"},
					"base_path":   "/v2",
					"description": "1735184737",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocols.#": "2",
						"base_path":   "/v2",
						"description": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocols": []string{
						"${var.protocol}"},
					"base_path": "/v1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocols.#": "1",
						"base_path":   "/v1",
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

var AlicloudApigHttpApiMap6880 = map[string]string{}

func AlicloudApigHttpApiBasicDependence6880(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "protocol" {
  default = "HTTP"
}

variable "protocol_https" {
  default = "HTTPS"
}


`, name)
}

// Test Apig HttpApi. <<< Resource test cases, automatically generated.
