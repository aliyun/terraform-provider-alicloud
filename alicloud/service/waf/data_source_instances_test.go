package waf_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudWafInstancesDataSource(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.WAFSupportRegions)
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudWafInstanceDataSourceConfig(rand, map[string]string{
			"ids": fmt.Sprintf(`["%s"]`, os.Getenv("ALICLOUD_WAF_INSTANCE_ID")),
		}),
		FakeConfig: testAccCheckAlicloudWafInstanceDataSourceConfig(rand, map[string]string{
			"ids": `["fake"]`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudWafInstanceDataSourceConfig(rand, map[string]string{
			"status": `"1"`,
		}),
		FakeConfig: testAccCheckAlicloudWafInstanceDataSourceConfig(rand, map[string]string{
			"status": `"0"`,
		}),
	}

	instanceSourceConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudWafInstanceDataSourceConfig(rand, map[string]string{
			"instance_source": `"waf-cloud"`,
		}),
		FakeConfig: testAccCheckAlicloudWafInstanceDataSourceConfig(rand, map[string]string{
			"ids":             `["fake"]`,
			"instance_source": `"waf-cloud"`,
		}),
	}

	resourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudWafInstanceDataSourceConfig(rand, map[string]string{
			"resource_group_id": "data.alicloud_resource_manager_resource_groups.default.groups.0.id",
		}),
		FakeConfig: testAccCheckAlicloudWafInstanceDataSourceConfig(rand, map[string]string{
			"ids":               `["fake"]`,
			"resource_group_id": "data.alicloud_resource_manager_resource_groups.default.groups.0.id",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudWafInstanceDataSourceConfig(rand, map[string]string{
			"ids":               fmt.Sprintf(`["%s"]`, os.Getenv("ALICLOUD_WAF_INSTANCE_ID")),
			"status":            `"1"`,
			"instance_source":   `"waf-cloud"`,
			"resource_group_id": "data.alicloud_resource_manager_resource_groups.default.groups.0.id",
		}),
		FakeConfig: testAccCheckAlicloudWafInstanceDataSourceConfig(rand, map[string]string{
			"ids":               `["fake"]`,
			"status":            `"1"`,
			"instance_source":   `"waf-cloud"`,
			"resource_group_id": "data.alicloud_resource_manager_resource_groups.default.groups.0.id",
		}),
	}

	var existDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"instances.#":                   "1",
			"instances.0.id":                acctest.CHECKSET,
			"instances.0.instance_id":       acctest.CHECKSET,
			"instances.0.end_date":          acctest.CHECKSET,
			"instances.0.in_debt":           acctest.CHECKSET,
			"instances.0.remain_day":        acctest.CHECKSET,
			"instances.0.status":            "1",
			"instances.0.subscription_type": "Subscription",
			"instances.0.trial":             acctest.CHECKSET,
		}
	}

	var fakeDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var wafInstancesRecordsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_waf_instances.default",
		ExistMapFunc: existDnsRecordsMapFunc,
		FakeMapFunc:  fakeDnsRecordsMapFunc,
	}

	var perCheck = func() {
		acctest.PreCheck(t)
	}

	wafInstancesRecordsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, perCheck, idsConf, statusConf, instanceSourceConf, resourceGroupIdConf, allConf)

}

func testAccCheckAlicloudWafInstanceDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
data "alicloud_resource_manager_resource_groups" "default"{
	name_regex = "^default$"
}

data "alicloud_waf_instances" "default" {
  %s
}
`, strings.Join(pairs, "\n  "))
	return config
}
