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

// Test Oss BucketTransferAcceleration. >>> Resource test cases, automatically generated.
// Case TransferAcceleration测试 6439
// lintignore: AT001
func TestAccAliCloudOssBucketTransferAcceleration_basic6439(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_transfer_acceleration.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOssBucketTransferAccelerationMap6439)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oss.NewOssServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOssBucketTransferAcceleration")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbuckettransferacceleration%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOssBucketTransferAccelerationBasicDependence6439)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": acctest.CHECKSET,
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
						"bucket":  acctest.CHECKSET,
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
// lintignore: AT001
func TestAccAliCloudOssBucketTransferAcceleration_basic6439_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_transfer_acceleration.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOssBucketTransferAccelerationMap6439)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oss.NewOssServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOssBucketTransferAcceleration")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbuckettransferacceleration%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOssBucketTransferAccelerationBasicDependence6439)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":  "${alicloud_oss_bucket.CreateBucket.bucket}",
					"enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":  acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOssBucketTransferAccelerationMap6439)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oss.NewOssServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOssBucketTransferAcceleration")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbuckettransferacceleration%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOssBucketTransferAccelerationBasicDependence6439)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":  "${alicloud_oss_bucket.CreateBucket.bucket}",
					"enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":  acctest.CHECKSET,
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
