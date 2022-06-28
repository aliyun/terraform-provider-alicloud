---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_secret_parameters"
sidebar_current: "docs-alicloud-datasource-oos-secret-parameters"
description: |-
  Provides a list of Oos Secret Parameters to the user.
---

# alicloud\_oos\_secret\_parameters

This data source provides the Oos Secret Parameters of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.147.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_oos_secret_parameters" "nameRegex" {
  ids = ["example_value"]
}
output "oos_secret_parameter_id_1" {
  value = data.alicloud_oos_secret_parameters.nameRegex.parameters.0.id
}

data "alicloud_oos_secret_parameters" "nameRegex" {
  name_regex = "^my-Parameter"
}
output "oos_parameter_id_2" {
  value = data.alicloud_oos_secret_parameters.nameRegex.parameters.0.id
}

data "alicloud_oos_secret_parameters" "resourceGroupId" {
  ids               = ["example_value"]
  resource_group_id = "example_value"
}
output "oos_parameter_id_3" {
  value = data.alicloud_oos_secret_parameters.resourceGroupId.parameters.0.id
}

data "alicloud_oos_secret_parameters" "tags" {
  ids = ["example_value"]
  tags = {
    Created = "TF"
    For     = "OosSecretParameter"
  }
}
output "oos_parameter_id_4" {
  value = data.alicloud_oos_secret_parameters.tags.parameters.0.id
}

```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Secret Parameter name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew) The ID of the Resource Group.
* `secret_parameter_name` - (Optional, ForceNew) The name of the secret parameter.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Secret Parameter IDs.
* `names` - A list of Secret Parameter names.
* `parameters` - A list of Oos Secret Parameters. Each element contains the following attributes:
	* `constraints` - The constraints of the encryption parameter.
	* `create_time` - The time when the encryption parameter was created.
	* `created_by` - The user who created the encryption parameter.
	* `description` - The description of the encryption parameter.
	* `id` - The ID of the Secret Parameter.
	* `key_id` - KeyId of KMS used for encryption.
	* `parameter_version` - The version number of the encryption parameter.
	* `resource_group_id` - The ID of the Resource Group.
	* `secret_parameter_id` - The ID of the encryption parameter.
	* `secret_parameter_name` - The name of the encryption parameter.
	* `share_type` - The share type of the encryption parameter.
	* `tags` - The tag of the resource.
	* `type` - The data type of the encryption parameter.
	* `updated_by` - The user who updated the encryption parameter.
	* `updated_date` - The time when the encryption parameter was updated.