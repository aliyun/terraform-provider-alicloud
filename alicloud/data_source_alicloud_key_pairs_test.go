package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestAccAlicloudKeyPairsDataSourceBasic(t *testing.T) {
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudKeyPairsDataSourceConfig(map[string]string{
			"name_regex": `"${alicloud_key_pair.default.key_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudKeyPairsDataSourceConfig(map[string]string{
			"name_regex": `"${alicloud_key_pair.default.key_name}_fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudKeyPairsDataSourceConfig(map[string]string{
			"ids": `["${alicloud_key_pair.default.key_name}"]`,
		}),
		fakeConfig: testAccCheckAlicloudKeyPairsDataSourceConfig(map[string]string{
			"ids": `["${alicloud_key_pair.default.key_name}_fake"]`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudKeyPairsDataSourceConfig(map[string]string{
			"name_regex":        `"${alicloud_key_pair.default.key_name}"`,
			"resource_group_id": `"${var.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudKeyPairsDataSourceConfig(map[string]string{
			"name_regex":        `"${alicloud_key_pair.default.key_name}"`,
			"resource_group_id": `"${var.resource_group_id}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudKeyPairsDataSourceConfig(map[string]string{
			"name_regex":        `"${alicloud_key_pair.default.key_name}"`,
			"resource_group_id": `"${var.resource_group_id}"`,
			"ids":               `["${alicloud_key_pair.default.key_name}"]`,
		}),
		fakeConfig: testAccCheckAlicloudKeyPairsDataSourceConfig(map[string]string{
			"name_regex":        `"${alicloud_key_pair.default.key_name}"`,
			"resource_group_id": `"${var.resource_group_id}"`,
			"ids":               `["${alicloud_key_pair.default.key_name}_fake"]`,
		}),
	}
	keyPairsCheckInfo.dataSourceTestCheck(t, 0, nameRegexConf, idsConf, resourceGroupIdConf, allConf)
}

func testAccCheckAlicloudKeyPairsDataSourceConfig(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "resource_group_id" {
	default = "%s"
}


resource "alicloud_key_pair" "default" {
	key_name = "tf-testAcc-key-pair-datasource"
	resource_group_id = "${var.resource_group_id}"
}
data "alicloud_key_pairs" "default" {
	%s
}`, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"), strings.Join(pairs, "\n  "))
	return config
}

var existKeyPairsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"names.#":                       "1",
		"ids.#":                         "1",
		"key_pairs.#":                   "1",
		"key_pairs.0.id":                CHECKSET,
		"key_pairs.0.key_name":          "tf-testAcc-key-pair-datasource",
		"key_pairs.0.resource_group_id": CHECKSET,
		"key_pairs.0.instances.#":       "0",
	}
}

var fakeKeyPairsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"names.#":     "0",
		"ids.#":       "0",
		"key_pairs.#": "0",
	}
}

var keyPairsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_key_pairs.default",
	existMapFunc: existKeyPairsMapFunc,
	fakeMapFunc:  fakeKeyPairsMapFunc,
}
