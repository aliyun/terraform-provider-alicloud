---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_access_key"
sidebar_current: "docs-alicloud-ephemeral-ram-access-key"
description: |-
    Creates a RAM AccessKey and exposes its secret without persisting it to the state or plan.
---

# Ephemeral: alicloud_ram_access_key

This ephemeral resource creates a RAM AccessKey through the `CreateAccessKey` API and exposes the resulting
AccessKey ID and AccessKey secret without persisting them to the state or plan. The credential is available
for the duration of a single Terraform operation and can be referenced in provider configuration, locals,
provisioners, and other ephemeral resources, but never appears in any Terraform artifact.

-> **NOTE:** Available since v2.1.0. Ephemeral resources require Terraform 1.10 or later.

-> **NOTE:** Every open of this ephemeral resource creates a new, permanent AccessKey. Terraform opens
ephemeral resources during both the plan and the apply phase, so a single `terraform apply` normally creates
two AccessKeys, and `terraform plan` alone also creates one. Each RAM user can hold at most two AccessKeys,
so an open can fail with a quota error when the user already has keys. Delete unused AccessKeys before
running Terraform.

-> **NOTE:** Terraform does not delete the created AccessKey when the ephemeral resource closes; the key stays
valid after the operation finishes. Because the RAM API only returns the secret at creation time and offers no
way to read it again, an AccessKey created by this resource cannot be recovered from Terraform. Record the
credential in its destination (for example Vault) and delete orphaned keys manually with
`DeleteAccessKey` or in the RAM console.

-> **NOTE:** The ephemeral resource is opened during the plan phase when `user_name` is known at plan time, so
the RAM user must already exist. When `user_name` references another resource in the same configuration, its
value is unknown during plan and the open is deferred to the apply phase automatically.

-> **NOTE:** When `user_name` is omitted, the AccessKey is created for the identity the provider itself
authenticates as. Prefer an explicit dedicated RAM user so the key can be identified and cleaned up.

## Example Usage

Create an AccessKey for a dedicated RAM user and hand it to a second provider instance in the same
configuration:

```terraform
variable "user_name" {
  default = "terraform-operator"
}

resource "alicloud_ram_user" "operator" {
  name  = var.user_name
  force = true
}

ephemeral "alicloud_ram_access_key" "operator" {
  user_name = alicloud_ram_user.operator.name
}

provider "alicloud" {
  alias      = "generated"
  access_key = ephemeral.alicloud_ram_access_key.operator.access_key_id
  secret_key = ephemeral.alicloud_ram_access_key.operator.access_key_secret
}

data "alicloud_caller_identity" "generated" {
  provider = alicloud.generated
}

output "generated_account_id" {
  value = data.alicloud_caller_identity.generated.account_id
}
```

## Argument Reference

The following arguments are supported:

* `user_name` - (Optional) The name of the RAM user the AccessKey belongs to. When omitted, the AccessKey is created for the identity the provider authenticates as. The name is not returned by the API when omitted.

## Attributes Reference

The following attributes are exported in the ephemeral result:

* `access_key_id` - The ID of the created AccessKey.
* `access_key_secret` - (Sensitive) The secret of the created AccessKey. It is returned only at creation time and cannot be read back through the RAM API.
* `status` - The status of the created AccessKey, for example `Active`.
* `create_time` - The time at which the AccessKey was created, in UTC, for example `2026-10-20T06:10:08Z`.
