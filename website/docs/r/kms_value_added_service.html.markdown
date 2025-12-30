---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_value_added_service"
description: |-
  Provides a Alicloud KMS Value Added Service resource.
---

# alicloud_kms_value_added_service

Provides a KMS Value Added Service resource.

Value Added Service.

For information about KMS Value Added Service and how to use it, see [What is Value Added Service](https://next.api.alibabacloud.com/document/BssOpenApi/2017-12-14/CreateInstance).

-> **NOTE:** Available since v1.267.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_kms_value_added_service" "default" {
  value_added_service = "2"
  period              = "1"
  payment_type        = "Subscription"
  renew_period        = "1"
  renew_status        = "AutoRenewal"
}
```

## Argument Reference

The following arguments are supported:
* `payment_type` - (Required, ForceNew) The payment type of the resource
* `period` - (Optional, Int) Prepaid cycle. Unit for year

-> **NOTE:** This parameter only applies during resource creation, update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `renew_period` - (Optional, Int) Automatic renewal period, in years.

-> **NOTE:**  When setting `RenewalStatus` to `AutoRenewal`, it must be set.

* `renew_status` - (Optional) The renewal status of the specified instance. Valid values:

  - AutoRenewal: The instance is automatically renewed.
  - ManualRenewal: The instance is manually renewed.
  - NotRenewal: The instance is not renewed.
* `value_added_service` - (Optional) value added service type, Instance Backup 1 default key rotation 2 Expert service 3

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource
* `region_id` - The region ID of the resource
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Value Added Service.
* `delete` - (Defaults to 5 mins) Used when delete the Value Added Service.
* `update` - (Defaults to 5 mins) Used when update the Value Added Service.

## Import

KMS Value Added Service can be imported using the id, e.g.

```shell
$ terraform import alicloud_kms_value_added_service.example <id>
```