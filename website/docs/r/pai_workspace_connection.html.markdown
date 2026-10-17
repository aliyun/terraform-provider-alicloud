---
subcategory: "PAI Workspace"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_workspace_connection"
description: |-
  Provides a Alicloud PAI Workspace Connection resource.
---

# alicloud_pai_workspace_connection

Provides a PAI Workspace Connection resource.

For information about PAI Workspace Connection and how to use it, see [What is Connection](https://next.api.alibabacloud.com/document/AIWorkSpace/2021-02-04/CreateConnection).

-> **NOTE:** Available since v1.295.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform_example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_pai_workspace_workspace" "default" {
  workspace_name = var.name
  description    = var.name
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_connection" "default" {
  connection_name = var.name
  connection_type = "MilvusConnection"
  workspace_id    = alicloud_pai_workspace_workspace.default.id
  accessibility   = "PRIVATE"
  description     = var.name

  configs = {
    database = "default"
    uri      = "example.milvus.com"
  }

  secrets = {
    token = "your-token"
  }

  models = [
    {
      model                   = "m1"
      display_name            = "model-1"
      model_type              = "OPEN_SOURCE"
      tool_call               = false
      support_reasoning       = true
      support_response_schema = true
      support_vision          = false
      max_model_length        = 4096
    }
  ]

  resource_meta = {
    instance_name = "instance-1"
    instance_id   = "inst-1"
    extra         = "extra-info"
  }
}
```

## Argument Reference

The following arguments are supported:

* `connection_name` - (Required, ForceNew) The name of the connection.
* `connection_type` - (Required, ForceNew) The type of the connection, such as `MilvusConnection`.
* `workspace_id` - (Required, ForceNew) The ID of the PAI workspace to which the connection belongs.
* `description` - (Optional) The description of the connection.
* `accessibility` - (Optional, ForceNew) The accessibility of the connection, such as `PRIVATE` or `PUBLIC`.
* `configs` - (Required, Map) The configuration info of the connection, a map of string key-value pairs.
* `secrets` - (Optional, Map) The secrets of the connection, a map of string key-value pairs.
* `models` - (Optional) The model configuration list of the connection. See [`models`](#models) below.
* `resource_meta` - (Optional, ForceNew) The resource meta info of the connection. See [`resource_meta`](#resource_meta) below.

### models

The `models` supports the following:

* `model` - (Optional) The model identifier.
* `display_name` - (Optional) The display name of the model.
* `model_type` - (Optional) The type of the model.
* `tool_call` - (Optional) Whether the model supports tool call.
* `support_reasoning` - (Optional) Whether the model supports reasoning. Create-only: written at creation and not applied on update.
* `support_response_schema` - (Optional) Whether the model supports response schema. Create-only: written at creation and not applied on update.
* `support_vision` - (Optional) Whether the model supports vision. Create-only: written at creation and not applied on update.
* `max_model_length` - (Optional) The max model length. Create-only: written at creation and not applied on update.

### resource_meta

The `resource_meta` supports the following:

* `instance_name` - (Optional) The instance name.
* `instance_id` - (Optional) The instance ID.
* `extra` - (Optional) The extra information of the resource.

-> **NOTE:** The `support_reasoning`, `support_response_schema`, `support_vision`, `max_model_length` fields under `models` can only be set when the connection is created; the API ignores them on update. To change one of them, recreate the resource (for example, by tainting it). The `resource_meta` block is `ForceNew` as a whole: changing any of its fields (`instance_name`, `instance_id`, `extra`) forces resource recreation.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. Same as `connection_id`.
* `connection_id` - The first ID of the resource, the connection ID.
* `creator` - The creator of the connection.
* `gmt_create_time` - The create time of the connection.
* `gmt_modified_time` - The modified time of the connection.

## Import

PAI Workspace Connection can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_workspace_connection.example <connection-id>
```
