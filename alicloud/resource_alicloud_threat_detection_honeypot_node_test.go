package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAlicloudThreatDetectionHoneypotNode_basic1988(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_honeypot_node.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionHoneypotNodeMap1988)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionHoneypotNode")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sThreatDetectionHoneypotNode%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionHoneypotNodeBasicDependence1988)
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
