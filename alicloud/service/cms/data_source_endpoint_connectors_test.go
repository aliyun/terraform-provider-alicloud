package cms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudCmsEndpointConnectorsDataSource_nameRegex covers name_regex and ids filtering.
func TestAccAliCloudCmsEndpointConnectorsDataSource_nameRegex(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccecds%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckCmsEndpointConnectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsEndpointConnectorsDataSourceNameRegex(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cms_endpoint_connectors.default", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_cms_endpoint_connectors.default", "connectors.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_endpoint_connectors.default", "connectors.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_endpoint_connectors.default", "connectors.0.connector_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_endpoint_connectors.default", "connectors.0.name"),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_endpoint_connectors.default", "connectors.0.workspace"),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_endpoint_connectors.default", "connectors.0.type"),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_endpoint_connectors.default", "connectors.0.endpoint"),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_endpoint_connectors.default", "connectors.0.created_at"),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_endpoint_connectors.default", "connectors.0.updated_at"),
				),
			},
			{
				Config: testAccCmsEndpointConnectorsDataSourceIds(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cms_endpoint_connectors.default", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_cms_endpoint_connectors.default", "connectors.#", "1"),
				),
			},
			{
				Config: testAccCmsEndpointConnectorsDataSourceFakeNameRegex(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cms_endpoint_connectors.default", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_cms_endpoint_connectors.default", "connectors.#", "0"),
				),
			},
		},
	})
}

func testAccCmsEndpointConnectorsDataSourceDependence(name string) string {
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

resource "alicloud_cms_endpoint_connector" "default" {
  workspace  = alicloud_cms_workspace.default.workspace_name
  type       = "model_service"
  name       = var.name
  endpoint   = "https://example.com/api"
  credential = {
    token = "secret-%s"
  }
}
`, name, name)
}

func testAccCmsEndpointConnectorsDataSourceNameRegex(name string) string {
	return fmt.Sprintf(`%s

data "alicloud_cms_endpoint_connectors" "default" {
  workspace  = alicloud_cms_workspace.default.workspace_name
  name_regex = alicloud_cms_endpoint_connector.default.name
}
`, testAccCmsEndpointConnectorsDataSourceDependence(name))
}

func testAccCmsEndpointConnectorsDataSourceIds(name string) string {
	return fmt.Sprintf(`%s

data "alicloud_cms_endpoint_connectors" "default" {
  workspace = alicloud_cms_workspace.default.workspace_name
  ids       = [alicloud_cms_endpoint_connector.default.id]
}
`, testAccCmsEndpointConnectorsDataSourceDependence(name))
}

func testAccCmsEndpointConnectorsDataSourceFakeNameRegex(name string) string {
	return fmt.Sprintf(`%s

data "alicloud_cms_endpoint_connectors" "default" {
  workspace  = alicloud_cms_workspace.default.workspace_name
  name_regex = "${alicloud_cms_endpoint_connector.default.name}_fake"
}
`, testAccCmsEndpointConnectorsDataSourceDependence(name))
}
