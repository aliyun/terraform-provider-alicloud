package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Oss BucketReferer. >>> Resource test cases, automatically generated.
// Case 4937
func TestAccAliCloudOssBucketReferer_basic4937(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_referer.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketRefererMap4937)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketReferer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketreferer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketRefererBasicDependence4937)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name":         "qd-api-test",
					"allow_empty_referer": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name":         "qd-api-test",
						"allow_empty_referer": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"referer_list": []map[string]interface{}{
						{
							"referer": []string{
								"www.abc.com", "www.aliyun.com"},
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
				Config: testAccConfig(map[string]interface{}{}),
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
								"*.aliyuncs.com", "*.alibaba-inc.com"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
								"abc.test", "*.abc.test"},
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
					"referer_list": []map[string]interface{}{
						{
							"referer": []string{
								"*aliyun.com", "*aliyuncs.com"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name":         "qd-api-test",
					"allow_empty_referer": "true",
					"referer_list": []map[string]interface{}{
						{
							"referer": []string{
								"www.abc.com", "www.aliyun.com"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name":         "qd-api-test",
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

var AlicloudOssBucketRefererMap4937 = map[string]string{}

func AlicloudOssBucketRefererBasicDependence4937(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 4937  twin
func TestAccAliCloudOssBucketReferer_basic4937_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_referer.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketRefererMap4937)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketReferer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketreferer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketRefererBasicDependence4937)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name":         "qd-api-test",
					"allow_empty_referer": "true",
					"referer_list": []map[string]interface{}{
						{
							"referer": []string{
								"*aliyun.com", "*aliyuncs.com", "https://example.com"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name":         "qd-api-test",
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

// Test Oss BucketReferer. <<< Resource test cases, automatically generated.
