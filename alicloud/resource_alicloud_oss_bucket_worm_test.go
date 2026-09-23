package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Oss BucketWorm. >>> Resource test cases, automatically generated.
// Case 测试BucketWorm锁定 9223
func TestAccAliCloudOssBucketWorm_basic9223(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_worm.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketWormMap9223)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketWorm")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketworm%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketWormBasicDependence9223)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":                   "${alicloud_oss_bucket.defaulthNMfIF.bucket}",
					"retention_period_in_days": "3",
					"status":                   "InProgress",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                   CHECKSET,
						"retention_period_in_days": "3",
						"status":                   "InProgress",
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
					"status":                   "Locked",
					"retention_period_in_days": "4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                   "Locked",
						"retention_period_in_days": "4",
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

var AlicloudOssBucketWormMap9223 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"worm_id":     CHECKSET,
}

func AlicloudOssBucketWormBasicDependence9223(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "defaulthNMfIF" {
  storage_class = "Standard"
  bucket = var.name
}


`, name)
}

// Case 测试BucketWorm取消 9745
func TestAccAliCloudOssBucketWorm_basic9745(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_worm.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketWormMap9745)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketWorm")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketworm%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketWormBasicDependence9745)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":                   "${alicloud_oss_bucket.defaulthNMfIF.bucket}",
					"retention_period_in_days": "1",
					"status":                   "Locked",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                   CHECKSET,
						"retention_period_in_days": "1",
						"status":                   "Locked",
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

var AlicloudOssBucketWormMap9745 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"worm_id":     CHECKSET,
}

func AlicloudOssBucketWormBasicDependence9745(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "defaulthNMfIF" {
  storage_class = "Standard"
  bucket = var.name
}


`, name)
}

// Case 测试BucketWorm锁定2 9746
func TestAccAliCloudOssBucketWorm_basic9746(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_worm.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketWormMap9746)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketWorm")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketworm%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketWormBasicDependence9746)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":                   "${alicloud_oss_bucket.defaulthNMfIF.bucket}",
					"retention_period_in_days": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":                   CHECKSET,
						"retention_period_in_days": "1",
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

var AlicloudOssBucketWormMap9746 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"worm_id":     CHECKSET,
}

func AlicloudOssBucketWormBasicDependence9746(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "defaulthNMfIF" {
  storage_class = "Standard"
  bucket = var.name
}


`, name)
}

// Test Oss BucketWorm. <<< Resource test cases, automatically generated.
