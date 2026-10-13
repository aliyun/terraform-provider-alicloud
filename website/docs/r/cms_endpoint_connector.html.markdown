---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_endpoint_connector"
description: |-
  Provides a Alicloud Cms Endpoint Connector resource.
---

# alicloud_cms_endpoint_connector

Provides a Cms Endpoint Connector resource.

For information about Cms Endpoint Connector and how to use it, see [What is Endpoint Connector](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateEndpointConnector).

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
  workspace   = alicloud_cms_workspace.default.workspace_name
  type        = "model_service"
  name        = var.name
  endpoint    = "https://example.com/api"
  alias       = var.name
  description = var.name
  credential = {
    token = "secret-token"
  }
  headers {
    key   = "X-Custom-Header"
    value = "foo"
  }
  properties = {
    modelId = "model-001"
  }
}
```

## Argument Reference

The following arguments are supported:
* `alias` - (Optional) The alias of the endpoint connector, used for display.
* `credential` - (Required, Sensitive) The credential configuration of the endpoint connector.
* `description` - (Optional) The description of the endpoint connector.
* `endpoint` - (Required) The service endpoint URL of the endpoint connector.
* `headers` - (Optional) The custom HTTP request headers of the endpoint connector. See [`headers`](#headers) below.
* `name` - (Required) The name of the endpoint connector, unique within the workspace and type.
* `properties` - (Optional) The type-specific properties of the endpoint connector. The required keys are validated by the backend based on the `type`.
* `type` - (Required, ForceNew) The type of the endpoint connector. Valid values: `model_service`, `agent_app`.
* `workspace` - (Required, ForceNew) The name of the workspace to which the endpoint connector belongs.

### `headers`

The `headers` supports the following:
* `key` - (Optional) The header key.
* `value` - (Optional) The header value.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. It is formatted as `<workspace>:<connector_id>`.
* `connector_id` - The ID of the endpoint connector.
* `created_at` - The creation timestamp of the endpoint connector.
* `region_id` - The region ID of the resource.
* `updated_at` - The last update timestamp of the endpoint connector.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Endpoint Connector.
* `delete` - (Defaults to 5 mins) Used when delete the Endpoint Connector.
* `update` - (Defaults to 5 mins) Used when update the Endpoint Connector.

## Import

Cms Endpoint Connector can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_endpoint_connector.example <workspace>:<connector_id>
```
