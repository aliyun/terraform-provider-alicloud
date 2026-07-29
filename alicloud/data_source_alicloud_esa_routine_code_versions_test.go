package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudEsaRoutineCodeVersionsDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_esa_routine_code_versions.default"
	name := fmt.Sprintf("tftestaccesacv%d", rand)
	codeFile := esaRoutineWriteCodeFixtureForDataSource(t, name)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, func(n string) string {
		return dataSourceEsaRoutineCodeVersionsConfig(n, codeFile)
	})

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name": "${alicloud_esa_routine.default.name}",
			"ids":  []string{"${alicloud_esa_routine.default.latest_code_version}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name": "${alicloud_esa_routine.default.name}",
			"ids":  []string{"${alicloud_esa_routine.default.latest_code_version}_fake"},
		}),
	}

	var existMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   CHECKSET,
			"versions.#":              CHECKSET,
			"versions.0.code_version": CHECKSET,
			"versions.0.status":       CHECKSET,
			"versions.0.create_time":  CHECKSET,
		}
	}

	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"versions.#": "0",
		}
	}

	var info = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	}

	info.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)
}

func esaRoutineWriteCodeFixtureForDataSource(t *testing.T, name string) string {
	return esaRoutineWriteCodeFixture(t, name, "addEventListener('fetch', e => e.respondWith(new Response('ds')))")
}

func dataSourceEsaRoutineCodeVersionsConfig(name, filename string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%[1]s"
}

resource "alicloud_esa_routine" "default" {
  name             = var.name
  description      = "tf-testacc esa routine code versions"
  filename         = "%[2]s"
  code_description = "code version ds"
}
`, name, filename)
}
