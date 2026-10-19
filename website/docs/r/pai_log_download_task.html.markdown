---
subcategory: "PAI"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_log_download_task"
description: |-
  Provides a Alicloud PAI EAS Log Download Task resource.
---

# alicloud_pai_log_download_task

Provides a PAI EAS Log Download Task resource.

A log download task collects the service logs of a specified EAS service within a
time range and produces a downloadable archive. The task is regional: its cluster
is the provider region, and it references a public EAS service by name.

-> **NOTE:** Available since v1.295.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_pai_service" "default" {
  develop = "false"
  service_config = jsonencode({
    metadata = {
      cpu      = 1
      gpu      = 0
      instance = 1
      memory   = 2000
      name     = "tfexample"
      rpc      = { keepalive = 70000 }
    }
    model_path      = "http://eas-data.oss-cn-shanghai.aliyuncs.com/processors/echo_processor_release.tar.gz"
    processor_entry = "libecho.so"
    processor_path  = "http://eas-data.oss-cn-shanghai.aliyuncs.com/processors/echo_processor_release.tar.gz"
    processor_type  = "cpp"
  })
}

resource "alicloud_pai_log_download_task" "default" {
  service_name  = alicloud_pai_service.default.id
  alias_name    = "tf-example"
  instance_name = "main"
  container     = "main"
}
```

## Argument Reference

The following arguments are supported:
* `service_name` - (Required, ForceNew) The name of the EAS service whose logs are downloaded.
* `alias_name` - (Optional, Computed, ForceNew) The alias name of the log download task.
* `from_time` - (Optional, Computed, ForceNew) The start time of the log collection range, formatted as `yyyy-MM-dd HH:mm:ss`. Defaults to seven days before `to_time`.
* `to_time` - (Optional, Computed, ForceNew) The end time of the log collection range, formatted as `yyyy-MM-dd HH:mm:ss`. Defaults to the current time.
* `instance_name` - (Optional, Computed, ForceNew) The name of the EAS instance to filter logs by. An empty value means all instances.
* `container` - (Optional, Computed, ForceNew) The name of the business container to filter logs by. An empty value means all business containers.

  -> **NOTE:** The EAS Create API enforces a server-side seven-day window on the log collection range. When `from_time` and `to_time` are omitted the API applies the defaults (`from_time` = `to_time` minus seven days, `to_time` = the current time); explicit values outside the window are rejected.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource. Same as `task_id`.
* `task_id` - The ID of the log download task assigned by EAS.
* `create_time` - The creation time of the log download task.
* `finish_time` - The finish time of the log download task.
* `status` - The status of the log download task, such as `Running`, `Finished` or `Failed`.
* `download_url` - The download URL of the produced log archive, populated when the task finishes.
* `region_id` - The region ID of the EAS cluster, same as the provider region.

## Import

PAI EAS Log Download Task can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_log_download_task.example <id>
```
