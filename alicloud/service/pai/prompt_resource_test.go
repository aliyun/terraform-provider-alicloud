package pai_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The PAI Workspace APIs used here are only offered in cn-hangzhou. Each test
// config pins provider "alicloud" { region = "cn-hangzhou" } so the cases run
// there regardless of the ALICLOUD_REGION the runner injects; the deleted
// root-package tests pinned the same region in their configs.
const testAccPaiWorkspaceRegionProvider = `
provider "alicloud" {
  region = "cn-hangzhou"
}
`

// The dependency workspace name keeps the letters/digits/underscore form of
// the tf-testAcc prefix: CreateWorkspace rejects hyphens and caps the name at
// 23 characters, so the hyphenated spelling cannot be used for that resource.

// Case 提示词测试用例 10621: create with one framework type, update to the
// other, then import. depRand stays a function-scoped variable so the update
// step reuses the same dependence workspace instead of regenerating it.
func TestAccAliCloudPaiWorkspacePrompt_basic10621(t *testing.T) {
	rand := tfacctest.RandIntRange(1, 999)
	depRand := tfacctest.RandIntRange(10000, 99999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPaiWorkspacePromptBasic10621Create(rand, depRand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_prompt.default", "id"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_prompt.default", "prompt_name", fmt.Sprintf("tf-testAccprompt-%d", rand)),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_prompt.default", "description", "这是一个提示词测试模版"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_prompt.default", "accessibility", "PRIVATE"),
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_prompt.default", "framework_content"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_prompt.default", "framework_type", "CRISPE"),
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_prompt.default", "workspace_id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_prompt.default", "prompt_id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_prompt.default", "create_time"),
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_prompt.default", "modify_time"),
				),
			},
			{
				Config: testAccPaiWorkspacePromptBasic10621Update(rand, depRand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_workspace_prompt.default", "prompt_name", fmt.Sprintf("tf-testAccprompt-%d", rand)),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_prompt.default", "description", "测试修改描述"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_prompt.default", "accessibility", "PRIVATE"),
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_prompt.default", "framework_content"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_prompt.default", "framework_type", "ICIO"),
				),
			},
			{
				// UpdatePrompt ignores JSON null and empty strings, so optional
				// fields cannot be cleared. Removing description from the config
				// must keep the prior value and leave no pending diff.
				Config: testAccPaiWorkspacePromptBasic10621RemoveDescription(rand, depRand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_workspace_prompt.default", "description", "测试修改描述"),
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_prompt.default", "framework_content"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_prompt.default", "framework_type", "ICIO"),
				),
			},
			{
				ResourceName:      "alicloud_pai_workspace_prompt.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Case 提示词测试用例_副本1747970886918 10831: same flow in the reverse
// framework-type order.
func TestAccAliCloudPaiWorkspacePrompt_basic10831(t *testing.T) {
	rand := tfacctest.RandIntRange(1, 999)
	depRand := tfacctest.RandIntRange(10000, 99999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPaiWorkspacePromptBasic10831Create(rand, depRand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_prompt.default", "id"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_prompt.default", "prompt_name", fmt.Sprintf("tf-testAccprompt-%d", rand)),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_prompt.default", "description", "这是一个提示词测试模版"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_prompt.default", "accessibility", "PRIVATE"),
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_prompt.default", "framework_content"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_prompt.default", "framework_type", "ICIO"),
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_prompt.default", "workspace_id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_prompt.default", "prompt_id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_prompt.default", "create_time"),
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_prompt.default", "modify_time"),
				),
			},
			{
				Config: testAccPaiWorkspacePromptBasic10831Update(rand, depRand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_workspace_prompt.default", "prompt_name", fmt.Sprintf("tf-testAccprompt-%d", rand)),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_prompt.default", "description", "测试修改描述"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_prompt.default", "accessibility", "PRIVATE"),
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_prompt.default", "framework_content"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_prompt.default", "framework_type", "CRISPE"),
				),
			},
			{
				ResourceName:      "alicloud_pai_workspace_prompt.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// The values carry backslash-escaped quotes so that, once substituted into
// framework_content = "<value>" in the HCL configs, they parse as the JSON
// strings the API expects; a bare quote would break HCL parsing at plan time.
const testAccPaiWorkspacePromptLongContentCreate = `{     \"PromptContext\":\"你是一个拥有十年驾龄的老司机，请你针对以下图片场景做出你的分析判断。\",     \"Tags\":{     \"侧翻的车辆\":\"车辆侧翻在地,4个车轮至少有两个离开地面\",     \"匝道\":\"只有明确看到高速路上的大弯道，一般匝道都在高速路干道的右侧，进出收费站才可判定存在。\"     } }`

const testAccPaiWorkspacePromptLongContentUpdate = `{     \"PromptContext\":\"你老司机，请你针对以下图片场景做出你的分析判断。\",     \"Tags\":{     \"侧翻的车辆\":\"车辆侧翻在地,4个车轮至少有两个离开地面\",     \"匝道\":\"只有明确看到高速路上的大弯道，一般匝道都在高速路干道的右侧，进出收费站才可判定存在。\"     } }`

func testAccPaiWorkspacePromptBasic10621Create(rand, depRand int) string {
	return fmt.Sprintf(`%[1]s
resource "alicloud_pai_workspace_workspace" "default7Adve7" {
  description    = "test_prompt_%[3]d"
  display_name   = "用来测试提示词"
  workspace_name = "tf_testAccWs%[3]d"
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_prompt" "default" {
  prompt_name       = "tf-testAccprompt-%[2]d"
  description       = "这是一个提示词测试模版"
  accessibility     = "PRIVATE"
  framework_content = "%[4]s"
  framework_type    = "CRISPE"
  workspace_id      = alicloud_pai_workspace_workspace.default7Adve7.id
}
`, testAccPaiWorkspaceRegionProvider, rand, depRand, testAccPaiWorkspacePromptLongContentCreate)
}

func testAccPaiWorkspacePromptBasic10621Update(rand, depRand int) string {
	return fmt.Sprintf(`%[1]s
resource "alicloud_pai_workspace_workspace" "default7Adve7" {
  description    = "test_prompt_%[3]d"
  display_name   = "用来测试提示词"
  workspace_name = "tf_testAccWs%[3]d"
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_prompt" "default" {
  prompt_name       = "tf-testAccprompt-%[2]d"
  description       = "测试修改描述"
  accessibility     = "PRIVATE"
  framework_content = "%[4]s"
  framework_type    = "ICIO"
  workspace_id      = alicloud_pai_workspace_workspace.default7Adve7.id
}
`, testAccPaiWorkspaceRegionProvider, rand, depRand, testAccPaiWorkspacePromptLongContentUpdate)
}

func testAccPaiWorkspacePromptBasic10621RemoveDescription(rand, depRand int) string {
	return fmt.Sprintf(`%[1]s
resource "alicloud_pai_workspace_workspace" "default7Adve7" {
  description    = "test_prompt_%[3]d"
  display_name   = "用来测试提示词"
  workspace_name = "tf_testAccWs%[3]d"
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_prompt" "default" {
  prompt_name       = "tf-testAccprompt-%[2]d"
  accessibility     = "PRIVATE"
  framework_content = "%[4]s"
  framework_type    = "ICIO"
  workspace_id      = alicloud_pai_workspace_workspace.default7Adve7.id
}
`, testAccPaiWorkspaceRegionProvider, rand, depRand, testAccPaiWorkspacePromptLongContentUpdate)
}

func testAccPaiWorkspacePromptBasic10831Create(rand, depRand int) string {
	return fmt.Sprintf(`%[1]s
resource "alicloud_pai_workspace_workspace" "default7Adve7" {
  description    = "test_prompt_%[3]d"
  display_name   = "用来测试提示词"
  workspace_name = "tf_testAccWs%[3]d"
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_prompt" "default" {
  prompt_name       = "tf-testAccprompt-%[2]d"
  description       = "这是一个提示词测试模版"
  accessibility     = "PRIVATE"
  framework_content = "%[4]s"
  framework_type    = "ICIO"
  workspace_id      = alicloud_pai_workspace_workspace.default7Adve7.id
}
`, testAccPaiWorkspaceRegionProvider, rand, depRand, testAccPaiWorkspacePromptLongContentCreate)
}

func testAccPaiWorkspacePromptBasic10831Update(rand, depRand int) string {
	return fmt.Sprintf(`%[1]s
resource "alicloud_pai_workspace_workspace" "default7Adve7" {
  description    = "test_prompt_%[3]d"
  display_name   = "用来测试提示词"
  workspace_name = "tf_testAccWs%[3]d"
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_prompt" "default" {
  prompt_name       = "tf-testAccprompt-%[2]d"
  description       = "测试修改描述"
  accessibility     = "PRIVATE"
  framework_content = "%[4]s"
  framework_type    = "CRISPE"
  workspace_id      = alicloud_pai_workspace_workspace.default7Adve7.id
}
`, testAccPaiWorkspaceRegionProvider, rand, depRand, testAccPaiWorkspacePromptLongContentUpdate)
}
