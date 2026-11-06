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

-> **NOTE:** Available since v2.2.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_arms_webhook_contact&exampleId=04057dbb-8a37-ba32-24f0-1b4f69a97758c08982e7&activeTab=example&spm=docs.r.arms_webhook_contact.0.04057dbb8a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_arms_webhook_contact" "default" {
  webhook_contact_name = var.name
  method               = "Post"
  url                  = "https://example.com/webhook"
  body                 = jsonencode({ "message" : "alert fired" })
  recover_body         = jsonencode({ "message" : "alert recovered" })
  biz_headers = {
    "Content-Type" = "application/json"
  }
  biz_params = {
    "foo" = "bar"
  }
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_arms_webhook_contact&spm=docs.r.arms_webhook_contact.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `biz_headers` - (Optional, ForceNew, Map) The HTTP request headers. Each value must be a string.
* `biz_params` - (Optional, ForceNew, Map) The HTTP request parameters. Each value must be a string.
* `body` - (Optional, ForceNew, String) The alert notification template. **Note: This parameter is required when `method` is set to `Post`.**
* `method` - (Required, ForceNew, String) The HTTP request method. Valid values: `Get`, `Post`. **Note: The parameter is immutable after resource creation.**
* `recover_body` - (Optional, ForceNew, String) The alert recover template. **Note: This parameter is required when `method` is set to `Post`.**
* `url` - (Required, ForceNew, String) The webhook URL.
* `webhook_contact_name` - (Required, ForceNew, String) The name of the webhook contact.

## Attributes Reference

The following attributes are exported:
* `id` - (Computed, String) The ID of the resource supplied above.
* `webhook_contact_id` - (Computed, String) The ID of the webhook contact.

## Import

Application Real-Time Monitoring Service (ARMS) Webhook Contact can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_webhook_contact.example <webhook_contact_id>
```
