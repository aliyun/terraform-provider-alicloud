// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlicloudCrArtifactLifecycleRuleDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCrArtifactLifecycleRuleSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_cr_artifact_lifecycle_rule.default.id}"]`,
			"instance_id": `"${alicloud_cr_ee_instance.resourceCase_20260526_Mmd6on_1.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCrArtifactLifecycleRuleSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_cr_artifact_lifecycle_rule.default.id}_fake"]`,
			"instance_id": `"${alicloud_cr_ee_instance.resourceCase_20260526_Mmd6on_1.id}"`,
		}),
	}

	InstanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCrArtifactLifecycleRuleSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_cr_artifact_lifecycle_rule.default.id}"]`,
			"instance_id": `"${alicloud_cr_ee_instance.resourceCase_20260526_Mmd6on_1.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCrArtifactLifecycleRuleSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_cr_artifact_lifecycle_rule.default.id}_fake"]`,
			"instance_id": `"${alicloud_cr_ee_instance.resourceCase_20260526_Mmd6on_1.id}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCrArtifactLifecycleRuleSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_cr_artifact_lifecycle_rule.default.id}"]`,
			"instance_id": `"${alicloud_cr_ee_instance.resourceCase_20260526_Mmd6on_1.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCrArtifactLifecycleRuleSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_cr_artifact_lifecycle_rule.default.id}_fake"]`,
			"instance_id": `"${alicloud_cr_ee_instance.resourceCase_20260526_Mmd6on_1.id}"`,
		}),
	}

	CrArtifactLifecycleRuleCheckInfo.dataSourceTestCheck(t, rand, idsConf, InstanceIdConf, allConf)
}

var existCrArtifactLifecycleRuleMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"rules.#":                            "1",
		"rules.0.id":                         CHECKSET,
		"rules.0.instance_id":                CHECKSET,
		"rules.0.create_time":                CHECKSET,
		"rules.0.repo_name":                  CHECKSET,
		"rules.0.auto":                       CHECKSET,
		"rules.0.namespace_name":             CHECKSET,
		"rules.0.retention_tag_count":        CHECKSET,
		"rules.0.schedule_time":              CHECKSET,
		"rules.0.modified_time":              CHECKSET,
		"rules.0.scope":                      CHECKSET,
		"rules.0.tag_regexp":                 CHECKSET,
		"rules.0.artifact_lifecycle_rule_id": CHECKSET,
	}
}

var fakeCrArtifactLifecycleRuleMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"rules.#": "0",
	}
}

var CrArtifactLifecycleRuleCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cr_artifact_lifecycle_rules.default",
	existMapFunc: existCrArtifactLifecycleRuleMapFunc,
	fakeMapFunc:  fakeCrArtifactLifecycleRuleMapFunc,
}

func testAccCheckAlicloudCrArtifactLifecycleRuleSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tfacc-cr-alr-%d"
}
resource "alicloud_cr_ee_instance" "resourceCase_20260526_Mmd6on_1" {
  default_oss_bucket = "true"
  instance_name      = var.name
  renewal_status     = "ManualRenewal"
  image_scanner      = "DISABLE"
  period             = "1"
  payment_type       = "Subscription"
  instance_type      = "Economy"
}

resource "alicloud_cr_ee_namespace" "namespaceCase_20260611_ArtifactLifecycleRule_1" {
  instance_id        = alicloud_cr_ee_instance.resourceCase_20260526_Mmd6on_1.id
  name               = var.name
  auto_create        = false
  default_visibility = "PRIVATE"
}

resource "alicloud_cr_ee_repo" "repoCase_20260611_ArtifactLifecycleRule_1" {
  instance_id = alicloud_cr_ee_instance.resourceCase_20260526_Mmd6on_1.id
  namespace   = alicloud_cr_ee_namespace.namespaceCase_20260611_ArtifactLifecycleRule_1.name
  name        = var.name
  repo_type   = "PRIVATE"
  summary     = "test repository for lifecycle rule"
}

resource "alicloud_cr_artifact_lifecycle_rule" "default" {
  auto                = true
  namespace_name      = alicloud_cr_ee_namespace.namespaceCase_20260611_ArtifactLifecycleRule_1.name
  retention_tag_count = "30"
  schedule_time       = "WEEK"
  scope               = "REPO"
  instance_id         = alicloud_cr_ee_instance.resourceCase_20260526_Mmd6on_1.id
  tag_regexp          = ".*"
  repo_name           = alicloud_cr_ee_repo.repoCase_20260611_ArtifactLifecycleRule_1.name
}

data "alicloud_cr_artifact_lifecycle_rules" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
