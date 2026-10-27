// Package actiontrail_test.
package actiontrail_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudActiontrailTrailsDataSource(t *testing.T) {
	resourceId := "data.alicloud_actiontrail_trails.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccactiontrail-%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceActiontrailTrailsDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_actiontrail_trail.default.id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_actiontrail_trail.default.id}-fake",
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_actiontrail_trail.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_actiontrail_trail.default.id}-fake"},
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_actiontrail_trail.default.id}"},
			"status": "Disable",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_actiontrail_trail.default.id}"},
			"status": "Enable",
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_actiontrail_trail.default.id}",
			"ids":        []string{"${alicloud_actiontrail_trail.default.id}"},
			"status":     "Disable",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_actiontrail_trail.default.id}-fake",
			"ids":        []string{"${alicloud_actiontrail_trail.default.id}"},
			"status":     "Enable",
		}),
	}
	var existActiontrailTrailMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"ids.0":                       acctest.CHECKSET,
			"trails.#":                    "1",
			"trails.0.event_rw":           "Write",
			"trails.0.oss_bucket_name":    acctest.CHECKSET,
			"trails.0.oss_key_prefix":     "",
			"trails.0.sls_project_arn":    "",
			"trails.0.sls_write_role_arn": "",
			"trails.0.status":             "Disable",
			"trails.0.id":                 acctest.CHECKSET,
			"trails.0.trail_name":         acctest.CHECKSET,
			"trails.0.trail_region":       "All",
		}
	}

	var fakeActiontrailTrailMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"trails.#": "0",
		}
	}

	var actionTrailsInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existActiontrailTrailMapFunc,
		FakeMapFunc:  fakeActiontrailTrailMapFunc,
	}

	actionTrailsInfo.DataSourceTestCheck(t, 0, nameRegexConf, idsConf, statusConf, allConf)
}

func dataSourceActiontrailTrailsDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_oss_bucket" "default" {
		bucket  =  "%[1]s"
	}

	resource "alicloud_actiontrail_trail" "default" {
	  trail_name = "%[1]s"
	  oss_bucket_name = alicloud_oss_bucket.default.id
	  event_rw = "Write"
	  status= "Disable"
	}`, name)
}
