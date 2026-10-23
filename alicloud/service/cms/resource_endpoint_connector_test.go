package cms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cms"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAliCloudCmsEndpointConnector_basic covers create, full update, clearing
// optional fields and import.
func TestAccAliCloudCmsEndpointConnector_basic(t *testing.T) {
	resourceId := "alicloud_cms_endpoint_connector.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccecbasic%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckCmsEndpointConnectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsEndpointConnectorBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "type", "model_service"),
					resource.TestCheckResourceAttr(resourceId, "name", name),
					resource.TestCheckResourceAttr(resourceId, "endpoint", "https://example.com/api"),
					resource.TestCheckResourceAttr(resourceId, "alias", name),
					resource.TestCheckResourceAttr(resourceId, "description", name),
					resource.TestCheckResourceAttrSet(resourceId, "connector_id"),
					resource.TestCheckResourceAttrSet(resourceId, "region_id"),
					resource.TestCheckResourceAttrSet(resourceId, "created_at"),
					resource.TestCheckResourceAttrSet(resourceId, "updated_at"),
				),
			},
			{
				Config: testAccCmsEndpointConnectorBasicUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "name", fmt.Sprintf("%s-v2", name)),
					resource.TestCheckResourceAttr(resourceId, "endpoint", "https://example.com/api/v2"),
					resource.TestCheckResourceAttr(resourceId, "alias", fmt.Sprintf("%s-updated", name)),
					resource.TestCheckResourceAttr(resourceId, "description", fmt.Sprintf("%s-updated", name)),
				),
			},
			{
				Config: testAccCmsEndpointConnectorBasicClear(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "description", ""),
					resource.TestCheckResourceAttr(resourceId, "alias", ""),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credential"},
			},
		},
	})
}

// TestAccAliCloudCmsEndpointConnector_agentApp covers the agent_app connector type.
func TestAccAliCloudCmsEndpointConnector_agentApp(t *testing.T) {
	resourceId := "alicloud_cms_endpoint_connector.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccecagent%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckCmsEndpointConnectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsEndpointConnectorAgentApp(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "type", "agent_app"),
					resource.TestCheckResourceAttr(resourceId, "name", name),
					resource.TestCheckResourceAttr(resourceId, "endpoint", "https://agent.example.com"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credential"},
			},
		},
	})
}

func testAccCheckCmsEndpointConnectorDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cms_endpoint_connector" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		_, err := cms.NewService(client).DescribeEndpointConnector(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing endpoint connector %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("endpoint connector %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccCmsEndpointConnectorDependence(name string) string {
	return fmt.Sprintf(`
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "%s"
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

func testAccCmsEndpointConnectorBasic(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_cms_endpoint_connector" "default" {
  workspace   = alicloud_cms_workspace.default.workspace_name
  type        = "model_service"
  name        = var.name
  endpoint    = "https://example.com/api"
  alias       = var.name
  description = var.name
  credential = {
    token = "secret-%s"
  }
  headers {
    key   = "X-Custom-Header"
    value = "foo"
  }
  properties = {
    modelId = "model-001"
  }
}
`, testAccCmsEndpointConnectorDependence(name), name)
}

func testAccCmsEndpointConnectorBasicUpdate(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_cms_endpoint_connector" "default" {
  workspace   = alicloud_cms_workspace.default.workspace_name
  type        = "model_service"
  name        = "%s-v2"
  endpoint    = "https://example.com/api/v2"
  alias       = "%s-updated"
  description = "%s-updated"
  credential = {
    token = "secret-%s-v2"
  }
  headers {
    key   = "X-Other-Header"
    value = "bar"
  }
  properties = {
    modelId = "model-002"
  }
}
`, testAccCmsEndpointConnectorDependence(name), name, name, name, name)
}

func testAccCmsEndpointConnectorBasicClear(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_cms_endpoint_connector" "default" {
  workspace   = alicloud_cms_workspace.default.workspace_name
  type        = "model_service"
  name        = var.name
  endpoint    = "https://example.com/api/v2"
  description = ""
  alias       = ""
  credential = {
    token = "secret-%s-v2"
  }
  properties = {}
}
`, testAccCmsEndpointConnectorDependence(name), name)
}

func testAccCmsEndpointConnectorAgentApp(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_cms_endpoint_connector" "default" {
  workspace  = alicloud_cms_workspace.default.workspace_name
  type       = "agent_app"
  name       = var.name
  endpoint   = "https://agent.example.com"
  credential = {
    apiKey = "ak-%s"
  }
  properties = {
    appId = "app-001"
  }
}
`, testAccCmsEndpointConnectorDependence(name), name)
}
