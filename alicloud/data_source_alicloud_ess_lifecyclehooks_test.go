package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEssLifecycleHooksDataSource(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssLifecycleHooksDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_lifecycle_hook.default.scaling_group_id}"`,
			"name_regex":       `"${alicloud_ess_lifecycle_hook.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEssLifecycleHooksDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_lifecycle_hook.default.scaling_group_id}"`,
			"name_regex":       `"${alicloud_ess_lifecycle_hook.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssLifecycleHooksDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_lifecycle_hook.default.scaling_group_id}"`,
			"ids":              `["${alicloud_ess_lifecycle_hook.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEssLifecycleHooksDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_lifecycle_hook.default.scaling_group_id}"`,
			"ids":              `["${alicloud_ess_lifecycle_hook.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssLifecycleHooksDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_lifecycle_hook.default.scaling_group_id}"`,
			"ids":              `["${alicloud_ess_lifecycle_hook.default.id}"]`,
			"name_regex":       `"${alicloud_ess_lifecycle_hook.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEssLifecycleHooksDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_lifecycle_hook.default.scaling_group_id}"`,
			"ids":              `["${alicloud_ess_lifecycle_hook.default.id}_fake"]`,
			"name_regex":       `"${alicloud_ess_lifecycle_hook.default.name}"`,
		}),
	}

	var existEsslifecyclehooksMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"hooks.#":                       "1",
			"hooks.0.name":                  fmt.Sprintf("tf-testAccDataSourceEssLifecycleHooks-%d", rand),
			"hooks.0.scaling_group_id":      CHECKSET,
			"hooks.0.default_result":        CHECKSET,
			"hooks.0.heartbeat_timeout":     "400",
			"hooks.0.lifecycle_transition":  "SCALE_OUT",
			"hooks.0.notification_arn":      CHECKSET,
			"hooks.0.notification_metadata": "helloworld",
		}
	}

	var fakeEsslifecyclehooksMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"hooks.#": "0",
			"ids.#":   "0",
		}
	}

	var essLifecyclehooksCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ess_lifecycle_hooks.default",
		existMapFunc: existEsslifecyclehooksMapFunc,
		fakeMapFunc:  fakeEsslifecyclehooksMapFunc,
	}

	essLifecyclehooksCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlicloudEssLifecycleHooksDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
%s

variable "name" {
	default = "tf-testAccDataSourceEssLifecycleHooks-%d"
}

resource "alicloud_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	scaling_group_name = "${var.name}"
	vswitch_ids = ["${alicloud_vswitch.default.id}"]
}
resource "alicloud_ess_lifecycle_hook" "default" {
  scaling_group_id      = "${alicloud_ess_scaling_group.default.id}"
  name                  = "${var.name}"
  lifecycle_transition  = "SCALE_OUT"
  heartbeat_timeout     = 400
  notification_metadata = "helloworld"
}

data "alicloud_ess_lifecycle_hooks" "default"{
  %s
}
`, EcsInstanceCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}
