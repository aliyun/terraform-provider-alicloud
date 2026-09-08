---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_context_store"
description: |-
  Provides a Alicloud Cms Context Store resource.
---

# alicloud_cms_context_store

Provides a Cms Context Store resource.

For information about Cms Context Store and how to use it, see [What is Context Store](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateContextStore).

-> **NOTE:** Available since v1.292.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cms_workspace" "default" {
  sls_project    = "example-project"
  workspace_name = "example-workspace"
}

resource "alicloud_cms_context_store" "example" {
  workspace          = alicloud_cms_workspace.default.workspace_name
  context_store_name = "example-context-store"
  context_type       = "private"
  description        = "example context store"
}
```

With config

```terraform
resource "alicloud_cms_context_store" "example" {
  workspace          = alicloud_cms_workspace.default.workspace_name
  context_store_name = "example-context-store"
  context_type       = "private"
  description        = "example context store"

  config {
    source {
      project    = "example-project"
      logstore   = "example-logstore"
      start_time = "2024-01-01T00:00:00Z"
    }
    metadata_field = {
      key1 = "value1"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required, ForceNew) The name of the workspace to which the context store belongs.
* `context_store_name` - (Required, ForceNew) The name of the context store.
* `context_type` - (Required) The type of the context store.
* `description` - (Optional) The description of the context store.
* `config` - (Optional) The configuration of the context store. See [`config`](#config) below.

### `config`

The config supports the following:

* `source` - (Optional) The data source configuration. See [`source`](#config-source) below.
* `metadata_field` - (Optional) The metadata field map. Key-value pairs.

### `config-source`

The source supports the following:

* `project` - (Optional) The name of the Simple Log Service (SLS) project.
* `logstore` - (Optional) The name of the SLS logstore.
* `start_time` - (Optional) The start time of the data source, in ISO 8601 format.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform. It is formatted as `<workspace>:<context_store_name>`.
* `status` - The status of the context store.
* `region_id` - The region ID of the context store.
* `create_time` - The creation time of the context store.
* `update_time` - The last update time of the context store.
* `dataset` - The dataset information of the context store. See [`dataset`](#dataset) below.

### `dataset`

The dataset supports the following:

* `name` - The name of the dataset.

## Import

Cms Context Store can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_context_store.example <workspace>:<context_store_name>
```
