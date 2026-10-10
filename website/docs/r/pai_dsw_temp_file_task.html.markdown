---
subcategory: "PAI DSW"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_dsw_temp_file_task"
description: |-
  Provides a PAI DSW Temp File Task resource.
---

# alicloud_pai_dsw_temp_file_task

Provides a PAI DSW Temp File Task resource.

For information about PAI DSW Temp File Task and how to use it, see [What is TempFileTask](https://next.api.alibabacloud.com/document/pai-dsw/2022-01-01/CreateTempFileTask).

-> **NOTE:** Available since v1.286.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_pai_dsw_temp_file_task" "default" {
  instance_id      = "dsw-xxxxxxxxxxxxxxxx"
  gmt_expired_time = "2026-12-31T23:59:59Z"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the PAI DSW instance. The Temp File Task is created under this instance.
* `gmt_expired_time` - (Optional, Computed) The expiration time of the temp file task, in ISO 8601 format (e.g. `2026-12-31T23:59:59Z`). Once expired, the task is automatically cleaned up.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Temp File Task, which is the `temp_file_task_id`.
* `temp_file_task_id` - The first-level ID of the resource.
* `owner_id` - The owner ID of the temp file task.
* `user_id` - The user ID of the temp file task.
* `create_time` - The creation time of the resource.
* `gmt_modified_time` - The modification time of the resource.
* `region_id` - The region ID of the resource.

## Import

PAI DSW Temp File Task can be imported using the id, e.g.

```shell
terraform import alicloud_pai_dsw_temp_file_task.example <temp_file_task_id>
```
