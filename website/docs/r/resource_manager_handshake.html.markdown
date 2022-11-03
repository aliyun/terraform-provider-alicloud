---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_handshake"
sidebar_current: "docs-alicloud-resource-resource-manager-handshake"
description: |-
  Provides a Resource Manager handshake resource.
---

# alicloud\_resource\_manager\_handshake

Provides a Resource Manager handshake resource. You can invite accounts to join a resource directory for unified management.
For information about Resource Manager handshake and how to use it, see [What is Resource Manager handshake](https://www.alibabacloud.com/help/en/doc-detail/135287.htm).

-> **NOTE:** Available in v1.82.0+.

## Example Usage

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

* `target_entity` - (Required, ForceNew) Invited account ID or login email.
* `target_type` - (Required, ForceNew) Type of account being invited. Valid values: `Account`, `Email`.
* `note` - (Optional, ForceNew) Remarks. The maximum length is 1024 characters.

## Attributes Reference

The following attributes are exported:

* `id` - This ID of Resource Manager handshake.  
* `expire_time` - The expiration time of the invitation.
* `master_account_id` - Resource account master account ID.
* `master_account_name` - The name of the main account of the resource directory.
* `modify_time` - The modification time of the invitation.
* `resource_directory_id` - Resource directory ID.
* `status` - Invitation status. Valid values: `Pending` waiting for confirmation, `Accepted`, `Cancelled`, `Declined`, `Expired`. 

## Import

Resource Manager handshake can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_handshake.example h-QmdexeFm1kE*****
```
