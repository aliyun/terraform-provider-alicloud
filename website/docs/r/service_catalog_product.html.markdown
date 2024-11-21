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

For information about Service Catalog Product and how to use it, see [What is Product](https://www.alibabacloud.com/help/en/service-catalog/developer-reference/api-servicecatalog-2021-09-01-createproduct).

-> **NOTE:** Available since v1.230.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_service_catalog_product&exampleId=972690f2-3292-1fb8-14fc-83c356de4eb5cd75e70b&activeTab=example&spm=docs.r.service_catalog_product.0.972690f232&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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