---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_user_attachment"
sidebar_current: "docs-alicloud-resource-bastionhost-user-attachment"
description: |-
  Provides a Alicloud Bastion Host User Attachment resource.
---

# alicloud\_bastionhost\_user\_attachment

Provides a Bastion Host User Attachment resource to add user to one user group.

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_bastionhost_user_attachment" "example" {
  instance_id   = "bastionhost-cn-tl3xxxxxxx"
  user_group_id = "10"
  user_id       = "100"
}

```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) Specifies the user group to add the user's bastion host ID of.
* `user_group_id` - (Required, ForceNew) Specifies the user group to which you want to add the user ID.
* `user_id` - (Required, ForceNew) Specify that you want to add to the policy attached to the user group ID. This includes response parameters in a Json-formatted string supports up to set up 100 USER ID.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of User Attachment. The value formats as `<instance_id>:<user_group_id>:<user_id>`.

## Import

Bastion Host User Attachment can be imported using the id, e.g.

```
$ terraform import alicloud_bastionhost_user_attachment.example <instance_id>:<user_group_id>:<user_id>
```
