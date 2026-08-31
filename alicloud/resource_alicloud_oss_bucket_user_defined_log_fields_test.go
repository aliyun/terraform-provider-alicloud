package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Oss BucketUserDefinedLogFields. >>> Resource test cases, automatically generated.
// Case UserDefinedLogFields测试 6647
func TestAccAliCloudOssBucketUserDefinedLogFields_basic6647(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_user_defined_log_fields.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketUserDefinedLogFieldsMap6647)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketUserDefinedLogFields")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketuserdefinedlogfields%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketUserDefinedLogFieldsBasicDependence6647)
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
					"param_set": []string{
						"oss-test", "test-para", "abc"},
					"header_set": []string{
						"def", "test-header"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":       CHECKSET,
						"param_set.#":  "3",
						"header_set.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"param_set": []string{
						"oss-test", "test-para", "abc"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"param_set.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"header_set": []string{
						"def", "test-header"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"header_set.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"param_set": []string{
						"oss-test"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"param_set.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"header_set": []string{
						"test-abc", "def", "xozy"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"header_set.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"param_set": []string{
						"abc"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"param_set.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"header_set": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"header_set.#": "0",
					}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"param_set": []string{},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"param_set.#": "0",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"header_set": []string{
						"def"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"header_set.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
					"param_set": []string{
						"oss-test", "test-para", "abc"},
					"header_set": []string{
						"def", "test-header"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":       CHECKSET,
						"param_set.#":  "3",
						"header_set.#": "2",
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

var AlicloudOssBucketUserDefinedLogFieldsMap6647 = map[string]string{}

func AlicloudOssBucketUserDefinedLogFieldsBasicDependence6647(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = var.name
}


`, name)
}

// Case UserDefinedLogFields测试 6647  twin
func TestAccAliCloudOssBucketUserDefinedLogFields_basic6647_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_user_defined_log_fields.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketUserDefinedLogFieldsMap6647)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketUserDefinedLogFields")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketuserdefinedlogfields%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketUserDefinedLogFieldsBasicDependence6647)
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
					"param_set": []string{
						"oss-test", "test-para", "abc"},
					"header_set": []string{
						"def", "test-header"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":       CHECKSET,
						"param_set.#":  "3",
						"header_set.#": "2",
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

// Case UserDefinedLogFields测试 6647  raw
func TestAccAliCloudOssBucketUserDefinedLogFields_basic6647_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_user_defined_log_fields.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketUserDefinedLogFieldsMap6647)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketUserDefinedLogFields")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketuserdefinedlogfields%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketUserDefinedLogFieldsBasicDependence6647)
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
					"param_set": []string{
						"oss-test", "test-para", "abc"},
					"header_set": []string{
						"def", "test-header"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":       CHECKSET,
						"param_set.#":  "3",
						"header_set.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"param_set": []string{
						"oss-test"},
					"header_set": []string{
						"test-abc", "def", "xozy"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"param_set.#":  "1",
						"header_set.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"param_set": []string{
						"abc"},
					"header_set": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"param_set.#":  "1",
						"header_set.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"param_set": []string{},
					"header_set": []string{
						"def"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"param_set.#":  "0",
						"header_set.#": "1",
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

// Test Oss BucketUserDefinedLogFields. <<< Resource test cases, automatically generated.
