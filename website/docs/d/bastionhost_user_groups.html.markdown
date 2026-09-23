---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_user_groups"
sidebar_current: "docs-alicloud-datasource-bastionhost-user-groups"
description: |-
  Provides a list of Bastionhost User Groups to the user.
---

# alicloud\_bastionhost\_user\_groups

This data source provides the Bastionhost User Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.132.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_bastionhost_user_groups" "ids" {
  instance_id = "bastionhost-cn-xxxx"
  ids         = ["1", "2"]
}
output "bastionhost_user_group_id_1" {
  value = data.alicloud_bastionhost_user_groups.ids.groups.0.id
}

data "alicloud_bastionhost_user_groups" "nameRegex" {
  instance_id = "bastionhost-cn-xxxx"
  name_regex  = "^my-UserGroup"
}
output "bastionhost_user_group_id_2" {
  value = data.alicloud_bastionhost_user_groups.nameRegex.groups.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of User Group self IDs.
* `instance_id` - (Required, ForceNew) Specify the New Group of the Bastion Host of Instance Id.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by User Group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `user_group_name` - (Optional, ForceNew) Specify the New Group Name. Supports up to 128 Characters.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of User Group names.
* `groups` - A list of Bastionhost User Groups. Each element contains the following attributes:
	* `comment` - Specify the New Group of Remark Information. Supports up to 500 Characters.
	* `id` - The ID of the User Group.
	* `instance_id` - Specify the New Group of the Bastion Host of Instance Id.
	* `user_group_id` - The User Group ID.
	* `user_group_name` - Specify the New Group Name. Supports up to 128 Characters.
