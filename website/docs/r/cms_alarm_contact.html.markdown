---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alarm_contact"
sidebar_current: "docs-alicloud-resource-cms-alarm-contact"
description: |-
  Provides a resource to add a alarm contact for cloud monitor.
---

# alicloud\_cms\_alarm\_contact

Creates or modifies an alarm contact. For information about alarm contact and how to use it, see [What is alarm contact](https://www.alibabacloud.com/help/en/doc-detail/114923.htm).

-> **NOTE:** Available in v1.99.0+.

## Example Usage

Basic Usage

```terraform
# If you use this template, you need to activate the link before you can return to the alarm contact information, otherwise diff will appear in terraform. So please confirm the activation link as soon as possible.
resource "alicloud_cms_alarm_contact" "example" {
  alarm_contact_name = "zhangsan"
  describe           = "For Test"
  channels_mail      = "terraform.test.com"
}
```

```terraform
# If you use this template, you can ignore the diff of the alarm contact information by `lifestyle`. We recommend the above usage and activate the link in time.
resource "alicloud_cms_alarm_contact" "example" {
  alarm_contact_name = "zhangsan"
  describe           = "For Test"
  channels_mail      = "terraform.test.com"
  lifecycle {
    ignore_changes = [channels_mail]
  }
}
```

## Argument Reference

The following arguments are supported:

* `alarm_contact_name` - (Required, ForceNew) The name of the alarm contact. The length should between 2 and 40 characters.
* `channels_aliim` - (Optional) The TradeManager ID of the alarm contact.
* `channels_ding_web_hook` - (Optional) The webhook URL of the DingTalk chatbot.
* `channels_mail` - (Optional) The email address of the alarm contact. After you add or modify an email address, the recipient receives an email that contains an activation link. The system adds the recipient to the list of alarm contacts only after the recipient activates the email address.
* `channels_sms` - (Optional) The phone number of the alarm contact. After you add or modify an email address, the recipient receives an email that contains an activation link. The system adds the recipient to the list of alarm contacts only after the recipient activates the email address.
* `describe` - (Required) The description of the alarm contact.
* `lang` - (Optional) The language type of the alarm. Valid values: `en`, `zh-cn`.

-> **NOTE:** Specify at least one of the following alarm notification targets: `channels_aliim`, `channels_ding_web_hook`, `channels_mail`, `channels_sms`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the alarm contact. The value is same with `alarm_contact_name`.

## Import

Alarm contact can be imported using the id, e.g.

```
$ terraform import alicloud_cms_alarm_contact.example abc12345
```
