package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudEsaListsDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_esa_lists.default"
	name := fmt.Sprintf("tf_testacc_esa_list_%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceEsaListsConfig)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_esa_list.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_esa_list.default.id}_fake"},
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_esa_list.default.name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_esa_list.default.name}_fake",
		}),
	}

	nameSubstring := `${replace(alicloud_esa_list.default.name, "tf_testacc_", "")}`
	itemSubstring := "${substr(alicloud_esa_list.default.items.0, 0, 6)}"
	nameQueryArgsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"query_args": []map[string]interface{}{
				{
					"name_like": nameSubstring,
					"order_by":  "name",
					"desc":      "true",
				},
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"query_args": []map[string]interface{}{
				{
					"name_like": nameSubstring + "_fake",
					"order_by":  "name",
					"desc":      "false",
				},
			},
		}),
	}

	kindQueryArgsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_esa_list.default.id}"},
			"query_args": []map[string]interface{}{
				{
					"kind": "ip",
				},
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_esa_list.default.id}"},
			"query_args": []map[string]interface{}{
				{
					"kind": "header",
				},
			},
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_esa_list.default.id}"},
			"name_regex": "${alicloud_esa_list.default.name}",
			"query_args": []map[string]interface{}{
				{
					"id_like":   "${alicloud_esa_list.default.id}",
					"name_like": "${alicloud_esa_list.default.name}",
					"kind":      "ip",
					"order_by":  "id",
					"desc":      "true",
				},
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_esa_list.default.id}"},
			"name_regex": "${alicloud_esa_list.default.name}",
			"query_args": []map[string]interface{}{
				{
					"id_like":   "${alicloud_esa_list.default.id}",
					"name_like": "${alicloud_esa_list.default.name}",
					"kind":      "header",
					"order_by":  "id",
					"desc":      "false",
				},
			},
		}),
	}

	configs := []acctest.DataSourceTestAccConfig{idsConf, nameRegexConf, nameQueryArgsConf, kindQueryArgsConf, allConf}
	for _, tc := range []struct {
		field string
		match string
		miss  string
	}{
		{"name_like", nameSubstring, nameSubstring + "_fake"},
		{"description_like", "${alicloud_esa_list.default.name}", "${alicloud_esa_list.default.name}_fake"},
		{"name_item_like", nameSubstring, nameSubstring + "_fake"},
		{"name_item_like", itemSubstring, "192.0.2."},
		{"id_like", "${substr(alicloud_esa_list.default.id, 1, -1)}", "${alicloud_esa_list.default.id}0"},
	} {
		// Isolate this fixture without masking the server-side filter with local filters.
		match := map[string]interface{}{"name_like": "${alicloud_esa_list.default.name}"}
		miss := map[string]interface{}{"name_like": "${alicloud_esa_list.default.name}"}
		match[tc.field], miss[tc.field] = tc.match, tc.miss
		configs = append(configs, acctest.DataSourceTestAccConfig{
			ExistConfig: testAccConfig(map[string]interface{}{"query_args": []map[string]interface{}{match}}),
			FakeConfig:  testAccConfig(map[string]interface{}{"query_args": []map[string]interface{}{miss}}),
		})
	}

	var existAliCloudEsaListsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"names.#":             "1",
			"lists.#":             "1",
			"lists.0.id":          acctest.CHECKSET,
			"lists.0.list_id":     acctest.CHECKSET,
			"lists.0.name":        name,
			"lists.0.kind":        "ip",
			"lists.0.description": "description_" + name + "_suffix",
			"lists.0.length":      "2",
			"lists.0.update_time": acctest.CHECKSET,
		}
	}

	var fakeAliCloudEsaListsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"lists.#": "0",
		}
	}

	var aliCloudEsaListsInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_esa_lists.default",
		ExistMapFunc: existAliCloudEsaListsMapFunc,
		FakeMapFunc:  fakeAliCloudEsaListsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	}

	aliCloudEsaListsInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, configs...)
}

func dataSourceEsaListsConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_esa_list" "default" {
  kind        = "ip"
  name        = var.name
  description = "description_${var.name}_suffix"
  items       = ["10.1.1.1", "10.1.1.2"]
}
`, name)
}
