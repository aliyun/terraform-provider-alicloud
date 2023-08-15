---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_ee_instance"
sidebar_current: "docs-alicloud-resource-cr-ee-instance"
description: |-
  Provides a Alicloud resource to manage Container Registry Enterprise Edition instances.
---

# alicloud_cr_ee_instance

This resource will help you to manager Container Registry Enterprise Edition instances.

For information about Container Registry Enterprise Edition instances and how to use it, see [Create a Instance](https://www.alibabacloud.com/help/en/doc-detail/208144.htm)

-> **NOTE:** Available since v1.124.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cr_ee_instance" "default" {
  payment_type   = "Subscription"
  period         = 1
  renew_period   = 0
  renewal_status = "ManualRenewal"
  instance_type  = "Advanced"
  instance_name  = "terraform-example"
}
```

### Deleting `alicloud_cr_ee_instance` or removing it from your configuration

The `alicloud_cr_ee_instance` resource allows you to manage `payment_type = "Subscription"` cr instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the CR Instance.
You can resume managing the subscription cr instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:

* `payment_type` - (Optional, ForceNew) Subscription of Container Registry Enterprise Edition instance. Default value: `Subscription`. Valid values: `Subscription`.
* `period` - (Optional, Int) Service time of Container Registry Enterprise Edition instance. Default value: `12`. Valid values: `1`, `2`, `3`, `6`, `12`, `24`, `36`, `48`, `60`. Unit: `month`.
* `renew_period` - (Optional, ForceNew, Int) Renewal period of Container Registry Enterprise Edition instance. Unit: `month`.
* `renewal_status` - (Optional, ForceNew) Renewal status of Container Registry Enterprise Edition instance. Valid values: `AutoRenewal`, `ManualRenewal`.
* `instance_type` - (Required, ForceNew) Type of Container Registry Enterprise Edition instance. Valid values: `Basic`, `Standard`, `Advanced`. **NOTE:** International Account doesn't supports `Standard`.
* `instance_name` - (Required, ForceNew) Name of Container Registry Enterprise Edition instance.
* `custom_oss_bucket` - (Optional) Name of your customized oss bucket. Use this bucket as instance storage if set.
* `password`- (Optional, Sensitive, Available since v1.132.0) The password of the Instance. The password is a string of 8 to 30 characters and must contain uppercase letters, lowercase letters, and numbers.
* `kms_encrypted_password` - (Optional, Available since v1.132.0) An KMS encrypts password used to an instance. If the `password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional, MapString, Available since v1.132.0) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating instance with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.

## Attributes Reference

The following attributes are exported:

* `id` - ID of Container Registry Enterprise Edition instance.
* `status` - Status of Container Registry Enterprise Edition instance.
* `created_time` - Time of Container Registry Enterprise Edition instance creation.
* `end_time` - Time of Container Registry Enterprise Edition instance expiration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 min) Used when create the Instance.


## Import

Container Registry Enterprise Edition instance can be imported using the `id`, e.g.

```shell
$ terraform import alicloud_cr_ee_instance.default cri-test
```
