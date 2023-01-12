---
subcategory: "Service Catalog"
layout: "alicloud"
page_title: "Alicloud: alicloud_service_catalog_product_versions"
sidebar_current: "docs-alicloud-datasource-service_catalog-product-versions"
description: |-
  Provides a list of Service Catalog Product Version owned by an Alibaba Cloud account.
---

# alicloud_service_catalog_product_versions

This data source provides Service Catalog Product Version available to the user.[What is Product Version](https://www.alibabacloud.com/help/en/servicecatalog/latest/api-doc-servicecatalog-2021-09-01-api-doc-listproductversions)

-> **NOTE:** Available in 1.196.0+

## Example Usage

```
data "alicloud_service_catalog_product_versions" "default" {
  name_regex = "1.0.0"
  product_id = "prod-bp125x4k29wb7q"
}

output "alicloud_service_catalog_product_version_example_id" {
  value = data.alicloud_service_catalog_product_versions.default.product_versions.0.id
}
```

## Argument Reference

The following arguments are supported:
* `product_id` - (Required,ForceNew) Product ID
* `ids` - (Optional, ForceNew, Computed) A list of Product Version IDs.
* `product_version_names` - (Optional, ForceNew) The name of the Product Version. You can specify at most 10 names.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Product Version IDs.
* `names` - A list of name of Product Versions.
* `product_versions` - A list of Product Version Entries. Each element contains the following attributes:
    * `id` - ID of product version.
    * `active` - Whether the version is activated
    * `create_time` - The creation time of the resource
    * `description` - Version description
    * `guidance` - Administrator guidance
    * `product_version_id` - The first ID of the resource
    * `product_version_name` - The name of the resource
    * `template_type` - Template Type
    * `template_url` - Template URL
