package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudOOSServiceSetting_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_service_setting.default"
	checkoutSupportedRegions(t, true, connectivity.OOSSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudOOSServiceSettingMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOosServiceSetting")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%soos%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOOSServiceSettingBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_oss_enabled":      "true",
					"delivery_oss_key_prefix":   "path1/",
					"delivery_oss_bucket_name":  "${alicloud_oss_bucket.default.0.bucket}",
					"delivery_sls_enabled":      "true",
					"delivery_sls_project_name": "${alicloud_log_project.default.0.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_oss_enabled":      "true",
						"delivery_oss_key_prefix":   "path1/",
						"delivery_oss_bucket_name":  CHECKSET,
						"delivery_sls_enabled":      "true",
						"delivery_sls_project_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_oss_key_prefix": "path2/",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_oss_key_prefix": "path2/",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_oss_bucket_name": "${alicloud_oss_bucket.default.1.bucket}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_oss_bucket_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_oss_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_oss_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_sls_project_name": "${alicloud_log_project.default.1.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_sls_project_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_sls_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_sls_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_oss_enabled":      "true",
					"delivery_oss_key_prefix":   "path1/",
					"delivery_oss_bucket_name":  "${alicloud_oss_bucket.default.0.bucket}",
					"delivery_sls_enabled":      "true",
					"delivery_sls_project_name": "${alicloud_log_project.default.0.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_oss_enabled":      "true",
						"delivery_oss_key_prefix":   "path1/",
						"delivery_oss_bucket_name":  CHECKSET,
						"delivery_sls_enabled":      "true",
						"delivery_sls_project_name": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudOOSServiceSettingMap0 = map[string]string{}

func AlicloudOOSServiceSettingBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

resource "alicloud_oss_bucket" "default" {
  count  = 2
  bucket = join("", [var.name, count.index])
  acl    = "public-read-write"
}

resource "alicloud_log_project" "default" {
  count = 2
  name  = join("", [var.name, count.index])
}
`, name)
}
