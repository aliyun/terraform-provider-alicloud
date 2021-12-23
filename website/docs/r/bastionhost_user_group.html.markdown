---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_user_group"
sidebar_current: "docs-alicloud-resource-bastionhost-user-group"
description: |-
  Provides a Alicloud Bastion Host User Group resource.
---

# alicloud\_bastionhost\_user\_group

Provides a Bastion Host User Group resource.

For information about Bastion Host User Group and how to use it, see [What is User Group](https://www.alibabacloud.com/help/doc-detail/204596.htm).

-> **NOTE:** Available in v1.132.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_bastionhost_user_group" "example" {
  instance_id     = "example_value"
  user_group_name = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `comment` - (Optional) Specify the New Group of Remark Information. Supports up to 500 Characters.
* `instance_id` - (Required, ForceNew) Specify the New Group of the Bastion Host of Instance Id.
* `user_group_name` - (Required) Specify the New Group Name. Supports up to 128 Characters.

## Attributes Reference

The following attributes are exported:

* `user_group_id` - The User Group self ID.
* `id` - The resource ID of User Group. The value formats as `<instance_id>:<user_group_id>`.

## Import

Bastion Host User Group can be imported using the id, e.g.

```
$ terraform import alicloud_bastionhost_user_group.example <instance_id>:<user_group_id>
```
