---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_alert_contact_group"
sidebar_current: "docs-alicloud-resource-arms-alert-contact-group"
description: |-
  Provides a Alicloud Application Real-Time Monitoring Service (ARMS) Alert Contact Group resource.
---

# alicloud_arms_alert_contact_group

Provides a Application Real-Time Monitoring Service (ARMS) Alert Contact Group resource.

For information about Application Real-Time Monitoring Service (ARMS) Alert Contact Group and how to use it, see [What is Alert Contact Group](https://next.api.aliyun.com/api/ARMS/2019-08-08/CreateAlertContactGroup).

-> **NOTE:** Available since v1.131.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_arms_alert_contact_group&exampleId=bb18d5b9-0088-3372-d9c0-fc554d0d65c5fecbdba1&activeTab=example&spm=docs.r.arms_alert_contact_group.0.bb18d5b900&intl_lang=EN_US" target="_blank">
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
resource "alicloud_arms_alert_contact_group" "example" {
  alert_contact_group_name = "example_value"
  contact_ids              = [alicloud_arms_alert_contact.example.id]
}

```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_arms_alert_contact_group&spm=docs.r.arms_alert_contact_group.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `alert_contact_group_name` - (Required) The name of the resource.
* `contact_ids` - (Optional) The list id of alert contact.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Alert Contact Group.

## Import

Application Real-Time Monitoring Service (ARMS) Alert Contact Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_alert_contact_group.example <id>
```
