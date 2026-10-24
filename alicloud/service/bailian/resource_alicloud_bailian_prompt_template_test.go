package bailian_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/bailian"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAlicloudBailianPromptTemplate_basic(t *testing.T) {
	resourceId := "alicloud_bailian_prompt_template.default"
	workspaceResourceId := "alicloud_model_studio_workspace.default"
	r := tfacctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		IDRefreshName:            resourceId,
		CheckDestroy:             testAccCheckBailianPromptTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudBailianPromptTemplateConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(workspaceResourceId, "id"),
					resource.TestCheckResourceAttrSet(workspaceResourceId, "api_host"),
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrSet(resourceId, "prompt_template_id"),
					resource.TestCheckResourceAttrPair(resourceId, "workspace_id", workspaceResourceId, "api_host"),
					resource.TestCheckResourceAttr(resourceId, "name", fmt.Sprintf("tf-testacc-prompt-%d", r)),
					resource.TestCheckResourceAttr(resourceId, "content", "You are a helpful assistant."),
				),
			},
			{
				Config: testAccAlicloudBailianPromptTemplateConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "name", fmt.Sprintf("tf-testacc-prompt-%d-upd", r)),
					resource.TestCheckResourceAttr(resourceId, "content", "You are an expert assistant."),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckBailianPromptTemplateDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_bailian_prompt_template" {
			continue
		}
		_, err := bailian.NewService(client).GetPromptTemplate(rs.Primary.Attributes["workspace_id"], rs.Primary.Attributes["prompt_template_id"])
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("checking prompt template destroy failed: %s", err)
		}
		return fmt.Errorf("prompt template %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccAlicloudBailianPromptTemplateConfig(r int) string {
	return fmt.Sprintf(`
resource "alicloud_model_studio_workspace" "default" {
  workspace_name = "tftestaccws%[1]d"
  service_site   = "global"
}

resource "alicloud_bailian_prompt_template" "default" {
  workspace_id = alicloud_model_studio_workspace.default.api_host
  name         = "tf-testacc-prompt-%[1]d"
  content      = "You are a helpful assistant."
}
`, r)
}

func testAccAlicloudBailianPromptTemplateConfigUpdate(r int) string {
	return fmt.Sprintf(`
resource "alicloud_model_studio_workspace" "default" {
  workspace_name = "tftestaccws%[1]d"
  service_site   = "global"
}

resource "alicloud_bailian_prompt_template" "default" {
  workspace_id = alicloud_model_studio_workspace.default.api_host
  name         = "tf-testacc-prompt-%[1]d-upd"
  content      = "You are an expert assistant."
}
`, r)
}
