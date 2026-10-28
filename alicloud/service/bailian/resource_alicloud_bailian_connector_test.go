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

func TestAccAlicloudBailianConnector_basic(t *testing.T) {
	resourceId := "alicloud_bailian_connector.default"
	workspaceResourceId := "alicloud_model_studio_workspace.default"
	r := tfacctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		IDRefreshName:            resourceId,
		CheckDestroy:             testAccCheckBailianConnectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudBailianConnectorConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(workspaceResourceId, "id"),
					resource.TestCheckResourceAttrSet(workspaceResourceId, "api_host"),
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrSet(resourceId, "connector_id"),
					resource.TestCheckResourceAttrPair(resourceId, "workspace_id", workspaceResourceId, "api_host"),
					resource.TestCheckResourceAttr(resourceId, "connector_name", fmt.Sprintf("tfc-%d", r)),
					resource.TestCheckResourceAttr(resourceId, "connector_type", "FILE"),
					resource.TestCheckResourceAttr(resourceId, "description", "tf-example-description"),
				),
			},
			{
				Config: testAccAlicloudBailianConnectorConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "connector_name", fmt.Sprintf("tfc-%d-u", r)),
					resource.TestCheckResourceAttr(resourceId, "description", "tf-example-description-upd"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connector_type"},
			},
		},
	})
}

func testAccCheckBailianConnectorDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_bailian_connector" {
			continue
		}
		_, err := bailian.NewService(client).GetConnector(rs.Primary.Attributes["workspace_id"], rs.Primary.Attributes["connector_id"])
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("checking connector destroy failed: %s", err)
		}
		return fmt.Errorf("connector %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccAlicloudBailianConnectorConfig(r int) string {
	return fmt.Sprintf(`
resource "alicloud_model_studio_workspace" "default" {
  workspace_name = "tftestaccws%[1]d"
  service_site   = "global"
}

resource "alicloud_bailian_connector" "default" {
  workspace_id   = alicloud_model_studio_workspace.default.api_host
  connector_name = "tfc-%[1]d"
  connector_type = "FILE"
  description    = "tf-example-description"
}
`, r)
}

func testAccAlicloudBailianConnectorConfigUpdate(r int) string {
	return fmt.Sprintf(`
resource "alicloud_model_studio_workspace" "default" {
  workspace_name = "tftestaccws%[1]d"
  service_site   = "global"
}

resource "alicloud_bailian_connector" "default" {
  workspace_id   = alicloud_model_studio_workspace.default.api_host
  connector_name = "tfc-%[1]d-u"
  connector_type = "FILE"
  description    = "tf-example-description-upd"
}
`, r)
}
