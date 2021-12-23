package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSimpleApplicationServerPlanDataSource(t *testing.T) {
	rand := acctest.RandInt()

	allConf := dataSourceTestAccConfig{
		existConfig: dataSourcesimpleApplicationServerPlansConfigDependence(rand, map[string]string{}),
		fakeConfig:  "",
	}
	var existSimpleApplicationServerPlansMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   CHECKSET,
			"plans.#": CHECKSET,
		}
	}

	var fakeSimpleApplicationServerPlansMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"plans.#": "0",
		}
	}

	var simpleApplicationServerPlansCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_simple_application_server_plans.default",
		existMapFunc: existSimpleApplicationServerPlansMapFunc,
		fakeMapFunc:  fakeSimpleApplicationServerPlansMapFunc,
	}

	simpleApplicationServerPlansCheckInfo.dataSourceTestCheck(t, rand, allConf)
}

func dataSourcesimpleApplicationServerPlansConfigDependence(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
data "alicloud_simple_application_server_plans" "default" {
  memory    = 1
  bandwidth = 3
  disk_size = 40
  flow      = 6
  core      = 2
  %s
}

`, strings.Join(pairs, "\n   "))
	return config
}
