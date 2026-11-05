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
func TestAccAliCloudThreatDetectionHoneyPot_basic1994(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_honey_pot.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudThreatDetectionHoneyPotMap1994)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return threatdetection.NewSasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeThreatDetectionHoneyPot")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tFacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sHoneyPot%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudThreatDetectionHoneyPotBasicDependence1994)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"honeypot_image_name": acctest.CHECKSET,
						"honeypot_image_id":   acctest.CHECKSET,
						"honeypot_name":       acctest.CHECKSET,
						"node_id":             acctest.CHECKSET,
						"status":              acctest.CHECKSET,
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"honeypot_name": "${var.name}-u",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"honeypot_name": acctest.CHECKSET,
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
}`, name)
}
