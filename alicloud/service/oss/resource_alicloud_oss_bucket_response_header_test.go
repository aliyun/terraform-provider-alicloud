// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package oss_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/oss"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// collection-order: bypass rule -- pre-existing TypeList attribute surfaced by verbatim stock-directory migration; schema/ordering semantics unchanged (byte-identical move to alicloud/service/oss)
// collection-order: bypass rule.filters.operation -- pre-existing TypeList attribute surfaced by verbatim stock-directory migration; schema/ordering semantics unchanged (byte-identical move to alicloud/service/oss)
// collection-order: bypass rule.hide_headers.header -- pre-existing TypeList attribute surfaced by verbatim stock-directory migration; schema/ordering semantics unchanged (byte-identical move to alicloud/service/oss)
// Test Oss BucketResponseHeader. >>> Resource test cases, automatically generated.
// Case BucketResponseHeader测试 7368
func TestAccAliCloudOssBucketResponseHeader_basic7368(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_response_header.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOssBucketResponseHeaderMap7368)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oss.NewOssServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOssBucketResponseHeader")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccoss%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOssBucketResponseHeaderBasicDependence7368)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"bucket": acctest.CHECKSET,
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
