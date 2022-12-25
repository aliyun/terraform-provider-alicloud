package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 2
func TestAccAlicloudThreatDetectionAntiBruteForceRule_basic1979(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_anti_brute_force_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionAntiBruteForceRuleMap1979)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionAntiBruteForceRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sThreatDetectionAntiBruteForceRule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionAntiBruteForceRuleBasicDependence1979)
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
					"default_rule":               "false",
					"anti_brute_force_rule_name": "${var.name}",
					"forbidden_time":             "360",
					"uuid_list": []string{
						"${data.alicloud_threat_detection_assets.default.assets.0.uuid}",
					},
					"fail_count": "80",
					"span":       "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_rule":               "false",
						"anti_brute_force_rule_name": CHECKSET,
						"forbidden_time":             "360",
						"uuid_list.#":                "1",
						"fail_count":                 "80",
						"span":                       "10",
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"anti_brute_force_rule_name": "${var.name}_update",
					"forbidden_time":             "300",
					"uuid_list": []string{
						"${data.alicloud_threat_detection_assets.default.assets.0.uuid}",
					},
					"fail_count":   "70",
					"span":         "9",
					"default_rule": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"anti_brute_force_rule_name": CHECKSET,
						"forbidden_time":             "300",
						"uuid_list.#":                "1",
						"fail_count":                 "70",
						"span":                       "9",
						"default_rule":               "true",
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

var AlicloudThreatDetectionAntiBruteForceRuleMap1979 = map[string]string{}

func AlicloudThreatDetectionAntiBruteForceRuleBasicDependence1979(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
data "alicloud_threat_detection_assets" "default" {
    machine_types = "ecs"
    ids = ["79d76eac-055a-492a-a5c8-eef3bac80c15"]
}

`, name)
}
