---
subcategory: "CDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_cdn_fc_trigger"
sidebar_current: "docs-alicloud-resource-cdn-fc-trigger"
description: |-
  Provides a Alicloud CDN Fc Trigger resource.
---

# alicloud\_cdn\_fc\_trigger

Provides a CDN Fc Trigger resource.

For information about CDN Fc Trigger and how to use it, see [What is Fc Trigger](https://www.alibabacloud.com/help/zh/alibaba-cloud-cdn/latest/add-function-calculation-trigger).

-> **NOTE:** Available in v1.165.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_account" "default" {}

data "alicloud_regions" "default" {
  current = true
}
resource "alicloud_cdn_fc_trigger" "example" {
  event_meta_name    = "LogFileCreated"
  event_meta_version = "1.0.0"
  notes              = "example_value"
  role_arn           = "acs:ram::${data.alicloud_account.default.id}:role/aliyuncdneventnotificationrole"
  source_arn         = "acs:cdn:*:${data.alicloud_account.default.id}:domain/example.com"
  trigger_arn        = "acs:fc:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:services/FCTestService/functions/printEvent/triggers/testtrigger"
}
```

## Argument Reference

The following arguments are supported:

* `event_meta_name` - (Required, ForceNew) The name of the Event.
* `event_meta_version` - (Required, ForceNew) The version of the Event.
* `function_arn` - (Optional) The function arn. The value formats as `acs:fc:{RegionID}:{AccountID}:{Filter}`.
* `notes` - (Required) The Note information.
* `role_arn` - (Required) The role authorized by RAM. The value formats as `acs:ram::{AccountID}:role/{RoleName}`.
* `source_arn` - (Required) Resources and filters for event listening. The value formats as `acs:cdn:{RegionID}:{AccountID}:{Filter}`.
* `trigger_arn` - (Required, ForceNew) The trigger corresponding to the function Compute Service. The value formats as `acs:fc:{RegionID}:{AccountID}:{Filter}`. See [Create a CDN Fc Trigger](https://www.alibabacloud.com/help/zh/alibaba-cloud-cdn/latest/add-function-calculation-trigger) for more details.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Fc Trigger. Its value is same as `trigger_arn`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Fc Trigger.
* `update` - (Defaults to 1 mins) Used when updating the Fc Trigger.
* `delete` - (Defaults to 1 mins) Used when deleting the Fc Trigger.

## Import

CDN Fc Trigger can be imported using the id, e.g.

```
$ terraform import alicloud_cdn_fc_trigger.example <trigger_arn>
```