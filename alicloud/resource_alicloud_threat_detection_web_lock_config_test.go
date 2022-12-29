package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAlicloudThreatDetectionWebLockConfig_basic1875(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_web_lock_config.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionWebLockConfigMap1875)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionWebLockConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sThreatDetectionWebLockConfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionWebLockConfigBasicDependence1875)
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
					"inclusive_file_type": "php;jsp;asp;aspx;js;cgi;html;htm;xml;shtml;shtm;jpg",
					"uuid":                "${data.alicloud_threat_detection_assets.default.ids.0}",
					"mode":                "whitelist",
					"local_backup_dir":    "/usr/local/aegis/bak",
					"dir":                 "/tmp/",
					"defence_mode":        "audit",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"inclusive_file_type": "php;jsp;asp;aspx;js;cgi;html;htm;xml;shtml;shtm;jpg",
						"uuid":                CHECKSET,
						"mode":                "whitelist",
						"local_backup_dir":    "/usr/local/aegis/bak",
						"dir":                 "/tmp/",
						"defence_mode":        "audit",
					}),
				),
			}, {
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudThreatDetectionWebLockConfigMap1875 = map[string]string{}

func AlicloudThreatDetectionWebLockConfigBasicDependence1875(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_threat_detection_assets" "default" {
  machine_types = "ecs"
}

`, name)
}
