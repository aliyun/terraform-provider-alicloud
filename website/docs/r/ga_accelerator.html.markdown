---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_accelerator"
sidebar_current: "docs-alicloud-resource-ga-accelerator"
description: |-
  Provides a Alicloud Global Accelerator (GA) Accelerator resource.
---

# alicloud\_ga\_accelerator

Provides a Global Accelerator (GA) Accelerator resource.

For information about Global Accelerator (GA) Accelerator and how to use it, see [What is Accelerator](https://help.aliyun.com/document_detail/153235.html).

-> **NOTE:** At present, The `alicloud_ga_accelerator` cannot be deleted. you need to wait until the resource is outdated and released automatically.

-> **NOTE:** Available in v1.111.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ga_accelerator" "example" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_name` - (Optional) The Name of the GA instance.
* `auto_use_coupon` - (Optional) Use coupons to pay bills automatically. Default value is `false`. Valid value: `true`: Use, `false`: Not used.
* `description` - (Optional) Descriptive information of the global acceleration instance.
* `duration` - (Required) The subscription duration. **NOTE:** Starting from v1.150.0+, the `duration` and  `pricing_cycle` are both required.
  * If the `pricing_cycle` parameter is set to `Month`, the valid values for the `duration` parameter are 1 to 9.
  * If the `pricing_cycle` parameter is set to `Year`, the valid values for the `duration` parameter are 1 to 3.
* `pricing_cycle`- (Optional, Available in v1.150.0+) The billing cycle of the GA instance. Valid values: `Month`,`Year`. The default value: `Month`.
  * `Month`: billed on a monthly basis.
  * `Year`: billed on an annual basis.
* `spec` - (Required) The instance type of the GA instance. Specification of global acceleration instance, value:
    `1`: Small 1.
    `2`: Small 2.
    `3`: Small 3.
    `5`: Medium 1.
    `8`: Medium 2.
    `10`: Medium 3.

* `auto_renew_duration` - (Optional, Available in v1.146.0+) Auto renewal period of an instance, in the unit of month. The value range is 1-12.
* `renewal_status` - (Optional, Available in v1.146.0+) Whether to renew an accelerator automatically or not. Default to "Normal". Valid values:
    - `AutoRenewal`: Enable auto renewal.
    - `Normal`: Disable auto renewal.
    - `NotRenewal`: No renewal any longer. After you specify this value, Alibaba Cloud stop sending notification of instance expiry, and only gives a brief reminder on the third day before the instance expiry.
    
### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Ga Accelerator.
* `update` - (Defaults to 6 mins) Used when updating the Ga Accelerator.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Accelerator. Value as `accelerator_id`.
* `status` - The status of the GA instance.

## Import

Ga Accelerator can be imported using the id, e.g.

```
$ terraform import alicloud_ga_accelerator.example <accelerator_id>
```
