---
subcategory: "Market Place"
layout: "alicloud"
page_title: "Alicloud: alicloud_market_product"
sidebar_current: "docs-alicloud-datasource-market-product"
description: |-
    Provides details of a Market product item.
---

# alicloud\_market\_product

This data source provides the Market product item details of Alibaba Cloud.

-> **NOTE:** Available in 1.69.0+

## Example Usage

```
data "alicloud_market_product" "default" {
  product_code = "cmapi022206"
}

output "product_name" {
  value = "${data.alicloud_market_product.default.product.0.name}"
}

output "first_product_sku_code" {
  value = "${data.alicloud_market_product.default.product.0.skus.0.sku_code}"
}

output "first_product_package_version" {
  value = "${data.alicloud_market_product.default.product.0.skus.0.package_versions.0.package_version}"
}
```

## Argument Reference

The following arguments are supported:

* `product_code` - (Required) The product code of the market product.
* `available_region` - (Available in 1.71.1+) A available region id used to filter market place Ecs images.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `product` - A product. It contains the following attributes:
  * `code` - The code of the product.
  * `name` - The name of the product.
  * `description` - The description of the product.
  * `skus` - A list of one element containing sku attributes of an object. Each element contains the following attributes:
    * `sku_code` - The sku code of this product sku.
    * `sku_name` - The sku name of this product sku.
    * `package_versions` - The list of package version details of this product sku, Each element contains the following attributes:
      * `package_name` - The package name of this product sku package.
      * `package_version` - The package version of this product sku package. Currently, the API products can return package_version, but others can not for ensure.
    * `images` - The list of custom ECS images, Each element contains the following attributes:
      * `image_id` - The Ecs image id.
      * `image_name` - The Ecs image display name.
      * `region_id` - The Ecs image region.

 
