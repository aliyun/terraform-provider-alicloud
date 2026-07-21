---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_hybrid_monitor_fc_tasks"
sidebar_current: "docs-alicloud-datasource-cms-hybrid-monitor-fc-tasks"
description: |-
  Provides a list of Cms Hybrid Monitor Fc Tasks to the user.
---

# alicloud\_cms\_hybrid\_monitor\_fc\_tasks

This data source provides the Cms Hybrid Monitor Fc Tasks of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.179.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_hybrid_monitor_fc_tasks" "ids" {
  ids = ["example_value"]
}
output "cms_hybrid_monitor_fc_task_id_1" {
  value = data.alicloud_cms_hybrid_monitor_fc_tasks.ids.tasks.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Hybrid Monitor Fc Task IDs.
* `keyword` - (Optional, ForceNew) The keyword.
* `namespace` - (Optional, ForceNew) The name of the namespace.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `tasks` - A list of Cms Hybrid Monitor Fc Tasks. Each element contains the following attributes:
 * `create_time` - Create the timestamp of the monitoring task. Unit: milliseconds.
 * `hybrid_monitor_fc_task_id` - The ID of the monitoring task.
 * `id` - The ID of the Hybrid Monitor Fc Task. The value formats as `<hybrid_monitor_fc_task_id>:<namespace>`.
 * `namespace` - The index warehouse where the host belongs.
 * `target_user_id` - The ID of the member account.
 * `yarm_config` - The configuration file of the Alibaba Cloud service that you want to monitor by using Hybrid Cloud Monitoring.