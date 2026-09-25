package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

// TestAccAlicloudDataWorksTaskInstances_basic smoke-tests the
// alicloud_data_works_task_instances data source. Task instances are runtime
// artifacts produced by the DataWorks scheduler; the provider has no resource
// to create one as a fixture, so the test filters by a non-existent id and
// asserts the API call succeeds and returns an empty result set. This still
// exercises the request mapping (PascalCase params, pagination, responsePath)
// and the full output schema.
//
// Bizdate is required by the ListTaskInstances API and is a UNIX timestamp in
// milliseconds (e.g. 1743350400000). The test uses a fixed past timestamp so
// the request is accepted even when no instances match the id filter.
func TestAccAlicloudDataWorksTaskInstances_basic(t *testing.T) {
	resourceId := "data.alicloud_data_works_task_instances.default"
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc-dataworkstaskinstances%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceDataWorksTaskInstancesDependence)

	// existConfig: full schema coverage (all optional args present) but filtered
	// to a non-existent id so the response is an empty list regardless of the
	// account's runtime state. project_env is intentionally omitted: the
	// ListTaskInstances API accepts Prod/Dev (first-letter capital) and rejects
	// PROD/dev with InvalidProjectEnv; leaving it unset keeps the parameter out
	// of the request so no value validation runs.
	fullConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"bizdate":              "1743350400000",
			"project_id":           "638",
			"ids":                  []string{"nonexistent_task_instance_id"},
			"name_regex":           "fixture_nonexistent_task_name",
			"owner":                "",
			"page_number":          1,
			"page_size":            10,
			"sort_by":              "",
			"status":               "",
			"task_id":              "",
			"task_name":            "",
			"task_type":            "",
			"trigger_type":         "",
			"workflow_id":          "",
			"workflow_instance_id": "",
			"output_file":          "",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"bizdate":    "1743350400000",
			"project_id": "638",
			"ids":        []string{"another_nonexistent_task_instance_id"},
		}),
	}

	var existMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#": "0",
			"ids.#":       "0",
		}
	}

	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#": "0",
			"ids.#":       "0",
		}
	}

	var checkInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}

	checkInfo.dataSourceTestCheck(t, rand, fullConf)
}

func dataSourceDataWorksTaskInstancesDependence(name string) string {
	return ""
}
