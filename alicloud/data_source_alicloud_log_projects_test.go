package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudLogProjectsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudLogProjectsDataSourceName(rand, map[string]string{
			"ids": `["terraform"]`,
		}),
		fakeConfig: testAccCheckAlicloudLogProjectsDataSourceName(rand, map[string]string{
			"ids": `["fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudLogProjectsDataSourceName(rand, map[string]string{
			"name_regex": `"terraform"`,
		}),
		fakeConfig: testAccCheckAlicloudLogProjectsDataSourceName(rand, map[string]string{
			"name_regex": `"fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudLogProjectsDataSourceName(rand, map[string]string{
			"ids":    `["terraform"]`,
			"status": `"Normal"`,
		}),
		fakeConfig: testAccCheckAlicloudLogProjectsDataSourceName(rand, map[string]string{
			"ids":    `["terraform"]`,
			"status": `"Disable"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudLogProjectsDataSourceName(rand, map[string]string{
			"ids":        `["terraform"]`,
			"name_regex": `"terraform"`,
			"status":     `"Normal"`,
		}),
		fakeConfig: testAccCheckAlicloudLogProjectsDataSourceName(rand, map[string]string{
			"ids":        `["fake"]`,
			"name_regex": `"fake"`,
			"status":     `"Disable"`,
		}),
	}
	var existAlicloudLogProjectsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"names.#":                     "1",
			"projects.#":                  "1",
			"projects.0.id":               CHECKSET,
			"projects.0.description":      CHECKSET,
			"projects.0.project_name":     "terraform",
			"projects.0.region":           CHECKSET,
			"projects.0.owner":            CHECKSET,
			"projects.0.last_modify_time": CHECKSET,
			"projects.0.status":           "Normal",
		}
	}
	var fakeAlicloudLogProjectsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudLogProjectsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_log_projects.default",
		existMapFunc: existAlicloudLogProjectsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudLogProjectsDataSourceNameMapFunc,
	}
	// After the project is successfully created, ListProject may not be able to query the newly created project, so here only the projects existing under the account are tested
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithRegions(t, true, connectivity.LogProjectsRegions)
	}
	alicloudLogProjectsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudLogProjectsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
data "alicloud_log_projects" "default" {
	%s	
}
`, strings.Join(pairs, " \n "))
	return config
}
