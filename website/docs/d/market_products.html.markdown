---
subcategory: "Market Place"
layout: "alicloud"
page_title: "Alicloud: alicloud_market_products"
sidebar_current: "docs-alicloud-datasource-market-products"
description: |-
    Provides a list of Market product items.
---

# alicloud\_market\_products

This data source provides the Market product items of Alibaba Cloud.

-> **NOTE:** Available in 1.64.0+

## Example Usage

```
data "alicloud_market_products" "default" {
  sort         = "created_on-desc"
  category_id  = "53690006"
  product_type = "SERVICE"
}

output "first_product_code" {
  value = "${data.alicloud_market_products.default.product_items.0.code}"
}

output "product_codes" {
  value = "${data.alicloud_market_products.default.ids}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, Available 1.66.0+) A regex string to apply to the product name.
* `ids` - (Optional, Available 1.66.0+) A list of product code.
* `sort` - (Optional, ForceNew) This field determines how to sort the filtered results, Valid values: `user_count-desc`, `created_on-desc`, `price-desc` and `score-desc`.
* `category_id` - (Optional, ForceNew) The Category ID of products. For more information, see [DescribeProducts](https://help.aliyun.com/document_detail/89834.htm). 
* `product_type` - (Optional, ForceNew) The type of products, Valid values: `APP`, `SERVICE`, `MIRROR`, `DOWNLOAD` and `API_SERVICE`.
* `search_term` - (Optional, ForceNew, Available 1.69.0+) Search term in this query.
* `supplier_id` - (Optional, ForceNew, Available 1.71.1+) The supplier id of the product.
* `supplier_name_keyword` - (Optional, ForceNew, Available 1.71.1+) The supplier name keyword of the product.
* `suggested_price` - (Optional, ForceNew, Available 1.71.1+) The suggested price of the product.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of product codes.
* `products` - A list of products. Each element contains the following attributes:
  * `code` - The code of the product.
  * `name` - The name of the product.
  * `category_id` - The category id of the product.
  * `supplier_id` - The supplier id of the product.
  * `supplier_name` - The supplier name of the product.
  * `short_description` - The short description of the product.
  * `tags` - The tags of the product.
  * `suggested_price` - The suggested price of the product.
  * `target_url` - The detail page URL of the product.
  * `image_url` - The image URL of the product.
  * `score` - The rating information of the product.
  * `operation_system` - The operation system of the product.
  * `warranty_date` - The warranty date of the product.
  * `delivery_date` - The delivery date of the product.
  * `delivery_way` - The delivery way of the product.
 
