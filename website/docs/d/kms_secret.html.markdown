---
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_secret"
sidebar_current: "docs-alicloud-datasource-kms-secret"
description: |-
    Decrypt a secret encrypted with KMS.
---

# alicloud\_kms\_secret

Decrypt a given ciphertext with KMS to use the resulting plaintext in resources.

~> **NOTE**: Using this data provider will allow you to conceal secret data within your resource definitions but does not take care of protecting that data in all Terraform logging and state output. Please take care to secure your secret data beyond just the Terraform configuration.

## Example Usage

```
resource "alicloud_kms_key" "key" {
  description             = "example key"
  is_enabled              = true
}

# Encrypt plaintext 'example'
resource "alicloud_kms_ciphertext" "encrypted" {
  key_id    = alicloud_kms_key.key.id
  plaintext = "example"
}

# Decrypt encrypted ciphertext
data "alicloud_kms_secret" "secret" {
  ciphertext_blob = alicloud_kms_ciphertext.encrypted.ciphertext_blob
}

# Output 'example' should match the plaintext encrypted in the beginning
output "decrypted" {
  value = data.alicloud_kms_secret.secret.plaintext
}
```

## Argument Reference

The following arguments are supported:

* `encryption_context` -
  (Optional) The Encryption context. If you specify this parameter in the Encrypt or GenerateDataKey API operation, it is also required when you call the Decrypt API operation. For more information, see [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm?spm=a2c63.p38356.b99.14.47562193BvC7Hu).
* `ciphertext_blob` - The ciphertext to be decrypted.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `plaintext` -  The decrypted plaintext.
* `key_id` - The globally unique ID of the CMK. It is the ID of the CMK used to decrypt ciphertext.
