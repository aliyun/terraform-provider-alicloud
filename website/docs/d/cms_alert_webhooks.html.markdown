---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alert_webhooks"
description: |-
  Provides a list of Cms Alert Webhooks to the user.
---

# alicloud_cms_alert_webhooks

This data source provides the Cms Alert Webhooks of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.292.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_cms_alert_webhook" "default" {
  alert_webhook_name = "${var.name}-${random_integer.default.result}"
  url                = "https://example.com/alert-webhook"
}

data "alicloud_cms_alert_webhooks" "ids" {
  ids = [alicloud_cms_alert_webhook.default.id]
}

output "cms_alert_webhook_id" {
  value = data.alicloud_cms_alert_webhooks.ids.webhooks.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of Alert Webhook IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Alert Webhook name.
* `alert_webhook_name` - (Optional, ForceNew) The name of the alert webhook used to filter results.
* `workspace` - (Optional, ForceNew) The workspace used to filter results.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Alert Webhook names.
* `webhooks` - A list of Alert Webhooks. Each element contains the following attributes:
  * `id` - The ID of the Alert Webhook.
  * `alert_webhook_id` - The ID of the Alert Webhook.
  * `alert_webhook_name` - The name of the Alert Webhook.
  * `content_type` - The data format of the webhook request body.
  * `headers` - The custom headers of the webhook request.
  * `lang` - The language of the notification content.
  * `method` - The request method used to call the webhook.
  * `url` - The URL of the Alert Webhook.
  * `workspace` - The workspace to which the Alert Webhook belongs.
