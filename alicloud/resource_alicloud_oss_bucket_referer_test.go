package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Oss BucketReferer. >>> Resource test cases, automatically generated.
// Case 新版Bucket资源测试 6182
func TestAccAliCloudOssBucketReferer_basic6182(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_referer.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketRefererMap6182)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketReferer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketreferer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketRefererBasicDependence6182)
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
					"bucket":              "${alicloud_oss_bucket.CreateBucket.bucket}",
					"allow_empty_referer": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":              CHECKSET,
						"allow_empty_referer": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"truncate_path": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"truncate_path": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allow_truncate_query_string": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allow_truncate_query_string": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"referer_list": []string{
						"www.test.com", "*.test2.com"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"referer_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"truncate_path": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"truncate_path": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allow_truncate_query_string": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allow_truncate_query_string": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"referer_list": []string{
						"*.test.com", "*.aliyun.com"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"referer_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"referer_blacklist": []string{
						"*.forbidden.com"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"referer_blacklist.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allow_empty_referer": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allow_empty_referer": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":                      "${alicloud_oss_bucket.CreateBucket.bucket}",
					"truncate_path":               "true",
					"allow_truncate_query_string": "true",
					"referer_list": []string{
						"www.test.com", "*.test2.com"},
					"allow_empty_referer": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                      CHECKSET,
						"truncate_path":               "true",
						"allow_truncate_query_string": "true",
						"referer_list.#":              "2",
						"allow_empty_referer":         "true",
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

var AlicloudOssBucketRefererMap6182 = map[string]string{
	"allow_truncate_query_string": CHECKSET,
}

func AlicloudOssBucketRefererBasicDependence6182(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = var.name
  lifecycle {
    ignore_changes = [
      referer_config,
    ]
  }
}


`, name)
}

// Case 新版Bucket资源测试 6182  twin
func TestAccAliCloudOssBucketReferer_basic6182_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_referer.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketRefererMap6182)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketReferer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketreferer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketRefererBasicDependence6182)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":                      "${alicloud_oss_bucket.CreateBucket.bucket}",
					"truncate_path":               "false",
					"allow_truncate_query_string": "true",
					"referer_list": []string{
						"*.test.com", "*.aliyun.com"},
					"allow_empty_referer": "false",
					"referer_blacklist": []string{
						"*.forbidden.com"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                      CHECKSET,
						"truncate_path":               "false",
						"allow_truncate_query_string": "true",
						"referer_list.#":              "2",
						"allow_empty_referer":         "false",
						"referer_blacklist.#":         "1",
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

// Test Oss BucketReferer. <<< Resource test cases, automatically generated.
