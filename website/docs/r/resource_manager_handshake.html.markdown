---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_handshake"
description: |-
  Provides a Alicloud Resource Manager Handshake resource.
---

# alicloud_resource_manager_handshake

Provides a Resource Manager Handshake resource.



For information about Resource Manager Handshake and how to use it, see [What is Handshake](https://www.alibabacloud.com/help/en/doc-detail/135287.htm).

-> **NOTE:** Available since v1.259.0.

## Example Usage

Basic Usage

```terraform
# Add a Resource Manager handshake.
resource "alicloud_resource_manager_handshake" "example" {
  target_entity = "1182775234******"
  target_type   = "Account"
  note          = "test resource manager handshake"
}
```

## Argument Reference

The following arguments are supported:
* `note` - (Optional, ForceNew) The description of the invitation.
The description can be up to 1,024 characters in length.
* `target_entity` - (Required, ForceNew) The ID or logon email address of the account that you want to invite.
* `target_type` - (Required, ForceNew) The type of the invited account. Valid values:

  - Account: indicates the ID of the account.
  - Email: indicates the logon email address of the account.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the invitation was created. The time is displayed in UTC.
* `expire_time` - The time when the invitation expires. The time is displayed in UTC.
* `master_account_id` - The ID of the management account of the resource directory.
* `master_account_name` - The name of the management account of the resource directory.
* `modify_time` - The time when the invitation was modified. The time is displayed in UTC.
* `resource_directory_id` - The ID of the resource directory.
* `status` - The status of the invitation. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Handshake.
* `delete` - (Defaults to 5 mins) Used when delete the Handshake.

## Import

Resource Manager Handshake can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_handshake.example <id>
```