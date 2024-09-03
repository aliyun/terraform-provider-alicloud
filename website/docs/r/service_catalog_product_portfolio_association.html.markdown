---
subcategory: "Service Catalog"
layout: "alicloud"
page_title: "Alicloud: alicloud_service_catalog_product_portfolio_association"
description: |-
  Provides a Alicloud Service Catalog Product Portfolio Association resource.
---

# alicloud_service_catalog_product_portfolio_association

Provides a Service Catalog Product Portfolio Association resource.

Product portfolio association.

For information about Service Catalog Product Portfolio Association and how to use it, see [What is Product Portfolio Association](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.230.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_service_catalog_portfolio" "default0yAgJ8" {
  provider_name  = var.name
  description    = "desc"
  portfolio_name = var.name
}

resource "alicloud_service_catalog_product" "defaultRetBJw" {
  provider_name = var.name
  product_name  = format("%s1", var.name)
  product_type  = "Ros"
}


resource "alicloud_service_catalog_product_portfolio_association" "default" {
  portfolio_id = alicloud_service_catalog_portfolio.default0yAgJ8.id
  product_id   = alicloud_service_catalog_product.defaultRetBJw.id
}
```

## Argument Reference

The following arguments are supported:
* `portfolio_id` - (Required, ForceNew) Product Portfolio ID
* `product_id` - (Required, ForceNew) Product ID

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<product_id>:<portfolio_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Product Portfolio Association.
* `delete` - (Defaults to 5 mins) Used when delete the Product Portfolio Association.

## Import

Service Catalog Product Portfolio Association can be imported using the id, e.g.

```shell
$ terraform import alicloud_service_catalog_product_portfolio_association.example <product_id>:<portfolio_id>
```