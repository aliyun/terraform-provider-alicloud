package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudBssOpenApiBudget_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bss_open_api_budget.default"
	ra := resourceAttrInit(resourceId, AliCloudBssOpenApiBudgetMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &BssOpenApiServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBssOpenApiBudget")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc-budget-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudBssOpenApiBudgetBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"budget_name":        name,
					"budget_type":        "cost",
					"metric":             "Cost",
					"cycle_type":         "Month",
					"cycle_start_period": "2026-09",
					"cycle_end_period":   "2026-12",
					"quota_type":         "quota",
					"quota":              "80",
					"comment":            "tf test budget",
					"nbid":               "",
					"cycle_quota": []map[string]interface{}{
						{"cycle_period": "2026-09", "quota": "80"},
						{"cycle_period": "2026-10", "quota": "90"},
					},
					"query_filter": []map[string]interface{}{
						{"code": "productCode", "select_type": "equal", "values": []string{"ECS"}},
					},
					"warn_confs": []map[string]interface{}{
						{"name": "tf-warn", "warn_target": "Msc", "threshold_type": "percent", "threshold_value": "80", "msc_contacts": []string{}, "msc_channels": []string{}, "event_bridge": false, "comment": "tf warn conf"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"budget_name":        name,
						"budget_type":        "cost",
						"metric":             "Cost",
						"cycle_type":         "Month",
						"cycle_start_period": "2026-09",
						"cycle_end_period":   "2026-12",
						"quota_type":         "quota",
						"quota":              "80",
						"comment":            "tf test budget",
						"cycle_quota.#":      "2",
						"query_filter.#":     "1",
						"warn_confs.#":       "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"budget_name":        name,
					"budget_type":        "usage",
					"metric":             "Usage",
					"cycle_type":         "Year",
					"cycle_start_period": "2026",
					"cycle_end_period":   "2027",
					"quota_type":         "percent",
					"quota":              "60",
					"comment":            "tf test budget updated",
					"nbid":               "",
					"cycle_quota": []map[string]interface{}{
						{"cycle_period": "2026", "quota": "60"},
					},
					"query_filter": []map[string]interface{}{
						{"code": "region", "select_type": "contains", "values": []string{"cn-hangzhou", "cn-beijing"}},
					},
					"warn_confs": []map[string]interface{}{
						{"name": "tf-warn-updated", "warn_target": "Dingtalk", "threshold_type": "absolute", "threshold_value": "60", "msc_contacts": []string{"contact1"}, "msc_channels": []string{"channel1"}, "event_bridge": true, "comment": "tf warn conf updated"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"budget_type":        "usage",
						"metric":             "Usage",
						"cycle_type":         "Year",
						"cycle_start_period": "2026",
						"cycle_end_period":   "2027",
						"quota_type":         "percent",
						"quota":              "60",
						"comment":            "tf test budget updated",
						"cycle_quota.#":      "1",
						"query_filter.#":     "1",
						"warn_confs.#":       "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"nbid"},
			},
		},
	})
}

// AliCloudBssOpenApiBudgetMap0 asserts the persisted attributes after apply.
var AliCloudBssOpenApiBudgetMap0 = map[string]string{
	"budget_name":        CHECKSET,
	"budget_type":        CHECKSET,
	"comment":            CHECKSET,
	"cycle_end_period":   CHECKSET,
	"cycle_start_period": CHECKSET,
	"cycle_type":         CHECKSET,
	"metric":             CHECKSET,
	"quota":              CHECKSET,
	"quota_type":         CHECKSET,
	"cycle_quota.#":      CHECKSET,
	"query_filter.#":     CHECKSET,
	"warn_confs.#":       CHECKSET,
}

// AliCloudBssOpenApiBudgetBasicDependence0 returns the dependency HCL (none for a global BSS resource).
var AliCloudBssOpenApiBudgetBasicDependence0 = func(name string) string {
	return ""
}
