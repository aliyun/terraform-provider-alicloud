---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_data_quality_alert_rule"
description: |-
  Provides a Alicloud Data Works Data Quality Alert Rule resource.
---

# alicloud_data_works_data_quality_alert_rule

Provides a Data Works Data Quality Alert Rule resource.

A data quality alert rule notifies recipients through the configured channels
when a DataQualityScan result matches the rule condition.

For information about Data Works Data Quality Alert Rule and how to use it, see [What is Data Quality Alert Rule](https://www.alibabacloud.com/help/en/dataworks/developer-reference/api-dataworks-public-2024-05-18-createdataqualityalertrule).

-> **NOTE:** Available since v2.2.0.

-> **NOTE:** The `target.ids` list references a DataQualityScan. DataQualityScan
has no Terraform resource yet; the scan id must reference a pre-existing scan
that belongs to the same project as `project_id`.

## Example Usage

Basic Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "terraform-example"
}

# The DataQualityScan referenced by target.ids has no Terraform resource yet;
# replace this with the id of a pre-existing DataQualityScan that belongs to
# the same DataWorks project as project_id.
variable "data_quality_scan_id" {
  type = number
}

provider "alicloud" {
  region = "cn-chengdu"
}

resource "alicloud_data_works_project" "default" {
  project_name     = var.name
  display_name     = var.name
  description      = var.name
  pai_task_enabled = true
}

resource "alicloud_data_works_data_quality_alert_rule" "default" {
  project_id = alicloud_data_works_project.default.id
  condition  = "results.any { r -> r.status == 'fail' && r.rule.severity == 'High' }"
  target {
    type = "DataQualityScan"
    ids  = [var.data_quality_scan_id]
  }
  notification {
    channels = ["Mail"]
    receivers {
      receiver_type = "TaskOwner"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, ForceNew, Int) The id of the DataWorks project that owns the alert rule and its monitoring target.
* `condition` - (Required, String) The alert condition expression evaluated against DataQualityScan results, e.g. `results.any { r -> r.status == 'fail' && r.rule.severity == 'High' }`.
* `target` - (Required) The monitoring target of the alert rule. See [`target`](#target) below.
* `notification` - (Required) The alert notification configuration. See [`notification`](#notification) below.

### `target`

The target supports the following:

* `type` - (Required, String) The monitoring target type. Only `DataQualityScan` is supported.
* `ids` - (Required, List of Int) The monitoring target id list. Only one id is supported; the referenced DataQualityScan must belong to the same project as `project_id`.

### `notification`

The notification supports the following:

* `channels` - (Required, List) The alert channel list. Except that `Mail` and `Sms` may be set together, only one channel is supported. Valid values: `Dingding`, `Weixin`, `Mail`, `Sms`, `Phone`, `Webhook`, `Feishu`.
* `receivers` - (Required, List) The alert receivers. See [`receivers`](#notification-receivers) below.

### `notification-receivers`

The receivers supports the following:

* `receiver_type` - (Required, String) The receiver type. Valid values: `AliUid`, `WebhookUrl`, `DingdingUrl`, `WeixinUrl`, `FeishuUrl`, `TaskOwner`, `DataQualityScanOwner`, `ShiftSchedule`.
* `receiver_values` - (Optional, List) The receiver values. Required when `receiver_type` is `AliUid` (an Alibaba Cloud account UID), `ShiftSchedule` (a shift schedule id) or one of the webhook URL types (`WebhookUrl`, `DingdingUrl`, `FeishuUrl`, `WeixinUrl`); leave empty for `TaskOwner` and `DataQualityScanOwner`.
* `extension` - (Optional, String) Extra receiver configuration. When `receiver_type` is `DingdingUrl`, set `{"atAll":true}` to @ all members.

## Attributes Reference

The following attributes are exported:

* `id` - (String) The id of the data quality alert rule.
