---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_entity_store"
description: |-
  Provides a Alicloud Cms EntityStore resource.
---

# alicloud_cms_entity_store

Provides a Cms EntityStore resource.

EntityStore is a singleton sub-resource of a Cms Workspace. It is identified by
the workspace name and does not support in-place updates; any change to the
identity forces a new resource.

For information about Cms EntityStore and how to use it, see [What is EntityStore](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateEntityStore).

-> **NOTE:** Available since v1.277.0.

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

resource "alicloud_cms_entity_store" "default" {
  workspace_name = alicloud_cms_workspace.default.workspace_name
}
```

## Argument Reference

The following arguments are supported:
* `workspace_name` - (Required, ForceNew) The name of the workspace that the EntityStore belongs to. EntityStore is a singleton resource identified by the workspace name.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. It is the workspace name.
* `region_id` - The region ID of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the EntityStore.
* `delete` - (Defaults to 5 mins) Used when delete the EntityStore.

## Import

Cms EntityStore can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_entity_store.example <workspace_name>
```
