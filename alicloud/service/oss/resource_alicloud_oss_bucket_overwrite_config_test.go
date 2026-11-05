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
// collection-order: bypass rule.principals.principal -- pre-existing TypeList attribute surfaced by verbatim stock-directory migration; schema/ordering semantics unchanged (byte-identical move to alicloud/service/oss)
// Test Oss BucketOverwriteConfig. >>> Resource test cases, automatically generated.
// Case 测试BucketOverwriteConfig 12526
func TestAccAliCloudOssBucketOverwriteConfig_basic12526(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_overwrite_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOssBucketOverwriteConfigMap12526)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oss.NewOssServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOssBucketOverwriteConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccoss%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOssBucketOverwriteConfigBasicDependence12526)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"bucket": acctest.CHECKSET,
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
