---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_basic_accelerator"
sidebar_current: "docs-alicloud-resource-ga-basic-accelerator"
description: |-
  Provides a Alicloud Global Accelerator (GA) Basic Accelerator resource.
---

# alicloud\_ga\_basic\_accelerator

Provides a Global Accelerator (GA) Basic Accelerator resource.

For information about Global Accelerator (GA) Basic Accelerator and how to use it, see [What is Basic Accelerator](https://www.alibabacloud.com/help/en/global-accelerator/latest/createbasicaccelerator).

-> **NOTE:** Available in v1.194.0+.

## Example Usage

Basic Usage

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

* `duration` - (Optional) The subscription duration. Default value: `1`.
  * If the `pricing_cycle` parameter is set to `Month`, the valid values for the `duration` parameter are `1` to `9`.
  * If the `pricing_cycle` parameter is set to `Year`, the valid values for the `duration` parameter are `1` to `3`.
* `pricing_cycle` - (Optional) The billing cycle. Default value: `Month`. Valid values: `Month`, `Year`.
* `basic_accelerator_name` - (Optional) The name of the Global Accelerator Basic Accelerator instance.
* `description` - (Optional) The description of the Global Accelerator Basic Accelerator instance.
* `bandwidth_billing_type` - (Optional, ForceNew) The bandwidth billing method. Valid values: `BandwidthPackage`, `CDT`, `CDT95`.
* `auto_pay` - (Optional) Specifies whether to enable automatic payment. Default value: `false`. Valid values:
  - `true`: enables automatic payment. Payments are automatically completed.
  - `false`: disables automatic payment. If you select this option, you must go to the Order Center to complete the payment after an order is generated.
* `auto_use_coupon` - (Optional) Specifies whether to automatically pay bills by using coupons. Default value: `false`. **NOTE:** This parameter is required only if `auto_pay` is set to `true`.
* `auto_renew` - (Optional) Specifies whether to enable auto-renewal for the GA Basic Accelerator instance. Default value: `false`. Valid values:
  - `true`: enables auto-renewal.
  - `false`: disables auto-renewal.
* `auto_renew_duration` - (Optional) The auto-renewal period. Unit: months. Valid values: `1` to `12`. Default value: `1`. **NOTE:** This parameter is required only if `auto_renew` is set to `true`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Basic Accelerator.
* `status` - The status of the Basic Accelerator instance.

#### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Basic Accelerator.
* `update` - (Defaults to 3 mins) Used when update the Basic Accelerator.
* `delete` - (Defaults to 3 mins) Used when delete the Basic Accelerator.

## Import

Global Accelerator (GA) Basic Accelerator can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_basic_accelerator.example <id>
```
