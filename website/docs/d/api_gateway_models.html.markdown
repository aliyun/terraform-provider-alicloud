---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_models"
sidebar_current: "docs-alicloud-datasource-api-gateway-models"
description: |-
  Provides a list of Api Gateway Models to the user.
---

# alicloud\_api\_gateway\_models

This data source provides the Api Gateway Models of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.187.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_api_gateway_models" "ids" {
  ids      = ["example_id"]
  group_id = "example_group_id"
}

output "api_gateway_model_id_1" {
  value = data.alicloud_api_gateway_models.ids.models.0.id
}

data "alicloud_api_gateway_models" "group_id" {
  group_id = "example_group_id"
}

output "api_gateway_model_id_2" {
  value = data.alicloud_api_gateway_models.group_id.models.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Model IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Model name.
* `group_id` - (Optional, ForceNew) The ID of the api group.
* `model_name` - (Optional, ForceNew) The name of the Model.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Model names.
* `models` - A list of Api Gateway Models. Each element contains the following attributes:
	* `id` - The ID of the Api Gateway Model.
	* `group_id` - The group of the model belongs to.
	* `model_name` - The name of the Model.
	* `schema` - The schema of the model.
	* `description` - The description of the model.
	* `model_id` - The id of the model.
	* `model_ref` - The reference of the model.
	* `modified_time` - The modified time of the model.
	* `create_time` - The creation time of the model.