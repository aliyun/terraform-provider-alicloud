package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAlicloudBssOpenApiBudgetsDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc-budget-ds-%d", rand)
	resourceId := "alicloud_bss_open_api_budget.default"
	dataSourceId := "data.alicloud_bss_open_api_budgets.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBssOpenApiBudgetDestroyWith(name),
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudBssOpenApiBudgetsDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "budgets.#", "1"),
					resource.TestCheckResourceAttr(dataSourceId, "ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceId, "budget_name", dataSourceId, "budgets.0.budget_name"),
					resource.TestCheckResourceAttrPair(resourceId, "budget_type", dataSourceId, "budgets.0.budget_type"),
					resource.TestCheckResourceAttrPair(resourceId, "metric", dataSourceId, "budgets.0.metric"),
					resource.TestCheckResourceAttrPair(resourceId, "quota", dataSourceId, "budgets.0.quota"),
				),
			},
		},
	})
}

func testAccCheckBssOpenApiBudgetDestroyWith(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		svc := BssOpenApiServiceV2{client}
		_, err := svc.DescribeBssOpenApiBudget(name)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("budget %s still exists", name)
	}
}

func testAccAlicloudBssOpenApiBudgetsDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_bss_open_api_budget" "default" {
  budget_name        = var.name
  budget_type        = "cost"
  metric             = "Cost"
  cycle_type         = "Month"
  cycle_start_period = "2026-09"
  cycle_end_period   = "2026-12"
  quota_type         = "quota"
  quota              = "80"
  comment            = "tf datasource test"
  cycle_quota {
    cycle_period = "2026-09"
    quota        = "80"
  }
  query_filter {
    code        = "productCode"
    select_type = "equal"
    values      = ["ECS"]
  }
  warn_confs {
    name             = "tf-warn"
    warn_target      = "Msc"
    threshold_type   = "percent"
    threshold_value  = "80"
    event_bridge     = false
    comment          = "tf warn conf"
  }
}

data "alicloud_bss_open_api_budgets" "default" {
  ids = [alicloud_bss_open_api_budget.default.budget_name]
}
`, name)
}
