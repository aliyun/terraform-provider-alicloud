---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_bandwidth_package"
sidebar_current: "docs-alicloud-resource-ga-bandwidth-package"
description: |-
  Provides a Alicloud Global Accelerator (GA) Bandwidth Package resource.
---

# alicloud\_ga\_bandwidth\_package

Provides a Global Accelerator (GA) Bandwidth Package resource.

For information about Global Accelerator (GA) Bandwidth Package and how to use it, see [What is Bandwidth Package](https://www.alibabacloud.com/help/en/doc-detail/153241.htm).

-> **NOTE:** At present, The `alicloud_ga_bandwidth_package` created with `Subscription` cannot be deleted. you need to wait until the resource is outdated and released automatically.

-> **NOTE:** Available in v1.112.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ga_bandwidth_package" "example" {
  bandwidth      = 20
  type           = "Basic"
  bandwidth_type = "Basic"
  duration       = 1
  auto_pay       = true
  ratio          = 30
}

```

## Argument Reference

The following arguments are supported:

* `auto_pay` - (Optional) Whether to pay automatically. Valid values:
    `false`: If automatic payment is not enabled, you need to go to the order center to complete the payment after the order is generated.
    `true`: Enable automatic payment, automatic payment order.
* `auto_use_coupon` - (Optional) Whether use vouchers. Default value is `false`. Valid values: `false`: Not used, `true`: Use.
* `bandwidth` - (Required) The bandwidth value of bandwidth packet.
* `bandwidth_package_name` - (Optional) The name of the bandwidth packet.
* `bandwidth_type` - (Optional) The bandwidth type of the bandwidth. Valid values: `Advanced`, `Basic`, `Enhanced`. If `type` is set to `Basic`, this parameter is required.

-> **NOTE:** At present, only basic can be configured to enhanced, but not enhanced and advanced to other types of accelerated bandwidth.

* `billing_type` - (Optional, ForceNew) The billing type. Valid values: `PayBy95`, `PayByTraffic`.
* `cbn_geographic_region_ida` - (Optional, ForceNew) Interworking area A of cross domain acceleration package. Only international stations support returning this parameter. Default value is `China-mainland`.
* `cbn_geographic_region_idb` - (Optional, ForceNew) Interworking area B of cross domain acceleration package. Only international stations support returning this parameter. Default value is `Global`.
* `description` - (Optional) The description of bandwidth package.
* `duration` - (Optional) The subscription duration. **NOTE:** The ForceNew attribute has be removed from version 1.148.0. If `payment_type` is set to `Subscription`, this parameter is required.
* `payment_type` - (Optional, ForceNew) The payment type of the bandwidth. Valid values: `PayAsYouGo`, `Subscription`. Default value is `Subscription`.
* `ratio` - (Optional, ForceNew) The minimum percentage for the pay-by-95th-percentile metering method. Valid values: 30 to 100.
* `type` - (Required, ForceNew) The type of the bandwidth packet. China station only supports return to basic. Valid values: `Basic`, `CrossDomain`.
* `auto_renew_duration` - (Optional, Available in 1.169.0+) Auto renewal period of a bandwidth packet, in the unit of month. The value range is 1-12.
* `renewal_status` - (Optional, Computed, Available in 1.169.0+) Whether to renew a bandwidth packet. automatically or not. Valid values:
  - `AutoRenewal`: Enable auto renewal.
  - `Normal`: Disable auto renewal.
  - `NotRenewal`: No renewal any longer. After you specify this value, Alibaba Cloud stop sending notification of instance expiry, and only gives a brief reminder on the third day before the instance expiry.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Bandwidth Package.
* `status` - The status of the bandwidth plan.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Bandwidth Package.
* `update` - (Defaults to 2 mins) Used when updating the Bandwidth Package.

## Import

Ga Bandwidth Package can be imported using the id, e.g.

```
$ terraform import alicloud_ga_bandwidth_package.example <id>
```
