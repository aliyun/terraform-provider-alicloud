package cloudmonitorservice_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCmsSlsGroupsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsSlsGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_sls_group.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudCmsSlsGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_sls_group.default.id}_fake"]`,
		}),
	}
	keywordConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsSlsGroupsDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_cms_sls_group.default.id}"]`,
			"keyword": `"${alicloud_cms_sls_group.default.sls_group_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudCmsSlsGroupsDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_cms_sls_group.default.id}"]`,
			"keyword": `"${alicloud_cms_sls_group.default.sls_group_name}_fake"`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsSlsGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cms_sls_group.default.sls_group_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudCmsSlsGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cms_sls_group.default.sls_group_name}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsSlsGroupsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cms_sls_group.default.id}"]`,
			"keyword":    `"${alicloud_cms_sls_group.default.sls_group_name}"`,
			"name_regex": `"${alicloud_cms_sls_group.default.sls_group_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudCmsSlsGroupsDataSourceName(rand, map[string]string{
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
			"groups.0.id":                 acctest.CHECKSET,
			"groups.0.create_time":        acctest.CHECKSET,
			"groups.0.sls_group_config.#": "1",
			"groups.0.sls_group_config.0.sls_user_id":  acctest.CHECKSET,
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
	var alicloudCmsSlsGroupsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cms_sls_groups.default",
		ExistMapFunc: existAlicloudCmsSlsGroupsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudCmsSlsGroupsDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudCmsSlsGroupsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, keywordConf, nameRegexConf, allConf)
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
