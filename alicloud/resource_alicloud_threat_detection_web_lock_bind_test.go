package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudThreatDetectionWebLockBind_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_web_lock_bind.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionWebLockBindMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionWebLockBind")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sThreatDetectionWebLockBind%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionWebLockBindBasicDependence)
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
					"status":              "on",
					"lang":                "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"inclusive_file_type": "php;jsp;asp;aspx;js;cgi;html;htm;xml;shtml;shtm;jpg",
						"uuid":                CHECKSET,
						"mode":                "whitelist",
						"local_backup_dir":    "/usr/local/aegis/bak",
						"dir":                 "/tmp/",
						"defence_mode":        "audit",
						"status":              "on",
						"lang":                "zh",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"inclusive_file_type": "php;jsp;asp;aspx;js;cgi;html;htm;xml;shtml;shtm;jpg",
					"uuid":                "${data.alicloud_threat_detection_assets.default.ids.0}",
					"mode":                "whitelist",
					"local_backup_dir":    "/usr/local/aegis/bak",
					"dir":                 "/tmp/",
					"defence_mode":        "audit",
					"status":              "off",
					"lang":                "en",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"inclusive_file_type": "php;jsp;asp;aspx;js;cgi;html;htm;xml;shtml;shtm;jpg",
						"uuid":                CHECKSET,
						"mode":                "whitelist",
						"local_backup_dir":    "/usr/local/aegis/bak",
						"dir":                 "/tmp/",
						"defence_mode":        "audit",
						"status":              "off",
						"lang":                "en",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"uuid":                "${data.alicloud_threat_detection_assets.default.ids.0}",
					"dir":                 "/tmp/",
					"local_backup_dir":    "/usr/local/aegis/bak",
					"defence_mode":        "audit",
					"mode":                "blacklist",
					"exclusive_dir":       "/tmp/logs/",
					"exclusive_file_type": "log;cache;tmp",
					"exclusive_file":      "/tmp/protected.db",
					"status":              "on",
					"lang":                "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"uuid":                CHECKSET,
						"dir":                 "/tmp/",
						"local_backup_dir":    "/usr/local/aegis/bak",
						"defence_mode":        "audit",
						"mode":                "blacklist",
						"exclusive_dir":       "/tmp/logs/",
						"exclusive_file_type": "log;cache;tmp",
						"exclusive_file":      "/tmp/protected.db",
						"status":              "on",
						"lang":                "zh",
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

var AlicloudThreatDetectionWebLockBindMap = map[string]string{}

func AlicloudThreatDetectionWebLockBindBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_threat_detection_assets" "default" {
  machine_types = "ecs"
}

`, name)
}
