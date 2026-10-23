// Package cms. This file is generated automatically. Please do not modify it manually, thank you!
// Nested attributes require Terraform Plugin Protocol 6.
package cms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// TestAccAliCloudCmsBizTrace_basic creates a dedicated workspace and log project
// dependence, then exercises create/update/import with a random biz_trace_code.
func TestAccAliCloudCmsBizTrace_basic(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsBizTrace_basic_step0(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_biz_trace.default", "biz_trace_code", fmt.Sprintf("tf_acc_biztrace_%d", r)),
					resource.TestCheckResourceAttr("alicloud_cms_biz_trace.default", "biz_trace_name", fmt.Sprintf("tfacccms%d", r)),
					resource.TestCheckResourceAttrSet("alicloud_cms_biz_trace.default", "rule_config"),
					resource.TestCheckResourceAttrSet("alicloud_cms_biz_trace.default", "advanced_config"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_biz_trace.default", tfjsonpath.New("create_time"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_biz_trace.default", tfjsonpath.New("region_id"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsBizTrace_basic_step1(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_biz_trace.default", "biz_trace_code", fmt.Sprintf("tf_acc_biztrace_%d", r)),
					resource.TestCheckResourceAttr("alicloud_cms_biz_trace.default", "biz_trace_name", fmt.Sprintf("tfacccms%d-updated", r)),
					resource.TestCheckResourceAttrSet("alicloud_cms_biz_trace.default", "rule_config"),
					resource.TestCheckResourceAttrSet("alicloud_cms_biz_trace.default", "advanced_config"),
				),
			},
			{
				ResourceName:      "alicloud_cms_biz_trace.default",
				ImportState:       true,
				ImportStateVerify: true,
				// rule_config is a JSON document the CMS server enriches on read with
				// server-resolved keys (callType, entranceServiceId, dyingScope); the
				// configured value is the canonical subset. ImportStateVerify does an
				// exact attribute comparison, but the imported (server-enriched) value
				// is a superset that can never equal the configured subset, so the field
				// is excluded here. fwJSONSubsetSuppress reconciles it on plan/refresh.
				ImportStateVerifyIgnore: []string{"rule_config"},
			},
		},
	})
}

func testAccAliCloudCmsBizTrace_basic_step0(r int) string {
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

func testAccAliCloudCmsBizTrace_basic_step1(r int) string {
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
  biz_trace_name  = "${var.name}-updated"
  rule_config     = "[{\"entrancePid\":\"xxxxx@b57c44xx6e86\",\"rpcMatcher\":{\"matchType\":\"EQUALS\",\"pattern\":\"/updateApp\"},\"characteristics\":{\"operation\":\"OR\",\"rules\":[{\"target\":\"CUSTOM_EXTRACT\",\"matcher\":{\"matchType\":\"CONTAINS\",\"pattern\":[]}}]}}]"
  advanced_config = "{\"sample\":{\"strategy\":\"BY_TRACE\"}}"
  workspace       = alicloud_cms_workspace.default.workspace_name
}
`, r)
}
