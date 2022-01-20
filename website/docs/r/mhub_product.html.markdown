---
subcategory: "Enterprise Mobile Application Studio (MHUB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mhub_product"
sidebar_current: "docs-alicloud-resource-mhub-product"
description: |-
  Provides a Alicloud MHUB Product resource.
---

# alicloud\_mhub\_product

Provides a MHUB Product resource.

For information about MHUB Product and how to use it, see [What is Product](https://help.aliyun.com/product/65109.html).

-> **NOTE:** Available in v1.138.0+.

-> **NOTE:** At present, the resource only supports cn-shanghai region.

## Example Usage

Basic Usage

```terraform
resource "alicloud_mhub_product" "example" {
  product_name = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `product_name` - (Required) ProductName.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Product.

## Import

MHUB Product can be imported using the id, e.g.

```
$ terraform import alicloud_mhub_product.example <id>
```
