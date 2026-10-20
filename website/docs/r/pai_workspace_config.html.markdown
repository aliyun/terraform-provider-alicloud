---
subcategory: "PAI Workspace"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_workspace_config"
description: |-
  Provides a Alicloud PAI Workspace Config resource.
---

# alicloud_pai_workspace_config

Provides a PAI Workspace Config resource.

A workspace configuration is a key/value pair scoped to a PAI workspace, grouped under a category name. There is no dedicated create API: the `UpdateConfig` operation upserts a configuration, so this resource creates and updates through the same operation.

For information about PAI Workspace Config and how to use it, see [GetConfig](https://www.alibabacloud.com/help/en/pai/developer-reference/api-aiworkspace-2021-02-04-getconfig).

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
  description    = "workspace for the config resource"
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_config" "default" {
  workspace_id  = alicloud_pai_workspace_workspace.default.id
  category_name = "CommonResourceConfig"
  config_key    = "tempStoragePath"
  config_value  = "oss://bucket-name.oss-cn-hangzhou.aliyuncs.com/example/"
}
```

## Argument Reference

The following arguments are supported:
* `workspace_id` - (Required, ForceNew) The ID of the PAI Workspace that owns the configuration.
* `category_name` - (Required, ForceNew) The category name of the configuration. Valid values: `CommonResourceConfig`, `DLCAutoRecycle`, `DLCPriorityConfig`, `DSWPriorityConfig`, `QuotaMaximumDuration`, `CommonTagConfig`. Each `config_key` is allowed only under a specific category (for example, `tempStoragePath` requires `CommonResourceConfig`). Changing this creates a new resource.
* `config_key` - (Required, ForceNew) The key of the workspace configuration. Changing this creates a new resource. Documented config keys include `tempStoragePath` (category `CommonResourceConfig`), `isAutoRecycle` (`DLCAutoRecycle`), `priorityConfig` (`DLCPriorityConfig`/`DSWPriorityConfig`), `quotaMaximumDuration` (`QuotaMaximumDuration`) and `predefinedTags` (`CommonTagConfig`).
* `config_value` - (Required) The value of the workspace configuration.
* `labels` - (Optional) The labels of the workspace configuration. PAI does not honor a client-requested label order, so labels are modeled as an unordered set; adding, removing or changing a label value is detected as a change, while reordering alone is not. See [`labels`](#labels) below.

### labels

The `labels` block supports the following:
* `key` - (Optional) The key of the label.
* `value` - (Optional) The value of the label.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource. The value is formulated as `<workspace_id>:<category_name>:<config_key>`.
* `gmt_create_time` - The creation time of the configuration.
* `gmt_modified_time` - The last modification time of the configuration. It is refreshed on every `UpdateConfig`, so it may differ between plan and apply.

## Import

PAI Workspace Config can be imported using the id, which consists of workspace_id, category_name and config_key, e.g.

```shell
$ terraform import alicloud_pai_workspace_config.example <workspace_id>:<category_name>:<config_key>
```
