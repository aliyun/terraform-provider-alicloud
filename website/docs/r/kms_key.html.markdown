---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_key"
description: |-
  Provides a Alicloud KMS Key resource.
---

# alicloud_kms_key

Provides a KMS Key resource.



For information about KMS Key and how to use it, see [What is Key](https://www.alibabacloud.com/help/en/kms/developer-reference/api-createkey).

-> **NOTE:** Available since v1.85.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_kms_key&exampleId=88de40bc-6464-f764-767b-2c8f68796a3c60b3415b&activeTab=example&spm=docs.r.kms_key.0.88de40bc64&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_kms_key" "default" {
  description            = "Hello KMS"
  status                 = "Enabled"
  pending_window_in_days = "7"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_kms_key&spm=docs.r.kms_key.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `automatic_rotation` - (Optional) Specifies whether to enable automatic key rotation. Default value: `Disabled`. Valid values: `Enabled`, `Disabled`.
* `deletion_protection` - (Optional, Available since v1.240.0) Specifies whether to enable deletion protection. Default value: `Disabled`. Valid values: `Enabled`, `Disabled`.
* `deletion_protection_description` - (Optional, Available since v1.240.0) The description of deletion protection. **NOTE:** `deletion_protection_description` takes effect only if `deletion_protection` is set to `Enabled`.
* `description` - (Optional) The description of the key.
* `dkms_instance_id` - (Optional, ForceNew, Available since v1.183.0) The ID of the KMS instance. If specified, the key is created in the specified KMS instance; if omitted, a default key (master key) is created in the current region.

  **NOTE:** The two types of keys differ in capability: only keys in KMS instances support key policies (`policy`), and the valid values of `key_spec` and `rotation_interval` are determined by the key management type. For more information, see [CreateKey](https://www.alibabacloud.com/help/en/kms/developer-reference/api-createkey).
* `key_spec` - (Optional, ForceNew) The specification of the key. Default value: `Aliyun_AES_256`. Valid values: `Aliyun_AES_256`, `Aliyun_AES_128`, `Aliyun_AES_192`, `Aliyun_SM4`, `RSA_2048`, `RSA_3072`, `EC_P256`, `EC_P256K`, `EC_SM2`.

  **NOTE:** The valid values vary based on the key management type. For keys in KMS instances, the supported specifications are determined by the instance type.
* `key_usage` - (Optional, ForceNew) The usage of the key. Default value: `ENCRYPT/DECRYPT`. Valid values:
  - `ENCRYPT/DECRYPT`: Encrypts or decrypts data.
  - `SIGN/VERIFY`: Generates or verifies a digital signature.
* `origin` - (Optional, ForceNew) The key material origin. Default value: `Aliyun_KMS`. Valid values: `Aliyun_KMS`, `EXTERNAL`.
* `pending_window_in_days` - (Optional, Int) The number of days before the CMK is deleted. During this period, the CMK is in the PendingDeletion state. After this period ends, you cannot cancel the deletion. Unit: days. Valid values: `7` to `366`. **NOTE:** From version 1.184.0, `pending_window_in_days` can be set to `366`.
* `policy` - (Optional, Available since v1.224.0) The content of the key policy. The value is in the JSON format. The value can be up to 32,768 bytes in length. For more information, see [How to use it](https://www.alibabacloud.com/help/en/kms/developer-reference/api-setkeypolicy).

  **NOTE:** This parameter takes effect only on keys in KMS instances (`dkms_instance_id` specified). Default keys do not support key policies, and this parameter does not take effect on them.
* `protection_level` - (Optional, ForceNew) The protection level of the key. Default value: `SOFTWARE`. Valid values: `SOFTWARE`, `HSM`.

  **NOTE:** The default value does not take effect when `dkms_instance_id` is specified. In this case, this parameter is ignored and the instance type determines the protection level: keys in a software key management instance are `SOFTWARE`, and keys in a hardware key management instance are `HSM`. For keys whose actual protection level is `HSM` (such as hardware instance keys), explicitly set `protection_level` to `HSM`. Otherwise, the configuration (default `SOFTWARE`) never matches the actual level and Terraform forces key recreation on every apply.
* `rotation_interval` - (Optional) The period of automatic key rotation. The following units are supported: d (day), h (hour), m (minute), and s (second). For example, you can use either 7d or 604800s to specify a seven-day interval.

  **NOTE:** This parameter is required if `automatic_rotation` is set to `Enabled`, and takes effect only when the key management type supports automatic rotation: default keys support only a fixed period of 365 days, software-protected keys support 7 to 365 days, and hardware-protected keys do not support automatic rotation.
* `status` - (Optional, Available since v1.123.1) The status of key. Default value: `Enabled`. Valid values: `Enabled`, `Disabled`, `PendingDeletion`.
* `tags` - (Optional, Map, Available since v1.207.0) A mapping of tags to assign to the resource.
* `deletion_window_in_days` - (Optional, Int, Deprecated since v1.85.0) Field `deletion_window_in_days` has been deprecated from provider version 1.85.0. New field `pending_window_in_days` instead.
* `key_state` - (Deprecated since v1.123.1) Field `key_state` has been deprecated from provider version 1.123.1. New field `status` instead.
* `is_enabled` - (Optional, Bool, Deprecated since v1.85.0) Field `is_enabled` has been deprecated from provider version 1.85.0. New field `status` instead.

-> **NOTE:** If you set the origin parameter to EXTERNAL or the key_spec parameter to an asymmetric CMK type, automatic key rotation is unavailable.

-> **NOTE:** The default type of the CMK is `Aliyun_AES_256`. Only keys in KMS instances support `Aliyun_AES_128` and `Aliyun_AES_192`.

-> **NOTE:** When the pre-deletion days elapses, the key is permanently deleted and cannot be recovered.

## Attributes Reference

* `id` - The resource ID in terraform of Key.
* `arn` - The ARN of the key.
* `primary_key_version` - The ID of the current primary key version of the symmetric CMK.
* `last_rotation_date` - The time when the last rotation was performed.
* `next_rotation_date` - The time when the next rotation will be performed.
* `material_expire_time` - The time when the key material expires.
* `creator` - The creator of the CMK.
* `creation_date` - The time when the CMK was created.
* `delete_date` - The time at which the CMK is scheduled for deletion.

## Timeouts

-> **NOTE:** Available since v1.224.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Key.
* `update` - (Defaults to 5 mins) Used when update the Key.
* `delete` - (Defaults to 5 mins) Used when delete the Key.

## Import

KMS Key can be imported using the id, e.g.

```shell
$ terraform import alicloud_kms_key.example <id>
```
