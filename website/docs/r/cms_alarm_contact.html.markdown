---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alarm_contact"
sidebar_current: "docs-alicloud-resource-cms-alarm-contact"
description: |-
  Provides a resource to add a alarm contact for cloud monitor.
---

# alicloud_cms_alarm_contact

Creates or modifies an alarm contact. For information about alarm contact and how to use it, see [What is alarm contact](https://www.alibabacloud.com/help/en/cloudmonitor/latest/putcontact).

-> **NOTE:** Available since v1.99.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cms_alarm_contact&exampleId=b511846b-0994-7e9f-b042-04a42a552d33a0f0c35e&activeTab=example&spm=docs.r.cms_alarm_contact.0.b511846b09&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
# You need to activate the link before you can return to the alarm contact information, otherwise diff will appear in terraform. So please confirm the activation link as soon as possible. Besides, you can ignore the diff of the alarm contact information by `lifestyle`. 
resource "alicloud_cms_alarm_contact" "example" {
  alarm_contact_name = "tf-example"
  describe           = "For example"
  channels_mail      = "terraform@test.com"
  lifecycle {
    ignore_changes = [channels_mail]
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cms_alarm_contact&spm=docs.r.cms_alarm_contact.example&intl_lang=EN_US)

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

```shell
$ terraform import alicloud_cms_alarm_contact.example abc12345
```
