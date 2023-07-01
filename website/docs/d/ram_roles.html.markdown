---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_roles"
sidebar_current: "docs-alicloud-datasource-ram-roles"
description: |-
  Provides a list of ram roles available to the user.
---

# alicloud_ram_roles

This data source provides a list of RAM Roles in an Alibaba Cloud account according to the specified filters.

## Example Usage

```terraform
data "alicloud_ram_roles" "roles_ds" {
  output_file = "roles.txt"
  name_regex  = ".*test.*"
  policy_name = "AliyunACSDefaultAccess"
  policy_type = "Custom"
}

output "first_role_id" {
  value = "${data.alicloud_ram_roles.roles_ds.roles.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, ForceNew) A regex string to filter results by the role name.
* `ids` (Optional, Available since v1.53.0) - A list of ram role IDs. 
* `policy_type` - (Optional, ForceNew) Filter results by a specific policy type. Valid values are `Custom` and `System`. If you set this parameter, you must set `policy_name` as well.
* `policy_name` - (Optional, ForceNew) Filter results by a specific policy name. If you set this parameter without setting `policy_type`, the later will be automatically set to `System`. The resulting roles will be attached to the specified policy.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of ram role IDs. 
* `names` - A list of ram role names. 
* `roles` - A list of roles. Each element contains the following attributes:
  * `id` - ID of the role.
  * `name` - Name of the role.
  * `arn` - Resource descriptor of the role.
  * `description` - Description of the role.
  * `assume_role_policy_document` - Authorization strategy of the role. This parameter is deprecated and replaced by `document`.
  * `document` - Authorization strategy of the role.
  * `create_date` - Creation date of the role.
  * `update_date` - Update date of the role.
