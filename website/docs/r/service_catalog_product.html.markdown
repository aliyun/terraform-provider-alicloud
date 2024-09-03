---
subcategory: "Service Catalog"
layout: "alicloud"
page_title: "Alicloud: alicloud_service_catalog_product"
description: |-
  Provides a Alicloud Service Catalog Product resource.
---

# alicloud_service_catalog_product

Provides a Service Catalog Product resource.

Service catalog product, IaC template encapsulation concept.

For information about Service Catalog Product and how to use it, see [What is Product](https://www.alibabacloud.com/help/en/).

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


resource "alicloud_service_catalog_product" "default" {
  provider_name = var.name
  description   = "desc"
  product_name  = var.name
  product_type  = "Ros"
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) The description of the product
* `product_name` - (Required) The name of the product
* `product_type` - (Required, ForceNew) The type of the product
* `provider_name` - (Required) The provider name of the product

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the product

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Product.
* `delete` - (Defaults to 5 mins) Used when delete the Product.
* `update` - (Defaults to 5 mins) Used when update the Product.

## Import

Service Catalog Product can be imported using the id, e.g.

```shell
$ terraform import alicloud_service_catalog_product.example <id>
```