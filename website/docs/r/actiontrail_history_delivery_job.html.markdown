---
subcategory: "Actiontrail"
layout: "alicloud"
page_title: "Alicloud: alicloud_actiontrail_history_delivery_job"
sidebar_current: "docs-alicloud-resource-actiontrail-history-delivery-job"
description: |-
  Provides a Alicloud Actiontrail History Delivery Job resource.
---

# alicloud\_actiontrail\_history\_delivery\_job

Provides a Actiontrail History Delivery Job resource.

For information about Actiontrail History Delivery Job and how to use it, see [What is History Delivery Job](https://www.alibabacloud.com/help/doc-detail/199999.htm).

-> **NOTE:** Available in v1.139.0+.

-> **NOTE:** You are authorized to use the historical event delivery task feature. To use this feature, [submit a ticket](https://workorder-intl.console.aliyun.com/?spm=a2c63.p38356.0.0.e29f552bb6odNZ#/ticket/createIndex) or ask the sales manager to add you to the whitelist.

-> **NOTE:** Make sure that you have called the `alicloud_actiontrail_trail` to create a single-account or multi-account trace that delivered to Log Service SLS.

-> **NOTE:** An Alibaba cloud account can only have one running delivery history job at the same time.



## Example Usage

Basic Usage

```terraform
data "alicloud_regions" "example" {
  current = true
}

data "alicloud_account" "example" {}

resource "alicloud_log_project" "example" {
  name        = "example_value"
  description = "tf actiontrail test"
}


resource "alicloud_actiontrail_trail" "example" {
  trail_name      = "example_value"
  sls_project_arn = "acs:log:${data.alicloud_regions.example.regions.0.id}:${data.alicloud_account.example.id}:project/${alicloud_log_project.example.name}"
}

resource "alicloud_actiontrail_history_delivery_job" "example" {
  trail_name = alicloud_actiontrail_trail.example.id
}
```

## Argument Reference

The following arguments are supported:

* `trail_name` - (Required, ForceNew) The name of the trail for which you want to create a historical event delivery task.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of History Delivery Job.
* `status` - The status of the task. Valid values: `0`, `1`, `2`, `3`. `0`: The task is initializing. `1`: The task is delivering historical events. `2`: The delivery of historical events is complete. `3`: The task fails.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Used when create the History Delivery Job.
* `delete` - (Defaults to 2 mins) Used when delete the History Delivery Job.

## Import

Actiontrail History Delivery Job can be imported using the id, e.g.

```
$ terraform import alicloud_actiontrail_history_delivery_job.example <id>
```
