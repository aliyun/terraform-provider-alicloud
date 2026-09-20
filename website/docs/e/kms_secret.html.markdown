---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: Ephemeral: alicloud_kms_secret"
sidebar_current: "docs-alicloud-ephemeral-kms-secret"
description: |-
    Retrieves the value of a KMS secret without persisting it to the state or plan.
---

# Ephemeral: alicloud_kms_secret

This ephemeral resource retrieves the value of a KMS secret through the `GetSecretValue` API without
persisting it to the state or plan. The value is available for the duration of a single Terraform
operation and can be referenced in provider configuration, write-only attributes, and other ephemeral
resources, but never appears in any Terraform artifact.

-> **NOTE:** Available since v2.0.0-beta5.

-> **NOTE:** Ephemeral resources require Terraform 1.10 or later.

## Example Usage

Basic Usage

```terraform
ephemeral "alicloud_kms_secret" "example" {
  secret_name = "tf-example-secret"
}
```

Reference the secret value in provider configuration:

```terraform
ephemeral "alicloud_kms_secret" "database" {
  secret_name = "prod/database/password"
}

provider "postgresql" {
  host     = "pg.example.com"
  username = "app"
  password = ephemeral.alicloud_kms_secret.database.secret_data
}
```

Access a secret owned by another Alibaba Cloud account by passing its ARN, built with the
provider-defined `arn_build` function:

```terraform
ephemeral "alicloud_kms_secret" "shared" {
  secret_name = provider::alicloud::arn_build("kms", "cn-hangzhou", "123456789012****", "secret/prod-database-password")
}
```

-> **NOTE:** Cross-account access is available only for secrets in a KMS instance and requires two-sided
authorization: the owning account must grant the caller's RAM user or role `kms:GetSecretValue` in the
[secret policy](https://www.alibabacloud.com/help/en/kms/key-management-service/security-and-compliance/modify-a-secret-policy),
and the caller's account must grant that RAM user or role the corresponding KMS permissions.

Feed the secret value into a write-only attribute (requires Terraform v1.11.0 or later), so the plaintext is never persisted anywhere:

```terraform
resource "alicloud_kms_ciphertext" "default" {
  key_id               = alicloud_kms_key.default.id
  plaintext_wo         = ephemeral.alicloud_kms_secret.example.secret_data
  plaintext_wo_version = 1
}
```

-> **NOTE:** The ephemeral resource is opened during the plan phase, so the secret it references must already exist. If the secret value is rotated, the rotated value reaches a write-only attribute only after its `_wo_version` is bumped in a subsequent apply.

## Argument Reference

The following arguments are supported:

* `secret_name` - (Required) The name of the secret, or the ARN of a secret in another Alibaba Cloud
  account (`acs:kms:${region}:${account}:secret/${secret-name}`). The ARN can be built with the
  provider-defined
  [`arn_build`](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/functions/arn_build)
  function, e.g.
  `provider::alicloud::arn_build("kms", "cn-hangzhou", "123456789012****", "secret/prod-database-password")`.
* `version_id` - (Optional) The version of the secret value to retrieve. If not specified, the
  version marked `ACSCurrent` is returned. The version actually retrieved is reported back in this
  attribute. Credential-type secrets (RDS, PolarDB, Redis/Tair, RAM, ECS) do not support version
  IDs; this parameter is ignored for them.
* `version_stage` - (Optional) The version stage of the secret value to retrieve. Defaults to
  `ACSCurrent`. Credential-type secrets (RDS, PolarDB, Redis/Tair, RAM, ECS) only support
  `ACSCurrent` and `ACSPrevious`.

## Attributes Reference

The following attributes are exported in the ephemeral result:

* `secret_data` - (Sensitive) The value of the secret. KMS decrypts the stored value before
  returning it. For `binary` values the string is base64-encoded.
* `version_id` - The version of the secret value actually retrieved.
* `version_stages` - The version stages attached to the retrieved version.
* `secret_data_type` - The data type of the secret value. Valid values: `text`, `binary`.
* `secret_type` - The type of the secret. Valid values: `Generic`, `Rds`, `Redis`, `RAMCredentials`,
  `ECS`, `PolarDB`.
