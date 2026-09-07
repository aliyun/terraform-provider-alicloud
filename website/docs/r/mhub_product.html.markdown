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

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_mhub_product&exampleId=72f57c0b-4c03-59e4-a79d-a9073353139c29b2ad49&activeTab=example&spm=docs.r.mhub_product.0.72f57c0b4c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-shanghai"
}

resource "alicloud_mhub_product" "example" {
  product_name = "example_value"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_mhub_product&spm=docs.r.mhub_product.example&intl_lang=EN_US)

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
