---
subcategory: "Service Catalog"
layout: "alicloud"
page_title: "Alicloud: alicloud_service_catalog_launch_options"
sidebar_current: "docs-alicloud-datasource-service-catalog-launch-options"
description: |-
  Provides a list of Service Catalog Launch Option owned by an Alibaba Cloud account.
---

# alicloud_service_catalog_launch_options

This data source provides Service Catalog Launch Option available to the user.[What is Launch Option](https://www.alibabacloud.com/help/en/servicecatalog/latest/api-doc-servicecatalog-2021-09-01-api-doc-listlaunchoptions)

-> **NOTE:** Available in 1.196.0+

## Example Usage

```terraform
data "alicloud_service_catalog_end_user_products" "default" {
  name_regex = "ram模板创建"
}
data "alicloud_service_catalog_launch_options" "default" {
  product_id = "data.alicloud_service_catalog_end_user_products.default.end_user_products.0.id"
}

output "alicloud_service_catalog_launch_option_example_id" {
  value = data.alicloud_service_catalog_launch_options.default.launch_options.0.id
}
```

## Argument Reference

The following arguments are supported:
* `product_id` - (Required,ForceNew) Product ID.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `name_regex` - (Required,ForceNew) A regex string to filter results by portfolio name.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `launch_options` - A list of Launch Option Entries. Each element contains the following attributes:
    * `id` - ID of Service Catalog Launch Option.
    * `constraint_summaries` - List of constraints.
        * `constraint_type` - Constraint type.The value is Launch, which indicates that the constraint is started.
        * `description` - Constraint description.
    * `portfolio_id` - Product mix ID.
    * `portfolio_name` - Product portfolio name.
