---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_parameters"
sidebar_current: "docs-alicloud-datasource-oos-parameters"
description: |-
  Provides a list of Oos Parameters to the user.
---

# alicloud_oos_parameters

This data source provides the Oos Parameters of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.147.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_oos_parameter" "default" {
  parameter_name = var.name
  value          = "tf-testacc-oos_parameter"
  type           = "String"
  description    = var.name
  constraints    = <<EOF
  {
    "AllowedValues": [
        "tf-testacc-oos_parameter"
    ],
    "AllowedPattern": "tf-testacc-oos_parameter",
    "MinLength": 1,
    "MaxLength": 100
  }
  EOF
  tags = {
    Created = "TF"
    For     = "Parameter"
  }
}

data "alicloud_oos_parameters" "ids" {
  ids = [alicloud_oos_parameter.default.id]
}

output "oos_secret_parameter_id_0" {
  value = data.alicloud_oos_parameters.ids.parameters.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of Parameter IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Parameter name.
* `parameter_name` - (Optional, ForceNew) The name of the common parameter. You can enter a keyword to query parameter names in fuzzy match mode.
* `type` - (Optional, ForceNew) The data type of the common parameter. Valid values: `String`, `StringList`.
* `resource_group_id` - (Optional, ForceNew) The ID of the Resource Group.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.
* `sort_field` - (Optional, ForceNew) The field used to sort the query results. Valid values: `Name`, `CreatedDate`.
* `sort_order` - (Optional, ForceNew) The order in which the entries are sorted. Default value: `Descending`. Valid values: `Ascending`, `Descending`.
* `enable_details` - (Optional, Bool) Whether to query the detailed list of resource attributes. Default value: `false`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Parameter names.
* `parameters` - A list of Oos Parameters. Each element contains the following attributes:
  * `id` - The ID of the Parameter. Its value is same as `parameter_name`.
  * `parameter_name` - The name of the common parameter.
  * `parameter_id` - The ID of the common parameter.
  * `type` - The data type of the common parameter.
  * `parameter_version` - The version number of the common parameter.
  * `share_type` - The share type of the common parameter.
  * `resource_group_id` - The ID of the Resource Group.
  * `description` - The description of the common parameter.
  * `tags` - The tags added to the common parameter.
  * `constraints` - The constraints of the common parameter. **Note:** `constraints` takes effect only if `enable_details` is set to `true`.
  * `value` - (Available since v1.231.0) The value of the common parameter. **Note:** `value` takes effect only if `enable_details` is set to `true`.
  * `created_by` - The user who created the common parameter.
  * `create_time` - The time when the common parameter was created.
  * `updated_by` - The user who updated the common parameter.
  * `updated_date` - The time when the common parameter was updated.
