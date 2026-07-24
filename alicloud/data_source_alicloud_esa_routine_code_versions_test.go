package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudEsaRoutineCodeVersionsDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEsaRoutineCodeVersionsSourceConfig(rand, map[string]string{
			"routine_name": `"${alicloud_esa_routine.default.name}"`,
			"ids":          `["${alicloud_esa_routine.default.latest_code_version}"]`,
		}),
		fakeConfig: testAccCheckAliCloudEsaRoutineCodeVersionsSourceConfig(rand, map[string]string{
			"routine_name": `"${alicloud_esa_routine.default.name}"`,
			"ids":          `["${alicloud_esa_routine.default.latest_code_version}_fake"]`,
		}),
	}

	EsaRoutineCodeVersionsCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}

var existEsaRoutineCodeVersionsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"versions.#":              "1",
		"versions.0.code_version": CHECKSET,
		"versions.0.status":       CHECKSET,
		"versions.0.create_time":  CHECKSET,
	}
}

var fakeEsaRoutineCodeVersionsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"versions.#": "0",
	}
}

var EsaRoutineCodeVersionsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_esa_routine_code_versions.default",
	existMapFunc: existEsaRoutineCodeVersionsMapFunc,
	fakeMapFunc:  fakeEsaRoutineCodeVersionsMapFunc,
}

func testAccCheckAliCloudEsaRoutineCodeVersionsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tftestacc%d"
}

resource "alicloud_esa_routine" "default" {
  name             = var.name
  description      = "tf-test-routine"
  code             = "addEventListener('fetch', e => e.respondWith(new Response('hello')))"
  code_description = "version 1"
}

data "alicloud_esa_routine_code_versions" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
