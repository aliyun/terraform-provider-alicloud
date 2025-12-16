---
subcategory: "Actiontrail"
layout: "alicloud"
page_title: "Alicloud: alicloud_actiontrail_history_delivery_job"
description: |-
  Provides a Alicloud Action Trail History Delivery Job resource.
---

# alicloud_actiontrail_history_delivery_job

Provides a Action Trail History Delivery Job resource.

Delivery History Tasks.

For information about Action Trail History Delivery Job and how to use it, see [What is History Delivery Job](https://www.alibabacloud.com/help/en/actiontrail/latest/api-actiontrail-2020-07-06-createdeliveryhistoryjob).

-> **NOTE:** Available since v1.139.0.

-> **NOTE:** You are authorized to use the historical event delivery task feature. To use this feature, [submit a ticket](https://workorder-intl.console.aliyun.com/?spm=a2c63.p38356.0.0.e29f552bb6odNZ#/ticket/createIndex) or ask the sales manager to add you to the whitelist.

-> **NOTE:** Make sure that you have called the `alicloud_actiontrail_trail` to create a single-account or multi-account trace that delivered to Log Service SLS.

-> **NOTE:** An Alibaba cloud account can only have one running delivery history job at the same time.



## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_regions" "default" {
  current = true
}

data "alicloud_account" "default" {}

data "alicloud_ram_roles" "default" {
  name_regex = "AliyunServiceRoleForActionTrail"
}

resource "alicloud_log_project" "default" {
  description  = var.name
  project_name = var.name
}

resource "alicloud_actiontrail_trail" "default" {
  event_rw                = "Write"
  sls_project_arn         = "acs:log:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:project/${alicloud_log_project.default.project_name}"
  trail_name              = var.name
  sls_write_role_arn      = data.alicloud_ram_roles.default.roles.0.arn
  trail_region            = "All"
  is_organization_trail   = false
  status                  = "Enable"
  event_selectors         = jsonencode([{ "ServiceName" : "PDS" }])
  data_event_trail_region = "cn-hangzhou"
}


resource "alicloud_actiontrail_history_delivery_job" "default" {
  trail_name = alicloud_actiontrail_trail.default.id
}
```

## Argument Reference

The following arguments are supported:
* `trail_name` - (Required, ForceNew) The Track Name.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the History Delivery Job.
* `delete` - (Defaults to 5 mins) Used when delete the History Delivery Job.

## Import

Action Trail History Delivery Job can be imported using the id, e.g.

```shell
$ terraform import alicloud_actiontrail_history_delivery_job.example <id>
```