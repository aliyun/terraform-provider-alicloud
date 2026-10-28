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

func TestAccAlicloudBailianMemory_basic(t *testing.T) {
	resourceId := "alicloud_bailian_memory.default"
	workspaceResourceId := "alicloud_model_studio_workspace.default"
	r := tfacctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		IDRefreshName:            resourceId,
		CheckDestroy:             testAccCheckBailianMemoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudBailianMemoryConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(workspaceResourceId, "id"),
					resource.TestCheckResourceAttrSet(workspaceResourceId, "api_host"),
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrSet(resourceId, "memory_id"),
					resource.TestCheckResourceAttrPair(resourceId, "workspace_id", workspaceResourceId, "api_host"),
					resource.TestCheckResourceAttr(resourceId, "description", fmt.Sprintf("tf-testacc-memory-%d", r)),
				),
			},
			{
				Config: testAccAlicloudBailianMemoryConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "description", fmt.Sprintf("tf-testacc-memory-%d-upd", r)),
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

func testAccCheckBailianMemoryDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_bailian_memory" {
			continue
		}
		_, err := bailian.NewService(client).GetMemory(rs.Primary.Attributes["workspace_id"], rs.Primary.Attributes["memory_id"])
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("checking memory destroy failed: %s", err)
		}
		return fmt.Errorf("memory %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccAlicloudBailianMemoryConfig(r int) string {
	return fmt.Sprintf(`
resource "alicloud_model_studio_workspace" "default" {
  workspace_name = "tftestaccws%[1]d"
  service_site   = "global"
}

resource "alicloud_bailian_memory" "default" {
  workspace_id = alicloud_model_studio_workspace.default.api_host
  description  = "tf-testacc-memory-%[1]d"
}
`, r)
}

func testAccAlicloudBailianMemoryConfigUpdate(r int) string {
	return fmt.Sprintf(`
resource "alicloud_model_studio_workspace" "default" {
  workspace_name = "tftestaccws%[1]d"
  service_site   = "global"
}

resource "alicloud_bailian_memory" "default" {
  workspace_id = alicloud_model_studio_workspace.default.api_host
  description  = "tf-testacc-memory-%[1]d-upd"
}
`, r)
}
