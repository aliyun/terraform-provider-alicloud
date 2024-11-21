---
subcategory: "Message Center (MscSub)"
layout: "alicloud"
page_title: "Alicloud: alicloud_msc_sub_webhook"
sidebar_current: "docs-alicloud-resource-msc-sub-webhook"
description: |-
  Provides a Alicloud Msc Sub Webhook resource.
---

# alicloud_msc_sub_webhook

Provides a Msc Sub Webhook resource.

-> **NOTE:** Available since v1.141.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_msc_sub_webhook&exampleId=50c682fd-ad2a-60b7-b141-135c461765a4628f6c0e&activeTab=example&spm=docs.r.msc_sub_webhook.0.50c682fdad&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tfexample"
}
variable "token" {
  default = "abcd****"
}
resource "alicloud_msc_sub_webhook" "example" {
  server_url   = format("https://oapi.dingtalk.com/robot/send?access_token=%s", var.token)
  webhook_name = var.name
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

```shell
$ terraform import alicloud_msc_sub_webhook.example <id>
```
