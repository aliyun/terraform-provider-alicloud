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

// NOTE: Test depends on data source or hardcoded are not stable and may fail at any time
// Case 1
func TestAccAliCloudThreatDetectionHoneypotProbe_basic2046(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_honeypot_probe.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudThreatDetectionHoneypotProbeMap2046)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return threatdetection.NewSasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeThreatDetectionHoneypotProbe")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tFacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sThreatDetectionHoneypotProbe%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudThreatDetectionHoneypotProbeBasicDependence2046)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"uuid":                 acctest.CHECKSET,
						"probe_type":           acctest.CHECKSET,
						"control_node_id":      acctest.CHECKSET,
						"ping":                 "true",
						"honeypot_bind_list.#": "1",
						"display_name":         acctest.CHECKSET,
						"arp":                  "true",
						"honeypot_probe_id":    acctest.CHECKSET,
						"probe_version":        acctest.CHECKSET,
						"service_ip_list.#":    "1",
						"status":               acctest.CHECKSET,
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
						"arp":          acctest.CHECKSET,
						"ping":         acctest.CHECKSET,
						"display_name": acctest.CHECKSET,
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
  //ids = ["038dcf91-3168-436c-a8cd-e5d36d80d872"]
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
