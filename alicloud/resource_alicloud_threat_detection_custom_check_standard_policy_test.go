package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1: SECTION-level policy (no parent dependency)
func TestAccAliCloudThreatDetectionCustomCheckStandardPolicy_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_custom_check_standard_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionCustomCheckStandardPolicyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionCustomCheckStandardPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-td-ccsp-sec-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionCustomCheckStandardPolicyBasicDependence)
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
					"policy_show_name": "${var.name}",
					"policy_type":      "SECTION",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_show_name": CHECKSET,
						"policy_type":      "SECTION",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_show_name": "${var.name}-update",
					"policy_type":      "SECTION",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_show_name": CHECKSET,
						"policy_type":      "SECTION",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Case 2: STANDARD-level policy with dependent_policy_id and type
func TestAccAliCloudThreatDetectionCustomCheckStandardPolicy_standard(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_custom_check_standard_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionCustomCheckStandardPolicyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionCustomCheckStandardPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-td-ccsp-std-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionCustomCheckStandardPolicyStandardDependence)
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
					"policy_show_name":    "${var.name}",
					"policy_type":         "STANDARD",
					"dependent_policy_id": "${alicloud_threat_detection_custom_check_standard_policy.requirement.policy_id}",
					"type":                "AISPM",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_show_name":    CHECKSET,
						"policy_type":         "STANDARD",
						"dependent_policy_id": CHECKSET,
						"type":                "AISPM",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_show_name":    "${var.name}-update",
					"policy_type":         "STANDARD",
					"dependent_policy_id": "${alicloud_threat_detection_custom_check_standard_policy.requirement.policy_id}",
					"type":                "RISK",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_show_name":    CHECKSET,
						"policy_type":         "STANDARD",
						"dependent_policy_id": CHECKSET,
						"type":                "RISK",
					}),
				),
			},
		},
	})
}

var AlicloudThreatDetectionCustomCheckStandardPolicyMap = map[string]string{}

func AlicloudThreatDetectionCustomCheckStandardPolicyBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}`, name)
}

func AlicloudThreatDetectionCustomCheckStandardPolicyStandardDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_threat_detection_custom_check_standard_policy" "section" {
  policy_show_name = "%s-section"
  policy_type      = "SECTION"
}

resource "alicloud_threat_detection_custom_check_standard_policy" "requirement" {
  policy_show_name    = "%s-requirement"
  policy_type         = "REQUIREMENT"
  dependent_policy_id = "${alicloud_threat_detection_custom_check_standard_policy.section.policy_id}"
}`, name, name, name)
}
