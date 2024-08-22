---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alarm_contact_group"
sidebar_current: "docs-alicloud-resource-cms-alarm-contact-group"
description: |-
  Provides a Alicloud CMS Alarm Contact Group resource.
---

# alicloud_cms_alarm_contact_group

Provides a CMS Alarm Contact Group resource.

For information about CMS Alarm Contact Group and how to use it, see [What is Alarm Contact Group](https://www.alibabacloud.com/help/en/cloudmonitor/latest/putcontactgroup).

-> **NOTE:** Available since v1.101.0.

## Example Usage
<div class="oics-button" style="float: right;margin: 0 0 -40px 0;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_cms_alarm_contact_group&exampleId=a45f08a7-074d-736f-0881-2bbb7b81f5385c3a8fc7&activeTab=example&spm=docs.r.cms_alarm_contact_group.0.a45f08a707" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>

Basic Usage

```terraform
resource "alicloud_cms_alarm_contact_group" "example" {
  alarm_contact_group_name = "tf-example"
}
```

## Argument Reference

The following arguments are supported:

* `alarm_contact_group_name` - (Required, ForceNew) The name of the alarm group.
* `contacts` - (Optional) The name of the alert contact.
* `describe` - (Optional) The description of the alert group.
* `enable_subscribed` - (Optional) Whether to open weekly subscription.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Alarm Contact Group.

## Import

CMS Alarm Contact Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_alarm_contact_group.example tf-testacc123
```
