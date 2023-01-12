---
subcategory: "Service Catalog"
layout: "alicloud"
page_title: "Alicloud: alicloud_service_catalog_product_as_end_users"
sidebar_current: "docs-alicloud-datasource-service-catalog-product-as-end-users"
description: |-
  Provides a list of Service Catalog Product As End User owned by an Alibaba Cloud account.
---

# alicloud_service_catalog_product_as_end_users

-> **DEPRECATED:** It has been deprecated from version `1.197.0`.
Please use new datasource [alicloud_service_catalog_end_user_products](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/service_catalog_end_user_products) instead.

This data source provides Service Catalog Product As End User available to the user.[What is Product As End User](https://www.alibabacloud.com/help/en/servicecatalog/latest/api-doc-servicecatalog-2021-09-01-api-doc-listproductsasenduser)

-> **NOTE:** Available in 1.196.0+

## Example Usage

```terraform
data "alicloud_service_catalog_product_as_end_users" "default" {
  name_regex = "ram模板创建"
}

output "alicloud_service_catalog_product_as_end_user_example_id" {
  value = data.alicloud_service_catalog_product_as_end_users.default.users.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of Product As End User IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by product name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Product As End User IDs.
* `users` - A list of Product As End User Entries. Each element contains the following attributes:
    * `id` - ID of product, Its value is the same as `product_id`.
    * `create_time` - Product creation time.According to ISO8601 standard, UTC time is used in the format: YYYY-MM-DDThh:mm:ssZ.
    * `description` - Product description.
    * `has_default_launch_option` - Whether there is a default Startup option. Value:-true: There is a default Startup option, and there is no need to fill in the portfolio when starting the product or updating the instance.-false: there is no default Startup option. You must fill in the portfolio when starting the product or updating the instance. > If the product is added to only one product portfolio, there will be a default Startup option. If the product is added to multiple product combinations, there will be multiple startup options at the same time, but there is no default Startup option at this time.
    * `product_arn` - Product ARN.
    * `product_id` - Product ID.
    * `product_name` - Product name.
    * `product_type` - Type of product.The value is Ros, which indicates the resource orchestration service (ROS).
    * `provider_name` - Product provider.
