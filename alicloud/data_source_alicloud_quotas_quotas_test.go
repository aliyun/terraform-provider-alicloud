package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudQuotasQuotasDataSource(t *testing.T) {
	resourceId := "data.alicloud_quotas_quotas.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccQuotasQuotas%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceQuotasQuotasDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":   "专有宿主机总数量上限",
			"product_code": "ecs",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":   "专有宿主机总数量上限-fake",
			"product_code": "ecs",
		}),
	}
	actionCodeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"product_code":      "ecs",
			"quota_action_code": "q_dedicated-hosts",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"product_code":      "ecs",
			"quota_action_code": "q_dedicated-hosts-fake",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"product_code":      "ecs",
			"name_regex":        "专有宿主机总数量上限",
			"quota_action_code": "q_dedicated-hosts",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"product_code":      "ecs",
			"name_regex":        "专有宿主机总数量上限-fake",
			"quota_action_code": "q_dedicated-hosts-fake",
		}),
	}
	var existQuotasQuotasMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"ids.0":                        CHECKSET,
			"names.#":                      "1",
			"names.0":                      CHECKSET,
			"quotas.#":                     "1",
			"quotas.0.id":                  CHECKSET,
			"quotas.0.adjustable":          CHECKSET,
			"quotas.0.applicable_type":     CHECKSET,
			"quotas.0.consumable":          CHECKSET,
			"quotas.0.quota_action_code":   "q_dedicated-hosts",
			"quotas.0.quota_description":   CHECKSET,
			"quotas.0.quota_name":          CHECKSET,
			"quotas.0.quota_type":          "",
			"quotas.0.quota_unit":          "",
			"quotas.0.total_quota":         CHECKSET,
			"quotas.0.total_usage":         CHECKSET,
			"quotas.0.unadjustable_detail": "",
		}
	}

	var fakeQuotasQuotasMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"quotas.#": "0",
		}
	}

	var QuotasQuotasInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existQuotasQuotasMapFunc,
		fakeMapFunc:  fakeQuotasQuotasMapFunc,
	}

	QuotasQuotasInfo.dataSourceTestCheck(t, 0, nameRegexConf, actionCodeConf, allConf)
}

func dataSourceQuotasQuotasDependence(name string) string {
	return ""
}
