---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_user_alarm_config"
description: |-
  Provides a Alicloud Cloud Firewall User Alarm Config resource.
---

# alicloud_cloud_firewall_user_alarm_config

Provides a Cloud Firewall User Alarm Config resource.

Configure alarm notifications and contacts.

For information about Cloud Firewall User Alarm Config and how to use it, see [What is User Alarm Config](https://next.api.alibabacloud.com/document/Cloudfw/2017-12-07/DescribeUserAlarmConfig).

-> **NOTE:** Available since v1.271.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cloud_firewall_user_alarm_config" "default" {
  alarm_config {
    alarm_value    = "on"
    alarm_type     = "bandwidth"
    alarm_period   = "1"
    alarm_hour     = "0"
    alarm_notify   = "0"
    alarm_week_day = "0"
  }
  use_default_contact = "1"
  notify_config {
    notify_value = "13000000000"
    notify_type  = "sms"
  }
  alarm_lang = "zh"
  lang       = "zh"
}
```

### Deleting `alicloud_cloud_firewall_user_alarm_config` or removing it from your configuration

Terraform cannot destroy resource `alicloud_cloud_firewall_user_alarm_config`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `alarm_config` - (Required, List) The alarm configuration. More details see [`alarm_config`](#alarm_config) below.
* `alarm_lang` - (Optional) The alarm language. Possible values are `zh`, `en`.
* `contact_config` - (Optional, Computed, List) Conflict with `notify_config`. The contact configuration. More details see [`contact_config`](#contact_config) below.
* `notify_config` - (Optional, Computed, List) Conflict with `contact_config`. The notification configuration. More details see [`notify_config`](#notify_config) below.
* `lang` - (Optional) The language type. Possible values are `zh`, `en`.

  ~> **NOTE:** This parameter only applies during resource creation, update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `use_default_contact` - (Optional) Whether to Use the default contact.

  ~> **NOTE:** This parameter only applies during resource creation, update. If modified in isolation without other property changes, Terraform will not trigger any action.

---
### `alarm_config` 

The alarm_config supports the following:
* `alarm_hour` - (Optional) The time of the day when the alarm is triggered. The range is `0 ~ 24`.
* `alarm_notify` - (Optional) The alarm notification type. Possible values are: `0`(sms/email), `1`(sms), `2`(email), `3`(none)
* `alarm_period` - (Optional) The alarm period. Possible values are: `0` (8:00 ~ 20:00), `1` 24 hours.
* `alarm_type` - (Optional) The alarm type. Possible values are: `weeklyReport`, `trafficPreAlert`, `outgoingRiskAll`, `ipsMiddlethreat`, `bandwidth`, `ipsHighthreat`, `outgoingRiskNonWhite`, `ipsIgnoreResolved` etc. 
* `alarm_value` - (Optional) The alarm notification message.
* `alarm_week_day` - (Optional) The day of the week when the alarm is triggered. The range is `1 ~ 7`.

---
### `contact_config`

The contact_config supports the following:
* `email` - (Optional) The email address of the contact.
* `mobile_phone` - (Optional) The mobile phone number of the contact.
* `name` - (Optional) The name of the contact.
* `status` - (Optional) The status of the contact configuration. Possible values are: `0` disable, `1` enable.

---
### `notify_config`

The notify_config supports the following:
* `notify_type` - (Optional) The notification type. Possible values are `sms`, `mail`.
* `notify_value` - (Optional) The notification value. Depending on the value of `notify_type`, it can be a mobile phone number or an email address.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<Alibaba Cloud Account ID>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the User Alarm Config.
* `update` - (Defaults to 5 mins) Used when update the User Alarm Config.
* `read` - (Defaults to 5 mins) Used when read the User Alarm Config.

## Import

Cloud Firewall User Alarm Config can be imported using the `Account ID`, e.g.

```shell
$ terraform import alicloud_cloud_firewall_user_alarm_config.example <Alibaba Cloud Account ID>
```