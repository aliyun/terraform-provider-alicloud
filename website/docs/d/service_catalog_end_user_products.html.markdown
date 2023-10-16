---
subcategory: "Service Catalog"
layout: "alicloud"
page_title: "Alicloud: alicloud_service_catalog_end_user_products"
sidebar_current: "docs-alicloud-datasource-service_catalog-end-user-products"
description: |-
  Provides a list of Service Catalog End User Product owned by an Alibaba Cloud account.
---

# alicloud_service_catalog_end_user_products

This data source provides Service Catalog End User Product available to the user.[What is End User Product](https://www.alibabacloud.com/help/en/servicecatalog/latest/api-servicecatalog-2021-09-01-listproductsasenduser)

-> **NOTE:** Available since v1.197.0.

## Example Usage

```terraform
data "alicloud_service_catalog_end_user_products" "default" {
  name_regex = "ram模板创建"
}

output "alicloud_service_catalog_end_user_product_example_id" {
  value = data.alicloud_service_catalog_end_user_products.default.end_user_products.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of End User Product IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by product name.
* `sort_by` - (Optional, ForceNew) The field that is used to sort the queried data. The value is fixed as CreateTime, which specifies the creation time of products.
* `sort_order` - (Optional, ForceNew) The order in which you want to sort the queried data. Valid values: `Asc`, `Desc`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of End User Product IDs.
* `end_user_products` - A list of End User Product Entries. Each element contains the following attributes:
  * `id` - ID of product, Its value is the same as `product_id`.
  * `create_time` - Product creation time.According to ISO8601 standard, UTC time is used in the format: YYYY-MM-DDThh:mm:ssZ.
  * `description` - Product description.
  * `has_default_launch_option` - Whether there is a default Startup option. Value:-true: There is a default Startup option, and there is no need to fill in the portfolio when starting the product or updating the instance.-false: there is no default Startup option. You must fill in the portfolio when starting the product or updating the instance. For more information about how to obtain the portfolio, see [ListLaunchOptions](~~ ListLaunchOptions ~~).> If the product is added to only one product portfolio, there will be a default Startup option. If the product is added to multiple product combinations, there will be multiple startup options at the same time, but there is no default Startup option at this time.
  * `product_arn` - Product ARN.
  * `product_id` - Product ID.
  * `product_name` - Product name.
  * `product_type` - Type of product.The value is Ros, which indicates the resource orchestration service (ROS).
  * `provider_name` - Product provider.
