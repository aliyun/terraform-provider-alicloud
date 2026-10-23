// Package cms. Framework-native data source test for alicloud_cms_biz_traces.
package cms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCmsBizTracesDataSource_basic(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCmsBizTracesDataSource_exist(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cms_biz_traces.default", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_cms_biz_traces.default", "biz_traces.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_biz_traces.default", "biz_traces.0.biz_trace_id"),
					resource.TestCheckResourceAttr("data.alicloud_cms_biz_traces.default", "biz_traces.0.biz_trace_code", fmt.Sprintf("tf_acc_biztrace_%d", r)),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_biz_traces.default", "biz_traces.0.rule_config"),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_biz_traces.default", "biz_traces.0.advanced_config"),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_biz_traces.default", "biz_traces.0.workspace"),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_biz_traces.default", "biz_traces.0.create_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_biz_traces.default", "biz_traces.0.region_id"),
				),
			},
			{
				Config: testAccCheckAlicloudCmsBizTracesDataSource_fake(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cms_biz_traces.default", "ids.#", "0"),
				),
			},
		},
	})
}

func bizTracesDataSourceDependence(r int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tfacccms%[1]d"
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
  biz_trace_code  = "tf_acc_biztrace_%[1]d"
  biz_trace_name  = var.name
  rule_config     = "[{\"entrancePid\":\"xxxxx@b57c44xx6e86\",\"rpcMatcher\":{\"matchType\":\"EQUALS\",\"pattern\":\"/createApp\"},\"characteristics\":{\"operation\":\"AND\",\"rules\":[{\"target\":\"CUSTOM_EXTRACT\",\"matcher\":{\"matchType\":\"CONTAINS\",\"pattern\":[]}}]}}]"
  advanced_config = "{\"sample\":{\"strategy\":\"BY_APP\"}}"
  workspace       = alicloud_cms_workspace.default.workspace_name
}
`, r)
}

func testAccCheckAlicloudCmsBizTracesDataSource_exist(r int) string {
	return fmt.Sprintf(`
%s

data "alicloud_cms_biz_traces" "default" {
  ids = ["${alicloud_cms_biz_trace.default.id}"]
}
`, bizTracesDataSourceDependence(r))
}

func testAccCheckAlicloudCmsBizTracesDataSource_fake(r int) string {
	return fmt.Sprintf(`
%s

data "alicloud_cms_biz_traces" "default" {
  ids = ["${alicloud_cms_biz_trace.default.id}_fake"]
}
`, bizTracesDataSourceDependence(r))
}
