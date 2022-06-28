---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_key"
sidebar_current: "docs-alicloud-resource-kms-key"
description: |-
  Provides a Alikms key resource.
---

# alicloud\_kms\_key

A kms key can help user to protect data security in the transmission process. For information about Alikms Key and how to use it, see [What is Resource Alikms Key](https://www.alibabacloud.com/help/doc-detail/28947.htm).

-> **NOTE:** Available in v1.85.0+.

## Example Usage

Basic Usage

```
resource "alicloud_kms_key" "key" {
  description             = "Hello KMS"
  pending_window_in_days  = "7"
  status                  = "Enabled"
}
```
## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the CMK. The description can be 0 to 8,192 characters in length.
* `key_usage` - (Optional, ForceNew) The usage of the CMK. Valid values:
  - ENCRYPT/DECRYPT(default value): encrypts or decrypts data. 
  - SIGN/VERIFY: generates or verifies a digital signature.
* `deletion_window_in_days` - (Optional) Field `deletion_window_in_days` has been deprecated from provider version 1.85.0. New field `pending_window_in_days` instead.
* `is_enabled` - (Optional) Field `is_enabled` has been deprecated from provider version 1.85.0. New field `key_state` instead.
* `automatic_rotation` - (Optional) Specifies whether to enable automatic key rotation. Valid values: 
  - Enabled
  - Disabled (default value)
  **NOTE**: If you set the origin parameter to EXTERNAL or the key_spec parameter to an asymmetric CMK type, automatic key rotation is unavailable.
    
* `key_spec`   - (Optional, ForceNew) The type of the CMK. Valid values: 
  "Aliyun_AES_256", "Aliyun_AES_128", "Aliyun_AES_192", "Aliyun_SM4", "RSA_2048", "RSA_3072", "EC_P256", "EC_P256K", "EC_SM2".
  Note: The default type of the CMK is Aliyun_AES_256. Only Dedicated KMS supports Aliyun_AES_128 and Aliyun_AES_192.
* `key_state` - (Optional) Field `key_state` has been deprecated from provider version 1.123.1. New field `status` instead.
* `status` - (Optional, Available in 1.123.1+) The status of CMK. Valid Values: 
  - Disabled
  - Enabled (default value)
  - PendingDeletion
  
* `origin` - (Optional, ForceNew) The source of key material. Valid values: 
  - Aliyun_KMS (default value)
  - EXTERNAL
  **NOTE**: The value of this parameter is case-sensitive. If you set the `key_spec` to an asymmetric CMK type, 
    you are not allowed to set the `origin` to EXTERNAL. If you set the `origin` to EXTERNAL, you must import key material. 
    For more information, see [import key material](https://www.alibabacloud.com/help/en/doc-detail/68523.htm).
    
* `pending_window_in_days` - (Optional) The number of days before the CMK is deleted. 
  During this period, the CMK is in the PendingDeletion state. 
  After this period ends, you cannot cancel the deletion. Valid values: 7 to 30. Unit: days.
* `protection_level` - (Optional, ForceNew) The protection level of the CMK. Valid values:
  - SOFTWARE (default value)
  - HSM
  **NOTE**: The value of this parameter is case-sensitive. Assume that you set this parameter to HSM. 
    If you set the origin parameter to Aliyun_KMS, the CMK is created in a managed hardware security module (HSM). 
    If you set the origin parameter to EXTERNA, you can import an external key to the managed HSM.
    
* `rotation_interval` - (Optional) The interval for automatic key rotation. Specify the value in the integer[unit] format.
  The following units are supported: d (day), h (hour), m (minute), and s (second). 
  For example, you can use either 7d or 604800s to specify a seven-day interval. 
  The interval can range from 7 days to 730 days. 
  **NOTE**: It is Required when `automatic_rotation = "Enabled"`
                                           
-> **NOTE:** When the pre-deletion days elapses, the key is permanently deleted and cannot be recovered.


## Attributes Reference

* `id` - The ID of the key.
* `arn` - The Alicloud Resource Name (ARN) of the key.
* `creation_date` -The date and time when the CMK was created. The time is displayed in UTC.
* `creator` -The creator of the CMK.
* `delete_date` -The scheduled date to delete CMK. The time is displayed in UTC. This value is returned only when the KeyState value is PendingDeletion.
* `last_rotation_date` - The date and time the last rotation was performed. The time is displayed in UTC. 
* `material_expire_time` - The time and date the key material for the CMK expires. The time is displayed in UTC. If the value is empty, the key material for the CMK does not expire.
* `next_rotation_date` - The time the next rotation is scheduled for execution. 
* `primary_key_version` - The ID of the current primary key version of the symmetric CMK. 


## Import

Alikms key can be imported using the id, e.g.

```
$ terraform import alicloud_kms_key.example abc123456
```
