---
subcategory: "PAI DSW"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_dsw_temp_file_task"
sidebar_current: "docs-alicloud-resource-pai-dsw-temp-file-task"
description: |-
  Provides a PAI DSW Temp File Task resource.
---

# alicloud_pai_dsw_temp_file_task

Provides a PAI DSW Temp File Task resource.

For information about PAI DSW Temp File Task and how to use it, see [CreateTempFileTask](https://www.alibabacloud.com/help/en/pai/developer-reference/api-pai-dsw-2022-01-01-createtempfiletask).

-> **NOTE:** Available since v2.0.0-beta5.

## Example Usage

Basic Usage

```terraform
resource "alicloud_pai_dsw_temp_file_task" "example" {
  instance_id      = "dsw-xxxxxxxx"
  gmt_expired_time = "2099-12-31 23:59:59"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the DSW instance.
* `gmt_expired_time` - (Optional) The expiry time of the temp file task.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. Same as `temp_file_task_id`.
* `owner_id` - The owner ID of the temp file task.
* `user_id` - The user ID of the temp file task.
* `create_time` - The creation time of the temp file task.
* `gmt_modified_time` - The modification time of the temp file task.
* `region_id` - The region ID of the temp file task.

## Import

PAI DSW Temp File Task can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_dsw_temp_file_task.example <temp_file_task_id>
```
