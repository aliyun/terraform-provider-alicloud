---
subcategory: "IaC Service (IaCService)"
layout: "alicloud"
page_title: "Alicloud: alicloud_iac_service_module"
description: |-
  Provides a Alicloud IaC Service Module resource.
---

# alicloud_iac_service_module

Provides a IaC Service Module resource.

Module.

For information about IaC Service Module and how to use it, see [What is Module](https://next.api.alibabacloud.com/document/IaCService/2021-08-06/CreateModule).

-> **NOTE:** Available since v1.294.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-zhangjiakou"
}

resource "alicloud_iac_service_module" "default" {
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
* `source` - (Required, ForceNew) The source of the module. Valid values: `OSS`, `Registry`, `ExportTask`, `Upload`. When `source` is set to `Upload`, the module is created from a local module package uploaded through the `UploadModule` API (a local ZIP file of no more than 10 MB).
* `source_path` - (Optional) The path of the module source. The format depends on the value of `source`:
  - When `source` is set to `Registry`, the value is `<workspace name>/<module name>:<module version>`, for example `terraform-alicloud-modules/rds:1.0.0`.
  - When `source` is set to `OSS`, the value is `oss::<file link>`. The file must be a ZIP file, for example `oss::https://example-bucket.oss-cn-zhangjiakou.aliyuncs.com/code.zip`.
  - When `source` is set to `ExportTask`, the value is `<export task ID>:<exported version>`, for example `ex-example1:1.0.0`.
  - When `source` is set to `Upload`, `source_path` must be left empty. The local module package (a ZIP file of no more than 10 MB) is uploaded through the `UploadModule` API.
* `state_path` - (Optional) The path of the state file that corresponds to the module. It is currently valid only for the `OSS` source. The value is in the format of `oss::<OSS path of the state file>`, for example `oss::https://example-bucket.oss-cn-zhangjiakou.aliyuncs.com/terraform.tfstate`.
* `version_strategy` - (Optional, Computed) The version generation policy of the module. Valid values: `Manual`, `SourcePathUpdated`. Default value: `Manual`. `Manual` means that versions are generated manually. `SourcePathUpdated` means that a new version is generated when `source_path` is changed.
* `group_info` - (Optional, Computed) The group and project to which the module belongs. If not specified, the module is created without a group assignment and the server fills in the assignment information when available. See [`group_info`](#group_info) below.
* `tags` - (Optional, Set) A mapping of tags to assign to the module. See [`tags`](#tags) below.

### `group_info`

The group_info supports the following:
* `group_id` - (Optional, Computed) The ID of the group to which the module belongs. The group must already exist.
* `project_id` - (Optional, Computed) The ID of the project to which the module belongs. The project must already exist.
* `group_name` - (Computed) The name of the group to which the module belongs.
* `project_name` - (Computed) The name of the project to which the module belongs.

### `tags`

The tags supports the following:
* `tag_key` - (Required) The key of the tag.
* `tag_value` - (Required) The value of the tag.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the module.
* `create_time` - The time when the module was created. The time is in the `yyyy-MM-ddTHH:mm:ssZ` format, which indicates UTC.
* `latest_version` - The latest version number of the module.
* `output_path` - The storage path of the module.
* `status` - The status of the module. Valid values: `Creating`, `Created`, `Errored`. A module version can be published only after the module enters the `Created` state.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Module.
* `delete` - (Defaults to 5 mins) Used when delete the Module.
* `update` - (Defaults to 5 mins) Used when update the Module.

## Import

IaC Service Module can be imported using the id, e.g.

```shell
$ terraform import alicloud_iac_service_module.example <id>
```
