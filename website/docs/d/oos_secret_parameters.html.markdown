---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_secret_parameters"
sidebar_current: "docs-alicloud-datasource-oos-secret-parameters"
description: |-
  Provides a list of Oos Secret Parameters to the user.
---

# alicloud_oos_secret_parameters

This data source provides the Oos Secret Parameters of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.147.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_oos_secret_parameter" "default" {
  secret_parameter_name = var.name
  value                 = "tf-testacc-oos_secret_parameter"
  type                  = "Secret"
  description           = var.name
  constraints           = <<EOF
  {
    "AllowedValues": [
        "tf-testacc-oos_secret_parameter"
    ],
    "AllowedPattern": "tf-testacc-oos_secret_parameter",
    "MinLength": 1,
    "MaxLength": 100
  }
  EOF
  tags = {
    Created = "TF"
    For     = "SecretParameter"
  }
}

data "alicloud_oos_secret_parameters" "ids" {
  ids = [alicloud_oos_secret_parameter.default.id]
}

output "oos_secret_parameter_id_0" {
  value = data.alicloud_oos_secret_parameters.ids.parameters.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of Secret Parameter IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Secret Parameter name.
* `secret_parameter_name` - (Optional, ForceNew) The name of the Secret Parameter.
* `resource_group_id` - (Optional, ForceNew) The ID of the Resource Group.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.
* `sort_field` - (Optional, ForceNew) The field used to sort the query results. Valid values: `Name`, `CreatedDate`.
* `sort_order` - (Optional, ForceNew) The order in which the entries are sorted. Default value: `Descending`. Valid values: `Ascending`, `Descending`.
* `enable_details` - (Optional, Bool) Whether to query the detailed list of resource attributes. Default value: `false`.
* `with_decryption` - (Optional, ForceNew, Bool, Available since v1.231.0) Specifies whether to decrypt the parameter value. Default value: `false`. **Note:** `with_decryption` takes effect only if `enable_details` is set to `true`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Secret Parameter names.
* `parameters` - A list of Oos Secret Parameters. Each element contains the following attributes:
  * `id` - The ID of the Secret Parameter.
  * `secret_parameter_id` - The ID of the encryption parameter.
  * `type` - The type of the parameter.
  * `parameter_version` - The version number of the encryption parameter.
  * `share_type` - The share type of the encryption parameter.
  * `key_id` - The ID of the key of Key Management Service (KMS) that is used for encryption.
  * `resource_group_id` - The ID of the Resource Group.
  * `secret_parameter_name` - The name of the encryption parameter.
  * `description` - The description of the encryption parameter.
  * `tags` - The tags of the parameter.
  * `constraints` - The constraints of the encryption parameter. **Note:** `constraints` takes effect only if `enable_details` is set to `true`.
  * `value` - (Available since v1.231.0) The value of the encryption parameter. **Note:** `value` takes effect only if `with_decryption` is set to `true`.
  * `created_by` - The user who created the encryption parameter.
  * `create_time` - The time when the encryption parameter was created.
  * `updated_by` - The user who updated the encryption parameter.
  * `updated_date` - The time when the encryption parameter was updated.
