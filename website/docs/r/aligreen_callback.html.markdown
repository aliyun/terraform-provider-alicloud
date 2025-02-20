---
subcategory: "Aligreen"
layout: "alicloud"
page_title: "Alicloud: alicloud_aligreen_callback"
description: |-
  Provides a Alicloud Aligreen Callback resource.
---

# alicloud_aligreen_callback

Provides a Aligreen Callback resource.

Detection Result Callback.

For information about Aligreen Callback and how to use it, see [What is Callback](https://next.api.alibabacloud.com/document/Green/2017-08-23/CreateCallback).

-> **NOTE:** Available since v1.228.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_aligreen_callback&exampleId=eb2f06de-71b9-821a-e3d7-bf40f41ed53c301f5190&activeTab=example&spm=docs.r.aligreen_callback.0.eb2f06de71&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform_example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_aligreen_callback" "default" {
  callback_url         = "https://www.aliyun.com"
  crypt_type           = "0"
  callback_name        = var.name
  callback_types       = ["machineScan", "selfAudit", "example"]
  callback_suggestions = ["block", "review", "pass"]
}
```

## Argument Reference

The following arguments are supported:
* `callback_name` - (Required) The Callback name defined by the customer. It can contain no more than 20 characters in Chinese, English, underscore (_), and digits.
* `callback_suggestions` - (Required, List) List of audit results supported by message notification. Value: block: confirmed violation, review: Suspected violation, review: normal.
* `callback_types` - (Required, List) A list of Callback types. Value: machineScan: Machine audit result notification, selfAudit: self-service audit notification.
* `callback_url` - (Required) The detection result will be called back to the url.
* `crypt_type` - (Optional, ForceNew, Int) The encryption algorithm is used to verify that the callback request is sent by the Aliyun Green Service to your business service. Value: 0:SHA256,1: SM3.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the Callback.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Callback.
* `delete` - (Defaults to 5 mins) Used when delete the Callback.
* `update` - (Defaults to 5 mins) Used when update the Callback.

## Import

Aligreen Callback can be imported using the id, e.g.

```shell
$ terraform import alicloud_aligreen_callback.example <id>
```