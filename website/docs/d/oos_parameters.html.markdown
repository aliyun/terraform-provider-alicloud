---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_parameters"
sidebar_current: "docs-alicloud-datasource-oos-parameters"
description: |-
  Provides a list of Oos Parameters to the user.
---

# alicloud\_oos\_parameters

This data source provides the Oos Parameters of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.147.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_oos_parameters" "ids" {
  ids = ["my-Parameter"]
}
output "oos_parameter_id_1" {
  value = data.alicloud_oos_parameters.ids.parameters.0.id
}

data "alicloud_oos_parameters" "nameRegex" {
  name_regex = "^my-Parameter"
}
output "oos_parameter_id_2" {
  value = data.alicloud_oos_parameters.nameRegex.parameters.0.id
}

data "alicloud_oos_parameters" "resourceGroupId" {
  ids               = ["my-Parameter"]
  resource_group_id = "example_value"
}
output "oos_parameter_id_3" {
  value = data.alicloud_oos_parameters.resourceGroupId.parameters.0.id
}

data "alicloud_oos_parameters" "tags" {
  ids = ["my-Parameter"]
  tags = {
    Created = "TF"
    For     = "OosParameter"
  }
}
output "oos_parameter_id_4" {
  value = data.alicloud_oos_parameters.tags.parameters.0.id
}

```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Parameter IDs. Its element value is same as Parameter Name.
* `parameter_name` - (Optional, ForceNew) The name of the common parameter. You can enter a keyword to query parameter names in fuzzy match mode.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Parameter name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew) The ID of the Resource Group.
* `type` - (Optional, ForceNew) The data type of the common parameter. Valid values: `String` and `StringList`.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Parameter names.
* `parameters` - A list of Oos Parameters. Each element contains the following attributes:
	* `constraints` - The constraints of the common parameter.
	* `create_time` - The time when the common parameter was created.
	* `created_by` - The user who created the common parameter.
	* `description` - The description of the common parameter.
	* `id` - The ID of the Parameter. Its value is same as `parameter_name`.
	* `parameter_id` - The ID of the common parameter.
	* `parameter_name` - The name of the common parameter.
	* `parameter_version` - The version number of the common parameter.
	* `resource_group_id` - The ID of the Resource Group.
	* `share_type` - The share type of the common parameter.
	* `tags` - The tag of the resource.
	* `type` - The data type of the common parameter.
	* `updated_by` - The user who updated the common parameter.
	* `updated_date` - The time when the common parameter was updated.
	* `value` - The value of the common parameter.