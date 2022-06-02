package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECDPolicyGroupDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdPolicyGroupDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_policy_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcdPolicyGroupDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_policy_group.default.id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdPolicyGroupDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_policy_group.default.policy_group_name}"`,
			"status":     `"AVAILABLE"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdPolicyGroupDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_policy_group.default.policy_group_name}"`,
			"status":     `"CREATING"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdPolicyGroupDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_policy_group.default.policy_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdPolicyGroupDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_policy_group.default.policy_group_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdPolicyGroupDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecd_policy_group.default.id}"]`,
			"status":     `"AVAILABLE"`,
			"name_regex": `"${alicloud_ecd_policy_group.default.policy_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdPolicyGroupDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecd_policy_group.default.id}_fake"]`,
			"status":     `"CREATING"`,
			"name_regex": `"${alicloud_ecd_policy_group.default.policy_group_name}_fake"`,
		}),
	}
	var existAlicloudEcdPolicyGroupDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"names.#":                 "1",
			"groups.#":                "1",
			"groups.0.visual_quality": "medium",
			"groups.0.clipboard":      "readwrite",
			"groups.0.local_drive":    "read",
			"groups.0.authorize_access_policy_rules.#":   "1",
			"groups.0.authorize_security_policy_rules.#": "1",
			"groups.0.recording":                         "off",
			"groups.0.recording_start_time":              "",
			"groups.0.recording_end_time":                "",
			"groups.0.recording_fps":                     "0",
			"groups.0.camera_redirect":                   "on",
		}
	}
	var fakeAlicloudEcdPolicyGroupDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudEcdPolicyGroupCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecd_policy_groups.default",
		existMapFunc: existAlicloudEcdPolicyGroupDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcdPolicyGroupDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.EcdSupportRegions)
	}
	alicloudEcdPolicyGroupCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudEcdPolicyGroupDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccPolicyGroup-%d"
}

resource "alicloud_ecd_policy_group" "default" {
  policy_group_name = var.name
  clipboard = "readwrite"
  local_drive = "read"
  authorize_access_policy_rules{
    description= var.name
    cidr_ip=     "1.2.3.4/24"
  }
  authorize_security_policy_rules  {
    type=        "inflow"
    policy=      "accept"
    description=  var.name
    port_range= "80/80"
    ip_protocol= "TCP"
    priority=    "1"
    cidr_ip=     "0.0.0.0/0"
  }
}

data "alicloud_ecd_policy_groups" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
