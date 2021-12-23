---
subcategory: "Market Place"
layout: "alicloud"
page_title: "Alicloud: alicloud_market_order"
sidebar_current: "docs-alicloud-resource-market-order"
description: |-
    Provides a market order resource.
---

# alicloud\_market\_order

Provides a market order resource.

-> **NOTE:** Terraform will auto build a market order  while it uses `alicloud_market_order` to build a market order resource.

-> **NOTE:** Available in 1.69.0+

## Example Usage

Basic Usage

```
resource "alicloud_market_order" "order" {
  product_code    = "cmapi033136"
  pay_type        = "prepay"
  quantity        = 1
  duration        = 1
  pricing_cycle   = "Month"
  package_version = "yuncode2713600001"
  coupon_id       = ""
}
```

## Argument Reference

The following arguments are supported:

* `product_code` - (Required, ForceNew) The product_code of market place product.
* `pay_type` - (Optional, ForceNew) Valid values are `PrePaid`, `PostPaid`,System default to `PostPaid`.
* `duration` - (Optional, ForceNew) The number of purchase cycles.
* `pricing_cycle` - (Required, ForceNew) The purchase cycle of the product, valid values are `Day`, `Month` and `Year`.
* `package_version` - (Required, ForceNew) The package version of the market product.
* `quantity` - (Optional, ForceNew) The quantity of the market product will be purchased.
* `coupon_id` - (Optional, ForceNew) The coupon id of the market product.
* `components` - (Optional, ForceNew) Service providers customize additional components.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the market order.

## Import

Market order can be imported using the id, e.g.

```
$ terraform import alicloud_market_order.order your-order-id
```
 
