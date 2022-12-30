---
subcategory: "Bss Open Api"
layout: "alicloud"
page_title: "Alicloud: alicloud_bss_openapi_products"
sidebar_current: "docs-alicloud-datasource-bss-openapi-products"
description: |-
  Provides a list of Bss Open Api Product owned by an Alibaba Cloud account.
---

# alicloud_bssopenapi_products

This data source provides Bss Open Api Product available to the user.[What is Product](https://www.alibabacloud.com/help/zh/bss-openapi/latest/api-doc-bssopenapi-2017-12-14-api-doc-queryproductlist)

-> **NOTE:** Available in 1.195.0+

## Example Usage

```
data "alicloud_bss_openapi_products" "default" {
  name_regex = "内容分发网络CDN"
}

output "alicloud_bssopenapi_product_example_id" {
  value = data.alicloud_bss_openapi_products.default.products.0.product_code
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of product IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Product name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `names` - A list of name of Products.
* `products` - A list of Product Entries. Each element contains the following attributes:
    * `id` - The ID of the product. The value is formulated as `<product_code>:<product_type>:<subscription_type>`.
    * `product_code` - Product code.
    * `product_name` - Product name.
    * `product_type` - Type of product.
    * `subscription_type` - Subscription type. Value:
      * Subscription: Prepaid.
      * PayAsYouGo: postpaid.
