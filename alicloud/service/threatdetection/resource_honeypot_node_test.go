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
func TestAccAliCloudThreatDetectionHoneypotNode_basic1988(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_honeypot_node.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudThreatDetectionHoneypotNodeMap1988)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return threatdetection.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeThreatDetectionHoneypotNode")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tFacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sThreatDetectionHoneypotNode%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudThreatDetectionHoneypotNodeBasicDependence1988)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"node_name":                      "${var.name}",
					"available_probe_num":            "20",
					"allow_honeypot_access_internet": "true",
					"security_group_probe_ip_list":   []string{"0.0.0.0/0"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_name":                      name,
						"available_probe_num":            "20",
						"allow_honeypot_access_internet": "true",
						"security_group_probe_ip_list.#": "1",
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"available_probe_num":          "24",
					"security_group_probe_ip_list": []string{"0.0.0.0/0", "10.0.0.0/8"},
					"node_name":                    "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"available_probe_num":            "24",
						"security_group_probe_ip_list.#": "2",
						"node_name":                      name + "_update",
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

var AlicloudThreatDetectionHoneypotNodeMap1988 = map[string]string{}

func AlicloudThreatDetectionHoneypotNodeBasicDependence1988(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}
