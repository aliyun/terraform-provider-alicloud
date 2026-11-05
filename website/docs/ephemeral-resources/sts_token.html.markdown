---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_sts_token"
sidebar_current: "docs-alicloud-ephemeral-sts-token"
description: |-
    Fetches a temporary access credential by assuming a RAM role without persisting it to the state or plan.
---

# Ephemeral: alicloud_sts_token

This ephemeral resource fetches a temporary access credential (an STS token) by assuming a RAM role through the `AssumeRole` API without persisting the credential to the state or plan. The credential is available for the duration of a single Terraform operation and can be referenced in provider configuration, locals, provisioners, and other ephemeral resources, but never appears in any Terraform artifact.

-> **NOTE:** Available since v2.1.0. Ephemeral resources require Terraform 1.10 or later.

-> **NOTE:** The session is valid for 3600 seconds by default. The validity period can be set through `duration_seconds` and ranges from 900 seconds up to the role's `MaxSessionDuration` (configurable per role from 3600 to 43200 seconds, 3600 by default). When the provider itself is configured with STS credentials, the session cannot outlive the remaining validity of the parent token. The credential is re-fetched on every operation that opens it and expires at the reported `expiration`.

-> **NOTE:** The ephemeral resource is opened during the plan phase when `role_arn` is known at plan time, so the role must already exist. When `role_arn` references another resource in the same configuration, its value is unknown during plan and the open is deferred to the apply phase automatically.

-> **NOTE:** If the only goal is for the provider itself to operate under a RAM role, configure the provider's own [`assume_role`](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs#assume_role) block instead — no ephemeral resource is needed. Use this ephemeral resource when the credential is consumed by something else in the configuration: a second provider instance, a provisioner, or a write-only attribute.

## Example Usage

Assume a role and expose the credential to another consumer in the same configuration:

```terraform
ephemeral "alicloud_sts_token" "example" {
  role_arn     = "acs:ram::123456789012345678:role/operator"
  session_name = "deploy"
}
```

Assume a role with a shorter validity period and a downscoping policy:

```terraform
ephemeral "alicloud_sts_token" "example" {
  role_arn         = "acs:ram::123456789012345678:role/operator"
  session_name     = "deploy"
  duration_seconds = 900
  policy = jsonencode({
    Version = "1"
    Statement = [{
      Effect   = "Allow"
      Action   = ["oss:ListBuckets"]
      Resource = ["*"]
    }]
  })
}
```

Build the role ARN from the current account ID with the provider-defined [`arn_build`](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/functions/arn_build) function:

```terraform
data "alicloud_caller_identity" "current" {}

ephemeral "alicloud_sts_token" "example" {
  role_arn     = provider::alicloud::arn_build("ram", "", data.alicloud_caller_identity.current.account_id, "role/operator")
  session_name = "deploy"
}
```

## Argument Reference

The following arguments are supported:

* `role_arn` - (Required) The ARN of the RAM role to assume, for example `acs:ram::123456789012345678:role/operator`. The role's trust policy must allow the current caller to assume it. The ARN can be built with the provider-defined [`arn_build`](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/functions/arn_build) function, e.g. `provider::alicloud::arn_build("ram", "", "123456789012****", "role/operator")`.
* `session_name` - (Required) The name of the role session, used to identify the session in audit logs. The assumed role user ARN takes the form `<role_arn>/<session_name>`.
* `duration_seconds` - (Optional) The validity period of the credential, in seconds. Valid values are `900` to `43200`. The effective maximum is the role's `MaxSessionDuration`. When this field is omitted, the session is valid for 3600 seconds.
* `policy` - (Optional) An inline policy, as a JSON string, that further restricts the permissions of the resulting session to a subset of the role's permissions.
* `external_id` - (Optional) The external ID to include in the request, used when the role's trust policy requires a specific external ID.

## Attributes Reference

The following attributes are exported in the ephemeral result:

* `access_key_id` - (Sensitive) The AccessKeyId of the temporary credential.
* `access_key_secret` - (Sensitive) The AccessKeySecret of the temporary credential.
* `security_token` - (Sensitive) The SecurityToken of the temporary credential.
* `expiration` - The time at which the credential expires, in UTC RFC3339 format, for example `2026-10-12T04:21:34Z`.
* `assumed_role_arn` - The ARN of the assumed role user, which takes the form `<role_arn>/<session_name>`.
