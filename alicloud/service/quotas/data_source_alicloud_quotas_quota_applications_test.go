package quotas_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

// The quota product does not support deletion, so skip the test.
func TestAccAlicloudQuotasApplicationInfosDataSource(t *testing.T) {
	resourceId := "data.alicloud_quotas_quota_applications.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceQuotasQuotaApplicationsDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"product_code":   "vpc",
			"enable_details": "true",
			"quota_category": "${alicloud_quotas_quota_application.default.quota_category}",
			"ids":            []string{"${alicloud_quotas_quota_application.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"product_code":   "vpc",
			"enable_details": "true",
			"ids":            []string{"${alicloud_quotas_quota_application.default.id}-fake"},
		}),
	}
	dimenstionsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"product_code":   "vpc",
			"quota_category": "${alicloud_quotas_quota_application.default.quota_category}",
			"enable_details": "true",
			"ids":            []string{"${alicloud_quotas_quota_application.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"product_code":   "vpc",
			"quota_category": "${alicloud_quotas_quota_application.default.quota_category}",
			"enable_details": "true",
			"ids":            []string{"${alicloud_quotas_quota_application.default.id}-fake"},
			"dimensions": []map[string]interface{}{
				{
					"key":   "regionId",
					"value": "cn-hangzhou",
				},
			},
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"product_code":   "vpc",
			"quota_category": "${alicloud_quotas_quota_application.default.quota_category}",
			"enable_details": "true",
			"ids":            []string{"${alicloud_quotas_quota_application.default.id}"},
			"status":         "Agree",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"product_code":   "vpc",
			"enable_details": "true",
			"ids":            []string{"${alicloud_quotas_quota_application.default.id}"},
			"status":         "Disagree",
		}),
	}
	var existQuotasApplicationInfosMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "1",
		}
	}

	var fakeQuotasApplicationInfosMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}

	var QuotasApplicationInfosInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existQuotasApplicationInfosMapFunc,
		FakeMapFunc:  fakeQuotasApplicationInfosMapFunc,
	}

	QuotasApplicationInfosInfo.DataSourceTestCheck(t, 0, idsConf, dimenstionsConf, statusConf)
}

func dataSourceQuotasQuotaApplicationsDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_quotas_quota_application" "default" {
	  product_code = "vpc"
	  notice_type = "3"
	  effective_time = "%s"
	  expire_time = "%s"
	  desire_value = "1"
	  reason = ""
	  quota_action_code = "vpc_whitelist/ha_vip_whitelist"
	  audit_mode = "Sync"
	  env_language = "zh"
	  quota_category = "WhiteListLabel"
	}`, time.Now().Add(1*time.Minute).Format("2006-01-02T15:04:05Z"), time.Now().Add(1*time.Hour).Format("2006-01-02T15:04:05Z"))
}
