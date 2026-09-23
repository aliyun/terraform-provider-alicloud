// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Oss BucketResponseHeader. >>> Resource test cases, automatically generated.
// Case BucketResponseHeader测试 7368
func TestAccAliCloudOssBucketResponseHeader_basic7368(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_response_header.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketResponseHeaderMap7368)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketResponseHeader")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccoss%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketResponseHeaderBasicDependence7368)
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
					"bucket": "${alicloud_oss_bucket.CreateBucket.id}",
					"rule": []map[string]interface{}{
						{
							"name": "name1",
							"filters": []map[string]interface{}{
								{
									"operation": []string{
										"GetObject", "Put*"},
								},
							},
							"hide_headers": []map[string]interface{}{
								{
									"header": []string{
										"Last-Modified"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": CHECKSET,
						"rule.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule": []map[string]interface{}{
						{
							"name": "name2",
							"filters": []map[string]interface{}{
								{
									"operation": []string{
										"Get*"},
								},
							},
							"hide_headers": []map[string]interface{}{
								{
									"header": []string{
										"x-oss-abc"},
								},
							},
						},
						{
							"name": "name3",
							"filters": []map[string]interface{}{
								{
									"operation": []string{
										"Delete"},
								},
							},
							"hide_headers": []map[string]interface{}{
								{
									"header": []string{
										"x-oss-def"},
								},
							},
						},
						{
							"name": "name4",
							"filters": []map[string]interface{}{
								{
									"operation": []string{
										"Get"},
								},
							},
							"hide_headers": []map[string]interface{}{
								{
									"header": []string{
										"x-oss-hij"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule": []map[string]interface{}{
						{
							"name": "name5",
							"filters": []map[string]interface{}{
								{
									"operation": []string{
										"Get"},
								},
							},
							"hide_headers": []map[string]interface{}{
								{
									"header": []string{
										"x-oss-xyz"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule": []map[string]interface{}{
						{
							"filters": []map[string]interface{}{
								{
									"operation": []string{
										"GetObject", "GetObjectMeta", "PutObject"},
								},
							},
							"hide_headers": []map[string]interface{}{
								{
									"header": []string{
										"x-oss-1", "x-oss-2", "x-oss-3"},
								},
							},
							"name": "name6",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule.#": "1",
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

var AlicloudOssBucketResponseHeaderMap7368 = map[string]string{}

func AlicloudOssBucketResponseHeaderBasicDependence7368(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
}


`, name)
}

// Test Oss BucketResponseHeader. <<< Resource test cases, automatically generated.
