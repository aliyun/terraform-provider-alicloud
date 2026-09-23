package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudQuotasTemplateApplicationsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudQuotasTemplateApplicationsSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_quotas_template_applications.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudQuotasTemplateApplicationsSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_quotas_template_applications.default.id}_fake"]`,
		}),
	}

	ProductCodeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudQuotasTemplateApplicationsSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_quotas_template_applications.default.id}"]`,
			"product_code": `"vpc"`,
		}),
		fakeConfig: testAccCheckAlicloudQuotasTemplateApplicationsSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_quotas_template_applications.default.id}_fake"]`,
			"product_code": `"vpc_fake"`,
		}),
	}
	QuotaActionCodeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudQuotasTemplateApplicationsSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_quotas_template_applications.default.id}"]`,
			"quota_action_code": `"vpc_whitelist/ha_vip_whitelist"`,
		}),
		fakeConfig: testAccCheckAlicloudQuotasTemplateApplicationsSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_quotas_template_applications.default.id}_fake"]`,
			"quota_action_code": `"vpc_whitelist/ha_vip_whitelist_fake"`,
		}),
	}
	QuotaCategoryConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudQuotasTemplateApplicationsSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_quotas_template_applications.default.id}"]`,
			"quota_category": `"FlowControl"`,
		}),
		fakeConfig: testAccCheckAlicloudQuotasTemplateApplicationsSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_quotas_template_applications.default.id}_fake"]`,
			"quota_category": `"CommonQuota"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudQuotasTemplateApplicationsSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_quotas_template_applications.default.id}"]`,
			"product_code": `"vpc"`,

			"quota_action_code": `"vpc_whitelist/ha_vip_whitelist"`,

			"quota_category": `"FlowControl"`,
		}),
		fakeConfig: testAccCheckAlicloudQuotasTemplateApplicationsSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_quotas_template_applications.default.id}_fake"]`,
			"product_code": `"vpc_fake"`,

			"quota_action_code": `"vpc_whitelist/ha_vip_whitelist_fake"`,

			"quota_category": `"CommonQuota"`,
		}),
	}

	QuotasTemplateApplicationsCheckInfo.dataSourceTestCheck(t, rand, idsConf, ProductCodeConf, QuotaActionCodeConf, QuotaCategoryConf, allConf)
}

var existQuotasTemplateApplicationsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"applications.#":    "1",
		"applications.0.id": CHECKSET,
	}
}

var fakeQuotasTemplateApplicationsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"applications.#": "0",
	}
}

var QuotasTemplateApplicationsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_quotas_template_applications.default",
	existMapFunc: existQuotasTemplateApplicationsMapFunc,
	fakeMapFunc:  fakeQuotasTemplateApplicationsMapFunc,
}

func testAccCheckAlicloudQuotasTemplateApplicationsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccQuotasTemplateApplications%d"
}

data "alicloud_resource_manager_accounts" "default" {
  status  = "CreateSuccess"
}



resource "alicloud_quotas_template_applications" "default" {
  quota_action_code = "vpc_whitelist/ha_vip_whitelist"
  product_code      = "vpc"
  quota_category    = "FlowControl"
  aliyun_uids       = ["${data.alicloud_resource_manager_accounts.default.ids.0}"]
  desire_value      = 6
  notice_type       = "0"
  env_language      = "zh"
  reason            = "测试"
  dimensions {
    key   = "apiName"
    value = "GetProductQuotaDimension"
  }
  dimensions {
    key   = "apiVersion"
    value = "2020-05-10"
  }
  dimensions {
    key   = "regionId"
    value = "cn-hangzhou"
  }
}

data "alicloud_quotas_template_applications" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
