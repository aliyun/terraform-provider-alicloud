---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_application_monitor"
description: |-
  Provides a Alicloud Ga Application Monitor resource.
---

# alicloud_ga_application_monitor

Provides a Ga Application Monitor resource. Source station detection task.

For information about Ga Application Monitor and how to use it, see [What is Application Monitor](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.210.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}


resource "alicloud_ga_application_monitor" "default" {
  address        = "www.baidu.com"
  task_name      = "aaaa"
  accelerator_id = "ga-bp1l8bs1z8gw0u6y1ag5g"
  listener_id    = "lsr-bp1jovsyx377pboiel4p2"
}
```

## Argument Reference

The following arguments are supported:
* `accelerator_id` - (Required, ForceNew) The id of the accelerator.
* `address` - (Required) The address of task.
* `detect_enable` - (Optional) Automatic diagnostic switch.
* `detect_threshold` - (Optional) Diagnostic threshold.
* `detect_times` - (Optional) Diagnostic trigger period.
* `listener_id` - (Required) The id of the listener.
* `options_json` - (Optional) Advanced options.
* `silence_time` - (Optional) Diagnostic silence Time.
* `status` - (Optional, Computed) The status of the resource.
* `task_name` - (Required) The Name of the task.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Application Monitor.
* `delete` - (Defaults to 5 mins) Used when delete the Application Monitor.
* `update` - (Defaults to 5 mins) Used when update the Application Monitor.

## Import

Ga Application Monitor can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_application_monitor.example <id>
```