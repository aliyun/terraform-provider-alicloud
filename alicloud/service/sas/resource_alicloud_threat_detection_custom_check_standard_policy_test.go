package sas_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudThreatDetectionCustomCheckStandardPolicy_basic covers a
// SECTION-level policy (no parent dependency): create, update the show name,
// and import.
func TestAccAliCloudThreatDetectionCustomCheckStandardPolicy_basic(t *testing.T) {
	rand := sdkacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-td-ccsp-sec-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccThreatDetectionCustomCheckStandardPolicyBasicConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_standard_policy.default", "policy_type", "SECTION"),
					resource.TestCheckResourceAttrSet("alicloud_threat_detection_custom_check_standard_policy.default", "policy_id"),
					resource.TestCheckResourceAttrSet("alicloud_threat_detection_custom_check_standard_policy.default", "id"),
				),
			},
			{
				Config: testAccThreatDetectionCustomCheckStandardPolicyBasicConfig(name + "-update"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_standard_policy.default", "policy_type", "SECTION"),
					resource.TestCheckResourceAttrSet("alicloud_threat_detection_custom_check_standard_policy.default", "policy_id"),
				),
			},
			{
				ResourceName:      "alicloud_threat_detection_custom_check_standard_policy.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccAliCloudThreatDetectionCustomCheckStandardPolicy_standard covers a
// STANDARD-level policy built on a SECTION -> REQUIREMENT -> STANDARD chain,
// exercising dependent_policy_id and type, then updating type.
func TestAccAliCloudThreatDetectionCustomCheckStandardPolicy_standard(t *testing.T) {
	rand := sdkacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-td-ccsp-std-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccThreatDetectionCustomCheckStandardPolicyStandardConfig(name, "AISPM"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_standard_policy.default", "policy_type", "STANDARD"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_standard_policy.default", "type", "AISPM"),
					resource.TestCheckResourceAttrSet("alicloud_threat_detection_custom_check_standard_policy.default", "dependent_policy_id"),
				),
			},
			{
				Config: testAccThreatDetectionCustomCheckStandardPolicyStandardConfig(name, "RISK"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_standard_policy.default", "policy_type", "STANDARD"),
					resource.TestCheckResourceAttr("alicloud_threat_detection_custom_check_standard_policy.default", "type", "RISK"),
					resource.TestCheckResourceAttrSet("alicloud_threat_detection_custom_check_standard_policy.default", "dependent_policy_id"),
				),
			},
		},
	})
}

func testAccThreatDetectionCustomCheckStandardPolicyBasicConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_threat_detection_custom_check_standard_policy" "default" {
  policy_show_name = var.name
  policy_type      = "SECTION"
}
`, name)
}

func testAccThreatDetectionCustomCheckStandardPolicyStandardConfig(name, policyType string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_threat_detection_custom_check_standard_policy" "section" {
  policy_show_name = "${var.name}-section"
  policy_type      = "SECTION"
}

resource "alicloud_threat_detection_custom_check_standard_policy" "requirement" {
  policy_show_name    = "${var.name}-requirement"
  policy_type         = "REQUIREMENT"
  dependent_policy_id = alicloud_threat_detection_custom_check_standard_policy.section.policy_id
}

resource "alicloud_threat_detection_custom_check_standard_policy" "default" {
  policy_show_name    = "${var.name}-standard"
  policy_type         = "STANDARD"
  dependent_policy_id = alicloud_threat_detection_custom_check_standard_policy.requirement.policy_id
  type                = "%s"
}
`, name, policyType)
}
