---
subcategory: "Actiontrail"
layout: "alicloud"
page_title: "Alicloud: alicloud_actiontrail_history_delivery_job"
sidebar_current: "docs-alicloud-resource-actiontrail-history-delivery-job"
description: |-
  Provides a Alicloud Actiontrail History Delivery Job resource.
---

# alicloud_actiontrail_history_delivery_job

Provides a Actiontrail History Delivery Job resource.

For information about Actiontrail History Delivery Job and how to use it, see [What is History Delivery Job](https://www.alibabacloud.com/help/en/actiontrail/latest/api-actiontrail-2020-07-06-createdeliveryhistoryjob).

-> **NOTE:** Available since v1.139.0.

-> **NOTE:** You are authorized to use the historical event delivery task feature. To use this feature, [submit a ticket](https://workorder-intl.console.aliyun.com/?spm=a2c63.p38356.0.0.e29f552bb6odNZ#/ticket/createIndex) or ask the sales manager to add you to the whitelist.

-> **NOTE:** Make sure that you have called the `alicloud_actiontrail_trail` to create a single-account or multi-account trace that delivered to Log Service SLS.

-> **NOTE:** An Alibaba cloud account can only have one running delivery history job at the same time.



## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_actiontrail_history_delivery_job&exampleId=7d41ea75-d293-4450-9a50-bb747d7fab118cd7c6e5&activeTab=example&spm=docs.r.actiontrail_history_delivery_job.0.7d41ea75d2&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
resource "random_integer" "default" {
  min = 10000
  max = 99999
}
data "alicloud_regions" "example" {
  current = true
}
data "alicloud_account" "example" {}

resource "alicloud_log_project" "example" {
  project_name = "${var.name}-${random_integer.default.result}"
  description  = "tf actiontrail example"
}

resource "alicloud_actiontrail_trail" "example" {
  trail_name      = "${var.name}-${random_integer.default.result}"
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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Used when create the History Delivery Job.
* `delete` - (Defaults to 2 mins) Used when delete the History Delivery Job.

## Import

Actiontrail History Delivery Job can be imported using the id, e.g.

```shell
$ terraform import alicloud_actiontrail_history_delivery_job.example <id>
```
