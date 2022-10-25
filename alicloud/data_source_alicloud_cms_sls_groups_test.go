package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCmsSlsGroupsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsSlsGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_sls_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCmsSlsGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_sls_group.default.id}_fake"]`,
		}),
	}
	keywordConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsSlsGroupsDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_cms_sls_group.default.id}"]`,
			"keyword": `"${alicloud_cms_sls_group.default.sls_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsSlsGroupsDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_cms_sls_group.default.id}"]`,
			"keyword": `"${alicloud_cms_sls_group.default.sls_group_name}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsSlsGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cms_sls_group.default.sls_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsSlsGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cms_sls_group.default.sls_group_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsSlsGroupsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cms_sls_group.default.id}"]`,
			"keyword":    `"${alicloud_cms_sls_group.default.sls_group_name}"`,
			"name_regex": `"${alicloud_cms_sls_group.default.sls_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsSlsGroupsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cms_sls_group.default.id}_fake"]`,
			"keyword":    `"${alicloud_cms_sls_group.default.sls_group_name}_fake"`,
			"name_regex": `"${alicloud_cms_sls_group.default.sls_group_name}_fake"`,
		}),
	}
	var existAlicloudCmsSlsGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"names.#":                     "1",
			"groups.#":                    "1",
			"groups.0.id":                 CHECKSET,
			"groups.0.create_time":        CHECKSET,
			"groups.0.sls_group_config.#": "1",
			"groups.0.sls_group_config.0.sls_user_id":  CHECKSET,
			"groups.0.sls_group_config.0.sls_logstore": "Logstore-ECS",
			"groups.0.sls_group_config.0.sls_project":  "aliyun-project",
			"groups.0.sls_group_config.0.sls_region":   "cn-hangzhou",
			"groups.0.sls_group_description":           fmt.Sprintf("tf_testAccSlsGroup_%d", rand),
			"groups.0.sls_group_name":                  fmt.Sprintf("tf_testAccSlsGroup_%d", rand),
		}
	}
	var fakeAlicloudCmsSlsGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCmsSlsGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cms_sls_groups.default",
		existMapFunc: existAlicloudCmsSlsGroupsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCmsSlsGroupsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudCmsSlsGroupsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, keywordConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudCmsSlsGroupsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf_testAccSlsGroup_%d"
}
data "alicloud_account" "this" {}

resource "alicloud_cms_sls_group" "default" {
	sls_group_config {
		sls_user_id = data.alicloud_account.this.id
		sls_logstore = "Logstore-ECS"
		sls_project = "aliyun-project"
		sls_region = "cn-hangzhou"
	}
	sls_group_description = var.name
	sls_group_name = var.name
}

data "alicloud_cms_sls_groups" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
