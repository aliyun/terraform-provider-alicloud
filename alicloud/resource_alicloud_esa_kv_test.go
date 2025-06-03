package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA Kv. >>> Resource test cases, automatically generated.
// Case Kv_test2
func TestAccAliCloudESAKvKv_test2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_kv.default"
	ra := resourceAttrInit(resourceId, AliCloudESAKvKv_test2Map)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaKv")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAKv%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAKvKv_test2BasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"expiration_ttl": "360",
					"expiration":     "1690",
					"namespace":      "${alicloud_esa_kv_namespace.resource_Namespace_Kv_test2.id}",
					"key":            "test_key",
					"url":            "https://yandanpub.oss-cn-hangzhou.aliyuncs.com/ces1.txt",
					"isbase":         "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"isbase": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"expiration": "1691",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"expiration_ttl": "361",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"url": "https://yandanpub.oss-cn-hangzhou.aliyuncs.com/ces1.txt",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"expiration", "expiration_ttl", "isbase", "url"},
			},
		},
	})
}

var AliCloudESAKvKv_test2Map = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAKvKv_test2BasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


resource "alicloud_esa_kv_namespace" "resource_Namespace_Kv_test2" {
  description  = "this is a test namespace."
  kv_namespace = "namespace1"
}

`, name)
}

// Case Kv_test1
func TestAccAliCloudESAKvKv_test1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_kv.default"
	ra := resourceAttrInit(resourceId, AliCloudESAKvKv_test1Map)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaKv")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAKv%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAKvKv_test1BasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"expiration_ttl": "360",
					"value":          "example_value",
					"expiration":     "1690",
					"namespace":      "${alicloud_esa_kv_namespace.resource_Namespace_Kv_test1.id}",
					"key":            "example_key",
					"isbase":         "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"isbase": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"expiration": "1691",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"expiration_ttl": "361",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"value": "test_value1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"expiration", "expiration_ttl", "isbase", "url"},
			},
		},
	})
}

var AliCloudESAKvKv_test1Map = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAKvKv_test1BasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


resource "alicloud_esa_kv_namespace" "resource_Namespace_Kv_test1" {
  description  = "this is a example namespace."
  kv_namespace = "namespace1"
}

`, name)
}

// Test ESA Kv. <<< Resource test cases, automatically generated.
