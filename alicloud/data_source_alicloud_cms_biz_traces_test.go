package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCmsBizTracesDataSource_basic(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(10000, 99999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsBizTracesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_biz_trace.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCmsBizTracesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_biz_trace.default.id}_fake"]`,
		}),
	}
	var existMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"biz_traces.#":                 "1",
			"biz_traces.0.biz_trace_id":    CHECKSET,
			"biz_traces.0.biz_trace_code":  fmt.Sprintf("tf_acc_biztrace_%d", rand),
			"biz_traces.0.biz_trace_name":  fmt.Sprintf("tfacccms%d", rand),
			"biz_traces.0.rule_config":     REGEXMATCH + ".*createApp.*",
			"biz_traces.0.advanced_config": REGEXMATCH + ".*BY_APP.*",
			"biz_traces.0.workspace":       fmt.Sprintf("tfacccms%d", rand),
			"biz_traces.0.create_time":     CHECKSET,
			"biz_traces.0.region_id":       CHECKSET,
		}
	}
	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var checkInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cms_biz_traces.default",
		existMapFunc: existMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	checkInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)
}

func testAccCheckAlicloudCmsBizTracesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tfacccms%d"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  sls_project    = alicloud_log_project.default.project_name
  workspace_name = var.name
}

resource "alicloud_cms_biz_trace" "default" {
  biz_trace_code  = "tf_acc_biztrace_%d"
  biz_trace_name  = var.name
  rule_config     = "[{\"entrancePid\":\"xxxxx@b57c44xx6e86\",\"rpcMatcher\":{\"matchType\":\"EQUALS\",\"pattern\":\"/createApp\"},\"characteristics\":{\"operation\":\"AND\",\"rules\":[{\"target\":\"CUSTOM_EXTRACT\",\"matcher\":{\"matchType\":\"CONTAINS\",\"pattern\":[]}}]}}]"
  advanced_config = "{\"sample\":{\"strategy\":\"BY_APP\"}}"
  workspace       = alicloud_cms_workspace.default.workspace_name
}

data "alicloud_cms_biz_traces" "default" {
  %s
}
`, rand, rand, strings.Join(pairs, " \n "))
	return config
}
