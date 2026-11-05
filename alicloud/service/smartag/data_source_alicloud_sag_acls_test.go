package smartag_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudSagAclsDataSource_basic(t *testing.T) {
	rand := tfacctest.RandInt()
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSagAclsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_sag_acl.default.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudSagAclsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_sag_acl.default.name}_fake"`,
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSagAclsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_sag_acl.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudSagAclsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_sag_acl.default.id}_fake"]`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSagAclsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_sag_acl.default.id}"]`,
			"name_regex": `"${alicloud_sag_acl.default.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudSagAclsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_sag_acl.default.id}_fake"]`,
			"name_regex": `"${alicloud_sag_acl.default.name}"`,
		}),
	}

	var existSagAclsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"acls.#":      "1",
			"ids.#":       "1",
			"names.#":     "1",
			"acls.0.id":   acctest.CHECKSET,
			"acls.0.name": fmt.Sprintf("tf-testAccSagAclDataSourceBisic-%d", rand),
		}
	}

	var fakeSagAclsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"acls.#":  "0",
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var sagAclsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_sag_acls.default",
		ExistMapFunc: existSagAclsMapFunc,
		FakeMapFunc:  fakeSagAclsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
	}
	sagAclsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlicloudSagAclsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
		variable "name" {
			default = "tf-testAccSagAclDataSourceBisic-%d"
		}
		resource "alicloud_sag_acl" "default" {
			name = "${var.name}"
		}

		data "alicloud_sag_acls" "default" {
		  %s
		}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
