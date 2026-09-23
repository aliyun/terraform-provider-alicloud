---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_secret_parameter"
sidebar_current: "docs-alicloud-ephemeral-oos-secret-parameter"
description: |-
    Retrieves the decrypted value of an OOS secret parameter without persisting it to the state or plan.
---

# Ephemeral: alicloud_oos_secret_parameter

This ephemeral resource retrieves the decrypted value of an OOS secret parameter through the
`GetSecretParameter` API without persisting it to the state or plan. The value is available for the
duration of a single Terraform operation and can be referenced in provider configuration, write-only
attributes, and other ephemeral resources, but never appears in any Terraform artifact.

-> **NOTE:** Available since v2.0.0-beta5. Ephemeral resources require Terraform 1.10 or later.

-> **NOTE:** The principal must have both the `oos:GetSecretParameter` permission and the KMS
`GetSecretValue` permission for the KMS key that encrypts the parameter.

-> **NOTE:** The ephemeral resource is opened during the plan phase, so the parameter it references
must already exist. If the parameter is created by another resource in the same configuration, the
first plan fails with a not-found error such as `EntityNotExists.Parameter`; create the parameter out of band or in
an earlier apply first. If the parameter value is rotated, the rotated value reaches a write-only
attribute only after its `_wo_version` is bumped in a subsequent apply.

-> **NOTE:** Object metadata such as the encryption key id, tags, and description is not part of the ephemeral
result; use the
[`alicloud_oos_secret_parameters`](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/oos_secret_parameters)
data source to read it.

## Example Usage

Basic Usage

```terraform
ephemeral "alicloud_oos_secret_parameter" "example" {
  secret_parameter_name = "tf-example-secret-parameter"
}
```

Reference the parameter value in provider configuration:

```terraform
ephemeral "alicloud_oos_secret_parameter" "database" {
  secret_parameter_name = "prod/database/password"
}

provider "postgresql" {
  host     = alicloud_db_instance.example.connection_string
  username = "app"
  password = ephemeral.alicloud_oos_secret_parameter.database.value
}
```

Feed the parameter value into a write-only attribute (requires Terraform v1.11.0 or later), so the
plaintext is never persisted anywhere:

```terraform
ephemeral "alicloud_oos_secret_parameter" "password" {
  secret_parameter_name = "prod/rds/password"
}

resource "alicloud_rds_account" "example" {
  db_instance_id              = alicloud_db_instance.example.id
  account_name                = "app"
  account_password_wo         = ephemeral.alicloud_oos_secret_parameter.password.value
  account_password_wo_version = 1
}
```

## Argument Reference

The following arguments are supported:

* `secret_parameter_name` - (Required) The name of the secret parameter. The name can contain
  letters, digits, hyphens (-) and underscores (_), must be 1 to 180 characters in length, and
  cannot start with `ALIYUN`, `ACS`, `ALIBABA`, `ALICLOUD`, or `OOS`.
* `parameter_version` - (Optional) The version of the parameter value to retrieve. Every update of
  the parameter value creates a new version. If not specified, the latest version is returned. The
  version actually retrieved is reported back in this attribute.

## Attributes Reference

The following attributes are exported in the ephemeral result:

* `value` - (Sensitive) The decrypted value of the secret parameter. The provider always requests
  the decrypted value: without decryption the `GetSecretParameter` API returns no value at all.
* `parameter_version` - The version of the parameter value actually retrieved.
