---
subcategory: "Quotas"
layout: "alicloud"
page_title: "Alicloud: alicloud_quotas_quota_alarm"
description: |-
  Provides a Alicloud Quotas Quota Alarm resource.
---

# alicloud_quotas_quota_alarm

Provides a Quotas Quota Alarm resource. 

For information about Quotas Quota Alarm and how to use it, see [What is Quota Alarm](https://www.alibabacloud.com/help/en/quota-center/latest/api-doc-quotas-2020-05-10-api-doc-createquotaalarm).

-> **NOTE:** Available since v1.116.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}


resource "alicloud_quotas_quota_alarm" "default" {
  quota_action_code = "q_desktop-count"
  quota_dimensions {
    key   = "regionId"
    value = "cn-hangzhou"
  }
  threshold_percent = 80
  product_code      = "gws"
  quota_alarm_name  = var.name
  threshold_type    = "used"
}
```

## Argument Reference

The following arguments are supported:
* `product_code` - (Required, ForceNew) The Product Code.
* `quota_action_code` - (Required, ForceNew) The Quota Action Code.
* `quota_alarm_name` - (Required) The name of Quota Alarm.
* `quota_dimensions` - (Optional, ForceNew) The Quota Dimensions. See [`quota_dimensions`](#quota_dimensions) below.
* `threshold` - (Optional) The threshold of Quota Alarm.
* `threshold_percent` - (Optional) The threshold percent of Quota Alarm.
* `threshold_type` - (Optional, Computed, Available in v1.206.0+) Quota alarm type. Value:
  - used: Quota used alarm.
  - usable: alarm for the remaining available quota.
* `web_hook` - (Optional) The WebHook of Quota Alarm.


### `quota_dimensions`

The quota_dimensions supports the following:
* `key` - (Optional, ForceNew) The Key of quota_dimensions.
* `value` - (Optional, ForceNew) The Value of quota_dimensions.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Quota Alarm.
* `delete` - (Defaults to 5 mins) Used when delete the Quota Alarm.
* `update` - (Defaults to 5 mins) Used when update the Quota Alarm.

## Import

Quotas Quota Alarm can be imported using the id, e.g.

```shell
$ terraform import alicloud_quotas_quota_alarm.example <id>
```