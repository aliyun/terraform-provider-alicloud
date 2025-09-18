---
subcategory: "Message Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_message_service_service"
description: |-
  Provides a Alicloud Message Service Service resource.
---

# alicloud_message_service_service

Provides a Message Service Service resource.

MNS Service Open Status.

For information about Message Service Service and how to use it, see [What is Service](https://next.api.alibabacloud.com/document/BssOpenApi/2017-12-14/CreateInstance).

-> **NOTE:** Available since v1.252.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_message_service_service&exampleId=82ab4d0c-4144-3358-60b9-2cf3e33962745bd8ce38&activeTab=example&spm=docs.r.message_service_service.0.82ab4d0c41&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_message_service_service" "default" {
}
```

### Creating `alicloud_message_service_service`

The `alicloud_message_service_service` resource is unique per account; repeated creation attempts to activate only one instance.

### Deleting `alicloud_message_service_service` or removing it from your configuration

Terraform cannot destroy resource `alicloud_message_service_service`. Terraform will remove this resource from the state file, however resources may remain.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - The status of the Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Service.
