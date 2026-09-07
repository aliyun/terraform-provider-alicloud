---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alert_webhook"
description: |-
  Provides a Alicloud Cms Alert Webhook resource.
---

# alicloud_cms_alert_webhook

Provides a Cms Alert Webhook resource.

The webhook notification target used by CloudMonitor alerting.

For information about Cms Alert Webhook and how to use it, see [What is Alert Webhook](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateAlertWebhook).

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
  content_type       = "JSON"
  method             = "POST"
  lang               = "zh_CN"
  headers = {
    X-Custom-Token = "example-token"
  }
}
```

## Argument Reference

The following arguments are supported:
* `alert_webhook_name` - (Required) The name of the alert webhook.
* `url` - (Required) The URL of the alert webhook.
* `content_type` - (Optional) The data format of the webhook request body. Valid values: `JSON` and `FORM`. Default value: `JSON`.
* `method` - (Optional) The request method used to call the webhook. Valid values: `GET` and `POST`. Default value: `POST`.
* `lang` - (Optional) The language of the notification content. Valid values: `zh_CN` and `en_US`.
* `headers` - (Optional, Map) The custom headers of the webhook request. Each header is a key-value pair of strings.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `id` - The resource ID in terraform of the alert webhook. It is the same as the alert webhook ID.
* `workspace` - The workspace to which the alert webhook belongs.

## Timeouts

-> **NOTE:** Available since v1.292.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Alert Webhook.
* `update` - (Defaults to 5 mins) Used when update the Alert Webhook.
* `delete` - (Defaults to 5 mins) Used when delete the Alert Webhook.

## Import

Cms Alert Webhook can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_alert_webhook.example <id>
```
