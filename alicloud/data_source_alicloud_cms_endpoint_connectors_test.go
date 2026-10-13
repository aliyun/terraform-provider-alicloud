package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCmsEndpointConnectorsDataSource_nameRegex(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsEndpointConnectorsDataSourceConfig(rand, map[string]string{
			"workspace":  `"${alicloud_cms_workspace.default.workspace_name}"`,
			"name_regex": `"${alicloud_cms_endpoint_connector.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsEndpointConnectorsDataSourceConfig(rand, map[string]string{
			"workspace":  `"${alicloud_cms_workspace.default.workspace_name}"`,
			"name_regex": `"${alicloud_cms_endpoint_connector.default.name}_fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsEndpointConnectorsDataSourceConfig(rand, map[string]string{
			"workspace": `"${alicloud_cms_workspace.default.workspace_name}"`,
			"ids":       `["${alicloud_cms_endpoint_connector.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCmsEndpointConnectorsDataSourceConfig(rand, map[string]string{
			"workspace": `"${alicloud_cms_workspace.default.workspace_name}"`,
			"ids":       `["${alicloud_cms_endpoint_connector.default.id}_fake"]`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsEndpointConnectorsDataSourceConfig(rand, map[string]string{
			"workspace":  `"${alicloud_cms_workspace.default.workspace_name}"`,
			"name_regex": `"${alicloud_cms_endpoint_connector.default.name}"`,
			"ids":        `["${alicloud_cms_endpoint_connector.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCmsEndpointConnectorsDataSourceConfig(rand, map[string]string{
			"workspace":  `"${alicloud_cms_workspace.default.workspace_name}"`,
			"name_regex": `"${alicloud_cms_endpoint_connector.default.name}_fake"`,
			"ids":        `["${alicloud_cms_endpoint_connector.default.id}_fake"]`,
		}),
	}
	var existAlicloudCmsEndpointConnectorsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"connectors.#":              "1",
			"connectors.0.id":           CHECKSET,
			"connectors.0.connector_id": CHECKSET,
			"connectors.0.name":         CHECKSET,
			"connectors.0.workspace":    CHECKSET,
			"connectors.0.type":         CHECKSET,
			"connectors.0.endpoint":     CHECKSET,
			"connectors.0.created_at":   CHECKSET,
			"connectors.0.updated_at":   CHECKSET,
		}
	}
	var fakeAlicloudCmsEndpointConnectorsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        "0",
			"connectors.#": "0",
		}
	}
	var alicloudCmsEndpointConnectorsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cms_endpoint_connectors.default",
		existMapFunc: existAlicloudCmsEndpointConnectorsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCmsEndpointConnectorsDataSourceNameMapFunc,
	}
	alicloudCmsEndpointConnectorsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlicloudCmsEndpointConnectorsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tftestacccmsendpointconnector%d"
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
    token = "secret-value"
  }
}

data "alicloud_cms_endpoint_connectors" "default" {
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
