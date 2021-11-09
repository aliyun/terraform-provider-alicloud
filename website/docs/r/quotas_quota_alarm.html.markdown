---
subcategory: "Quotas"
layout: "alicloud"
page_title: "Alicloud: alicloud_quotas_quota_alarm"
sidebar_current: "docs-alicloud-resource-quotas-quota-alarm"
description: |-
  Provides a Alicloud Quotas Quota Alarm resource.
---

# alicloud\_quotas\_quota\_alarm

Provides a Quotas Quota Alarm resource.

For information about Quotas Quota Alarm and how to use it, see [What is Quota Alarm](https://help.aliyun.com/document_detail/184343.html).

-> **NOTE:** Available in v1.116.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_quotas_quota_alarm" "example" {
  quota_alarm_name  = "tf-testAcc"
  product_code      = "ecs"
  quota_action_code = "q_prepaid-instance-count-per-once-purchase"
  threshold         = "100"
  quota_dimensions {
    key   = "regionId"
    value = "cn-hangzhou"
  }
}

```

## Argument Reference

The following arguments are supported:

* `quota_alarm_name` - (Required) The name of Quota Alarm.
* `product_code` - (Required, ForceNew) The Product Code.
* `quota_action_code` - (Required, ForceNew) The Quota Action Code.
* `quota_dimensions` - (Optional, ForceNew) The Quota Dimensions.
* `threshold` - (Optional) The threshold of Quota Alarm.
* `threshold_percent` - (Optional) The threshold percent of Quota Alarm.
* `web_hook` - (Optional) The WebHook of Quota Alarm.

#### Block quota_dimensions

The quota_dimensions supports the following: 

* `key` - (Optional, ForceNew) The Key of quota_dimensions.
* `value` - (Optional, ForceNew) The Value of quota_dimensions.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Quota Alarm.

## Import

Quotas Quota Alarm can be imported using the id, e.g.

```
$ terraform import alicloud_quotas_quota_alarm.example <id>
```
