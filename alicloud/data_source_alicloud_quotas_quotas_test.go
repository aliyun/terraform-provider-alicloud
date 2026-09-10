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
			"name_regex":   "Maximum*",
			"product_code": "ecs",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":   "Maximum-fake",
			"product_code": "ecs",
		}),
	}
	actionCodeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"product_code":      "ecs",
			"quota_action_code": "q_cloud-assistant-activation-count",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"product_code":      "ecs",
			"quota_action_code": "q_cloud-assistant-activation-count-fake",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"product_code":      "ecs",
			"name_regex":        "Maximum*",
			"quota_action_code": "q_cloud-assistant-activation-count",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"product_code":      "ecs",
			"name_regex":        "Maximum-fake",
			"quota_action_code": "q_cloud-assistant-activation-count-fake",
		}),
	}
	var existQuotasQuotasMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    CHECKSET,
			"quotas.#": CHECKSET,
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
