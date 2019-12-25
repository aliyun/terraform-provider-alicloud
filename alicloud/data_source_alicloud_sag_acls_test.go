package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudSagAclsDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSagAclsDataSourceConfig(rand, map[string]string{
			"name_regex": `alicloud_sag_acl.default.name`,
		}),
		fakeConfig: testAccCheckAlicloudSagAclsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_sag_acl.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSagAclsDataSourceConfig(rand, map[string]string{
			"ids": `[alicloud_sag_acl.default.id]`,
		}),
		fakeConfig: testAccCheckAlicloudSagAclsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_sag_acl.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSagAclsDataSourceConfig(rand, map[string]string{
			"ids":        `[alicloud_sag_acl.default.id]`,
			"name_regex": `alicloud_sag_acl.default.name`,
		}),
		fakeConfig: testAccCheckAlicloudSagAclsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_sag_acl.default.id}_fake"]`,
			"name_regex": `alicloud_sag_acl.default.name`,
		}),
	}

	var existDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"acls.#":      "1",
			"ids.#":       "1",
			"names.#":     "1",
			"acls.0.id":   CHECKSET,
			"acls.0.name": fmt.Sprintf("tf-testAccSagAclDataSourceBisic-%d", rand),
		}
	}

	var fakeDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"acls.#":  "0",
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var sagAclsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_sag_acls.default",
		existMapFunc: existDnsRecordsMapFunc,
		fakeMapFunc:  fakeDnsRecordsMapFunc,
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
			name = var.name
		}

		data "alicloud_sag_acls" "default" {
		  %s
		}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
