---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_endpoint_connectors"
sidebar_current: "docs-alicloud-datasource-cms-endpoint-connectors"
description: |-
  Provides a list of Cms Endpoint Connectors to the user.
---

# alicloud_cms_endpoint_connectors

This data source provides the Cms Endpoint Connectors of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.292.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.name
  sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_endpoint_connector" "default" {
  workspace = alicloud_cms_workspace.default.workspace_name
  type      = "model_service"
  name      = var.name
  endpoint  = "https://example.com/api"
  credential = {
    token = "secret-token"
  }
}

data "alicloud_cms_endpoint_connectors" "default" {
  workspace = alicloud_cms_workspace.default.workspace_name
  ids       = [alicloud_cms_endpoint_connector.default.id]
}
output "cms_endpoint_connector_id_1" {
  value = data.alicloud_cms_endpoint_connectors.default.connectors.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Endpoint Connector IDs. Its element value is formatted as `<workspace>:<connector_id>`.
* `name` - (Optional, ForceNew) The name of the endpoint connector, used to filter results by exact name match.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by endpoint connector name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `type` - (Optional, ForceNew) The type of the endpoint connector, used to filter results by type. Valid values: `model_service`, `agent_app`.
* `workspace` - (Required, ForceNew) The name of the workspace to which the endpoint connectors belong.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `connectors` - A list of Cms Endpoint Connectors. Each element contains the following attributes:
  * `alias` - The alias of the endpoint connector.
  * `connector_id` - The ID of the endpoint connector.
  * `created_at` - The creation timestamp of the endpoint connector.
  * `description` - The description of the endpoint connector.
  * `endpoint` - The service endpoint URL of the endpoint connector.
  * `id` - The ID of the resource. It is formatted as `<workspace>:<connector_id>`.
  * `name` - The name of the endpoint connector.
  * `properties` - The type-specific properties of the endpoint connector.
  * `type` - The type of the endpoint connector.
  * `updated_at` - The last update timestamp of the endpoint connector.
  * `workspace` - The name of the workspace to which the endpoint connector belongs.
