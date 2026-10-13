---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alert_notify_templates"
description: |-
  Provides a list of Cloud Monitor Alert Notify Templates to the user according to the specified filters.
---

# alicloud_cms_alert_notify_templates

This data source provides a list of Cloud Monitor Alert Notify Templates in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available since v1.292.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_alert_notify_templates" "default" {
  ids = ["my-template-id"]
}

output "first_template_name" {
  value = data.alicloud_cms_alert_notify_templates.default.templates.0.alert_notify_template_name
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of alert notify template IDs. If this parameter is set, the data source will only return the templates whose IDs match the specified values.
* `alert_notify_template_name` - (Optional) The name of the alert notify template. The data source will only return templates whose names match the specified value.
* `program_lang` - (Optional) The programming language of the alert notify template. Valid values: `zh_CN`, `en_US`.
* `type` - (Optional) The channel type of the alert notify template. Valid values: `DING`, `WEIXIN`, `FEISHU`, `SLACK`, `CUSTOM`, `TEAMS`, `SMS`, `EMAIL`, `CALL`, `PAGER_DUTY`.
* `output_file` - (Optional) File path where to save the results (in JSON format). If not specified, the results are not saved to a file.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of alert notify template IDs.
* `templates` - A list of alert notify templates. Each element contains the following attributes:
  * `alert_notify_template_id` - The ID of the alert notify template.
  * `alert_notify_template_name` - The name of the alert notify template.
  * `type` - The channel type of the alert notify template.
  * `program_lang` - The programming language of the alert notify template.
  * `templates` - The notification channel templates. Each element contains the following attributes:
    * `channel` - The channel key of the alert notify template.
    * `title` - The title of the alert notification template.
    * `content` - The content of the alert notification template.
