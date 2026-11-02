---
subcategory: "Service Catalog"
layout: "alicloud"
page_title: "Alicloud: alicloud_service_catalog_provisioned_products"
sidebar_current: "docs-alicloud-datasource-service-catalog-provisioned-products"
description: |-
  Provides a list of Service Catalog Provisioned Product owned by an Alibaba Cloud account.
---

# alicloud_service_catalog_provisioned_products

This data source provides Service Catalog Provisioned Product available to the user. [What is Provisioned Product](https://www.alibabacloud.com/help/en/service-catalog/developer-reference/api-servicecatalog-2021-09-01-launchproduct)

-> **NOTE:** Available since v1.196.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_service_catalog_provisioned_products" "default" {
  ids        = ["IdExample"]
  name_regex = "NameRegexExample"
}

output "alicloud_service_catalog_provisioned_product_example_id" {
  value = data.alicloud_service_catalog_provisioned_products.default.provisioned_products.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of Provisioned Product IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Product name.
* `access_level_filter` - (Optional, ForceNew) The access filter.
* `sort_by` - (Optional, ForceNew) The field that is used to sort the queried data.
* `sort_order` - (Optional, ForceNew) The sorting method.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Provisioned Product IDs.
* `names` - A list of name of Provisioned Products.
* `provisioned_products` - (Available since v1.197.0) A list of Provisioned Product Entries. Each element contains the following attributes:
  * `id` - The ID of the provisioned product.
  * `create_time` - The creation time of the product instance
  * `last_provisioning_task_id` - The ID of the last instance operation task
  * `last_successful_provisioning_task_id` - The ID of the last successful instance operation task
  * `last_task_id` - The ID of the last task
  * `outputs` - The output value of the template.
    * `description` - Description of the output value defined in the template.
    * `output_key` - The name of the output value defined in the template.
    * `output_value` - The content of the output value defined in the template.
  * `owner_principal_id` - The RAM entity ID of the owner
  * `owner_principal_type` - The RAM entity type of the owner
  * `parameters` - The template parameters entered by the user.
    * `parameter_key` - The name of the parameter defined in the template.
    * `parameter_value` - The Template parameter value entered by the user.
  * `portfolio_id` - Product mix ID.
  * `product_id` - Product ID.
  * `product_name` - The name of the product
  * `product_version_id` - Product version ID.
  * `product_version_name` - The name of the product version
  * `provisioned_product_arn` - The ARN of the product instance
  * `provisioned_product_id` - The ID of the instance.
  * `provisioned_product_name` - The name of the instance.
  * `provisioned_product_type` - Instance type.
  * `stack_id` - The ID of the ROS stack
  * `stack_region_id` - The ID of the region to which the resource stack belongs.
  * `status` - Instance status
  * `status_message` - The status message of the product instance
  * `tags` - A mapping of tags of the provisioned product.
* `products` - (Deprecated since v1.197.0) A list of Provisioned Product Entries. Each element contains the following attributes:
  * `id` - The ID of the provisioned product.
  * `create_time` - The creation time of the product instance
  * `last_provisioning_task_id` - The ID of the last instance operation task
  * `last_successful_provisioning_task_id` - The ID of the last successful instance operation task
  * `last_task_id` - The ID of the last task
  * `outputs` - The output value of the template.
    * `description` - Description of the output value defined in the template.
    * `output_key` - The name of the output value defined in the template.
    * `output_value` - The content of the output value defined in the template.
  * `owner_principal_id` - The RAM entity ID of the owner
  * `owner_principal_type` - The RAM entity type of the owner
  * `parameters` - The template parameters entered by the user.
    * `parameter_key` - The name of the parameter defined in the template.
    * `parameter_value` - The Template parameter value entered by the user.
  * `portfolio_id` - Product mix ID.
  * `product_id` - Product ID.
  * `product_name` - The name of the product
  * `product_version_id` - Product version ID.
  * `product_version_name` - The name of the product version
  * `provisioned_product_arn` - The ARN of the product instance
  * `provisioned_product_id` - The ID of the instance.
  * `provisioned_product_name` - The name of the instance.
  * `provisioned_product_type` - Instance type.
  * `stack_id` - The ID of the ROS stack
  * `stack_region_id` - The ID of the region to which the resource stack belongs.
  * `status` - Instance status
  * `status_message` - The status message of the product instance
  * `tags` - A mapping of tags of the provisioned product.
