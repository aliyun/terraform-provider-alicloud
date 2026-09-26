package pai_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAlicloudPaiWorkspacePromptsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)

	existCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("data.alicloud_pai_workspace_prompts.default", "ids.#", "1"),
		resource.TestCheckResourceAttr("data.alicloud_pai_workspace_prompts.default", "names.#", "1"),
		resource.TestCheckResourceAttr("data.alicloud_pai_workspace_prompts.default", "prompts.#", "1"),
		resource.TestCheckResourceAttrSet("data.alicloud_pai_workspace_prompts.default", "prompts.0.id"),
		resource.TestCheckResourceAttrSet("data.alicloud_pai_workspace_prompts.default", "prompts.0.prompt_id"),
		resource.TestCheckResourceAttr("data.alicloud_pai_workspace_prompts.default", "prompts.0.prompt_name", fmt.Sprintf("tf-testAccprompt-%d", rand)),
		resource.TestCheckResourceAttr("data.alicloud_pai_workspace_prompts.default", "prompts.0.accessibility", "PRIVATE"),
		resource.TestCheckResourceAttr("data.alicloud_pai_workspace_prompts.default", "prompts.0.framework_type", "ICIO"),
		resource.TestCheckResourceAttr("data.alicloud_pai_workspace_prompts.default", "prompts.0.description", "datasource test prompt"),
		resource.TestCheckResourceAttrSet("data.alicloud_pai_workspace_prompts.default", "prompts.0.create_time"),
		resource.TestCheckResourceAttrSet("data.alicloud_pai_workspace_prompts.default", "prompts.0.modify_time"),
	)
	fakeCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("data.alicloud_pai_workspace_prompts.default", "ids.#", "0"),
		resource.TestCheckResourceAttr("data.alicloud_pai_workspace_prompts.default", "names.#", "0"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPaiWorkspacePromptsDataSourceConfig(rand, map[string]string{
					"ids": `["${alicloud_pai_workspace_prompt.default.id}"]`,
				}),
				Check: existCheck,
			},
			{
				Config: testAccPaiWorkspacePromptsDataSourceConfig(rand, map[string]string{
					"ids": `["${alicloud_pai_workspace_prompt.default.id}_fake"]`,
				}),
				Check: fakeCheck,
			},
			{
				Config: testAccPaiWorkspacePromptsDataSourceConfig(rand, map[string]string{
					"name_regex": `"${alicloud_pai_workspace_prompt.default.prompt_name}"`,
				}),
				Check: existCheck,
			},
			{
				Config: testAccPaiWorkspacePromptsDataSourceConfig(rand, map[string]string{
					"name_regex": `"${alicloud_pai_workspace_prompt.default.prompt_name}_fake"`,
				}),
				Check: fakeCheck,
			},
			{
				Config: testAccPaiWorkspacePromptsDataSourceConfig(rand, map[string]string{
					"ids":            `["${alicloud_pai_workspace_prompt.default.id}"]`,
					"name_regex":     `"${alicloud_pai_workspace_prompt.default.prompt_name}"`,
					"framework_type": `"ICIO"`,
				}),
				Check: existCheck,
			},
			{
				Config: testAccPaiWorkspacePromptsDataSourceConfig(rand, map[string]string{
					"ids":        `["${alicloud_pai_workspace_prompt.default.id}_fake"]`,
					"name_regex": `"${alicloud_pai_workspace_prompt.default.prompt_name}_fake"`,
				}),
				Check: fakeCheck,
			},
			{
				Config: testAccPaiWorkspacePromptsDataSourceConfig(rand, map[string]string{
					"ids":            `["${alicloud_pai_workspace_prompt.default.id}"]`,
					"framework_type": `"ICIO"`,
				}),
				Check: existCheck,
			},
			{
				Config: testAccPaiWorkspacePromptsDataSourceConfig(rand, map[string]string{
					"ids":            `["${alicloud_pai_workspace_prompt.default.id}"]`,
					"framework_type": `"CRISPE"`,
				}),
				Check: fakeCheck,
			},
		},
	})
}

func testAccPaiWorkspacePromptsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	return fmt.Sprintf(`%s
variable "name" {
  default = "tf-testAccprompt-%d"
}

resource "alicloud_pai_workspace_workspace" "default" {
  description    = "test_prompt_ds_%d"
  display_name   = "test prompt datasource"
  workspace_name = "tf_testAccWs%d"
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_prompt" "default" {
  prompt_name       = var.name
  description       = "datasource test prompt"
  accessibility     = "PRIVATE"
  framework_content = "{\"PromptContext\":\"test\",\"Tags\":{\"key\":\"value\"}}"
  framework_type    = "ICIO"
  workspace_id      = alicloud_pai_workspace_workspace.default.id
}

data "alicloud_pai_workspace_prompts" "default" {
  workspace_id = alicloud_pai_workspace_workspace.default.id
  %s
}
`, testAccPaiWorkspaceRegionProvider, rand, rand, rand, strings.Join(pairs, "\n  "))
}
