---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_alert_contact"
sidebar_current: "docs-alicloud-resource-arms-alert-contact"
description: |-
  Provides a Alicloud Application Real-Time Monitoring Service (ARMS) Alert Contact resource.
---

# alicloud_arms_alert_contact

Provides a Application Real-Time Monitoring Service (ARMS) Alert Contact resource.

For information about Application Real-Time Monitoring Service (ARMS) Alert Contact and how to use it, see [What is Alert Contact](https://next.api.aliyun.com/api/ARMS/2019-08-08/CreateAlertContact).

-> **NOTE:** Available since v1.129.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_arms_alert_contact&exampleId=87d466b0-3a49-65d7-33aa-377007af720e98c756b9&activeTab=example&spm=docs.r.arms_alert_contact.0.87d466b03a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_arms_alert_contact" "example" {
  alert_contact_name     = "example_value"
  ding_robot_webhook_url = "https://oapi.dingtalk.com/robot/send?access_token=91f2f6****"
  email                  = "someone@example.com"
  phone_num              = "1381111****"
}

```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_arms_alert_contact&spm=docs.r.arms_alert_contact.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `alert_contact_name` - (Optional) The name of the alert contact.
* `ding_robot_webhook_url` - (Optional) The webhook URL of the DingTalk chatbot. For more information about how to obtain the URL, see Configure a DingTalk chatbot to send alert notifications: https://www.alibabacloud.com/help/en/doc-detail/106247.htm. You must specify at least one of the following parameters: PhoneNum, Email, and DingRobotWebhookUrl.
* `email` - (Optional) The email address of the alert contact. You must specify at least one of the following parameters: PhoneNum, Email, and DingRobotWebhookUrl.
* `phone_num` - (Optional) The mobile number of the alert contact. You must specify at least one of the following parameters: PhoneNum, Email, and DingRobotWebhookUrl.
* `system_noc` - (Optional) Specifies whether the alert contact receives system notifications. Valid values:  true: receives system notifications. false: does not receive system notifications.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Alert Contact.

## Import

Application Real-Time Monitoring Service (ARMS) Alert Contact can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_alert_contact.example <id>
```
