---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_ciphertext"
sidebar_current: "docs-alicloud-resource-kms-ciphertext"
description: |-
  Encrypt data with KMS.
---

# alicloud\_kms\_ciphertext

Encrypt a given plaintext with KMS. The produced ciphertext stays stable across applies. If the plaintext should be re-encrypted on each apply use the [`alicloud_kms_ciphertext`](/docs/providers/alicloud/d/kms_ciphertext.html) data source.

~> **NOTE**: Using this data provider will allow you to conceal secret data within your resource definitions but does not take care of protecting that data in all Terraform logging and state output. Please take care to secure your secret data beyond just the Terraform configuration.

## Example Usage

```terraform
resource "alicloud_kms_key" "key" {
  description             = "example key"
  is_enabled              = true
}

resource "alicloud_kms_ciphertext" "encrypted" {
  key_id    = alicloud_kms_key.key.id
  plaintext = "example"
}
```

## Argument Reference

The following arguments are supported:

* `plaintext` - (ForceNew) The plaintext to be encrypted which must be encoded in Base64.
* `key_id` - (ForceNew) The globally unique ID of the CMK.
* `encryption_context` -
  (Optional, ForceNew) The Encryption context. If you specify this parameter here, it is also required when you call the Decrypt API operation. For more information, see [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ciphertext_blob` - The ciphertext of the data key encrypted with the primary CMK version.
