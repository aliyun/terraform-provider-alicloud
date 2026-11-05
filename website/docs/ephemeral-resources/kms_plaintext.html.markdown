---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_plaintext"
sidebar_current: "docs-alicloud-ephemeral-kms-plaintext"
description: |-
    Decrypts a KMS ciphertext without persisting the plaintext to the state or plan.
---

# Ephemeral: alicloud_kms_plaintext

This ephemeral resource decrypts a KMS ciphertext through the `Decrypt` API without persisting the
plaintext to the state or plan. The value is available for the duration of a single Terraform
operation and can be referenced in provider configuration, write-only attributes, and other ephemeral
resources, but never appears in any Terraform artifact.

-> **NOTE:** Available since v2.0.0-beta5. Ephemeral resources require Terraform 1.10 or later.

-> **NOTE:** The principal must have the `kms:Decrypt` permission for the key that encrypted the ciphertext.

-> **NOTE:** The ephemeral resource is opened during the plan phase when `ciphertext_blob` is known at plan
time, so the key that encrypted the ciphertext must still exist; if the key has been deleted, decryption
fails with a not-found error such as `Forbidden.KeyNotFound`. When `ciphertext_blob` references another
resource in the same configuration, its value is unknown during plan and the open is deferred to the
apply phase automatically.

-> **NOTE:** The
[`alicloud_kms_plaintext`](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/kms_plaintext)
data source persists the decrypted plaintext into the state; use this ephemeral resource when the value
must never appear in any Terraform artifact.

## Example Usage

Basic Usage

```terraform
ephemeral "alicloud_kms_plaintext" "example" {
  ciphertext_blob = "AQICAHh4WGJvT2Z...h3MjQ5Ng=="
}
```

Decrypt with an encryption context that matches the one supplied at encryption time:

```terraform
ephemeral "alicloud_kms_plaintext" "example" {
  ciphertext_blob = "AQICAHh4WGJvT2Z...h3MjQ5Ng=="

  encryption_context = {
    stage = "production"
  }
}
```

Reference the plaintext in provider configuration:

```terraform
ephemeral "alicloud_kms_plaintext" "database" {
  ciphertext_blob = "AQICAHh4WGJvT2Z...h3MjQ5Ng=="
}

provider "postgresql" {
  host     = alicloud_db_instance.example.connection_string
  username = "app"
  password = ephemeral.alicloud_kms_plaintext.database.plaintext
}
```

Feed the plaintext into a write-only attribute (requires Terraform v1.11.0 or later), so the plaintext
is never persisted anywhere:

```terraform
ephemeral "alicloud_kms_plaintext" "password" {
  ciphertext_blob = "AQICAHh4WGJvT2Z...h3MjQ5Ng=="
}

resource "alicloud_rds_account" "example" {
  db_instance_id              = alicloud_db_instance.example.id
  account_name                = "app"
  account_password_wo         = ephemeral.alicloud_kms_plaintext.password.plaintext
  account_password_wo_version = 1
}
```

## Argument Reference

The following arguments are supported:

* `ciphertext_blob` - (Required) The ciphertext to decrypt, produced by the KMS `Encrypt`,
  `GenerateDataKey`, or `GenerateDataKeyWithoutPlaintext` API.
* `encryption_context` - (Optional) The key/value pairs that were supplied as the encryption context
  when the ciphertext was produced. If the ciphertext was encrypted with an encryption context,
  decryption requires the exact same context.

## Attributes Reference

The following attributes are exported in the ephemeral result:

* `plaintext` - (Sensitive) The decrypted plaintext.
* `key_id` - The ID of the key that was used to decrypt the ciphertext.
