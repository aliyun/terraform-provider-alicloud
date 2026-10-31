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

// Test Oss BucketPolicy. >>> Resource test cases, automatically generated.
// Case BucketPolicy测试 6363
func TestAccAliCloudOssBucketPolicy_basic6363(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOssBucketPolicyMap6363)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oss.NewOssServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOssBucketPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketpolicy%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOssBucketPolicyBasicDependence6363)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy": "{\\\"Statement\\\":[{\\\"Action\\\":[\\\"oss:PutObject\\\",\\\"oss:GetObject\\\"],\\\"Effect\\\":\\\"Deny\\\",\\\"Principal\\\":[\\\"1234567890\\\"],\\\"Resource\\\":[\\\"acs:oss:*:1234567890:*/*\\\"]}],\\\"Version\\\":\\\"1\\\"}",
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": acctest.CHECKSET,
						"policy": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy": "{\\\"Statement\\\":[{\\\"Action\\\":[\\\"oss:PutObject\\\",\\\"oss:GetObject\\\"],\\\"Effect\\\":\\\"Deny\\\",\\\"Principal\\\":[\\\"1234567990\\\"],\\\"Resource\\\":[\\\"acs:oss:*:1234567890:*/*\\\"]}],\\\"Version\\\":\\\"1\\\"}",
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy": acctest.CHECKSET,
						"bucket": acctest.CHECKSET,
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

var AlicloudOssBucketPolicyMap6363 = map[string]string{}

func AlicloudOssBucketPolicyBasicDependence6363(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = var.name
  lifecycle {
    ignore_changes = [
      policy,
    ]
  }
}


`, name)
}

// Case BucketPolicy测试 6363  twin
func TestAccAliCloudOssBucketPolicy_basic6363_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOssBucketPolicyMap6363)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oss.NewOssServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOssBucketPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketpolicy%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOssBucketPolicyBasicDependence6363)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy": "{\\\"Statement\\\":[{\\\"Action\\\":[\\\"oss:PutObject\\\",\\\"oss:GetObject\\\"],\\\"Effect\\\":\\\"Deny\\\",\\\"Principal\\\":[\\\"1234567990\\\"],\\\"Resource\\\":[\\\"acs:oss:*:1234567890:*/*\\\"]}],\\\"Version\\\":\\\"1\\\"}",
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy": acctest.CHECKSET,
						"bucket": acctest.CHECKSET,
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

// Test Oss BucketPolicy. <<< Resource test cases, automatically generated.
