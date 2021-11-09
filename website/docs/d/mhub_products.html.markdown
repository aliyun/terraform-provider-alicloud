---
subcategory: "Enterprise Mobile Application Studio (MHUB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mhub_products"
sidebar_current: "docs-alicloud-datasource-mhub-products"
description: |-
  Provides a list of Mhub Products to the user.
---

# alicloud\_mhub\_products

This data source provides the Mhub Products of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.138.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "example_value"
}

resource "alicloud_mhub_product" "default" {
  product_name = var.name
}
data "alicloud_mhub_products" "ids" {}
output "mhub_product_id_1" {
  value = data.alicloud_mhub_products.ids.products.0.id
}

data "alicloud_mhub_products" "nameRegex" {
  name_regex = "^my-Product"
}
output "mhub_product_id_2" {
  value = data.alicloud_mhub_products.nameRegex.products.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Product IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Product name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Product names.
* `products` - A list of Mhub Products. Each element contains the following attributes:
	* `id` - The ID of the Product.
	* `product_id` - The ID of the Product.
	* `product_name` - The name of the Product.
