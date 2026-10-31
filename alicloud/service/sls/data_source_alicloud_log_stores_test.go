package sls_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudLogStoresDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudLogStoresDataSourceName(rand, map[string]string{
			"ids": `[alicloud_log_store.default.name]`,
		}),
		FakeConfig: testAccCheckAlicloudLogStoresDataSourceName(rand, map[string]string{
			"ids": `["fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudLogStoresDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_log_store.default.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudLogStoresDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_log_store.default.name}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudLogStoresDataSourceName(rand, map[string]string{
			"ids":        `[alicloud_log_store.default.name]`,
			"name_regex": `"${alicloud_log_store.default.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudLogStoresDataSourceName(rand, map[string]string{
			"ids":        `["fake"]`,
			"name_regex": `"${alicloud_log_store.default.name}_fake"`,
		}),
	}
	var existAlicloudLogStoresDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"names.#":             "1",
			"stores.#":            "1",
			"stores.0.id":         acctest.CHECKSET,
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
	var alicloudLogStoresCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_log_stores.default",
		ExistMapFunc: existAlicloudLogStoresDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudLogStoresDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudLogStoresCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
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
