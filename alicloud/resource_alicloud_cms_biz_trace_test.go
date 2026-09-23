package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCmsBizTraceBasic(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	var v map[string]interface{}
	resourceId := "alicloud_cms_biz_trace.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsBizTraceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsBizTrace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsBizTraceBasicDependence)
	ruleConfig := `[{"entrancePid":"xxxxx@b57c44xx6e86","rpcMatcher":{"matchType":"EQUALS","pattern":"/createApp"},"characteristics":{"operation":"AND","rules":[{"target":"CUSTOM_EXTRACT","matcher":{"matchType":"CONTAINS","pattern":[]}}]}}]`
	ruleConfigUpdated := `[{"entrancePid":"xxxxx@b57c44xx6e86","rpcMatcher":{"matchType":"EQUALS","pattern":"/updateApp"},"characteristics":{"operation":"OR","rules":[{"target":"CUSTOM_EXTRACT","matcher":{"matchType":"CONTAINS","pattern":[]}}]}}]`
	advancedConfig := `{"sample":{"strategy":"BY_APP"}}`
	advancedConfigUpdated := `{"sample":{"strategy":"BY_TRACE"}}`
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"biz_trace_code":  fmt.Sprintf("tf_acc_biztrace_%d", rand),
					"biz_trace_name":  name,
					"rule_config":     strings.ReplaceAll(ruleConfig, `"`, `\"`),
					"advanced_config": strings.ReplaceAll(advancedConfig, `"`, `\"`),
					"workspace":       "${alicloud_cms_workspace.default.workspace_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"biz_trace_code":  fmt.Sprintf("tf_acc_biztrace_%d", rand),
						"biz_trace_name":  name,
						"rule_config":     REGEXMATCH + ".*createApp.*",
						"advanced_config": REGEXMATCH + ".*BY_APP.*",
						"workspace":       name,
						"biz_trace_id":    CHECKSET,
						"create_time":     CHECKSET,
						"region_id":       CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"biz_trace_name":  fmt.Sprintf("%s-updated", name),
					"rule_config":     strings.ReplaceAll(ruleConfigUpdated, `"`, `\"`),
					"advanced_config": strings.ReplaceAll(advancedConfigUpdated, `"`, `\"`),
					"workspace":       "${alicloud_cms_workspace.default.workspace_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"biz_trace_name":  fmt.Sprintf("%s-updated", name),
						"rule_config":     REGEXMATCH + ".*updateApp.*",
						"advanced_config": REGEXMATCH + ".*BY_TRACE.*",
						"workspace":       name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudCmsBizTraceMap = map[string]string{
	"biz_trace_id": CHECKSET,
	"create_time":  CHECKSET,
	"region_id":    CHECKSET,
}

func AliCloudCmsBizTraceBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
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
`, name)
}
