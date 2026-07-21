---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_host_groups"
sidebar_current: "docs-alicloud-datasource-bastionhost-host-groups"
description: |-
  Provides a list of Bastionhost Host Groups to the user.
---

# alicloud\_bastionhost\_host\_groups

This data source provides the Bastionhost Host Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_bastionhost_host_groups" "ids" {
  instance_id = "bastionhost-cn-tl3xxxxxxx"
  ids         = ["example_value-1", "example_value-2"]
}
output "bastionhost_host_group_id_1" {
  value = data.alicloud_bastionhost_host_groups.ids.groups.0.id
}

data "alicloud_bastionhost_host_groups" "nameRegex" {
  instance_id = "bastionhost-cn-tl3xxxxxxx"
  name_regex  = "^my-HostGroup"
}
output "bastionhost_host_group_id_2" {
  value = data.alicloud_bastionhost_host_groups.nameRegex.groups.0.id
}

```

## Argument Reference

The following arguments are supported:

* `host_group_name` - (Optional, ForceNew) Specify the New Host Group Name, Supports up to 128 Characters.
* `ids` - (Optional, ForceNew, Computed)  A list of Host Group IDs.
* `instance_id` - (Required, ForceNew) Specify the New Host Group Where the Bastion Host ID of.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Host Group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Host Group names.
* `groups` - A list of Bastionhost Host Groups. Each element contains the following attributes:
	* `comment` - Specify the New Host Group of Notes, Supports up to 500 Characters.
	* `host_group_id` - Host Group ID.
	* `host_group_name` - Specify the New Host Group Name, Supports up to 128 Characters.
	* `id` - The ID of the Host Group.
	* `instance_id` - Specify the New Host Group Where the Bastion Host ID of.
