---
subcategory: "PAI DSW"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_dsw_temp_file_task"
description: |-
  Provides a data source for the PAI DSW Temp File Task.
---

# alicloud_pai_dsw_temp_file_task

Provides a data source to query a PAI DSW Temp File Task by its ID.

-> **NOTE:** Available since v1.286.0.

## Example Usage

```terraform
data "alicloud_pai_dsw_temp_file_task" "default" {
  temp_file_task_id = "tft-xxxxxxxxxxxxxxxx"
}

output "instance_id" {
  value = data.alicloud_pai_dsw_temp_file_task.default.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `temp_file_task_id` - (Required) The ID of the PAI DSW Temp File Task to query.
* `output_file` - (Optional) File path where the results will be saved after running the data source query.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Temp File Task (same as `temp_file_task_id`).
* `instance_id` - The ID of the PAI DSW instance.
* `owner_id` - The owner ID of the temp file task.
* `user_id` - The user ID of the temp file task.
* `create_time` - The creation time of the resource.
* `gmt_modified_time` - The modification time of the resource.
* `gmt_expired_time` - The expiration time of the resource.
* `region_id` - The region ID of the resource.
