package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Oss BucketReferer. >>> Resource test cases, automatically generated.
// Case 5764
func TestAccAliCloudOssBucketReferer_basic5764(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_referer.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketRefererMap5764)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketReferer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketreferer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketRefererBasicDependence5764)
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
					"bucket":              "${alicloud_oss_bucket.bucket.bucket}",
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
					"referer_blacklist": []map[string]interface{}{
						{
							"referer": []string{
								"http://www.abc.com", "https://*.abc.com"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"referer_blacklist": REMOVEKEY,
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
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"referer_list": []map[string]interface{}{
						{
							"referer": []string{
								"http://*.aliyun.com", "https://*.aliyun.com"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allow_empty_referer": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allow_empty_referer": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"truncate_path": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"truncate_path": "true",
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
					"referer_list": []map[string]interface{}{
						{
							"referer": []string{
								"https://allow.aliyun.com", "http://allow.aliyun.com"},
						},
					},
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
				Config: testAccConfig(map[string]interface{}{
					"referer_list": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allow_empty_referer": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allow_empty_referer": "true",
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
				Config: testAccConfig(map[string]interface{}{
					"referer_blacklist": []map[string]interface{}{
						{
							"referer": []string{
								"http://www.abc.com", "https://*.abc.com"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"referer_list": []map[string]interface{}{
						{
							"referer": []string{
								"http://www.aliyun.com", "https://*.aliyun.com"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":            "${alicloud_oss_bucket.bucket.bucket}",
					"referer_blacklist": REMOVEKEY,
					"referer_list": []map[string]interface{}{
						{
							"referer": []string{
								"http://*.aliyun.com", "https://*.aliyun.com"},
						},
					},
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudOssBucketRefererMap5764 = map[string]string{}

func AlicloudOssBucketRefererBasicDependence5764(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "bucket" {
  bucket = var.name
  acl    = "private"
}


`, name)
}

// Case 5764  twin
func TestAccAliCloudOssBucketReferer_basic5764_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_referer.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketRefererMap5764)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketReferer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketreferer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketRefererBasicDependence5764)
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
					"bucket":                      "${alicloud_oss_bucket.bucket.bucket_name}",
					"allow_truncate_query_string": "true",
					"referer_blacklist": []map[string]interface{}{
						{
							"referer": []string{
								"http://www.abc.com", "https://*.abc.com", "https://www.abc.com"},
						},
					},
					"referer_list": []map[string]interface{}{
						{
							"referer": []string{
								"http://www.aliyun.com", "https://*.aliyun.com", "https://www.*aliyuncs.com"},
						},
					},
					"allow_empty_referer": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                      CHECKSET,
						"allow_truncate_query_string": "true",
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

// Test Oss BucketReferer. <<< Resource test cases, automatically generated.
