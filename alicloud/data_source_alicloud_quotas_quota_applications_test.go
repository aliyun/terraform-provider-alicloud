package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

// The quota product does not support deletion, so skip the test.
func SkipTestAccAlicloudQuotasApplicationInfosDataSource(t *testing.T) {
	resourceId := "data.alicloud_quotas_quota_applications.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceQuotasQuotaApplicationsDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"product_code":   "ess",
			"enable_details": "true",
			"ids":            []string{"${alicloud_quotas_application_info.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"product_code":   "ess",
			"enable_details": "true",
			"ids":            []string{"${alicloud_quotas_application_info.default.id}-fake"},
		}),
	}
	dimenstionsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"product_code":   "ess",
			"enable_details": "true",
			"ids":            []string{"${alicloud_quotas_application_info.default.id}"},
			"dimensions": []map[string]interface{}{
				{
					"key":   "regionId",
					"value": "cn-hangzhou",
				},
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"product_code":   "ess",
			"enable_details": "true",
			"ids":            []string{"${alicloud_quotas_application_info.default.id}-fake"},
			"dimensions": []map[string]interface{}{
				{
					"key":   "regionId",
					"value": "cn-hangzhou",
				},
			},
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"product_code":   "ess",
			"enable_details": "true",
			"ids":            []string{"${alicloud_quotas_application_info.default.id}"},
			"status":         "Process",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"product_code":   "ess",
			"enable_details": "true",
			"ids":            []string{"${alicloud_quotas_application_info.default.id}"},
			"status":         "Disagree",
		}),
	}
	var existQuotasApplicationInfosMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"ids.0":                     CHECKSET,
			"infos.#":                   "1",
			"infos.0.approve_value":     CHECKSET,
			"infos.0.audit_reason":      "",
			"infos.0.desire_value":      "100",
			"infos.0.id":                CHECKSET,
			"infos.0.dimensions.#":      "1",
			"infos.0.effective_time":    "",
			"infos.0.expire_time":       "",
			"infos.0.notice_type":       "0",
			"infos.0.product_code":      "ess",
			"infos.0.quota_action_code": "q_db_instance",
			"infos.0.quota_description": CHECKSET,
			"infos.0.quota_name":        CHECKSET,
			"infos.0.quota_unit":        "",
			"infos.0.reason":            CHECKSET,
			"infos.0.status":            CHECKSET,
		}
	}

	var fakeQuotasApplicationInfosMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"infos.#": "0",
		}
	}

	var QuotasApplicationInfosInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existQuotasApplicationInfosMapFunc,
		fakeMapFunc:  fakeQuotasApplicationInfosMapFunc,
	}

	QuotasApplicationInfosInfo.dataSourceTestCheck(t, 0, idsConf, dimenstionsConf, statusConf)
}

func dataSourceQuotasQuotaApplicationsDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_quotas_quota_application" "default" {
		notice_type  =  "0"
		desire_value =  "100"
		product_code =  "ess"
		quota_action_code = "q_db_instance"
		reason       =   "For Terraform Test"
		dimensions {
                key   =   "regionId"
                value = "cn-hangzhou"
				}
	}`)
}
