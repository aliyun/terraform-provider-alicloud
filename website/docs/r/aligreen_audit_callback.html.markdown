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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_aligreen_audit_callback&exampleId=87797eac-9b31-e712-8730-a951d64e143e4f1939ef&activeTab=example&spm=docs.r.aligreen_audit_callback.0.87797eac9b&intl_lang=EN_US" target="_blank">
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


resource "alicloud_aligreen_audit_callback" "default" {
  crypt_type           = "SM3"
  audit_callback_name  = var.name
  url                  = "https://www.aliyun.com/"
  callback_types       = ["aliyunAudit", "selfAduit", "example"]
  callback_suggestions = ["block", "review", "pass"]
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_aligreen_audit_callback&spm=docs.r.aligreen_audit_callback.example&intl_lang=EN_US)

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

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Audit Callback.
* `delete` - (Defaults to 5 mins) Used when delete the Audit Callback.
* `update` - (Defaults to 5 mins) Used when update the Audit Callback.

## Import

Aligreen Audit Callback can be imported using the id, e.g.

```shell
$ terraform import alicloud_aligreen_audit_callback.example <id>
```