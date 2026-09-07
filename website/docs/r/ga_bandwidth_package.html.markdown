---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_bandwidth_package"
sidebar_current: "docs-alicloud-resource-ga-bandwidth-package"
description: |-
  Provides a Alicloud Global Accelerator (GA) Bandwidth Package resource.
---

# alicloud_ga_bandwidth_package

Provides a Global Accelerator (GA) Bandwidth Package resource.

For information about Global Accelerator (GA) Bandwidth Package and how to use it, see [What is Bandwidth Package](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-createbandwidthpackage).

-> **NOTE:** At present, The `alicloud_ga_bandwidth_package` created with `Subscription` cannot be deleted. you need to wait until the resource is outdated and released automatically.

-> **NOTE:** Available since v1.112.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ga_bandwidth_package&exampleId=f612bc89-4a69-713e-14a6-97f9b2248ab0e0f3e3bc&activeTab=example&spm=docs.r.ga_bandwidth_package.0.f612bc894a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_ga_bandwidth_package" "example" {
  bandwidth      = 20
  type           = "Basic"
  bandwidth_type = "Basic"
  duration       = 1
  auto_pay       = true
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ga_bandwidth_package&spm=docs.r.ga_bandwidth_package.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `bandwidth` - (Required, Int) The bandwidth value of bandwidth packet.
* `type` - (Required, ForceNew) The type of the bandwidth packet. China station only supports return to basic. Valid values: `Basic`, `CrossDomain`.
* `bandwidth_type` - (Optional) The bandwidth type of the bandwidth. Valid values: `Advanced`, `Basic`, `Enhanced`. If `type` is set to `Basic`, this parameter is required.
-> **NOTE:** At present, only basic can be configured to enhanced, but not enhanced and advanced to other types of accelerated bandwidth.
* `payment_type` - (Optional, ForceNew) The payment type of the bandwidth. Default value: `Subscription`. Valid values: `PayAsYouGo`, `Subscription`.
* `billing_type` - (Optional, ForceNew) The billing type. Valid values: `PayBy95`, `PayByTraffic`. **NOTE:** `billing_type` is valid only when `payment_type` is set to `PayAsYouGo`.
* `ratio` - (Optional, ForceNew, Int) The minimum percentage for the pay-by-95th-percentile metering method. Valid values: `30` to `100`. **NOTE:** `ratio` is valid only when `billing_type` is set to `PayBy95`.
* `cbn_geographic_region_ida` - (Optional, ForceNew) Interworking area A of cross domain acceleration package. Only international stations support returning this parameter. Default value: `China-mainland`.
* `cbn_geographic_region_idb` - (Optional, ForceNew) Interworking area B of cross domain acceleration package. Only international stations support returning this parameter. Default value: `Global`.
* `auto_pay` - (Optional, Bool) Whether to pay automatically. Valid values:
  - `false`: If automatic payment is not enabled, you need to go to the order center to complete the payment after the order is generated.
  - `true`: Enable automatic payment, automatic payment order.
* `duration` - (Optional) The subscription duration. **NOTE:** The ForceNew attribute has be removed from version 1.148.0. If `payment_type` is set to `Subscription`, this parameter is required.
* `auto_use_coupon` - (Optional, Bool) Whether use vouchers. Default value: `false`. Valid values:
  - `false`: Not used.
  - `true`: Use.
* `renewal_status` - (Optional, Available since v1.169.0) Whether to renew a bandwidth packet. automatically or not. Valid values:
  - `AutoRenewal`: Enable auto renewal.
  - `Normal`: Disable auto renewal.
  - `NotRenewal`: No renewal any longer. After you specify this value, Alibaba Cloud stop sending notification of instance expiry, and only gives a brief reminder on the third day before the instance expiry.
* `auto_renew_duration` - (Optional, Int, Available since v1.169.0) Auto renewal period of a bandwidth packet, in the unit of month. Valid values: `1` to `12`.
* `promotion_option_no` - (Optional, Available since v1.208.0) The code of the coupon. **NOTE:** The `promotion_option_no` takes effect only for accounts registered on the international site (alibabacloud.com).
* `resource_group_id` - (Optional, Available since v1.226.0) The ID of the resource group. **Note:** Once you set a value of this property, you cannot set it to an empty string anymore.
* `bandwidth_package_name` - (Optional) The name of the bandwidth packet.
* `description` - (Optional) The description of bandwidth package.
* `tags` - (Optional, Available since v1.207.1) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Bandwidth Package.
* `status` - The status of the Bandwidth Package.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Bandwidth Package.
* `update` - (Defaults to 5 mins) Used when update the Bandwidth Package.
* `delete` - (Defaults to 3 mins) Used when delete the Bandwidth Package.

## Import

Ga Bandwidth Package can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_bandwidth_package.example <id>
```
