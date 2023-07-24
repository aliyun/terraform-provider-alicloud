---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_accelerator"
description: |-
  Provides a Alicloud Ga Accelerator resource.
---

# alicloud_ga_accelerator

Provides a Ga Accelerator resource. 

For information about Ga Accelerator and how to use it, see [What is Accelerator](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.209.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}


resource "alicloud_ga_accelerator" "default" {
  instance_charge_type   = "POSTPAY"
  duration               = 1
  auto_pay               = true
  bandwidth_billing_type = "CDT"
  pricing_cycle          = "Month"
}
```

## Argument Reference

The following arguments are supported:
* `accelerator_name` - (Optional, Available since v1.111.0) The Name of the GA instance.
* `auto_pay` - (Optional) Whether to pay automatically, value:
  - **false** (default): automatic payment is not enabled. After generating an order, you need to complete the payment at the order center.
  - **true**: Enable automatic payment to automatically pay for orders.
* `auto_renew` - (Optional) Whether automatic renewal is turned on. Value:
  - **true**: Yes.
  - **false** (default): No.
* `auto_renew_duration` - (Optional, Available since v1.111.0) The duration of automatic renewal. Unit: Month.

Valid values: **1** to **12 * *.
-> **NOTE:**  This item takes effect only when **AutoRenew** is **true.
* `auto_use_coupon` - (Optional, Available since v1.111.0) 当前属性没有在镇元上录入属性描述，请补充后再生成代码。.
* `bandwidth_billing_type` - (Optional, ForceNew, Available since v1.111.0) Bandwidth billing method.
  - **BandwidthPackage**: Billed by bandwidth package.
  - **CDT**: Billing by traffic.
  - **CDT95**: billed at 95 and settled by CDT. This bandwidth billing method is available only to whitelist users.
* `cross_border_mode` - (Optional) cross border.
* `cross_border_status` - (Optional) Whether the cross-border line function is turned on.
* `ddos_id` - (Optional) DDoS high-defense instance ID that is unbound from the global acceleration instance.
* `ddos_region_id` - (Optional) The region where the DDoS pro instance is located. Value:
  - **cn-hangzhou**: Mainland China.
  - **ap-southeast-1**: Non-Mainland China.
* `description` - (Optional, Available since v1.111.0) Descriptive information of the global acceleration instance.
* `duration` - (Optional, Available since v1.111.0) Duration.
* `instance_charge_type` - (Optional) The billing type of the Global Acceleration instance. The default value is PREPAY.
  - PREPAY-Prepaid
  - POSTPAY-POSTPAY.
* `ip_set_config` - (Optional, ForceNew) Accelerate area configuration. See [`ip_set_config`](#ip_set_config) below.
* `pricing_cycle` - (Optional, Available since v1.111.0) PricingCycle.
* `promotion_option_no` - (Optional) Coupon number.
-> **NOTE:**  Only the international station involves this parameter.
* `renewal_status` - (Optional, Available since v1.111.0) Automatic renewal status.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `spec` - (Optional, Available since v1.111.0) The instance type of the GA instance.
* `tags` - (Optional, Map) The tag of the resource.

### `ip_set_config`

The ip_set_config supports the following:
* `access_mode` - (Optional, ForceNew) Accelerated zone access mode. Value:
  - **UserDefine**: Customize the nearest access mode. You can select an acceleration region and region based on your business needs. Global acceleration provides an independent EIP for each acceleration region.
  - **Anycast**: uses automatic nearest access mode. You do not need to configure the acceleration area. Global acceleration provides an Anycast EIP in multiple regions around the world. Users can access the Alibaba Cloud Acceleration Network from the nearest access point through Anycast EIP.
.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Time for global acceleration instance creation.
* `payment_type` - The payment type of the resource.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Accelerator.
* `delete` - (Defaults to 5 mins) Used when delete the Accelerator.
* `update` - (Defaults to 5 mins) Used when update the Accelerator.

## Import

Ga Accelerator can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_accelerator.example <id>
```