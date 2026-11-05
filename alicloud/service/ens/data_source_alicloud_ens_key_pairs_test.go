package ens_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudEnsKeyPairsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEnsKeyPairSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ens_key_pair.default.key_pair_name}"`,
			"version":    `"2017-11-10"`,
		}),
		FakeConfig: testAccCheckAlicloudEnsKeyPairSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ens_key_pair.default.key_pair_name}_fake"`,
			"version":    `"2017-11-10"`,
		}),
	}

	namesRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEnsKeyPairSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ens_key_pair.default.key_pair_name}"`,
			"version":    `"2017-11-10"`,
		}),
		FakeConfig: testAccCheckAlicloudEnsKeyPairSourceConfig(rand, map[string]string{
			"name_regex": `"TestAccAlicloudAmqpBindingsDataSource"`,
			"version":    `"2017-11-10"`,
		}),
	}

	keyNameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEnsKeyPairSourceConfig(rand, map[string]string{
			"key_pair_name": `"${alicloud_ens_key_pair.default.key_pair_name}"`,
			"version":       `"2017-11-10"`,
		}),
		FakeConfig: testAccCheckAlicloudEnsKeyPairSourceConfig(rand, map[string]string{
			"key_pair_name": `"${alicloud_ens_key_pair.default.key_pair_name}_fake"`,
			"version":       `"2017-11-10"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEnsKeyPairSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ens_key_pair.default.key_pair_name}"`,
			"version":    `"2017-11-10"`,
		}),
		FakeConfig: testAccCheckAlicloudEnsKeyPairSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ens_key_pair.default.key_pair_name}_fake"`,
			"version":    `"2017-11-10"`,
		}),
	}

	EnsKeyPairCheckInfo.DataSourceTestCheck(t, rand, nameRegexConf, namesRegexConf, keyNameRegexConf, allConf)
}

func testAccCheckAlicloudEnsKeyPairSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccEnsKeyPairsDataSource%d"
}

resource "alicloud_ens_key_pair" "default" {
  key_pair_name = var.name
  version       = "2017-11-10"
}

data "alicloud_ens_key_pairs" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

var existEnsKeyPairMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"pairs.#":               "1",
		"pairs.0.version":       "2017-11-10",
		"pairs.0.key_pair_name": fmt.Sprintf("tf-testAccEnsKeyPairsDataSource%d", rand),
	}
}

var fakeEnsKeyPairMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"pairs.#": "0",
	}
}

var EnsKeyPairCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_ens_key_pairs.default",
	ExistMapFunc: existEnsKeyPairMapFunc,
	FakeMapFunc:  fakeEnsKeyPairMapFunc,
}
