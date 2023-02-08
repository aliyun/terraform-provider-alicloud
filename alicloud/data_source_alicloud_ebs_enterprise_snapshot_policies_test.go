package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEbsEnterpriseSnapshotPolicyDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEbsEnterpriseSnapshotPolicySourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ebs_enterprise_snapshot_policy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEbsEnterpriseSnapshotPolicySourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ebs_enterprise_snapshot_policy.default.id}_fake"]`,
		}),
	}

	enterpriseSnapshotPolicyIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEbsEnterpriseSnapshotPolicySourceConfig(rand, map[string]string{
			"enterprise_snapshot_policy_ids": `["${alicloud_ebs_enterprise_snapshot_policy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEbsEnterpriseSnapshotPolicySourceConfig(rand, map[string]string{
			"enterprise_snapshot_policy_ids": `["${alicloud_ebs_enterprise_snapshot_policy.default.id}_fake"]`,
		}),
	}

	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEbsEnterpriseSnapshotPolicySourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_ebs_enterprise_snapshot_policy.default.id}"]`,
			"resource_group_id": `"${alicloud_ebs_enterprise_snapshot_policy.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEbsEnterpriseSnapshotPolicySourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_ebs_enterprise_snapshot_policy.default.id}"]`,
			"resource_group_id": `"${alicloud_ebs_enterprise_snapshot_policy.default.resource_group_id}_fake"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEbsEnterpriseSnapshotPolicySourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ebs_enterprise_snapshot_policy.default.enterprise_snapshot_policy_name}"`,
			"ids":        `["${alicloud_ebs_enterprise_snapshot_policy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEbsEnterpriseSnapshotPolicySourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ebs_enterprise_snapshot_policy.default.id}"]`,
			"name_regex": `"${alicloud_ebs_enterprise_snapshot_policy.default.enterprise_snapshot_policy_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEbsEnterpriseSnapshotPolicySourceConfig(rand, map[string]string{
			"ids":                            `["${alicloud_ebs_enterprise_snapshot_policy.default.id}"]`,
			"resource_group_id":              `"${alicloud_ebs_enterprise_snapshot_policy.default.resource_group_id}"`,
			"enterprise_snapshot_policy_ids": `["${alicloud_ebs_enterprise_snapshot_policy.default.id}"]`,
			"name_regex":                     `"${alicloud_ebs_enterprise_snapshot_policy.default.enterprise_snapshot_policy_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEbsEnterpriseSnapshotPolicySourceConfig(rand, map[string]string{
			"ids":                            `["${alicloud_ebs_enterprise_snapshot_policy.default.id}_fake"]`,
			"enterprise_snapshot_policy_ids": `["${alicloud_ebs_enterprise_snapshot_policy.default.id}_fake"]`,
			"resource_group_id":              `"${alicloud_ebs_enterprise_snapshot_policy.default.resource_group_id}_fake"`,
			"name_regex":                     `"${alicloud_ebs_enterprise_snapshot_policy.default.enterprise_snapshot_policy_name}_fake"`,
		}),
	}

	EbsEnterpriseSnapshotPolicyCheckInfo.dataSourceTestCheck(t, rand, idsConf, enterpriseSnapshotPolicyIdsConf, resourceGroupIdConf, nameRegexConf, allConf)
}

var existEbsEnterpriseSnapshotPolicyMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                                                         "1",
		"enterprise_snapshot_policies.#":                                                "1",
		"enterprise_snapshot_policies.0.id":                                             CHECKSET,
		"enterprise_snapshot_policies.0.create_time":                                    CHECKSET,
		"enterprise_snapshot_policies.0.cross_region_copy_info.#":                       "1",
		"enterprise_snapshot_policies.0.cross_region_copy_info.0.enabled":               "true",
		"enterprise_snapshot_policies.0.cross_region_copy_info.0.regions.#":             "1",
		"enterprise_snapshot_policies.0.cross_region_copy_info.0.regions.0.region_id":   defaultRegionToTest,
		"enterprise_snapshot_policies.0.cross_region_copy_info.0.regions.0.retain_days": "7",
		"enterprise_snapshot_policies.0.desc":                                           CHECKSET,
		"enterprise_snapshot_policies.0.enterprise_snapshot_policy_id":                  CHECKSET,
		"enterprise_snapshot_policies.0.enterprise_snapshot_policy_name":                CHECKSET,
		"enterprise_snapshot_policies.0.resource_group_id":                              CHECKSET,
		"enterprise_snapshot_policies.0.retain_rule.#":                                  "1",
		"enterprise_snapshot_policies.0.retain_rule.0.time_interval":                    "1",
		"enterprise_snapshot_policies.0.retain_rule.0.time_unit":                        "DAYS",
		"enterprise_snapshot_policies.0.schedule.#":                                     "1",
		"enterprise_snapshot_policies.0.schedule.0.cron_expression":                     "0 0 */12 * * *",
		"enterprise_snapshot_policies.0.status":                                         CHECKSET,
		"enterprise_snapshot_policies.0.target_type":                                    CHECKSET,
		"enterprise_snapshot_policies.0.tags.%":                                         "2",
		"enterprise_snapshot_policies.0.tags.Created":                                   "TF",
		"enterprise_snapshot_policies.0.tags.For":                                       "acceptance test",
		"enterprise_snapshot_policies.0.storage_rule.#":                                 "1",
		"enterprise_snapshot_policies.0.storage_rule.0.enable_immediate_access":         "false",
	}
}

var fakeEbsEnterpriseSnapshotPolicyMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                          "0",
		"enterprise_snapshot_policies.#": "0",
	}
}

var EbsEnterpriseSnapshotPolicyCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ebs_enterprise_snapshot_policies.default",
	existMapFunc: existEbsEnterpriseSnapshotPolicyMapFunc,
	fakeMapFunc:  fakeEbsEnterpriseSnapshotPolicyMapFunc,
}

func testAccCheckAlicloudEbsEnterpriseSnapshotPolicySourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccEbsEnterpriseSnapshotPolicy%d"
}


resource "alicloud_ebs_enterprise_snapshot_policy" "default" {
  status = "ENABLED"
  desc   = var.name
  schedule {
    cron_expression = "0 0 */12 * * *"
  }
  target_type = "DISK"
  retain_rule {
    time_interval = 1
    time_unit     = "DAYS"
  }
  cross_region_copy_info {
	enabled = true
	regions {
	  region_id = "%s"
	  retain_days = 7
	}
  }
	tags = {
		Created = "TF"
		For     = "acceptance test"	
	}
  enterprise_snapshot_policy_name = var.name
}

data "alicloud_ebs_enterprise_snapshot_policies" "default" {
%s
}
`, rand, os.Getenv("ALICLOUD_REGION"), strings.Join(pairs, "\n   "))
	return config
}
