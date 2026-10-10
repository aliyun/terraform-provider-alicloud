package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudPaiWorkspacePromptsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccAlicloudPaiWorkspacePromptsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_pai_workspace_prompt.default.id}"]`,
		}),
		fakeConfig: testAccAlicloudPaiWorkspacePromptsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_pai_workspace_prompt.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccAlicloudPaiWorkspacePromptsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_pai_workspace_prompt.default.prompt_name}"`,
		}),
		fakeConfig: testAccAlicloudPaiWorkspacePromptsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_pai_workspace_prompt.default.prompt_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccAlicloudPaiWorkspacePromptsDataSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_pai_workspace_prompt.default.id}"]`,
			"name_regex":     `"${alicloud_pai_workspace_prompt.default.prompt_name}"`,
			"framework_type": `"ICIO"`,
		}),
		fakeConfig: testAccAlicloudPaiWorkspacePromptsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_pai_workspace_prompt.default.id}_fake"]`,
			"name_regex": `"${alicloud_pai_workspace_prompt.default.prompt_name}_fake"`,
		}),
	}
	var existAlicloudPaiWorkspacePromptsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"names.#":                  "1",
			"prompts.#":                "1",
			"prompts.0.id":             CHECKSET,
			"prompts.0.prompt_id":      CHECKSET,
			"prompts.0.prompt_name":    fmt.Sprintf("tf-testaccprompt-%d", rand),
			"prompts.0.accessibility":  "PRIVATE",
			"prompts.0.framework_type": "ICIO",
			"prompts.0.description":    "datasource test prompt",
			"prompts.0.create_time":    CHECKSET,
			"prompts.0.modify_time":    CHECKSET,
		}
	}
	var fakeAlicloudPaiWorkspacePromptsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudPaiWorkspacePromptsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_pai_workspace_prompts.default",
		existMapFunc: existAlicloudPaiWorkspacePromptsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudPaiWorkspacePromptsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		testAccPreCheck(t)
	}
	alicloudPaiWorkspacePromptsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func testAccAlicloudPaiWorkspacePromptsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testaccprompt-%d"
}

resource "alicloud_pai_workspace_workspace" "default" {
  description    = "test_prompt_ds_%d"
  display_name   = "test prompt datasource"
  workspace_name = "prompt_ds_%d"
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_prompt" "default" {
  prompt_name        = var.name
  description        = "datasource test prompt"
  accessibility     = "PRIVATE"
  framework_content = "{\"PromptContext\":\"test\",\"Tags\":{\"key\":\"value\"}}"
  framework_type    = "ICIO"
  workspace_id      = alicloud_pai_workspace_workspace.default.id
}

data "alicloud_pai_workspace_prompts" "default" {
  workspace_id = alicloud_pai_workspace_workspace.default.id
  %s
}
`, rand, rand, rand, strings.Join(pairs, "\n  "))
	return config
}
