// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test PaiWorkspace Prompt. >>> Resource test cases, automatically generated.
// Case 提示词测试用例 10621
func TestAccAliCloudPaiWorkspacePrompt_basic10621(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_prompt.default"
	ra := resourceAttrInit(resourceId, AlicloudPaiWorkspacePromptMap10621)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspacePrompt")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPaiWorkspacePromptBasicDependence10621)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"prompt_name":       name,
					"description":       "这是一个提示词测试模版",
					"accessibility":     "PRIVATE",
					"framework_content": "{     \\\"PromptContext\\\":\\\"你是一个拥有十年驾龄的老司机，请你针对以下图片场景做出你的分析判断。\\\",     \\\"Tags\\\":{     \\\"侧翻的车辆\\\":\\\"车辆侧翻在地,4个车轮至少有两个离开地面\\\",     \\\"匝道\\\":\\\"只有明确看到高速路上的大弯道，一般匝道都在高速路干道的右侧，进出收费站才可判定存在。\\\"     } }",
					"framework_type":    "CRISPE",
					"workspace_id":      "${alicloud_pai_workspace_workspace.default7Adve7.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prompt_name":       name,
						"description":       "这是一个提示词测试模版",
						"accessibility":     "PRIVATE",
						"framework_content": CHECKSET,
						"framework_type":    "CRISPE",
						"workspace_id":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "测试修改描述",
					"framework_content": "{     \\\"PromptContext\\\":\\\"你老司机，请你针对以下图片场景做出你的分析判断。\\\",     \\\"Tags\\\":{     \\\"侧翻的车辆\\\":\\\"车辆侧翻在地,4个车轮至少有两个离开地面\\\",     \\\"匝道\\\":\\\"只有明确看到高速路上的大弯道，一般匝道都在高速路干道的右侧，进出收费站才可判定存在。\\\"     } }",
					"framework_type":    "ICIO",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "测试修改描述",
						"framework_content": CHECKSET,
						"framework_type":    "ICIO",
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

var AlicloudPaiWorkspacePromptMap10621 = map[string]string{
	"modify_time": CHECKSET,
	"create_time": CHECKSET,
	"prompt_id":   CHECKSET,
}

func AlicloudPaiWorkspacePromptBasicDependence10621(name string) string {
	rand := acctest.RandIntRange(10000, 99999)
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_pai_workspace_workspace" "default7Adve7" {
  description    = "test_prompt_%d"
  display_name   = "用来测试提示词"
  workspace_name = "prompt_%d"
  env_types      = ["prod"]
}


`, name, rand, rand)
}

// Case 提示词测试用例_副本1747970886918 10831
func TestAccAliCloudPaiWorkspacePrompt_basic10831(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_prompt.default"
	ra := resourceAttrInit(resourceId, AlicloudPaiWorkspacePromptMap10831)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspacePrompt")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPaiWorkspacePromptBasicDependence10831)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"prompt_name":       name,
					"description":       "这是一个提示词测试模版",
					"accessibility":     "PRIVATE",
					"framework_content": "{     \\\"PromptContext\\\":\\\"你是一个拥有十年驾龄的老司机，请你针对以下图片场景做出你的分析判断。\\\",     \\\"Tags\\\":{     \\\"侧翻的车辆\\\":\\\"车辆侧翻在地,4个车轮至少有两个离开地面\\\",     \\\"匝道\\\":\\\"只有明确看到高速路上的大弯道，一般匝道都在高速路干道的右侧，进出收费站才可判定存在。\\\"     } }",
					"framework_type":    "ICIO",
					"workspace_id":      "${alicloud_pai_workspace_workspace.default7Adve7.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prompt_name":       name,
						"description":       "这是一个提示词测试模版",
						"accessibility":     "PRIVATE",
						"framework_content": CHECKSET,
						"framework_type":    "ICIO",
						"workspace_id":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "测试修改描述",
					"framework_content": "{     \\\"PromptContext\\\":\\\"你老司机，请你针对以下图片场景做出你的分析判断。\\\",     \\\"Tags\\\":{     \\\"侧翻的车辆\\\":\\\"车辆侧翻在地,4个车轮至少有两个离开地面\\\",     \\\"匝道\\\":\\\"只有明确看到高速路上的大弯道，一般匝道都在高速路干道的右侧，进出收费站才可判定存在。\\\"     } }",
					"framework_type":    "CRISPE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "测试修改描述",
						"framework_content": CHECKSET,
						"framework_type":    "CRISPE",
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

var AlicloudPaiWorkspacePromptMap10831 = map[string]string{
	"modify_time": CHECKSET,
	"create_time": CHECKSET,
	"prompt_id":   CHECKSET,
}

func AlicloudPaiWorkspacePromptBasicDependence10831(name string) string {
	rand := acctest.RandIntRange(10000, 99999)
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_pai_workspace_workspace" "default7Adve7" {
  description    = "test_prompt_%d"
  display_name   = "用来测试提示词"
  workspace_name = "prompt_%d"
  env_types      = ["prod"]
}


`, name, rand, rand)
}

// Test PaiWorkspace Prompt. <<< Resource test cases, automatically generated.
