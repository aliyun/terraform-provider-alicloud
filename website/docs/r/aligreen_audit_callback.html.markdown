---
subcategory: "Aligreen"
layout: "alicloud"
page_title: "Alicloud: alicloud_aligreen_audit_callback"
description: |-
  Provides a Alicloud Aligreen Audit Callback resource.
---

# alicloud_aligreen_audit_callback

Provides a Aligreen Audit Callback resource.

Callback notification after detection is completed.

For information about Aligreen Audit Callback and how to use it, see [What is Audit Callback](https://next.api.alibabacloud.com/document/Green/2017-08-23/CreateAuditCallback).

-> **NOTE:** Available since v1.228.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform_example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_aligreen_audit_callback" "default" {
  crypt_type           = "SM3"
  audit_callback_name  = var.name
  url                  = "https://www.aliyun.com/"
  callback_types       = ["aliyunAudit", "selfAduit", "example"]
  callback_suggestions = ["block", "review", "pass"]
}
```

## Argument Reference

The following arguments are supported:
* `audit_callback_name` - (Required, ForceNew) The AuditCallback name defined by the customer. It can contain no more than 20 characters in Chinese, English, underscore (_), and digits.
* `callback_suggestions` - (Required, List) List of audit results supported by message notification. Value: block: confirmed violation, review: Suspected violation, review: normal.
* `callback_types` - (Required, List) A list of Callback types. Value: machineScan: Machine audit result notification, selfAudit: self-service audit notification.
* `crypt_type` - (Required) The encryption algorithm is used to verify that the callback request is sent by the content security service to your business service. The value is SHA256:SHA256 encryption algorithm and SM3: SM3 encryption algorithm.
* `url` - (Required) The detection result will be called back to the url.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Audit Callback.
* `delete` - (Defaults to 5 mins) Used when delete the Audit Callback.
* `update` - (Defaults to 5 mins) Used when update the Audit Callback.

## Import

Aligreen Audit Callback can be imported using the id, e.g.

```shell
$ terraform import alicloud_aligreen_audit_callback.example <id>
```