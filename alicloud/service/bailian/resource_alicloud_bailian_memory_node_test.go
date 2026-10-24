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

func TestAccAlicloudBailianMemoryNode_basic(t *testing.T) {
	resourceId := "alicloud_bailian_memory_node.default"
	workspaceResourceId := "alicloud_model_studio_workspace.default"
	r := tfacctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		IDRefreshName:            resourceId,
		CheckDestroy:             testAccCheckBailianMemoryNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudBailianMemoryNodeConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(workspaceResourceId, "id"),
					resource.TestCheckResourceAttrSet(workspaceResourceId, "api_host"),
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrSet(resourceId, "memory_node_id"),
					resource.TestCheckResourceAttrPair(resourceId, "workspace_id", workspaceResourceId, "api_host"),
					resource.TestCheckResourceAttr(resourceId, "content", fmt.Sprintf("tf-testacc-memory-node-content-%d", r)),
				),
			},
			{
				Config: testAccAlicloudBailianMemoryNodeConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "content", fmt.Sprintf("tf-testacc-memory-node-content-%d-upd", r)),
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

func testAccCheckBailianMemoryNodeDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_bailian_memory_node" {
			continue
		}
		_, err := bailian.NewService(client).GetMemoryNode(
			rs.Primary.Attributes["workspace_id"],
			rs.Primary.Attributes["memory_id"],
			rs.Primary.Attributes["memory_node_id"],
		)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("checking memory node destroy failed: %s", err)
		}
		return fmt.Errorf("memory node %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccAlicloudBailianMemoryNodeConfig(r int) string {
	return fmt.Sprintf(`
resource "alicloud_model_studio_workspace" "default" {
  workspace_name = "tftestaccws%[1]d"
  service_site   = "global"
}

resource "alicloud_bailian_memory" "default" {
  workspace_id = alicloud_model_studio_workspace.default.api_host
  description  = "tf-testacc-memory-%[1]d"
}

resource "alicloud_bailian_memory_node" "default" {
  workspace_id = alicloud_bailian_memory.default.workspace_id
  memory_id    = alicloud_bailian_memory.default.memory_id
  content      = "tf-testacc-memory-node-content-%[1]d"
}
`, r)
}

func testAccAlicloudBailianMemoryNodeConfigUpdate(r int) string {
	return fmt.Sprintf(`
resource "alicloud_model_studio_workspace" "default" {
  workspace_name = "tftestaccws%[1]d"
  service_site   = "global"
}

resource "alicloud_bailian_memory" "default" {
  workspace_id = alicloud_model_studio_workspace.default.api_host
  description  = "tf-testacc-memory-%[1]d"
}

resource "alicloud_bailian_memory_node" "default" {
  workspace_id = alicloud_bailian_memory.default.workspace_id
  memory_id    = alicloud_bailian_memory.default.memory_id
  content      = "tf-testacc-memory-node-content-%[1]d-upd"
}
`, r)
}
