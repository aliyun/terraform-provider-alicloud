package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudWAFInstancesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudWafInstanceDataSourceConfig(rand, map[string]string{
			"ids": fmt.Sprintf(`["%s"]`, os.Getenv("ALICLOUD_WAF_INSTANCE_ID")),
		}),
		fakeConfig: testAccCheckAlicloudWafInstanceDataSourceConfig(rand, map[string]string{
			"ids": `["fake"]`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudWafInstanceDataSourceConfig(rand, map[string]string{
			"status": `"1"`,
		}),
		fakeConfig: testAccCheckAlicloudWafInstanceDataSourceConfig(rand, map[string]string{
			"status": `"0"`,
		}),
	}

	instanceSourceConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudWafInstanceDataSourceConfig(rand, map[string]string{
			"instance_source": `"waf-cloud"`,
		}),
		fakeConfig: testAccCheckAlicloudWafInstanceDataSourceConfig(rand, map[string]string{
			"ids":             `["fake"]`,
			"instance_source": `"waf-cloud"`,
		}),
	}

	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudWafInstanceDataSourceConfig(rand, map[string]string{
			"resource_group_id": "data.alicloud_resource_manager_resource_groups.default.groups.0.id",
		}),
		fakeConfig: testAccCheckAlicloudWafInstanceDataSourceConfig(rand, map[string]string{
			"ids":               `["fake"]`,
			"resource_group_id": "data.alicloud_resource_manager_resource_groups.default.groups.0.id",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudWafInstanceDataSourceConfig(rand, map[string]string{
			"ids":               fmt.Sprintf(`["%s"]`, os.Getenv("ALICLOUD_WAF_INSTANCE_ID")),
			"status":            `"1"`,
			"instance_source":   `"waf-cloud"`,
			"resource_group_id": "data.alicloud_resource_manager_resource_groups.default.groups.0.id",
		}),
		fakeConfig: testAccCheckAlicloudWafInstanceDataSourceConfig(rand, map[string]string{
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
			"instances.0.id":                CHECKSET,
			"instances.0.instance_id":       CHECKSET,
			"instances.0.end_date":          CHECKSET,
			"instances.0.in_debt":           CHECKSET,
			"instances.0.remain_day":        CHECKSET,
			"instances.0.status":            "1",
			"instances.0.subscription_type": "Subscription",
			"instances.0.trial":             CHECKSET,
		}
	}

	var fakeDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var wafInstancesRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_waf_instances.default",
		existMapFunc: existDnsRecordsMapFunc,
		fakeMapFunc:  fakeDnsRecordsMapFunc,
	}

	var perCheck = func() {
		testAccPreCheck(t)
		testAccPreCheckWithEnvVariable(t, "ALICLOUD_WAF_INSTANCE_ID")
	}

	wafInstancesRecordsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, perCheck, idsConf, statusConf, instanceSourceConf, resourceGroupIdConf, allConf)

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
