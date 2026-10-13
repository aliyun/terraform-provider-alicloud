package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Oss AccessPointPolicy. >>> Resource test cases, automatically generated.
// Case AccessPointPolicy测试 6710
func TestAccAliCloudOssAccessPointPolicy_basic6710(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_access_point_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudOssAccessPointPolicyMap6710)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssAccessPointPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossaccesspointpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssAccessPointPolicyBasicDependence6710)
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
					"access_point_name": "${alicloud_oss_access_point.CreateAccessPoint.access_point_name}",
					"bucket":            "${alicloud_oss_bucket.CreateBucket.bucket}",
					"policy":            "{\\\"Statement\\\":[{\\\"Action\\\":[\\\"oss:PutObject\\\",\\\"oss:GetObject\\\"],\\\"Effect\\\":\\\"Allow\\\",\\\"Principal\\\":[\\\"${data.alicloud_caller_identity.current.account_id}\\\"],\\\"Resource\\\":[\\\"acs:oss:cn-hangzhou:${data.alicloud_caller_identity.current.account_id}:accesspoint/${alicloud_oss_access_point.CreateAccessPoint.access_point_name}/object/*\\\"]}],\\\"Version\\\":\\\"1\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_point_name": CHECKSET,
						"bucket":            CHECKSET,
						"policy":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_point_name": "${alicloud_oss_access_point.CreateAccessPoint.access_point_name}",
					"bucket":            "${alicloud_oss_bucket.CreateBucket.bucket}",
					"policy":            "{\\\"Statement\\\":[{\\\"Action\\\":[\\\"oss:PutObject\\\",\\\"oss:GetObject\\\"],\\\"Effect\\\":\\\"Deny\\\",\\\"Principal\\\":[\\\"${data.alicloud_caller_identity.current.account_id}\\\"],\\\"Resource\\\":[\\\"acs:oss:cn-hangzhou:${data.alicloud_caller_identity.current.account_id}:accesspoint/${alicloud_oss_access_point.CreateAccessPoint.access_point_name}/object/*\\\"]}],\\\"Version\\\":\\\"1\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy":            CHECKSET,
						"access_point_name": CHECKSET,
						"bucket":            CHECKSET,
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

var AlicloudOssAccessPointPolicyMap6710 = map[string]string{}

func AlicloudOssAccessPointPolicyBasicDependence6710(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = var.name
}

resource "alicloud_oss_access_point" "CreateAccessPoint" {
  access_point_name = "${var.name}-ap"
  bucket            = alicloud_oss_bucket.CreateBucket.bucket
  network_origin    = "internet"
}

data "alicloud_caller_identity" "current" {}


`, name)
}

// Case AccessPointPolicy测试 6710  twin
func TestAccAliCloudOssAccessPointPolicy_basic6710_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_access_point_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudOssAccessPointPolicyMap6710)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssAccessPointPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossaccesspointpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssAccessPointPolicyBasicDependence6710)
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
					"access_point_name": "${alicloud_oss_access_point.CreateAccessPoint.access_point_name}",
					"bucket":            "${alicloud_oss_bucket.CreateBucket.bucket}",
					"policy":            "{\\\"Statement\\\":[{\\\"Action\\\":[\\\"oss:PutObject\\\",\\\"oss:GetObject\\\"],\\\"Effect\\\":\\\"Deny\\\",\\\"Principal\\\":[\\\"${data.alicloud_caller_identity.current.account_id}\\\"],\\\"Resource\\\":[\\\"acs:oss:cn-hangzhou:${data.alicloud_caller_identity.current.account_id}:accesspoint/${alicloud_oss_access_point.CreateAccessPoint.access_point_name}/object/*\\\"]}],\\\"Version\\\":\\\"1\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy":            CHECKSET,
						"bucket":            CHECKSET,
						"access_point_name": CHECKSET,
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

// Test Oss AccessPointPolicy. <<< Resource test cases, automatically generated.
