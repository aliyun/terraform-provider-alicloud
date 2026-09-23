package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAlicloudThreatDetectionHoneyPot_basic1994(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_honey_pot.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionHoneyPotMap1994)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionHoneyPot")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sHoneyPot%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionHoneyPotBasicDependence1994)
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
					"honeypot_image_name": "ruoyi",
					"honeypot_image_id":   "${data.alicloud_threat_detection_honeypot_images.default.images.0.honeypot_image_id}",
					"honeypot_name":       "${var.name}",
					"node_id":             "${alicloud_threat_detection_honeypot_node.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"honeypot_image_name": CHECKSET,
						"honeypot_image_id":   CHECKSET,
						"honeypot_name":       CHECKSET,
						"node_id":             CHECKSET,
						"status":              CHECKSET,
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"honeypot_name": "${var.name}-u",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"honeypot_name": CHECKSET,
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

var AlicloudThreatDetectionHoneyPotMap1994 = map[string]string{}

func AlicloudThreatDetectionHoneyPotBasicDependence1994(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_threat_detection_honeypot_images" "default" {
  name_regex = "^ruoyi"
}

resource "alicloud_threat_detection_honeypot_node" "default" {
  node_name           = var.name
  available_probe_num = 20
  security_group_probe_ip_list = ["0.0.0.0/0"]
}

`, name)
}
