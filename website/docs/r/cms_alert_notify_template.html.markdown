---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alert_notify_template"
description: |-
  Provides a Cloud Monitor Alert Notify Template resource.
---

# alicloud_cms_alert_notify_template

Provides a Cloud Monitor Alert Notify Template resource.

For information about Cloud Monitor Alert Notify Template and how to use it, see [Manage notification templates](https://www.alibabacloud.com/help/en/cms/user-guide/manage-notification-templates).

-> **NOTE:** Available since v1.292.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cms_alert_notify_template" "default" {
  alert_notify_template_id   = "my-template-id"
  alert_notify_template_name = "my-template-name"
  type                       = "DING"
  program_lang               = "zh_CN"
  templates {
    channel = "ding"
    title   = "Alert Title"
    content = "Alert Content"
  }
}
```

## Argument Reference

The following arguments are supported:

* `alert_notify_template_id` - (Required, ForceNew) The ID of the alert notify template. It must be unique within the region.
* `alert_notify_template_name` - (Optional) The name of the alert notify template.
* `templates` - (Optional, List) The notification channel templates. Each element is an object representing a single channel. See [`templates`](#templates) below.
* `type` - (Optional, ForceNew) The channel type of the alert notify template. Valid values: `DING`, `WEIXIN`, `FEISHU`, `SLACK`, `CUSTOM`, `TEAMS`, `SMS`, `EMAIL`, `CALL`, `PAGER_DUTY`.
* `program_lang` - (Optional, ForceNew) The programming language of the alert notify template. Valid values: `zh_CN`, `en_US`.

### `templates`

The templates supports the following:

* `channel` - (Required) The channel key of the template, which identifies the notification channel this template content applies to, e.g. `ding`.
* `title` - (Optional) The title of the alert notification template.
* `content` - (Optional) The content of the alert notification template.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of the alert notify template. The value is the same as `alert_notify_template_id`.
* `region_id` - The region ID of the alert notify template.

## Import

Cloud Monitor Alert Notify Template can be imported using the id, e.g.

```shell
terraform import alicloud_cms_alert_notify_template.default my-template-id
```
