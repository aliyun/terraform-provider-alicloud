---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_scheduled_preload_execution"
description: |-
  Provides a Alicloud ESA Scheduled Preload Execution resource.
---

# alicloud_esa_scheduled_preload_execution

Provides a ESA Scheduled Preload Execution resource.



For information about ESA Scheduled Preload Execution and how to use it, see [What is Scheduled Preload Execution](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateScheduledPreloadExecutions).

-> **NOTE:** Available since v1.248.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "default" {
  site_name   = "terraform.cn"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_scheduled_preload_job" "default" {
  insert_way                 = "textBox"
  site_id                    = alicloud_esa_site.default.id
  scheduled_preload_job_name = "example_scheduledpreloadexecution_job"
  url_list                   = "http://example.gositecdn.cn/example/example.txt"
}

resource "alicloud_esa_scheduled_preload_execution" "default" {
  slice_len                = "5"
  end_time                 = "2024-06-04T10:02:09.000+08:00"
  start_time               = "2024-06-04T00:00:00.000+08:00"
  scheduled_preload_job_id = alicloud_esa_scheduled_preload_job.default.scheduled_preload_job_id
  interval                 = "30"
}
```

## Argument Reference

The following arguments are supported:
* `end_time` - (Optional) The end time of the prefetch plan.
* `interval` - (Required, Int) The time interval between each batch execution. Unit: seconds.
* `scheduled_preload_job_id` - (Required, ForceNew) The ID of the prefetch task.
* `slice_len` - (Required, Int) The number of URLs prefetched in each batch.
* `start_time` - (Optional) The start time of the prefetch plan.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<scheduled_preload_job_id>:<scheduled_preload_execution_id>`.
* `scheduled_preload_execution_id` - The ID of the prefetch plan.
* `status` - The status of the prefetch plan, including the following statuses.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Scheduled Preload Execution.
* `delete` - (Defaults to 5 mins) Used when delete the Scheduled Preload Execution.
* `update` - (Defaults to 5 mins) Used when update the Scheduled Preload Execution.

## Import

ESA Scheduled Preload Execution can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_scheduled_preload_execution.example <scheduled_preload_job_id>:<scheduled_preload_execution_id>
```