---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_roles"
sidebar_current: "docs-alicloud-datasource-resource-manager-roles"
description: |-
    Provides a list of Resource Manager Roles to the user.
---

# alicloud\_resource\_manager\_roles

This data source provides the Resource Manager Roles of the current Alibaba Cloud user.

-> **NOTE:**  Available in 1.86.0+.

## Example Usage

```
data "alicloud_resource_manager_roles" "example" {
  name_regex = "tftest"
}

output "first_role_id" {
  value = "${data.alicloud_resource_manager_roles.example.roles.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Resource Manager Role IDs.
* `name_regex` - (Optional) A regex string to filter results by role name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of role IDs.
* `names` - A list of role names.
* `roles` - A list of roles. Each element contains the following attributes:
    * `id` - The ID of the role.
    * `role_id`- The ID of the role.
    * `role_name`- The name of the role.
    * `arn`- The Alibaba Cloud Resource Name (ARN) of the RAM role.
    * `create_date`- The time when the RAM role was created.
    * `update_date`- The time when the RAM role was updated.
    * `description`- The description of the RAM role.
    * `max_session_duration`- The maximum session duration of the RAM role.
    
