---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_parameter_groups"
sidebar_current: "docs-alicloud-datasource-rds-parameter-groups"
description: |-
  Provides a list of Rds Parameter Groups to the user.
---

# alicloud\_rds\_parameter\_groups

This data source provides the Rds Parameter Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.119.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_rds_parameter_groups" "example" {
  ids        = ["example_value"]
  name_regex = "the_resource_name"
}

output "first_rds_parameter_group_id" {
  value = data.alicloud_rds_parameter_groups.example.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Parameter Group IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Parameter Group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Parameter Group names.
* `groups` - A list of Rds Parameter Groups. Each element contains the following attributes:
	* `engine` - The database engine.
	* `engine_version` - The version of the database engine.
	* `force_restart` - Indicates whether applying the parameter template requires the instance to be restarted. Valid values: `0`: A restart is not required. `1`: A restart is required.
	* `id` - The ID of the Parameter Group.
	* `param_detail` - Parameter list.
		* `param_value` - The name of a parameter.
		* `param_name` - The value of a parameter.
	* `parameter_group_desc` - The description of the parameter template.
	* `parameter_group_id` - The ParameterGroupId of the Parameter Group.
	* `parameter_group_name` - The name of the parameter template.
	* `parameter_group_type` - The type of the parameter template. Valid values:
        `0`: indicates a default parameter template.
        `1`: indicates a custom parameter template.
        `2`: an automatic backup parameter template. After you apply this type of template to an RDS instance, the system automatically backs up the parameter settings of that instance and saves the data backup as a template.
