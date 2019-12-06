---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_key"
sidebar_current: "docs-alicloud-resource-kms-key"
description: |-
  Provides a Alicloud kms key resource.
---

# alicloud\_kms\_key

A kms key can help user to protect data security in the transmission process.

## Example Usage

Basic Usage

```
resource "alicloud_kms_key" "key" {
  description             = "Hello KMS"
  deletion_window_in_days = "7"
  is_enabled              = true
}
```
## Argument Reference

The following arguments are supported:

* `description` - (Optional, ForceNew) The description of the key as viewed in Alicloud console. Default to "From Terraform".
* `key_usage` - (Optional) Specifies the usage of CMK. Currently, default to 'ENCRYPT/DECRYPT', indicating that CMK is used for encryption and decryption.
* `deletion_window_in_days` - (Optional) Duration in days after which the key is deleted
	after destruction of the resource, must be between 7 and 30 days. Defaults to 30 days.
* `is_enabled` - (Optional) Specifies whether the key is enabled. Defaults to true.

-> **NOTE:** At present, the resource only supports to modify `is_enabled`.

-> **NOTE:** When the pre-deletion days elapses, the key is permanently deleted and cannot be recovered.


## Attributes Reference

* `id` - The ID of the key.
* `arn` - The Alicloud Resource Name (ARN) of the key.
* `description` - The description of the key.
* `key_usage` - (ForceNew) Specifies the usage of CMK.
* `deletion_window_in_days` - During pre-deletion days.
* `is_enabled` - Whether the key is enabled.


## Import

KMS key can be imported using the id, e.g.

```
$ terraform import alicloud_kms_key.example abc123456
```
