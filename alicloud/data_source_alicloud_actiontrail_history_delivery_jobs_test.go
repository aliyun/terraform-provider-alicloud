package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudActiontrailHistoryDeliveryJobsDataSource(t *testing.T) {
	resourceId := "data.alicloud_actiontrail_history_delivery_jobs.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccactiontrail-%d", rand)
	checkoutSupportedRegions(t, true, connectivity.ActiontrailSupportRegions)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceActiontrailHistoryDeliveryJobsDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_actiontrail_history_delivery_job.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_actiontrail_history_delivery_job.default.id}-fake"},
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_actiontrail_history_delivery_job.default.id}"},
			"status":         `2`,
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_actiontrail_history_delivery_job.default.id}"},
			"status": `0`,
		}),
	}
	var existActiontrailHistoryDeliveryJobMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":             "1",
			"ids.0":             CHECKSET,
			"jobs.#":            "1",
			"jobs.0.trail_name": fmt.Sprintf("tf-testaccactiontrail-%d", rand),
		}
	}

	var fakeActiontrailHistoryDeliveryJobMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":  "0",
			"jobs.#": "0",
		}
	}

	var actionHistoryDeliveryJobsInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existActiontrailHistoryDeliveryJobMapFunc,
		fakeMapFunc:  fakeActiontrailHistoryDeliveryJobMapFunc,
	}

	actionHistoryDeliveryJobsInfo.dataSourceTestCheck(t, rand, idsConf, statusConf)
}

func dataSourceActiontrailHistoryDeliveryJobsDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%v"
	}
	data "alicloud_regions" "default" {
	  current = true
	}
	
	data "alicloud_account" "default" {}
	
	resource "alicloud_log_project" "default" {
	  name = var.name
	  description = "tf actiontrail test"
	}

	resource "alicloud_actiontrail_trail" "default" {
	  trail_name = var.name
	  sls_project_arn = "acs:log:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:project/${alicloud_log_project.default.name}"
	}

	resource "alicloud_actiontrail_history_delivery_job" "default" {
	  trail_name = alicloud_actiontrail_trail.default.name
	}`, name)
}
