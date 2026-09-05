package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// lintignore: AT001
func TestAccAliCloudThreatDetectionGlobalConfig_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_global_config.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionGlobalConfigMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionGlobalConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-tdgc-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionGlobalConfigBasicDependence)
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
					"global_config_name":  "${var.name}",
					"global_config_value": "config-value-1",
					"lang":                "zh",
					"role_for":            0,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_config_name":  name,
						"global_config_value": "config-value-1",
						"lang":                "zh",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_config_name":  "${var.name}",
					"global_config_value": "config-value-2",
					"lang":                "zh",
					"role_for":            0,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_config_name":  name,
						"global_config_value": "config-value-2",
						"lang":                "zh",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_config_name":  "${var.name}",
					"global_config_value": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_config_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang", "role_for"},
			},
		},
	})
}

var AlicloudThreatDetectionGlobalConfigMap = map[string]string{}

func AlicloudThreatDetectionGlobalConfigBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

`, name)
}
