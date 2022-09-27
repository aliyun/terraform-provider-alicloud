---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_hybrid_monitor_fc_task"
sidebar_current: "docs-alicloud-resource-cms-hybrid-monitor-fc-task"
description: |-
  Provides a Alicloud Cloud Monitor Service Hybrid Monitor Fc Task resource.
---

# alicloud\_cms\_hybrid\_monitor\_fc\_task

Provides a Cloud Monitor Service Hybrid Monitor Fc Task resource.

For information about Cloud Monitor Service Hybrid Monitor Fc Task and how to use it, see [What is Hybrid Monitor Fc Task](https://www.alibabacloud.com/help/en/cloudmonitor/latest/createhybridmonitortask).

-> **NOTE:** Available in v1.179.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cms_hybrid_monitor_fc_task" "example" {
  namespace      = "example_value"
  yarm_config    = "example_value"
  target_user_id = "example_value"
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, ForceNew) The index warehouse where the host belongs.
* `target_user_id` - (Optional, ForceNew) The ID of the member account. If you call API operations by using a management account, you can connect the Alibaba Cloud services that are activated for a member account in Resource Directory to Hybrid Cloud Monitoring. You can use Resource Directory to monitor Alibaba Cloud services across enterprise accounts.
* `yarm_config` - (Required) The configuration file of the Alibaba Cloud service that you want to monitor by using Hybrid Cloud Monitoring.
  - `namespace`: the namespace of the Alibaba Cloud service.
  - `metric_list`: the metrics of the Alibaba Cloud service.
  
## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Hybrid Monitor Fc Task. The value formats as `<hybrid_monitor_fc_task_id>:<namespace>`.
* `hybrid_monitor_fc_task_id` - The ID of the monitoring task.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Hybrid Monitor Fc Task.
* `delete` - (Defaults to 2 mins) Used when delete the Hybrid Monitor Fc Task.
* `update` - (Defaults to 2 mins) Used when update the Hybrid Monitor Fc Task.

## Import

Cloud Monitor Service Hybrid Monitor Fc Task can be imported using the id, e.g.

```
$ terraform import alicloud_cms_hybrid_monitor_fc_task.example <hybrid_monitor_fc_task_id>:<namespace>
```