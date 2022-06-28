---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_secret"
sidebar_current: "docs-alicloud-resource-kms-secret"
description: |-
  Provides a Alibaba Cloud kms secret resource.
---

# alicloud\_kms\_secret

This resouce used to create a secret and store its initial version.

-> **NOTE:** Available in 1.76.0+.

## Example Usage

Basic Usage

```
resource "alicloud_kms_secret" "default" {
  secret_name                   = "secret-foo"
  description                   = "from terraform"
  secret_data                   = "Secret data."
  version_id                    = "000000000001"
  force_delete_without_recovery = true
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the secret.
* `encryption_key_id` - (Optional, ForceNew) The ID of the KMS CMK that is used to encrypt the secret value. If you do not specify this parameter, Secrets Manager automatically creates an encryption key to encrypt the secret.
* `force_delete_without_recovery` - (Optional) Specifies whether to forcibly delete the secret. If this parameter is set to true, the secret cannot be recovered. Valid values: true, false. Default to: false.
* `recovery_window_in_days` - (Optional) Specifies the recovery period of the secret if you do not forcibly delete it. Default value: 30. It will be ignored when `force_delete_without_recovery` is true.
* `secret_data` - (Required) The value of the secret that you want to create. Secrets Manager encrypts the secret value and stores it in the initial version.
* `secret_data_type` - (Optional) The type of the secret value. Valid values: text, binary. Default to "text".
* `secret_name` - (Required, ForceNew) The name of the secret.
* `version_id` - (Required) The version number of the initial version. Version numbers are unique in each secret object.
* `version_stages` - (Optional, List(string)) The stage labels that mark the new secret version. If you do not specify this parameter, Secrets Manager marks it with "ACSCurrent".
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `enable_automatic_rotation` - (Optional, ForceNew, Available in 1.124.0+) Whether to enable automatic key rotation.
* `rotation_interval` - (Optional, ForceNew, Available in 1.124.0+) The time period of automatic rotation. The format is integer[unit], where integer represents the length of time, and unit represents the unit of time. The legal unit units are: d (day), h (hour), m (minute), s (second). 7d or 604800s both indicate a 7-day cycle.

## Attributes Reference

* `id` - The ID of the secret. It same with `secret_name`.
* `arn` - The Alicloud Resource Name (ARN) of the secret.
* `planned_delete_time` - The time when the secret is scheduled to be deleted.

### Timeouts

-> **NOTE:** Available in 1.103.2+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the KMS Secret. 
* `update` - (Defaults to 1 mins, Available in 1.105.0+) Used when updating the KMS Secret. 
* `delete` - (Defaults to 1 mins, Available in 1.105.0+) Used when deleting the KMS Secret. 

## Import

KMS secret can be imported using the id, e.g.

```
$ terraform import alicloud_kms_secret.default secret-foo
```
