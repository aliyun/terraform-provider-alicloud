---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_access_key"
sidebar_current: "docs-alicloud-ephemeral-ram-access-key"
description: |-
    Creates a RAM access key for the duration of an ephemeral reference and revokes it on close, without persisting the secret to state.
---

# Ephemeral: alicloud_ram_access_key

This ephemeral resource creates a RAM access key (an `AccessKeyId` / `AccessKeySecret` pair) for a RAM user through the `CreateAccessKey` API and exposes it for the duration of a single Terraform operation. The secret never appears in the Terraform state or plan. When the ephemeral reference closes, the access key is deleted through the `DeleteAccessKey` API so that no residual credential is left behind.

This is intended for short-lived credentials that another part of the same Terraform operation consumes directly — for example an aliased provider configured with the key — without ever letting the access key secret touch the state file. In contrast, the long-lived [`alicloud_ram_access_key` managed resource](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/ram_access_key) writes `secret` to state.

-> **NOTE:** Available since v2.1.0. Ephemeral resources require Terraform 1.10 or later.

-> **NOTE:** The access key is created during the operation that opens the ephemeral resource and is deleted when that operation ends. A single RAM user may hold at most two access keys; if the limit is reached, `CreateAccessKey` fails until an existing key is revoked.

-> **NOTE:** When `user_name` is omitted, the access key is created for the caller itself (the RAM user whose credentials the provider is configured with). Set `user_name` to issue a key for a different RAM user; the caller must be permitted to manage that user's access keys.

-> **NOTE:** `CreateAccessKey` has no idempotency token, so the provider disables automatic retries for it. If a create request fails after it reached RAM, an access key may exist without having been tracked for revocation; list the user's access keys with `ListAccessKeys` and delete any unexpected key before retrying.

-> **NOTE:** Storing the access key in a secrets manager such as [HashiCorp Vault](https://registry.terraform.io/providers/hashicorp/vault/latest/docs) for later use is incompatible with this resource: Terraform closes the ephemeral resource during the same operation, Close deletes the key, and the stored copy stops working. Use the [`alicloud_ram_access_key` managed resource](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/ram_access_key) when the credential must outlive the Terraform operation; the ephemeral resource only serves consumers that use the credential before Close revokes it.

-> **NOTE:** Closing the reference is best-effort idempotent: if the access key has already been deleted out-of-band, Close succeeds without error. The key is only created once, on open; there is no renewal.

## Example Usage

Issue a short-lived access key for a RAM user and consume it within the same operation through an aliased provider, without the secret ever reaching the Terraform state:

```terraform
ephemeral "alicloud_ram_access_key" "deploy" {
  user_name = "ci-deployer"
}

provider "alicloud" {
  alias      = "deploy"
  region     = "cn-hangzhou"
  access_key = ephemeral.alicloud_ram_access_key.deploy.access_key_id
  secret_key = ephemeral.alicloud_ram_access_key.deploy.access_key_secret
}

data "alicloud_account" "deploy" {
  provider = alicloud.deploy
}
```

Issue an access key for the caller itself (the configured RAM user) by omitting `user_name`:

```terraform
ephemeral "alicloud_ram_access_key" "self" {}
```

## Argument Reference

The following arguments are supported:

* `user_name` - (Optional, String) The name of the RAM user for whom the access key is created. When omitted, the access key is created for the caller itself. The caller must be permitted to manage the target user's access keys.

## Attributes Reference

The following attributes are exported in the ephemeral result:

* `access_key_id` - (Sensitive, String) The AccessKeyId of the created access key.
* `access_key_secret` - (Sensitive, String) The AccessKeySecret of the created access key. Only available while the ephemeral resource is open; it is never persisted to state.
* `status` - (String) The status of the access key, `Active` or `Inactive`.
* `create_date` - (String) The creation time of the access key, in UTC, for example `2026-10-20T03:21:59Z`.
