package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Oss BucketTransferAcceleration. >>> Resource test cases, automatically generated.
// Case TransferAcceleration测试 6439
func TestAccAliCloudOssBucketTransferAcceleration_basic6439(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_transfer_acceleration.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketTransferAccelerationMap6439)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketTransferAcceleration")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbuckettransferacceleration%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketTransferAccelerationBasicDependence6439)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":  "${alicloud_oss_bucket.CreateBucket.bucket}",
					"enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":  CHECKSET,
						"enabled": "true",
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

var AlicloudOssBucketTransferAccelerationMap6439 = map[string]string{}

func AlicloudOssBucketTransferAccelerationBasicDependence6439(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = var.name
  lifecycle {
    ignore_changes = [
      transfer_acceleration,
    ]
  }
}


`, name)
}

// Case TransferAcceleration测试 6439  twin
func TestAccAliCloudOssBucketTransferAcceleration_basic6439_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_transfer_acceleration.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketTransferAccelerationMap6439)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketTransferAcceleration")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbuckettransferacceleration%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketTransferAccelerationBasicDependence6439)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":  "${alicloud_oss_bucket.CreateBucket.bucket}",
					"enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":  CHECKSET,
						"enabled": "true",
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

// Case TransferAcceleration测试 6439  raw
func TestAccAliCloudOssBucketTransferAcceleration_basic6439_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_transfer_acceleration.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketTransferAccelerationMap6439)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketTransferAcceleration")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbuckettransferacceleration%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketTransferAccelerationBasicDependence6439)
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
					"bucket":  "${alicloud_oss_bucket.CreateBucket.bucket}",
					"enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":  CHECKSET,
						"enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "false",
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

// Test Oss BucketTransferAcceleration. <<< Resource test cases, automatically generated.
