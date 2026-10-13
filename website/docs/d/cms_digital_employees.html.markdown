---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_digital_employees"
description: |-
  This data source provides the list of Cms Digital Employees.
---

# alicloud_cms_digital_employees

This data source provides a list of Cms Digital Employees in an Alibaba Cloud account.

-> **NOTE:** Available since v1.282.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_digital_employees" "default" {
  digital_employee_name = "tf-digital-employee"
  ids                   = ["tf-digital-employee"]
}

output "first_id" {
  value = data.alicloud_cms_digital_employees.default.digital_employees.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of digital employee names. Only digital employees whose name matches one of these IDs are returned.
* `name_regex` - (Optional) A regex string to filter results by the digital employee name.
* `digital_employee_name` - (Optional) Filter results by the exact digital employee name.
* `display_name` - (Optional) Filter results by the exact display name.
* `employee_type` - (Optional) Filter results by the exact employee type.
* `resource_group_id` - (Optional) Filter results by the resource group ID.
* `output_file` - (Optional) The name of the file to save the result set.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of digital employee names.
* `digital_employees` - A list of Cms Digital Employees.

### `digital_employees`

The `digital_employees` supports the following attributes:
* `create_time` - The creation time of the digital employee.
* `default_rule` - The default rule of the digital employee.
* `description` - The description of the digital employee.
* `digital_employee_name` - The name of the digital employee.
* `display_name` - The display name of the digital employee.
* `employee_type` - The type of the digital employee.
* `id` - The ID of the digital employee. It is the value of `digital_employee_name`.
* `region_id` - The region ID of the resource.
* `resource_group_id` - The ID of the resource group to which the digital employee belongs.
* `role_arn` - The ARN of the RAM role that the digital employee assumes.
* `tags` - A set of tags for the digital employee.
* `update_time` - The update time of the digital employee.

### `digital_employees.tags`

The `tags` supports the following attributes:
* `key` - The key of the tag.
* `value` - The value of the tag.
