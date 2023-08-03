---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_instance"
description: |-
  Provides a Alicloud KMS Instance resource.
---

# alicloud_kms_instance

Provides a KMS Instance resource. User-exclusive KMS instances.

For information about KMS Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.209.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}


resource "alicloud_kms_instance" "default" {
  vpc_num         = "1"
  key_num         = "1000"
  spec            = "1000"
  renew_status    = "AutoRenewal"
  product_type    = "kms_ddi_public_cn"
  product_version = "3"
  renew_period    = "3"
}
```

### Deleting `alicloud_kms_instance` or removing it from your configuration

Terraform cannot destroy resource `alicloud_kms_instance`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `key_num` - (Required) Maximum number of stored keys.
* `payment_type` - (Optional, ForceNew, Computed) Payment type is Prepaid.
* `product_type` - (Required) KMS instance type (China /International ).
* `product_version` - (Optional) KMS Instance commodity type (software/hardware).
* `renew_period` - (Optional) Automatic renewal period, in months.
* `renew_status` - (Optional) Renewal options (manual renewal, automatic renewal, no renewal).
* `secret_num` - (Required) Maximum number of Secrets.
* `spec` - (Required) The computation performance level of the KMS instance.
* `vpc_num` - (Required) The number of managed accesses. The maximum number of VPCs that can access this KMS instance.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.
* `instance_name` - The name of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Instance.
* `update` - (Defaults to 5 mins) Used when update the Instance.

## Import

KMS Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_kms_instance.example <id>
```