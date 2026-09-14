---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_webhook_contact"
description: |-
  Provides a Alicloud Application Real-Time Monitoring Service (ARMS) Webhook Contact resource.
---

# alicloud_arms_webhook_contact

Provides a Application Real-Time Monitoring Service (ARMS) Webhook Contact resource.



For information about Application Real-Time Monitoring Service (ARMS) Webhook Contact and how to use it, see [What is Webhook Contact](https://next.api.alibabacloud.com/document/ARMS/2019-08-08/CreateOrUpdateWebhookContact).

-> **NOTE:** Available since v1.294.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_arms_webhook_contact" "default" {
  webhook_contact_name = var.name

  webhook {
    method       = "Post"
    url          = "https://example.com/webhook"
    body         = jsonencode({ "message" : "alert fired" })
    recover_body = jsonencode({ "message" : "alert recovered" })
    biz_headers = {
      "Content-Type" = "application/json"
    }
  }
}
```

## Argument Reference

The following arguments are supported:
* `webhook` - (Required, ForceNew, List) The webhook object that contains the request method, URL, headers, parameters, alert notification template and recover template. See [`webhook`](#webhook) below.
* `webhook_contact_name` - (Required, ForceNew) The name of the webhook contact.

### `webhook`

The webhook supports the following:
* `biz_headers` - (Optional, ForceNew, Map) The HTTP request headers.
* `biz_params` - (Optional, ForceNew, Map) The HTTP request parameters.
* `body` - (Optional, ForceNew) The alert notification template.
* `method` - (Required, ForceNew) The HTTP request method.
* `recover_body` - (Optional, ForceNew) The alert recover template.
* `url` - (Required, ForceNew) The webhook URL.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Webhook Contact.
* `delete` - (Defaults to 5 mins) Used when delete the Webhook Contact.

## Import

Application Real-Time Monitoring Service (ARMS) Webhook Contact can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_webhook_contact.example <webhook_contact_id>
```
