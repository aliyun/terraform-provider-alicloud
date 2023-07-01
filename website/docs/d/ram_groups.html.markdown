---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_groups"
sidebar_current: "docs-alicloud-datasource-ram-groups"
description: |-
  Provides a list of ram groups available to the user.
---

# alicloud_ram_groups

This data source provides a list of RAM Groups in an Alibaba Cloud account according to the specified filters.

## Example Usage

```terraform
data "alicloud_ram_groups" "groups_ds" {
  output_file = "groups.txt"
  user_name   = "user1"
  name_regex  = "^group[0-9]*"
}

output "first_group_name" {
  value = "${data.alicloud_ram_groups.groups_ds.groups.0.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, ForceNew) A regex string to filter the returned groups by their names.
* `user_name` - (Optional, ForceNew) Filter the results by a specific the user name.
* `policy_type` - (Optional, ForceNew) Filter the results by a specific policy type. Valid items are `Custom` and `System`. If you set this parameter, you must set `policy_name` as well.
* `policy_name` - (Optional, ForceNew) Filter the results by a specific policy name. If you set this parameter without setting `policy_type`, it will be automatically set to `System`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of ram group names.
* `groups` - A list of groups. Each element contains the following attributes:
  * `name` - Name of the group.
  * `comments` - Comments of the group.
