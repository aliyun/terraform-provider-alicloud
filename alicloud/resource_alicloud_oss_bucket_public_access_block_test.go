package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Oss BucketPublicAccessBlock. >>> Resource test cases, automatically generated.
// Case Bucket阻止公共访问 6640
func TestAccAliCloudOssBucketPublicAccessBlock_basic6640(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_public_access_block.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketPublicAccessBlockMap6640)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketPublicAccessBlock")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketpublicaccessblock%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketPublicAccessBlockBasicDependence6640)
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
					"block_public_access": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":              CHECKSET,
						"block_public_access": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"block_public_access": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"block_public_access": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":              "${alicloud_oss_bucket.CreateBucket.bucket}",
					"block_public_access": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":              CHECKSET,
						"block_public_access": "true",
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

var AlicloudOssBucketPublicAccessBlockMap6640 = map[string]string{}

func AlicloudOssBucketPublicAccessBlockBasicDependence6640(name string) string {
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

// Case Bucket阻止公共访问 6640  twin
func TestAccAliCloudOssBucketPublicAccessBlock_basic6640_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_public_access_block.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketPublicAccessBlockMap6640)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketPublicAccessBlock")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketpublicaccessblock%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketPublicAccessBlockBasicDependence6640)
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
					"block_public_access": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":              CHECKSET,
						"block_public_access": "true",
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

// Case Bucket阻止公共访问 6640  raw
func TestAccAliCloudOssBucketPublicAccessBlock_basic6640_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_public_access_block.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketPublicAccessBlockMap6640)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketPublicAccessBlock")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketpublicaccessblock%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketPublicAccessBlockBasicDependence6640)
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
					"block_public_access": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":              CHECKSET,
						"block_public_access": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"block_public_access": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"block_public_access": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"block_public_access": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"block_public_access": "true",
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

// Test Oss BucketPublicAccessBlock. <<< Resource test cases, automatically generated.
