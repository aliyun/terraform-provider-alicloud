package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAlicloudThreatDetectionHoneypotProbe_basic2046(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_honeypot_probe.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionHoneypotProbeMap2046)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionHoneypotProbe")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sThreatDetectionHoneypotProbe%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionHoneypotProbeBasicDependence2046)
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
					"uuid":            "${data.alicloud_threat_detection_assets.default.assets.0.uuid}",
					"probe_type":      "host_probe",
					"control_node_id": "${alicloud_threat_detection_honeypot_node.default.id}",
					"ping":            "true",
					"honeypot_bind_list": []map[string]interface{}{
						{
							"bind_port_list": []map[string]interface{}{
								{
									"start_port": "80",
									"end_port":   "80",
								},
							},
							"honeypot_id": "${alicloud_threat_detection_honey_pot.default.id}",
						},
					},
					"display_name": "${var.name}",
					"arp":          "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"uuid":                 CHECKSET,
						"probe_type":           CHECKSET,
						"control_node_id":      CHECKSET,
						"ping":                 "true",
						"honeypot_bind_list.#": "1",
						"display_name":         CHECKSET,
						"arp":                  "true",
						"honeypot_probe_id":    CHECKSET,
						"probe_version":        CHECKSET,
						"service_ip_list.#":    "1",
						"status":               CHECKSET,
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"arp":          "false",
					"ping":         "false",
					"display_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"arp":          CHECKSET,
						"ping":         CHECKSET,
						"display_name": CHECKSET,
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

var AlicloudThreatDetectionHoneypotProbeMap2046 = map[string]string{}

func AlicloudThreatDetectionHoneypotProbeBasicDependence2046(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
data "alicloud_threat_detection_assets" "default" {
    machine_types = "ecs"
    ids = ["53926396-0690-49a4-aa08-a38ca4853cdf"]
}

resource "alicloud_threat_detection_honeypot_node" "default" {
  node_name           = var.name
  available_probe_num = 20
  security_group_probe_ip_list = ["0.0.0.0/0"]
}

data "alicloud_threat_detection_honeypot_images" "default" {
  name_regex = "^ruoyi"
}

resource "alicloud_threat_detection_honey_pot" "default" {
  honeypot_image_name = "ruoyi"
  honeypot_image_id   = data.alicloud_threat_detection_honeypot_images.default.images.0.honeypot_image_id
  honeypot_name       = var.name
  node_id             = alicloud_threat_detection_honeypot_node.default.id
}

`, name)
}
