---
subcategory: "Actiontrail"
layout: "alicloud"
page_title: "Alicloud: alicloud_actiontrail_history_delivery_jobs"
sidebar_current: "docs-alicloud-datasource-actiontrail-history-delivery-jobs"
description: |-
  Provides a list of Actiontrail History Delivery Jobs to the user.
---

# alicloud\_actiontrail\_history\_delivery\_jobs

This data source provides the Actiontrail History Delivery Jobs of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.139.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_actiontrail_history_delivery_jobs" "ids" {
  ids = ["example_id"]
}
output "actiontrail_history_delivery_job_id_1" {
  value = data.alicloud_actiontrail_history_delivery_jobs.ids.jobs.0.id
}

data "alicloud_actiontrail_history_delivery_jobs" "status" {
  ids    = ["example_id"]
  status = "2"
}
output "actiontrail_history_delivery_job_id_2" {
  value = data.alicloud_actiontrail_history_delivery_jobs.status.jobs.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of History Delivery Job IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the task. Valid values: `0`, `1`, `2`, `3`. `0`: The task is initializing. `1`: The task is delivering historical events. `2`: The delivery of historical events is complete. `3`: The task fails.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `jobs` - A list of Actiontrail History Delivery Jobs. Each element contains the following attributes:
    * `create_time` - The time when the task was created.
    * `end_time` - The time when the task ended.
    * `history_delivery_job_id` -  The resource ID in terraform of History Delivery Job.
    * `home_region` - The home region of the trail.
    * `id` - The ID of the History Delivery Job.
    * `start_time` - The time when the task started.
    * `status` - The status of the task. Valid values: `0`, `1`, `2`, `3`. `0`: The task is initializing. `1`: The task is delivering historical events. `2`: The delivery of historical events is complete. `3`: The task fails.
    * `trail_name` - The name of the trail.
    * `updated_time` - The time when the task was updated.
    * `job_status` - Detail status of delivery job.
      * `region` - The region of the delivery job.
      * `status` - The status of the task. Valid values: `0`, `1`, `2`, `3`. `0`: The task is initializing. `1`: The task is delivering historical events. `2`: The delivery of historical events is complete. `3`: The task fails.
