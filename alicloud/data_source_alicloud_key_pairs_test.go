package alicloud

import (
	"fmt"
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
	keyPairsCheckInfo.dataSourceTestCheck(t, 0, nameRegexConf)
}

func testAccCheckAlicloudKeyPairsDataSourceConfig(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
resource "alicloud_key_pair" "default" {
	key_name = "tf-testAcc-key-pair-datasource"
}
data "alicloud_key_pairs" "default" {
	%s
}`, strings.Join(pairs, "\n  "))
	return config
}

var existKeyPairsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"names.#":                  "1",
		"key_pairs.#":             "1",
		"key_pairs.0.id":          CHECKSET,
		"key_pairs.0.key_name":    "tf-testAcc-key-pair-datasource",
		"key_pairs.0.instances.#": "0",
	}
}

var fakeKeyPairsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"names.#" :   "0",
		"key_pairs.#": "0",
	}
}

var keyPairsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_key_pairs.default",
	existMapFunc: existKeyPairsMapFunc,
	fakeMapFunc:  fakeKeyPairsMapFunc,
}
