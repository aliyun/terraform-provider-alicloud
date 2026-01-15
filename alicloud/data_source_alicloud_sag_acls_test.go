package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSagAclsDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSagAclsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_sag_acl.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSagAclsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_sag_acl.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSagAclsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_sag_acl.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSagAclsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_sag_acl.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSagAclsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_sag_acl.default.id}"]`,
			"name_regex": `"${alicloud_sag_acl.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSagAclsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_sag_acl.default.id}_fake"]`,
			"name_regex": `"${alicloud_sag_acl.default.name}"`,
		}),
	}

	var existSagAclsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"acls.#":      "1",
			"ids.#":       "1",
			"names.#":     "1",
			"acls.0.id":   CHECKSET,
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

	var sagAclsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_sag_acls.default",
		existMapFunc: existSagAclsMapFunc,
		fakeMapFunc:  fakeSagAclsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
	}
	sagAclsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, allConf)
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
