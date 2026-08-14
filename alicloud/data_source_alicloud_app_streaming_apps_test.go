package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudAppStreamingAppsDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAppStreamingAppsSourceConfig(rand, map[string]string{}),
		fakeConfig: testAccCheckAlicloudAppStreamingAppsSourceConfig(rand, map[string]string{
			"ids": `["ca-tf-acc-fake-app-id"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAppStreamingAppsSourceConfig(rand, map[string]string{
			"name_regex": `".*"`,
		}),
		fakeConfig: testAccCheckAlicloudAppStreamingAppsSourceConfig(rand, map[string]string{
			"name_regex": `"tf-acc-never-match-x9z$"`,
		}),
	}

	outputFileConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAppStreamingAppsSourceConfig(rand, map[string]string{
			"output_file": `"app_streaming_apps.json"`,
		}),
	}

	var existAppStreamingAppsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":  CHECKSET,
			"apps.#": CHECKSET,
		}
	}
	var fakeAppStreamingAppsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":  "0",
			"apps.#": "0",
		}
	}
	var appStreamingAppsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_app_streaming_apps.default",
		existMapFunc: existAppStreamingAppsMapFunc,
		fakeMapFunc:  fakeAppStreamingAppsMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.AppStreamingSupportRegions)
	}
	appStreamingAppsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, outputFileConf)
}

func testAccCheckAlicloudAppStreamingAppsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccAppStreamingApp%d"
}

data "alicloud_app_streaming_apps" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
