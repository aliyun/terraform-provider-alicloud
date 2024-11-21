---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_basic_accelerator"
sidebar_current: "docs-alicloud-resource-ga-basic-accelerator"
description: |-
  Provides a Alicloud Global Accelerator (GA) Basic Accelerator resource.
---

# alicloud_ga_basic_accelerator

Provides a Global Accelerator (GA) Basic Accelerator resource.

For information about Global Accelerator (GA) Basic Accelerator and how to use it, see [What is Basic Accelerator](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-createbasicaccelerator).

-> **NOTE:** Available since v1.194.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ga_basic_accelerator&exampleId=c51d573b-a0bd-742d-bba5-aedca8d825b93a001f70&activeTab=example&spm=docs.r.ga_basic_accelerator.0.c51d573ba0&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_ga_basic_accelerator" "default" {
  duration               = 1
  pricing_cycle          = "Month"
  basic_accelerator_name = "tf-example-value"
  description            = "tf-example-value"
  bandwidth_billing_type = "BandwidthPackage"
  auto_pay               = true
  auto_use_coupon        = "true"
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth_billing_type` - (Optional, ForceNew) The bandwidth billing method. Valid values: `BandwidthPackage`, `CDT`, `CDT95`.
* `payment_type` - (Optional, ForceNew, Available since v1.208.1) The payment type. Default value: `Subscription`. Valid values: `PayAsYouGo`, `Subscription`.
* `cross_border_status` - (Optional, Bool, Available since v1.208.1) Indicates whether cross-border acceleration is enabled. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `auto_pay` - (Optional, Bool) Specifies whether to enable automatic payment. Default value: `false`. Valid values:
  - `true`: enables automatic payment. Payments are automatically completed.
  - `false`: disables automatic payment. If you select this option, you must go to the Order Center to complete the payment after an order is generated.
* `duration` - (Optional, Int) The subscription duration. Default value: `1`.
  * If the `pricing_cycle` parameter is set to `Month`, the valid values for the `duration` parameter are `1` to `9`.
  * If the `pricing_cycle` parameter is set to `Year`, the valid values for the `duration` parameter are `1` to `3`.
* `pricing_cycle` - (Optional) The billing cycle. Default value: `Month`. Valid values: `Month`, `Year`.
* `auto_use_coupon` - (Optional) Specifies whether to automatically pay bills by using coupons. Default value: `false`. **NOTE:** This parameter is required only if `auto_pay` is set to `true`.
* `auto_renew` - (Optional, Bool) Specifies whether to enable auto-renewal for the GA Basic Accelerator instance. Default value: `false`. Valid values:
  - `true`: enables auto-renewal.
  - `false`: disables auto-renewal.
* `auto_renew_duration` - (Optional, Int) The auto-renewal period. Unit: months. Default value: `1`. Valid values: `1` to `12`. **NOTE:** This parameter is required only if `auto_renew` is set to `true`.
* `promotion_option_no` - (Optional, Available since v1.208.1) The code of the coupon. **NOTE:** The `promotion_option_no` takes effect only for accounts registered on the international site (alibabacloud.com).
* `resource_group_id` - (Optional, Available since v1.226.0) The ID of the resource group. **Note:** Once you set a value of this property, you cannot set it to an empty string anymore.
* `basic_accelerator_name` - (Optional) The name of the Global Accelerator Basic Accelerator instance.
* `description` - (Optional) The description of the Global Accelerator Basic Accelerator instance.
* `tags` - (Optional, Available since v1.207.1) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Basic Accelerator.
* `status` - The status of the Basic Accelerator instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Basic Accelerator.
* `update` - (Defaults to 3 mins) Used when update the Basic Accelerator.
* `delete` - (Defaults to 3 mins) Used when delete the Basic Accelerator.

## Import

Global Accelerator (GA) Basic Accelerator can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_basic_accelerator.example <id>
```
