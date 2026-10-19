---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alert_webhook"
sidebar_current: "docs-alicloud-resource-cms-alert-webhook"
description: |-
  Provides a Cloud Monitor Service Alert Webhook resource.
---

# alicloud\_cms\_alert\_webhook

Provides a Cloud Monitor Service Alert Webhook resource.

For information about Cloud Monitor Service Alert Webhook and how to use it, see [What is Alert Webhook](https://www.alibabacloud.com/help/en/cloudmonitor/developer-reference/api-cms-2024-03-30-createalertwebhook).

-> **NOTE:** Available since v1.245.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

resource "alicloud_cms_alert_webhook" "default" {
  alert_webhook_name = var.name
  url                = "https://example.com/alert-webhook"
  content_type       = "JSON"
  method             = "POST"
  lang               = "en_US"
  headers = {
    X-Token = "example-token"
  }
}
```

## Argument Reference

The following arguments are supported:

* `alert_webhook_name` - (Required) The name of the alert webhook.
* `url` - (Required) The URL of the alert webhook.
* `content_type` - (Optional) The content type of the alert webhook request body. Valid values: `JSON`, `FORM`.
* `method` - (Optional) The HTTP method of the alert webhook request. Valid values: `GET`, `POST`.
* `lang` - (Optional) The language. Valid values: `zh_CN`, `en_US`.
* `headers` - (Optional) The custom request headers of the alert webhook.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the alert webhook.
* `workspace` - The workspace to which the alert webhook belongs.

## Import

Cloud Monitor Service Alert Webhook can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_alert_webhook.example <webhook-id>
```
