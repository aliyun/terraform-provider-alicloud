---
subcategory: "PAI Workspace"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_workspace_product"
description: |-
  Provides a Alicloud PAI Workspace Product resource.
---

# alicloud_pai_workspace_product

Provides a PAI Workspace Product resource. Creating the resource places a real product order through the `CreateProductOrders` API.

-> **NOTE:** Available since v1.236.0.

-> **NOTE:** `alicloud_pai_workspace_product` is an order-only resource. There is no Update or Delete API: a product order is irreversible. Any change to a configurable attribute forces replacement (a new order), and `terraform destroy` only removes the resource from Terraform state without canceling the purchased product. Use `auto_pay = false` (the default) and review the order in the console before paying.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_pai_workspace_product" "default" {
  product_code  = "PAI_share"
  charge_type   = "POSTPAY"
  pricing_cycle = "Hour"
  order_type    = "BUY"
  auto_pay      = false
}
```

## Argument Reference

The following arguments are supported:

* `product_code` - (Required, Forces new order) The code of the product to purchase. Valid value: `PAI_share`.
* `charge_type` - (Optional, Forces new order) The charge type of the order, such as `POSTPAY`, `PayAsYouGo` or `Subscription`.
* `pricing_cycle` - (Optional, Forces new order) The pricing cycle, such as `Hour`, `Month` or `Year`. The unit of `duration` depends on it.
* `duration` - (Optional, Forces new order) The purchase duration; the unit is determined by `pricing_cycle`.
* `auto_renew` - (Optional, Forces new order) Whether to auto renew the product.
* `order_type` - (Optional, Forces new order) The order type, such as `BUY`, `New` or `Renew`.
* `auto_pay` - (Optional, Forces new order) Whether to auto pay for all products in the order. Defaults to `false`.
* `instance_properties` - (Optional, Forces new order) The instance properties attached to the order. See [`instance_properties`](#instance_properties) below.

### instance_properties

The `instance_properties` block supports:

* `code` - (Optional) The property code.
* `name` - (Optional) The property name.
* `value` - (Optional) The property value.

-> **NOTE:** The `instance_properties` block is write-only. The `CreateProductOrders` API does not return it, so the values are preserved from configuration and cannot be read back.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. Same as `product_id`.
* `product_id` - The product ID returned by `CreateProductOrders`, the unique key of the purchased product.
* `is_purchased` - Whether the product has been purchased.
* `has_permission_to_purchase` - Whether the caller has permission to purchase the product.
* `purchase_url` - The URL to purchase the product.
* `region_id` - The region ID of the resource.

## Import

PAI Workspace Product can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_workspace_product.default <product_id>
```

-> **NOTE:** The create-only order attributes (`charge_type`, `pricing_cycle`, `duration`, `auto_renew`, `order_type`, `auto_pay` and `instance_properties`) are not returned by the read API. After import they are empty in state; set them in configuration only when placing a new order.
