package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Oss BucketCors. >>> Resource test cases, automatically generated.
// Case BucketCors测试 6362
func TestAccAliCloudOssBucketCors_basic6362(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_cors.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketCorsMap6362)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketCors")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketcors%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketCorsBasicDependence6362)
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
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
					"cors_rule": []map[string]interface{}{
						{
							"allowed_methods": []string{
								"GET"},
							"allowed_origins": []string{
								"*"},
							"allowed_headers": []string{
								"x-oss-test", "x-oss-abc"},
							"expose_header": []string{
								"x-oss-request-id"},
							"max_age_seconds": "1000",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cors_rule.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"response_vary": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"response_vary": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cors_rule": []map[string]interface{}{
						{
							"allowed_methods": []string{
								"POST", "HEAD", "GET", "DELETE"},
							"allowed_origins": []string{
								"oss.aliyuncs.com"},
							"allowed_headers": []string{
								"*"},
							"max_age_seconds": "100",
						},
						{
							"allowed_methods": []string{
								"PUT"},
							"allowed_origins": []string{
								"allow.aliyuncs.com.", "*.aliyuncs.com"},
							"allowed_headers": []string{
								"test-oss"},
							"expose_header": []string{
								"x-oss-meta"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cors_rule.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":        "${alicloud_oss_bucket.CreateBucket.bucket}",
					"response_vary": "true",
					"cors_rule": []map[string]interface{}{
						{
							"allowed_methods": []string{
								"GET"},
							"allowed_origins": []string{
								"*"},
							"allowed_headers": []string{
								"x-oss-test", "x-oss-abc"},
							"expose_header": []string{
								"x-oss-request-id"},
							"max_age_seconds": "1000",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":        CHECKSET,
						"response_vary": "true",
						"cors_rule.#":   "1",
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

var AlicloudOssBucketCorsMap6362 = map[string]string{
	"response_vary": "false",
}

func AlicloudOssBucketCorsBasicDependence6362(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = var.name
  lifecycle {
    ignore_changes = [
      cors_rule,
    ]
  }
}


`, name)
}

// Case BucketCors测试 6362  twin
func TestAccAliCloudOssBucketCors_basic6362_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_cors.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketCorsMap6362)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketCors")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketcors%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketCorsBasicDependence6362)
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
					"bucket":        "${alicloud_oss_bucket.CreateBucket.bucket}",
					"response_vary": "false",
					"cors_rule": []map[string]interface{}{
						{
							"allowed_methods": []string{
								"POST", "HEAD", "GET", "DELETE"},
							"allowed_origins": []string{
								"oss.aliyuncs.com"},
							"allowed_headers": []string{
								"*", "x-oss-abc"},
							"expose_header": []string{
								"x-oss-request-id"},
							"max_age_seconds": "100",
						},
						{
							"allowed_methods": []string{
								"PUT"},
							"allowed_origins": []string{
								"allow.aliyuncs.com.", "*.aliyuncs.com"},
							"allowed_headers": []string{
								"test-oss"},
							"expose_header": []string{
								"x-oss-meta"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":        CHECKSET,
						"response_vary": "false",
						"cors_rule.#":   "2",
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

// Test Oss BucketCors. <<< Resource test cases, automatically generated.
