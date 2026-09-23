package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Oss AccessPoint. >>> Resource test cases, automatically generated.
// Case AccessPoint测试 6642
func TestAccAliCloudOssAccessPoint_basic6642(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_access_point.default"
	ra := resourceAttrInit(resourceId, AlicloudOssAccessPointMap6642)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssAccessPoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossaccesspoint%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssAccessPointBasicDependence6642)
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
					"access_point_name": name,
					"bucket":            "${alicloud_oss_bucket.CreateBucket.bucket}",
					"vpc_configuration": []map[string]interface{}{
						{
							"vpc_id": "vpc-abctest",
						},
					},
					"network_origin": "vpc",
					"public_access_block_configuration": []map[string]interface{}{
						{
							"block_public_access": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_point_name": name,
						"bucket":            CHECKSET,
						"network_origin":    "vpc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"public_access_block_configuration": []map[string]interface{}{
						{
							"block_public_access": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlicloudOssAccessPointMap6642 = map[string]string{
	"status": CHECKSET,
}

func AlicloudOssAccessPointBasicDependence6642(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
}


`, name)
}

// Test Oss AccessPoint. <<< Resource test cases, automatically generated.
