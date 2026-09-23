package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Oss BucketServerSideEncryption. >>> Resource test cases, automatically generated.
// Case ServerSideEncryption测试 6458
func TestAccAliCloudOssBucketServerSideEncryption_basic6458(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_server_side_encryption.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketServerSideEncryptionMap6458)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketServerSideEncryption")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketserversideencryption%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketServerSideEncryptionBasicDependence6458)
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
					"sse_algorithm": "SM4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":        CHECKSET,
						"sse_algorithm": "SM4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sse_algorithm": "AES256",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sse_algorithm": "AES256",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kms_data_encryption": "SM4",
					"kms_master_key_id":   "${alicloud_kms_key.GetKMS.id}",
					"bucket":              "${alicloud_oss_bucket.CreateBucket.bucket}",
					"sse_algorithm":       "KMS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kms_data_encryption": "SM4",
						"kms_master_key_id":   CHECKSET,
						"bucket":              CHECKSET,
						"sse_algorithm":       "KMS",
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

var AlicloudOssBucketServerSideEncryptionMap6458 = map[string]string{}

func AlicloudOssBucketServerSideEncryptionBasicDependence6458(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = var.name
  lifecycle {
    ignore_changes = [
      server_side_encryption_rule,
    ]
  }
}

resource "alicloud_kms_key" "GetKMS" {
  origin             = "Aliyun_KMS"
  protection_level   = "SOFTWARE"
  description        = "用于测试OSS服务端加密"
  key_spec           = "Aliyun_AES_256"
  key_usage          = "ENCRYPT/DECRYPT"
  automatic_rotation = "Disabled"
  pending_window_in_days = 7
}


`, name)
}

// Case ServerSideEncryption测试 6458  twin
func TestAccAliCloudOssBucketServerSideEncryption_basic6458_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_server_side_encryption.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketServerSideEncryptionMap6458)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketServerSideEncryption")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketserversideencryption%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketServerSideEncryptionBasicDependence6458)
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
					"kms_data_encryption": "SM4",
					"kms_master_key_id":   "${alicloud_kms_key.GetKMS.id}",
					"bucket":              "${alicloud_oss_bucket.CreateBucket.bucket}",
					"sse_algorithm":       "KMS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kms_data_encryption": "SM4",
						"kms_master_key_id":   CHECKSET,
						"bucket":              CHECKSET,
						"sse_algorithm":       "KMS",
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

// Test Oss BucketServerSideEncryption. <<< Resource test cases, automatically generated.
