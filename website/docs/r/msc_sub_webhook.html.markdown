---
subcategory: "Message Center"
layout: "alicloud"
page_title: "Alicloud: alicloud_msc_sub_webhook"
sidebar_current: "docs-alicloud-resource-msc-sub-webhook"
description: |-
  Provides a Alicloud Msc Sub Webhook resource.
---

# alicloud\_msc\_sub\_webhook

Provides a Msc Sub Webhook resource.

-> **NOTE:** Available in v1.141.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_msc_sub_webhook" "example" {
  server_url   = "example_value"
  webhook_name = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `server_url` - (Required) The serverUrl of the Webhook. This url must start with `https://oapi.dingtalk.com/robot/send?access_token=`.
* `webhook_name` - (Required) The name of the Webhook. **Note:** The name must be `2` to `12` characters in length, and can contain uppercase and lowercase letters.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Webhook.

## Import

Msc Sub Webhook can be imported using the id, e.g.

```
$ terraform import alicloud_msc_sub_webhook.example <id>
```
