---
subcategory: "IaC Service (IaCService)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ia_c_service_module"
description: |-
  Provides a Alicloud IaC Service Module resource.
---

# alicloud_ia_c_service_module

Provides a IaC Service Module resource.

Module.

For information about IaC Service Module and how to use it, see [What is Module](https://next.api.alibabacloud.com/document/IaCService/2021-08-06/CreateModule).

-> **NOTE:** Available since v1.287.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-zhangjiakou"
}

resource "alicloud_ia_c_service_module" "default" {
  module_name      = var.name
  source           = "Registry"
  source_path      = "alibaba/security-group:2.4.1"
  version_strategy = "Manual"
  description      = var.name
  tags {
    tag_key   = "Created"
    tag_value = "TF"
  }
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) The description of the module.
* `module_name` - (Required) The name of the module.
* `source` - (Required, ForceNew) The source of the module. Valid values: `OSS`, `Registry`, `ExportTask`, `Upload`, `Shared`, `Editor`.
* `source_path` - (Optional) The path of the module source. The format depends on the value of `source`:
  - When `source` is set to `Registry`, the value is `<workspace name>/<module name>:<module version>`, for example `terraform-alicloud-modules/rds:1.0.0`.
  - When `source` is set to `OSS`, the value is `oss::<file link>`, for example `oss::https://example-bucket.oss-cn-zhangjiakou.aliyuncs.com/code.zip`.
  - When `source` is set to `ExportTask`, the value is `<export task ID>:<exported version>`, for example `ex-example1:1.0.0`.
* `state_path` - (Optional) The path of the state file that corresponds to the module. It is currently valid only for the `OSS` source. The value is in the format of `oss::<OSS path of the state file>`, for example `oss::https://example-bucket.oss-cn-zhangjiakou.aliyuncs.com/terraform.tfstate`.
* `version_strategy` - (Optional, Computed) The version generation policy of the module. Valid values: `Manual`, `SourcePathUpdated`. Default value: `Manual`. `Manual` means that versions are generated manually. `SourcePathUpdated` means that a new version is generated when `source_path` is changed.
* `tags` - (Optional, Set) A mapping of tags to assign to the module. See [`tags`](#tags) below.

### `tags`

The tags supports the following:
* `tag_key` - (Required) The key of the tag.
* `tag_value` - (Optional) The value of the tag.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the module.
* `create_time` - The time when the module was created. The time is in the `yyyy-MM-ddTHH:mm:ssZ` format, which indicates UTC.
* `group_info` - The group information of the module. It is read-only for now; assigning a module to a group will be supported after the corresponding group and project resources are available. See [`group_info`](#group_info) below.
* `latest_version` - The latest version number of the module.
* `output_path` - The storage path of the module.
* `status` - The status of the module. Valid values: `Creating`, `Created`, `Errored`. A module version can be published only after the module enters the `Created` state.

### `group_info`

The group_info supports the following:
* `group_id` - The ID of the group to which the module belongs.
* `group_name` - The name of the group to which the module belongs.
* `project_id` - The ID of the project to which the module belongs.
* `project_name` - The name of the project to which the module belongs.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Module.
* `delete` - (Defaults to 5 mins) Used when delete the Module.
* `update` - (Defaults to 5 mins) Used when update the Module.

## Import

IaC Service Module can be imported using the id, e.g.

```shell
$ terraform import alicloud_ia_c_service_module.example <id>
```
