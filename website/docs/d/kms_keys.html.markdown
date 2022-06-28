---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_keys"
sidebar_current: "docs-alicloud-datasource-kms-keys"
description: |-
    Provides a list of available KMS Keys.
---

# alicloud\_kms\_keys

This data source provides a list of KMS keys in an Alibaba Cloud account according to the specified filters.

## Example Usage

```
# Declare the data source
data "alicloud_kms_keys" "kms_keys_ds" {
  description_regex = "Hello KMS"
  output_file       = "kms_keys.json"
}

output "first_key_id" {
  value = "${data.alicloud_kms_keys.kms_keys_ds.keys.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of KMS key IDs.
* `description_regex` - (Optional) A regex string to filter the results by the KMS key description.
* `status` - (Optional) Filter the results by status of the KMS keys. Valid values: `Enabled`, `Disabled`, `PendingDeletion`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` -  A list of KMS key IDs.
* `keys` - A list of KMS keys. Each element contains the following attributes:
  * `id` - ID of the key.
  * `arn` - The Alibaba Cloud Resource Name (ARN) of the key.
  * `description` - Description of the key.
  * `status` - Status of the key. Possible values: `Enabled`, `Disabled` and `PendingDeletion`.
  * `creation_date` - Creation date of key.
  * `delete_date` - Deletion date of key.
  * `creator` - The owner of the key.
  * `automatic_rotation` -(Available in 1.123.1+) Specifies whether to enable automatic key rotation.
  * `key_id` -(Available in 1.123.1+)  ID of the key.
  * `key_spec` -(Available in 1.123.1+)  The type of the CMK.
  * `key_usage` -(Available in 1.123.1+)  The usage of CMK.
  * `last_rotation_date` -(Available in 1.123.1+)  The date and time the last rotation was performed.
  * `material_expire_time` -(Available in 1.123.1+)  The time and date the key material for the CMK expires.
  * `next_rotation_date` -(Available in 1.123.1+)  The time the next rotation is scheduled for execution. 
  * `origin` -(Available in 1.123.1+)  The source of the key material for the CMK.
  * `protection_level` -(Available in 1.123.1+)  The protection level of the CMK.
  * `rotation_interval` -(Available in 1.123.1+)  The period of automatic key rotation.
  * `primary_key_version` -(Available in 1.123.1+)  The ID of the current primary key version of the symmetric CMK. 
  
