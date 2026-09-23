package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Oss BucketPolicy. >>> Resource test cases, automatically generated.
// Case BucketPolicy测试 6363
func TestAccAliCloudOssBucketPolicy_basic6363(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketPolicyMap6363)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketPolicyBasicDependence6363)
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
					"policy": "{\\\"Statement\\\":[{\\\"Action\\\":[\\\"oss:PutObject\\\",\\\"oss:GetObject\\\"],\\\"Effect\\\":\\\"Deny\\\",\\\"Principal\\\":[\\\"1234567890\\\"],\\\"Resource\\\":[\\\"acs:oss:*:1234567890:*/*\\\"]}],\\\"Version\\\":\\\"1\\\"}",
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": CHECKSET,
						"policy": CHECKSET,
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
						"policy": CHECKSET,
						"bucket": CHECKSET,
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
	ra := resourceAttrInit(resourceId, AlicloudOssBucketPolicyMap6363)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketPolicyBasicDependence6363)
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
					"policy": "{\\\"Statement\\\":[{\\\"Action\\\":[\\\"oss:PutObject\\\",\\\"oss:GetObject\\\"],\\\"Effect\\\":\\\"Deny\\\",\\\"Principal\\\":[\\\"1234567990\\\"],\\\"Resource\\\":[\\\"acs:oss:*:1234567890:*/*\\\"]}],\\\"Version\\\":\\\"1\\\"}",
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy": CHECKSET,
						"bucket": CHECKSET,
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
