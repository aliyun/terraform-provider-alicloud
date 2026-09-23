package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudActiontrailTrailsDataSource(t *testing.T) {
	resourceId := "data.alicloud_actiontrail_trails.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccactiontrail-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceActiontrailTrailsDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_actiontrail_trail.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_actiontrail_trail.default.id}-fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_actiontrail_trail.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_actiontrail_trail.default.id}-fake"},
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_actiontrail_trail.default.id}"},
			"status": "Disable",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_actiontrail_trail.default.id}"},
			"status": "Enable",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_actiontrail_trail.default.id}",
			"ids":        []string{"${alicloud_actiontrail_trail.default.id}"},
			"status":     "Disable",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_actiontrail_trail.default.id}-fake",
			"ids":        []string{"${alicloud_actiontrail_trail.default.id}"},
			"status":     "Enable",
		}),
	}
	var existActiontrailTrailMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"ids.0":                       CHECKSET,
			"trails.#":                    "1",
			"trails.0.event_rw":           "Write",
			"trails.0.oss_bucket_name":    CHECKSET,
			"trails.0.oss_key_prefix":     "",
			"trails.0.sls_project_arn":    "",
			"trails.0.sls_write_role_arn": "",
			"trails.0.status":             "Disable",
			"trails.0.id":                 CHECKSET,
			"trails.0.trail_name":         CHECKSET,
			"trails.0.trail_region":       "All",
		}
	}

	var fakeActiontrailTrailMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"trails.#": "0",
		}
	}

	var actionTrailsInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existActiontrailTrailMapFunc,
		fakeMapFunc:  fakeActiontrailTrailMapFunc,
	}

	actionTrailsInfo.dataSourceTestCheck(t, 0, nameRegexConf, idsConf, statusConf, allConf)
}

func dataSourceActiontrailTrailsDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_oss_bucket" "default" {
		bucket  =  "%[1]s"
	}

	resource "alicloud_actiontrail_trail" "default" {
	  trail_name = "%[1]s"
	  oss_bucket_name = alicloud_oss_bucket.default.id
	  status= "Disable"
	}`, name)
}
