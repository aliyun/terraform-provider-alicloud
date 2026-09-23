package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAlicloudThreatDetectionHoneypotPreset_basic1974(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_honeypot_preset.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionHoneypotPresetMap1974)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionHoneypotPreset")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sThreatDetectionHoneypotPreset%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionHoneypotPresetBasicDependence1974)
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
					"honeypot_image_name": "shiro",
					"meta": []map[string]interface{}{
						{
							"portrait_option": "true",
							"burp":            "open",
							"trojan_git":      "open",
						},
					},
					"node_id":     "${alicloud_threat_detection_honeypot_node.default.id}",
					"preset_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"honeypot_image_name":    "shiro",
						"meta.#":                 "1",
						"meta.0.portrait_option": "true",
						"meta.0.burp":            "open",
						"meta.0.trojan_git":      "open",
						"node_id":                CHECKSET,
						"preset_name":            name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preset_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preset_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"honeypot_image_name": "shiro",
					"preset_name":         "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"honeypot_image_name": "shiro",
						"preset_name":         name,
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

var AlicloudThreatDetectionHoneypotPresetMap1974 = map[string]string{}

func AlicloudThreatDetectionHoneypotPresetBasicDependence1974(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_threat_detection_honeypot_node" "default" {
  node_name           = var.name
  available_probe_num = 20
  security_group_probe_ip_list = ["0.0.0.0/0"]
}
`, name)
}
