---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_accelerator"
sidebar_current: "docs-alicloud-resource-ga-accelerator"
description: |-
  Provides a Alicloud Global Accelerator (GA) Accelerator resource.
---

# alicloud_ga_accelerator

Provides a Global Accelerator (GA) Accelerator resource.

For information about Global Accelerator (GA) Accelerator and how to use it, see [What is Accelerator](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-createaccelerator).

-> **NOTE:** Available since v1.111.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ga_accelerator&exampleId=4c5ac150-b6e5-64a4-04b9-bbc59be166ab875af5e4&activeTab=example&spm=docs.r.ga_accelerator.0.4c5ac150b6&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_ga_accelerator" "example" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}
```

## Argument Reference

The following arguments are supported:

* `spec` - (Optional) The instance type of the GA instance. Specification of global acceleration instance. Valid values:
  - `1`: Small 1.
  - `2`: Small 2.
  - `3`: Small 3.
  - `5`: Medium 1.
  - `8`: Medium 2.
  - `10`: Medium 3.
* `bandwidth_billing_type` - (Optional, ForceNew, Available since v1.205.0) The bandwidth billing method. Default value: `BandwidthPackage`. Valid values:
  - `BandwidthPackage`: billed based on bandwidth plans.
  - `CDT`: billed based on data transfer.
* `payment_type` - (Optional, ForceNew, Available since v1.208.1) The payment type. Default value: `Subscription`. Valid values: `PayAsYouGo`, `Subscription`.
* `cross_border_status` - (Optional, Bool, Available since v1.208.1) Indicates whether cross-border acceleration is enabled. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.
**NOTE:** `cross_border_status` is valid only when `bandwidth_billing_type` is set to `CDT`.
* `cross_border_mode` - (Optional, Available since v1.208.1) The type of cross-border acceleration. Default value: `bgpPro`. Valid values: `bgpPro`, `private`. **NOTE:** `cross_border_mode` is valid only when `cross_border_status` is set to `true`.
* `duration` - (Optional, Int) The subscription duration.
  * If the `pricing_cycle` parameter is set to `Month`, the valid values for the `duration` parameter are 1 to 9.
  * If the `pricing_cycle` parameter is set to `Year`, the valid values for the `duration` parameter are 1 to 3.
* `pricing_cycle`- (Optional, Available since v1.150.0) The billing cycle of the GA instance. Default value: `Month`. Valid values:
  - `Month`: billed on a monthly basis.
  - `Year`: billed on an annual basis.
* `auto_use_coupon` - (Optional, Bool) Use coupons to pay bills automatically. Default value: `false`. Valid values:
  - `true`: Use.
  - `false`: Not used.
* `renewal_status` - (Optional, Available since v1.146.0) Whether to renew an accelerator automatically or not. Default value: `Normal`. Valid values:
  - `AutoRenewal`: Enable auto renewal.
  - `Normal`: Disable auto renewal.
  - `NotRenewal`: No renewal any longer. After you specify this value, Alibaba Cloud stop sending notification of instance expiry, and only gives a brief reminder on the third day before the instance expiry.
* `auto_renew_duration` - (Optional, Int, Available since v1.146.0) Auto renewal period of an instance, in the unit of month. The value range is 1-12.
* `promotion_option_no` - (Optional, Available since v1.208.1) The code of the coupon. **NOTE:** The `promotion_option_no` takes effect only for accounts registered on the international site (alibabacloud.com).
* `resource_group_id` - (Optional, Available since v1.226.0) The ID of the resource group. **Note:** Once you set a value of this property, you cannot set it to an empty string anymore.
* `accelerator_name` - (Optional) The Name of the GA instance.
* `description` - (Optional) Descriptive information of the global acceleration instance.
* `tags` - (Optional, Available since v1.207.1) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Accelerator. Value as `accelerator_id`.
* `status` - The status of the GA instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Ga Accelerator.
* `update` - (Defaults to 6 mins) Used when updating the Ga Accelerator.
* `delete` - (Defaults to 3 mins) Used when deleting the Ga Accelerator.

## Import

Ga Accelerator can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_accelerator.example <id>
```
