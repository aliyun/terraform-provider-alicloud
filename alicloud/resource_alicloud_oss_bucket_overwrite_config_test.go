// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Oss BucketOverwriteConfig. >>> Resource test cases, automatically generated.
// Case 测试BucketOverwriteConfig 12526
func TestAccAliCloudOssBucketOverwriteConfig_basic12526(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_overwrite_config.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketOverwriteConfigMap12526)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketOverwriteConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccoss%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketOverwriteConfigBasicDependence12526)
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
					"bucket": "${alicloud_oss_bucket.defaultrdrM3m.id}",
					"rule": []map[string]interface{}{
						{
							"id":     "rule1",
							"action": "forbid",
							"prefix": "rule1-prefix/",
							"suffix": "rule1-suffix/",
							"principals": []map[string]interface{}{
								{
									"principal": []string{
										"a", "b", "c"},
								},
							},
						},
						{
							"id":     "rule2",
							"action": "forbid",
							"prefix": "rule2-prefix/",
							"suffix": "rule2-suffix/",
							"principals": []map[string]interface{}{
								{
									"principal": []string{
										"d", "e", "f"},
								},
							},
						},
						{
							"id":     "rule3",
							"action": "forbid",
							"prefix": "rule3-prefix/",
							"suffix": "rule3-suffix/",
							"principals": []map[string]interface{}{
								{
									"principal": []string{
										"1", "2", "3"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": CHECKSET,
						"rule.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule": []map[string]interface{}{
						{
							"id":     "rule1",
							"action": "forbid",
							"prefix": "prefix/",
							"suffix": "suffix/",
							"principals": []map[string]interface{}{
								{
									"principal": []string{
										"x"},
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
							"id":     "rule1",
							"action": "forbid",
							"prefix": "a/",
							"suffix": "b/",
							"principals": []map[string]interface{}{
								{
									"principal": []string{},
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudOssBucketOverwriteConfigMap12526 = map[string]string{}

func AlicloudOssBucketOverwriteConfigBasicDependence12526(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "defaultrdrM3m" {
  storage_class = "Standard"
}


`, name)
}

// Test Oss BucketOverwriteConfig. <<< Resource test cases, automatically generated.
