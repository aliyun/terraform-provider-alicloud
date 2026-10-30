package quotas_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudQuotasQuotasDataSource(t *testing.T) {
	resourceId := "data.alicloud_quotas_quotas.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccQuotasQuotas%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceQuotasQuotasDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":   "Maximum*",
			"product_code": "ecs",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":   "Maximum-fake",
			"product_code": "ecs",
		}),
	}
	actionCodeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"product_code":      "ecs",
			"quota_action_code": "q_cloud-assistant-activation-count",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"product_code":      "ecs",
			"quota_action_code": "q_cloud-assistant-activation-count-fake",
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"product_code":      "ecs",
			"name_regex":        "Maximum*",
			"quota_action_code": "q_cloud-assistant-activation-count",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"product_code":      "ecs",
			"name_regex":        "Maximum-fake",
			"quota_action_code": "q_cloud-assistant-activation-count-fake",
		}),
	}
	var existQuotasQuotasMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    acctest.CHECKSET,
			"quotas.#": acctest.CHECKSET,
		}
	}

	var fakeQuotasQuotasMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"quotas.#": "0",
		}
	}

	var QuotasQuotasInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existQuotasQuotasMapFunc,
		FakeMapFunc:  fakeQuotasQuotasMapFunc,
	}

	QuotasQuotasInfo.DataSourceTestCheck(t, 0, nameRegexConf, actionCodeConf, allConf)
}

func dataSourceQuotasQuotasDependence(name string) string {
	return ""
}
