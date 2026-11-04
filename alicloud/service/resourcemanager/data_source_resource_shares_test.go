package resourcemanager_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudResourceManagerResourceSharesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudResourceManagerResourceSharesSourceConfig(rand, map[string]string{
			"name_regex": `"^${alicloud_resource_manager_resource_share.example.resource_share_name}$"`,
		}),
		FakeConfig: testAccCheckAlicloudResourceManagerResourceSharesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_resource_share.example.resource_share_name}_fake"`,
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudResourceManagerResourceSharesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_resource_manager_resource_share.example.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudResourceManagerResourceSharesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_resource_manager_resource_share.example.id}_fake"]`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudResourceManagerResourceSharesSourceConfig(rand, map[string]string{
			"name_regex": `"^${alicloud_resource_manager_resource_share.example.resource_share_name}$"`,
			"status":     `"Active"`,
		}),
		FakeConfig: testAccCheckAlicloudResourceManagerResourceSharesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_resource_share.example.resource_share_name}"`,
			"status":     `"Deleted"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudResourceManagerResourceSharesSourceConfig(rand, map[string]string{
			"name_regex": `"^${alicloud_resource_manager_resource_share.example.resource_share_name}$"`,
			"ids":        `["${alicloud_resource_manager_resource_share.example.id}"]`,
			"status":     `"Active"`,
		}),
		FakeConfig: testAccCheckAlicloudResourceManagerResourceSharesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_resource_share.example.resource_share_name}_fake"`,
			"ids":        `["${alicloud_resource_manager_resource_share.example.id}"]`,
			"status":     `"Deleted"`,
		}),
	}

	var existResourceManagerResourceSharesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"shares.#":                      "1",
			"names.#":                       "1",
			"ids.#":                         "1",
			"shares.0.id":                   acctest.CHECKSET,
			"shares.0.resource_share_name":  fmt.Sprintf("tf-testaccResourceManagerResourceShare%d", rand),
			"shares.0.resource_share_id":    acctest.CHECKSET,
			"shares.0.resource_share_owner": acctest.CHECKSET,
			"shares.0.status":               "Active",
		}
	}

	var fakeResourceManagerResourceSharesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"shares.#": "0",
			"ids.#":    "0",
			"names.#":  "0",
		}
	}

	var resourceGroupsRecordsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_resource_manager_resource_shares.example",
		ExistMapFunc: existResourceManagerResourceSharesMapFunc,
		FakeMapFunc:  fakeResourceManagerResourceSharesMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckEnterpriseAccountEnabled(t)
	}
	resourceGroupsRecordsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, statusConf, allConf)

}

func testAccCheckAlicloudResourceManagerResourceSharesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
resource "alicloud_resource_manager_resource_share" "example"{
	resource_share_name = "tf-testaccResourceManagerResourceShare%d"
}

data "alicloud_resource_manager_resource_shares" "example"{
	resource_share_owner = "Self"
%s
}

`, rand, strings.Join(pairs, "\n   "))
	return config
}
