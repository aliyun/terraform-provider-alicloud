---
subcategory: "Classic Load Balancer (SLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_load_balancer"
sidebar_current: "docs-alicloud-resource-slb-load-balancer"
description: |-
  Provides an Application Load Balancer resource.
---

# alicloud_slb_load_balancer

Provides an Application Load Balancer resource.

-> **NOTE:** Available in 1.123.1+

-> **NOTE:** At present, to avoid some unnecessary regulation confusion, SLB can not support alicloud international account to create `PayByBandwidth` instance.

-> **NOTE:** The supported specifications vary by region. Currently, not all regions support guaranteed-performance instances.
For more details about guaranteed-performance instance, see [Guaranteed-performance instances](https://www.alibabacloud.com/help/en/server-load-balancer/latest/createloadbalancer-2#t4182.html).

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_slb_load_balancer&exampleId=d3db5e79-78d0-fbc5-20d4-cb5cf51ed710a7863d05&activeTab=example&spm=docs.r.slb_load_balancer.0.d3db5e7978&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
# Create a intranet SLB instance
variable "slb_load_balancer_name" {
  default = "forSlbLoadBalancer"
}

data "alicloud_zones" "load_balancer" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "load_balancer" {
  vpc_name = var.slb_load_balancer_name
}

resource "alicloud_vswitch" "load_balancer" {
  vpc_id       = alicloud_vpc.load_balancer.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_zones.load_balancer.zones[0].id
  vswitch_name = var.slb_load_balancer_name
}

resource "alicloud_slb_load_balancer" "load_balancer" {
  load_balancer_name = var.slb_load_balancer_name
  address_type       = "intranet"
  load_balancer_spec = "slb.s2.small"
  vswitch_id         = alicloud_vswitch.load_balancer.id
  tags = {
    info = "create for internet"
  }
  instance_charge_type = "PayBySpec"
}
```

### Deleting `alicloud_slb_load_balancer` or removing it from your configuration

The `alicloud_slb_load_balancer` resource allows you to manage `payment_type = "Subscription"` load balancer, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Load Balancer.
You can resume managing the subscription load balancer via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:

* `load_balancer_name` - (Optional) The name of the SLB. This name must be unique within your AliCloud account, can have a maximum of 80 characters,
must contain only alphanumeric characters or hyphens, such as "-","/",".","_", and must not begin or end with a hyphen. If not specified,
Terraform will autogenerate a name beginning with `tf-lb`.
* `address_type` - (Optional, ForceNew) The network type of the SLB instance. Valid values: ["internet", "intranet"]. If load balancer launched in VPC, this value must be `intranet`.
    - internet: After an Internet SLB instance is created, the system allocates a public IP address so that the instance can forward requests from the Internet.
    - intranet: After an intranet SLB instance is created, the system allocates an intranet IP address so that the instance can only forward intranet requests.
* `internet_charge_type` - (Optional) Valid values are `PayByBandwidth`, `PayByTraffic`. If this value is `PayByBandwidth`, then argument `address_type` must be `internet`. Default is `PayByTraffic`. If load balancer launched in VPC, this value must be `PayByTraffic`. Before version 1.10.1, the valid values are `paybybandwidth` and `paybytraffic`.
* `bandwidth` - (Optional) Valid value is between 1 and 5120, If argument `internet_charge_type` is `PayByTraffic`, then this value will be ignored.
* `vswitch_id` - (Optional, ForceNew) The VSwitch ID to launch in. **Note:** Required for a VPC SLB. If `address_type` is internet, it will be ignored.
* `load_balancer_spec` - (Optional) The specification of the Server Load Balancer instance. Default to empty string indicating it is "Shared-Performance" instance.
 Launching "Performance-guaranteed" instance, it must be specified. Valid values: `slb.s1.small`, `slb.s2.small`, `slb.s2.medium`,
 `slb.s3.small`, `slb.s3.medium`, `slb.s3.large` and `slb.s4.large`. It will be ignored when `instance_charge_type = "PayByCLCU"`.
* `tags` - (Optional, Computed) A mapping of tags to assign to the resource. The `tags` can have a maximum of 10 tag for every load balancer instance. This filed mark as `Computed` since v1.217.1.
* `payment_type` - (Optional) The billing method of the load balancer. Valid values are `PayAsYouGo` and `Subscription`. Default to `PayAsYouGo`.
* `period` - (Optional) The duration that you will buy the resource, in month. It is valid when `PaymentType` is `Subscription`. Default to 1. Valid values: [1-9, 12, 24, 36]. This attribute is only used to create `Subscription` instance or modify the `PayAsYouGo` instance to `Subscription`. Once effect, it will not be modified that means running `terraform apply` will not affect the resource.
* `master_zone_id` - (Optional, ForceNew) The primary zone ID of the SLB instance. If not specified, the system will be randomly assigned. You can query the primary and standby zones in a region by calling the [DescribeZone](https://help.aliyun.com/document_detail/27585.htm) API.
* `slave_zone_id` - (Optional, ForceNew) The standby zone ID of the SLB instance. If not specified, the system will be randomly assigned. You can query the primary and standby zones in a region by calling the DescribeZone API.
* `delete_protection` - (Optional) Whether enable the deletion protection or not. on: Enable deletion protection. off: Disable deletion protection. Default to off. Only postpaid instance support this function.   
* `address_ip_version` - (Optional) The IP version of the SLB instance to be created, which can be set to `ipv4` or `ipv6` . Default to `ipv4`. Now, only internet instance support `ipv6` address.
* `address` - (Optional) Specify the IP address of the private network for the SLB instance, which must be in the destination CIDR block of the corresponding switch.
* `resource_group_id` - (Optional, ForceNew) The id of resource group which the SLB belongs.
* `modification_protection_reason` - (Optional) The reason of modification protection. It's effective when `modification_protection_status` is `ConsoleProtection`.
* `modification_protection_status` - (Optional) The status of modification protection. Valid values: `ConsoleProtection` and `NonProtection`. Default value is `NonProtection`.
* `status` - (Optional) The status of slb load balancer. Valid values: `active` and `inactice`. The system default value is `active`.
* `name` - (Optional, Deprecated from v1.123.1) Field `name` has been deprecated from provider version 1.123.1 New field `load_balancer_name` instead.
* `instance_charge_type` - (Optional, V1.193.0+) Support `PayBySpec` (default) and `PayByCLCU`, This parameter takes effect when the value of **payment_type** (instance payment mode) is **PayAsYouGo** (pay-as-you-go).
* `specification` - (Optional, Deprecated from v1.123.1) Field `specification` has been deprecated from provider version 1.123.1 New field `load_balancer_spec` instead.
* `internet` - (Optional, Deprecated from v1.124.0) Field `internet` has been deprecated from provider version 1.124.0 New field `address_type` instead.

-> **NOTE:** A "Shared-Performance" instance can be changed to "Performance-guaranteed", but the change is irreversible.

-> **NOTE:** To change a "Shared-Performance" instance to a "Performance-guaranteed" instance, the SLB will have a short probability of business interruption (10 seconds-30 seconds). Advise to change it during the business downturn, or migrate business to other SLB Instances by using GSLB before changing.

-> **NOTE:** Currently, the alibaba cloud international account does not support creating a `Subscription` SLB instance.

-> **NOTE:** This parameter `instance_charge_type` is only valid for China sites and only if the `payment_type` value is `PayAsYouGo`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the load balancer.
* `address` - The IP address of the load balancer.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the SLB load balancer.(until it reaches the initial `active` status). 
* `delete` - (Defaults to 9 mins) Used when terminating the SLB load balancer.

## Import

Load balancer can be imported using the id, e.g.

```shell
$ terraform import alicloud_slb_load_balancer.example lb-abc123456
```
