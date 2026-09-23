---
subcategory: "Security Center"
layout: "alicloud"
page_title: "Alicloud: alicloud_security_center_groups"
sidebar_current: "docs-alicloud-datasource-security-center-groups"
description: |-
  Provides a list of Security Center Groups to the user.
---

# alicloud\_security\_center\_groups

This data source provides the Security Center Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.133.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_security_center_groups" "nameRegex" {
  name_regex = "^my-Group"
}
output "security_center_groups" {
  value = data.alicloud_security_center_groups.nameRegex.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Group IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Group names.
* `groups` - A list of Security Center Groups. Each element contains the following attributes:
	* `id` - The ID of the Group(same as the group_id).
	* `group_id` - The ID of Group.
	* `group_name` - The name of Group.
	* `group_flag` - GroupFlag, '0' mean default group(created by system), '1' means customer defined group.
