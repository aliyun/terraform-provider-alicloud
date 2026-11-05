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

// Test Oss BucketServerSideEncryption. >>> Resource test cases, automatically generated.
// Case ServerSideEncryption测试 6458
func TestAccAliCloudOssBucketServerSideEncryption_basic6458(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_server_side_encryption.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOssBucketServerSideEncryptionMap6458)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oss.NewOssServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOssBucketServerSideEncryption")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketserversideencryption%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOssBucketServerSideEncryptionBasicDependence6458)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":        "${alicloud_oss_bucket.CreateBucket.bucket}",
					"sse_algorithm": "SM4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":        acctest.CHECKSET,
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
						"kms_master_key_id":   acctest.CHECKSET,
						"bucket":              acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOssBucketServerSideEncryptionMap6458)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oss.NewOssServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOssBucketServerSideEncryption")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketserversideencryption%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOssBucketServerSideEncryptionBasicDependence6458)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"kms_master_key_id":   acctest.CHECKSET,
						"bucket":              acctest.CHECKSET,
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
