package esa_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudEsaRoutineCodeVersionsDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(10000, 99999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudEsaRoutineCodeVersionsSourceConfig(rand, map[string]string{
			"routine_name": `"${alicloud_esa_routine.default.name}"`,
			"ids":          `["${alicloud_esa_routine.default.latest_code_version}"]`,
		}),
		FakeConfig: testAccCheckAliCloudEsaRoutineCodeVersionsSourceConfig(rand, map[string]string{
			"routine_name": `"${alicloud_esa_routine.default.name}"`,
			"ids":          `["${alicloud_esa_routine.default.latest_code_version}_fake"]`,
		}),
	}

	EsaRoutineCodeVersionsCheckInfo.DataSourceTestCheck(t, rand, idsConf)
}

var existEsaRoutineCodeVersionsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"versions.#":              "1",
		"versions.0.code_version": acctest.CHECKSET,
		"versions.0.status":       acctest.CHECKSET,
		"versions.0.create_time":  acctest.CHECKSET,
	}
}

var fakeEsaRoutineCodeVersionsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"versions.#": "0",
	}
}

var EsaRoutineCodeVersionsCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_esa_routine_code_versions.default",
	ExistMapFunc: existEsaRoutineCodeVersionsMapFunc,
	FakeMapFunc:  fakeEsaRoutineCodeVersionsMapFunc,
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
