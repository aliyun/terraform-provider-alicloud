---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_handshake_acceptance"
description: |-
  Provides a Alicloud Resource Manager Handshake Acceptance resource.
---

# alicloud_resource_manager_handshake_acceptance

Provides a Resource Manager Handshake Acceptance resource. It is used by the **invited account** to accept an invitation (handshake) to join a resource directory.

For information about Resource Manager Handshake Acceptance and how to use it, see [What is Handshake](https://www.alibabacloud.com/help/en/doc-detail/135287.htm) and the [AcceptHandshake](https://next.api.alibabacloud.com/document/ResourceDirectoryMaster/2022-04-19/AcceptHandshake) API.

-> **NOTE:** This resource must be applied with the credentials of the **invited account**, not the management (master) account that sent the invitation. The invitation itself is created by the management account via `alicloud_resource_manager_handshake`. Use a separate [provider alias](https://developer.hashicorp.com/terraform/language/providers/configuration#alias-multiple-provider-configurations) configured with the invited account's credentials.

-> **NOTE:** Accepting a handshake is irreversible. Destroying this resource only removes it from the Terraform state; it does not make the invited account leave the resource directory.

-> **NOTE:** Available since v1.282.0.

## Example Usage

Basic Usage

```terraform
# The management account sends the invitation.
resource "alicloud_resource_manager_handshake" "example" {
  target_entity = "1182775234******"
  target_type   = "Account"
  note          = "test resource manager handshake"
}

# The invited account accepts it (configure invited_account with the invited account's credentials).
resource "alicloud_resource_manager_handshake_acceptance" "example" {
  provider     = alicloud.invited_account
  handshake_id = alicloud_resource_manager_handshake.example.id
}
```

## Argument Reference

The following arguments are supported:
* `handshake_id` - (Required, ForceNew) The ID of the invitation.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource. The value is the same as `handshake_id`.
* `create_time` - The time when the invitation was created. The time is displayed in UTC.
* `expire_time` - The time when the invitation expires. The time is displayed in UTC.
* `invited_account_real_name` - The real-name verification information of the invited account.
* `master_account_id` - The ID of the management account of the resource directory.
* `master_account_name` - The name of the management account of the resource directory.
* `master_account_real_name` - The real-name verification information of the management account.
* `modify_time` - The time when the invitation was modified. The time is displayed in UTC.
* `note` - The note of the invitation.
* `resource_directory_id` - The ID of the resource directory.
* `status` - The status of the invitation. The acceptance resource exists only when the status is `Accepted`.
* `target_entity` - The invited account. The value is an account ID when `target_type` is `Account`, or a logon email address when `target_type` is `Email`.
* `target_type` - The type of the invited account. Valid values: `Account` and `Email`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when accept the Handshake.
* `delete` - (Defaults to 5 mins) Used when delete the Handshake Acceptance.

## Import

Resource Manager Handshake Acceptance can be imported using the id (handshake id), e.g.

```shell
$ terraform import alicloud_resource_manager_handshake_acceptance.example <id>
```
