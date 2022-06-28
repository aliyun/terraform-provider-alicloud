package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudLogStoresDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudLogStoresDataSourceName(rand, map[string]string{
			"ids": `[alicloud_log_store.default.name]`,
		}),
		fakeConfig: testAccCheckAlicloudLogStoresDataSourceName(rand, map[string]string{
			"ids": `["fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudLogStoresDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_log_store.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudLogStoresDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_log_store.default.name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudLogStoresDataSourceName(rand, map[string]string{
			"ids":        `[alicloud_log_store.default.name]`,
			"name_regex": `"${alicloud_log_store.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudLogStoresDataSourceName(rand, map[string]string{
			"ids":        `["fake"]`,
			"name_regex": `"${alicloud_log_store.default.name}_fake"`,
		}),
	}
	var existAlicloudLogStoresDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"names.#":             "1",
			"stores.#":            "1",
			"stores.0.id":         CHECKSET,
			"stores.0.store_name": fmt.Sprintf("tf-testacclogstores-%d", rand),
		}
	}
	var fakeAlicloudLogStoresDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"stores.#": "0",
		}
	}
	var alicloudLogStoresCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_log_stores.default",
		existMapFunc: existAlicloudLogStoresDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudLogStoresDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudLogStoresCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudLogStoresDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testacclogstores-%d"
}

resource "alicloud_log_project" "default" {
  name        = var.name
  description = "created by terraform"
}

resource "alicloud_log_store" "default" {
  project = alicloud_log_project.default.name
  name    = var.name
}

data "alicloud_log_stores" "default" {
	project= alicloud_log_store.default.project
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
