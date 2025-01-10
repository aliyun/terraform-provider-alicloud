---
subcategory: "Cloud Control"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_control_prices"
sidebar_current: "docs-alicloud-datasource-cloud-control-prices"
description: |-
  Provides a list of Cloud Control Price owned by an Alibaba Cloud account.
---

# alicloud_cloud_control_prices

This data source provides Cloud Control Price available to the user.[What is Price](https://www.alibabacloud.com/help/en/)

-> **NOTE:** Available since v1.241.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_cloud_control_prices" "default" {
  desire_attributes = {
    AddressType = "internet"
    PaymentType = "PayAsYouGo"
  }
  product       = "SLB"
  resource_code = "LoadBalancer"
}

output "alicloud_cloud_control_price_example_id" {
  value = data.alicloud_cloud_control_prices.default.prices.0.discount_price
}
```

## Argument Reference

The following arguments are supported:
* `desire_attributes` - (Optional, ForceNew) This property represent the detailed configuration of the Resource which you are going to get price.  Give same content as DesireAttributes of the 'Resource' Resource when start Create operation. 'PaymentType' is necessary when in DesireAttributes.  Here is a probably example when you get the price of SLB LoadBalancer:```json{"LoadBalancerName": "cc-test","Bandwidth": 6,"PaymentType": "PayAsYouGo","AddressType": "internet","LoadBalancerSpec": "slb.s3.small","InternetChargeType": "paybybandwidth"} See [`DesireAttributes`](#DesireAttributes) below.
* `product` - (Required, ForceNew) The product Code represents the product to be operated. Currently supported products and resources can be queried at the following link: [supported-services-and-resource-types](https://help.aliyun.com/zh/cloud-control-api/product-overview/supported-services-and-resource-types).
* `resource_code` - (Required, ForceNew) Resource Code, if there is a parent resource, split with `::`, such as VPC::VSwitch. The supported resource Code can be obtained from the following link: [supported-services-and-resource-types](https://help.aliyun.com/zh/cloud-control-api/product-overview/supported-services-and-resource-types).
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `prices` - A list of Price Entries. Each element contains the following attributes:
  * `currency` - Currency. Value range: CNY: RMB. USD: USD. JPY: Japanese yen.
  * `discount_price` - Discount
  * `module_details` - Pricing Module Price Details
    * `cost_after_discount` - Preferential price.
    * `invoice_discount` - Discount.
    * `module_code` - Valuation Module Identification.
    * `module_name` - Pricing Module Name.
    * `original_cost` - Original Price.
    * `price_type` - Price Type.
  * `original_price` - Original Price
  * `promotion_details` - Offer Details
    * `promotion_desc` - Offer Description.
    * `promotion_id` - Offer logo.
    * `promotion_name` - Offer Name.
  * `trade_price` - Preferential price