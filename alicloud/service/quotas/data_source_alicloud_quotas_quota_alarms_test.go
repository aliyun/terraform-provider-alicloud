package quotas_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudQuotasQuotaAlarmsDataSource(t *testing.T) {
	resourceId := "data.alicloud_quotas_quota_alarms.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccQuotasQuotaAlarmsTest%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceQuotasQuotaAlarmsDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_quotas_quota_alarm.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_quotas_quota_alarm.default.id}-fake"},
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":     name,
			"enable_details": "true",
			"ids":            []string{"${alicloud_quotas_quota_alarm.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     name,
			"enable_details": "true",
			"ids":            []string{"${alicloud_quotas_quota_alarm.default.id}-fake"},
		}),
	}
	dimensionsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_quotas_quota_alarm.default.id}"},
			"quota_dimensions": []map[string]interface{}{
				{
					"key":   "regionId",
					"value": "cn-hangzhou",
				},
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_quotas_quota_alarm.default.id}"},
			"quota_dimensions": []map[string]interface{}{
				{
					"key":   "regionId",
					"value": "cn-beijing",
				},
			},
		}),
	}
	var existQuotasQuotaAlarmsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"ids.0":                       acctest.CHECKSET,
			"names.#":                     "1",
			"names.0":                     name,
			"alarms.#":                    "1",
			"alarms.0.id":                 acctest.CHECKSET,
			"alarms.0.alarm_id":           acctest.CHECKSET,
			"alarms.0.product_code":       "ecs",
			"alarms.0.quota_action_code":  "q_prepaid-instance-count-per-once-purchase",
			"alarms.0.quota_alarm_name":   name,
			"alarms.0.threshold":          "100",
			"alarms.0.quota_dimensions.#": "1",
			"alarms.0.threshold_percent":  "0",
			"alarms.0.web_hook":           "",
		}
	}

	var fakeQuotasQuotaAlarmsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"alarms.#": "0",
		}
	}

	var QuotasApplicationInfosInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existQuotasQuotaAlarmsMapFunc,
		FakeMapFunc:  fakeQuotasQuotaAlarmsMapFunc,
	}

	QuotasApplicationInfosInfo.DataSourceTestCheck(t, 0, idsConf, nameRegexConf, dimensionsConf)
}

func dataSourceQuotasQuotaAlarmsDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_quotas_quota_alarm" "default" {
		quota_alarm_name  = "%s"
		product_code      = "ecs"
		quota_action_code = "q_prepaid-instance-count-per-once-purchase"
		threshold         = "100"
		quota_dimensions {
			key   = "regionId"
			value = "cn-hangzhou"
		}
	}`, name)
}
