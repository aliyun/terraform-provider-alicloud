---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_application_access_point"
description: |-
  Provides a Alicloud KMS Application Access Point resource.
---

# alicloud_kms_application_access_point

Provides a KMS Application Access Point resource. An application access point (AAP) is used to implement fine-grained access control for Key Management Service (KMS) resources. An application can access a KMS instance only after an AAP is created for the application. .

For information about KMS Application Access Point and how to use it, see [What is Application Access Point](https://www.alibabacloud.com/help/zh/key-management-service/latest/api-createapplicationaccesspoint).

-> **NOTE:** Available since v1.210.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_kms_application_access_point&exampleId=f1d4768e-178f-2eb4-6a66-b64d8b2668f0aeda76ba&activeTab=example&spm=docs.r.kms_application_access_point.0.f1d4768e17&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}


resource "alicloud_kms_application_access_point" "default" {
  description                   = "example aap"
  application_access_point_name = var.name
  policies                      = ["abc", "efg", "hfc"]
}
```

## Argument Reference

The following arguments are supported:
* `application_access_point_name` - (Required, ForceNew) Application Access Point Name.
* `description` - (Optional) Description .
* `policies` - (Required) The policies that have bound to the Application Access Point (AAP).

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Application Access Point.
* `delete` - (Defaults to 5 mins) Used when delete the Application Access Point.
* `update` - (Defaults to 5 mins) Used when update the Application Access Point.

## Import

KMS Application Access Point can be imported using the id, e.g.

```shell
$ terraform import alicloud_kms_application_access_point.example <id>
```