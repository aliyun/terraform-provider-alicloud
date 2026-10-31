package threatdetection_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/threatdetection"
	tFacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Case 1
func TestAccAliCloudThreatDetectionWebLockConfig_basic1875(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_web_lock_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudThreatDetectionWebLockConfigMap1875)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return threatdetection.NewSasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeThreatDetectionWebLockConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tFacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sThreatDetectionWebLockConfig%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudThreatDetectionWebLockConfigBasicDependence1875)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"uuid":                acctest.CHECKSET,
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
