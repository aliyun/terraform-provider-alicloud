package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudThreatDetectionBackupPoliciesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionBackupPoliciesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_backup_policy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionBackupPoliciesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_backup_policy.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionBackupPoliciesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_threat_detection_backup_policy.default.backup_policy_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionBackupPoliciesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_threat_detection_backup_policy.default.backup_policy_name}_fake"`,
		}),
	}
	nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionBackupPoliciesDataSourceName(rand, map[string]string{
			"name": `"${alicloud_threat_detection_backup_policy.default.backup_policy_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionBackupPoliciesDataSourceName(rand, map[string]string{
			"name": `"${alicloud_threat_detection_backup_policy.default.backup_policy_name}_fake"`,
		}),
	}
	machineRemarkConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionBackupPoliciesDataSourceName(rand, map[string]string{
			"machine_remark": `"launch-advisor-20220810"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionBackupPoliciesDataSourceName(rand, map[string]string{
			"machine_remark": `"tf-test"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionBackupPoliciesDataSourceName(rand, map[string]string{
			"name":   `"${alicloud_threat_detection_backup_policy.default.backup_policy_name}"`,
			"status": `"enabled"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionBackupPoliciesDataSourceName(rand, map[string]string{
			"name":   `"${alicloud_threat_detection_backup_policy.default.backup_policy_name}_fake"`,
			"status": `"closed"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionBackupPoliciesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_threat_detection_backup_policy.default.id}"]`,
			"name_regex":     `"${alicloud_threat_detection_backup_policy.default.backup_policy_name}"`,
			"name":           `"${alicloud_threat_detection_backup_policy.default.backup_policy_name}"`,
			"machine_remark": `"launch-advisor-20220810"`,
			"status":         `"enabled"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionBackupPoliciesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_threat_detection_backup_policy.default.id}_fake"]`,
			"name_regex":     `"${alicloud_threat_detection_backup_policy.default.backup_policy_name}_fake"`,
			"name":           `"${alicloud_threat_detection_backup_policy.default.backup_policy_name}_fake"`,
			"machine_remark": `"tf-test"`,
			"status":         `"closed"`,
		}),
	}
	var existAlicloudThreatDetectionBackupPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"policies.#":                    "1",
			"policies.0.id":                 CHECKSET,
			"policies.0.backup_policy_id":   CHECKSET,
			"policies.0.backup_policy_name": CHECKSET,
			"policies.0.policy":             CHECKSET,
			"policies.0.policy_version":     "2.0.0",
			"policies.0.uuid_list.#":        "1",
			"policies.0.status":             "enabled",
		}
	}
	var fakeAlicloudThreatDetectionBackupPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"policies.#": "0",
		}
	}
	var alicloudThreatDetectionBackupPoliciesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_threat_detection_backup_policies.default",
		existMapFunc: existAlicloudThreatDetectionBackupPoliciesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudThreatDetectionBackupPoliciesDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudThreatDetectionBackupPoliciesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, nameConf, machineRemarkConf, statusConf, allConf)
}

func testAccCheckAlicloudThreatDetectionBackupPoliciesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccThreatDetectionBackupPolicy-%d"
	}

	data "alicloud_threat_detection_assets" "default" {
	  machine_types = "ecs"
	}

	resource "alicloud_threat_detection_backup_policy" "default" {
  		backup_policy_name = var.name
  		policy             = "{\"Exclude\":[\"/bin/\",\"/usr/bin/\",\"/sbin/\",\"/boot/\",\"/proc/\",\"/sys/\",\"/srv/\",\"/lib/\",\"/selinux/\",\"/usr/sbin/\",\"/run/\",\"/lib32/\",\"/lib64/\",\"/lost+found/\",\"/var/lib/kubelet/\",\"/var/lib/ntp/proc\",\"/var/lib/container\"],\"ExcludeSystemPath\":true,\"Include\":[],\"IsDefault\":1,\"Retention\":7,\"Schedule\":\"I|1668703620|PT24H\",\"Source\":[],\"SpeedLimiter\":\"\",\"UseVss\":true}"
  		policy_version     = "2.0.0"
  		uuid_list          = [data.alicloud_threat_detection_assets.default.ids.0]
	}

	data "alicloud_threat_detection_backup_policies" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
