// Package actiontrail_test.
package actiontrail_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

// TestAccAlicloudActiontrailHistoryDeliveryJobsDataSource is the real ACC for
// the alicloud_actiontrail_history_delivery_jobs data source (previously a
// Skip-prefixed test that never executed). It lists delivery history jobs via
// ListDeliveryHistoryJobs and supports ids / status / enable_details filters.
//
// The id-based query (idsConf) is the deterministic representative ACC: the
// freshly-created job is matched by id with enable_details, and a
// non-matching id yields an empty result. The status filter is exercised only
// in its exclusion direction (statusConf Fake with status=0): after
// CreateDeliveryHistoryJob the resource waits for JobStatus in {1, 2, 3}, so
// the job is never status 0 and the filter reliably excludes it. The
// inclusion direction (status=2) is intentionally not asserted: the
// post-Create status is non-deterministic (the resource's state refresh
// targets {1, 2, 3}, and resource_history_delivery_job_test.go asserts
// status=CHECKSET rather than a fixed value), so a status=2 inclusion
// assertion would be flaky; SKIP/FAIL would not count as PASS per AC2.
//
// The trail precondition aligns with the proven
// TestAccAliCloudActionTrailHistoryDeliveryJob_basic12075 resource test: it
// uses the AliyunServiceRoleForActionTrail service-linked role
// (sls_write_role_arn) so the trail can deliver to SLS, and runs in
// cn-hangzhou where that role and the history delivery job are exercised.
func TestAccAlicloudActiontrailHistoryDeliveryJobsDataSource(t *testing.T) {
	resourceId := "data.alicloud_actiontrail_history_delivery_jobs.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccactiontrail-%d", rand)
	acctest.CheckoutSupportedRegions(t, true, connectivity.ActiontrailSupportRegions)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceActiontrailHistoryDeliveryJobsDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_actiontrail_history_delivery_job.default.id}"},
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_actiontrail_history_delivery_job.default.id}-fake"},
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		// status=0 reliably excludes the created job: the resource waits for
		// JobStatus in {1, 2, 3} after CreateDeliveryHistoryJob, so the job is
		// never status 0 and the filter yields an empty result. The inclusion
		// direction (status=2) is not asserted because the post-Create status
		// is non-deterministic.
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_actiontrail_history_delivery_job.default.id}"},
			"status": `0`,
		}),
	}
	var existActiontrailHistoryDeliveryJobMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":             "1",
			"ids.0":             acctest.CHECKSET,
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

	var actionHistoryDeliveryJobsInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existActiontrailHistoryDeliveryJobMapFunc,
		FakeMapFunc:  fakeActiontrailHistoryDeliveryJobMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		acctest.PreCheck(t)
	}
	actionHistoryDeliveryJobsInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf)
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

	data "alicloud_ram_roles" "default" {
	  name_regex = "AliyunServiceRoleForActionTrail"
	}

	resource "alicloud_log_project" "default" {
	  name = var.name
	  description = "tf actiontrail test"
	}

	resource "alicloud_actiontrail_trail" "default" {
	  trail_name = var.name
	  sls_project_arn = "acs:log:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:project/${alicloud_log_project.default.name}"
	  sls_write_role_arn = data.alicloud_ram_roles.default.roles.0.arn
	  status = "Enable"
	}

	resource "alicloud_actiontrail_history_delivery_job" "default" {
	  trail_name = alicloud_actiontrail_trail.default.name
	}`, name)
}
