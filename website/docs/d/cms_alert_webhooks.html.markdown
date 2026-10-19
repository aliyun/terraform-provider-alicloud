---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alert_webhooks"
sidebar_current: "docs-alicloud-datasource-cms-alert-webhooks"
description: |-
  Provides a list of Cloud Monitor Service Alert Webhooks to the user.
---

# alicloud\_cms\_alert\_webhooks

This data source provides the Cloud Monitor Service Alert Webhooks available to the user.

-> **NOTE:** Available since v1.245.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_alert_webhooks" "default" {
  name_regex = "tf-example"
}

output "first_webhook_id" {
  value = data.alicloud_cms_alert_webhooks.default.webhooks.0.alert_webhook_id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of alert webhook ids.
* `name_regex` - (Optional) A regex string to filter results by the alert webhook name.
* `alert_webhook_name` - (Optional) The name of the alert webhook to filter by.
* `workspace` - (Optional) The workspace to filter by.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of alert webhook ids.
* `names` - A list of alert webhook names.
* `webhooks` - A list of alert webhooks. Each element contains the following attributes:
  * `id` - The ID of the alert webhook.
  * `alert_webhook_id` - The ID of the alert webhook.
  * `alert_webhook_name` - The name of the alert webhook.
  * `content_type` - The content type of the alert webhook request body.
  * `headers` - The custom request headers of the alert webhook.
  * `lang` - The language of the alert webhook.
  * `method` - The HTTP method of the alert webhook request.
  * `url` - The URL of the alert webhook.
  * `workspace` - The workspace to which the alert webhook belongs.
