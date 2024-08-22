---
subcategory: "Enterprise Mobile Application Studio (MHUB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mhub_product"
sidebar_current: "docs-alicloud-resource-mhub-product"
description: |-
  Provides a Alicloud MHUB Product resource.
---

# alicloud_mhub_product

Provides a MHUB Product resource.

For information about MHUB Product and how to use it, see [What is Product](https://help.aliyun.com/product/65109.html).

-> **NOTE:** Available since v1.138.0+.

-> **NOTE:** At present, the resource only supports cn-shanghai region.

## Example Usage
<div class="oics-button" style="float: right;margin: 0 0 -40px 0;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_mhub_product&exampleId=72f57c0b-4c03-59e4-a79d-a9073353139c29b2ad49&activeTab=example&spm=docs.r.mhub_product.0.72f57c0b4c" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-shanghai"
}

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

```shell
$ terraform import alicloud_mhub_product.example <id>
```
