---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_context_store_api_key"
description: |-
  Provides a Alicloud Cms Context Store API Key resource.
---

# alicloud_cms_context_store_api_key

Provides a Cms Context Store API Key resource.

An API key of a context store in a Cloud Monitor 2.0 workspace.

For information about Cms Context Store API Key and how to use it, see [What is Context Store API Key](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateContextStoreAPIKey).

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

resource "alicloud_cms_context_store_api_key" "default" {
  workspace          = alicloud_cms_workspace.default.workspace_name
  context_store_name = "example-context-store"
  name               = var.name
}
```

## Argument Reference

The following arguments are supported:
* `workspace` - (Required, ForceNew) The name of the workspace to which the context store belongs.
* `context_store_name` - (Required, ForceNew) The name of the context store.
* `name` - (Required, ForceNew) The display name of the API key, which is used to identify the purpose of the key.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<workspace>:<context_store_name>:<name>`.
* `api_key` - The value of the API key.
* `create_time` - The creation time of the API key.
* `region_id` - The region ID of the API key.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Context Store API Key.
* `delete` - (Defaults to 5 mins) Used when delete the Context Store API Key.

## Import

Cms Context Store API Key can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_context_store_api_key.example <workspace>:<context_store_name>:<name>
```
