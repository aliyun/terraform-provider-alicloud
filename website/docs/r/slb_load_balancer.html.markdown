---
subcategory: "Classic Load Balancer (CLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_load_balancer"
sidebar_current: "docs-alicloud-resource-slb-load-balancer"
description: |-
  Provides an Application Load Banlancer resource.
---

# alicloud\_slb\_load\_balancer

Provides an Application Load Balancer resource.

-> **NOTE:** Available in 1.123.1+

-> **NOTE:** At present, to avoid some unnecessary regulation confusion, SLB can not support alicloud international account to create `PayByBandwidth` instance.

-> **NOTE:** The supported specifications vary by region. Currently, not all regions support guaranteed-performance instances.
For more details about guaranteed-performance instance, see [Guaranteed-performance instances](https://www.alibabacloud.com/help/doc-detail/27657.htm).

## Example Usage

```terraform
# Create a intranet SLB instance
variable "name" {
  default = "terraformtestslbconfig"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name = var.name
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = var.name
  address_type       = "intranet"
  load_balancer_spec = "slb.s2.small"
  vswitch_id         = alicloud_vswitch.default.id
  tags = {
    info = "create for internet"
  }
}
```

### Deleting `alicloud_slb_load_balancer` or removing it from your configuration

The `alicloud_slb_load_balancer` resource allows you to manage `payment_type = "Subscription"` or `instance_charge_type = "Prepaid"` load balancer, but Terraform cannot destroy it.
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
* `bandwidth` - (Optional) Valid value is between 1 and 1000, If argument `internet_charge_type` is `PayByTraffic`, then this value will be ignore.
* `vswitch_id` - (Optional, ForceNew) The VSwitch ID to launch in. **Note:** Required for a VPC SLB. If `address_type` is internet, it will be ignore.
* `load_balancer_spec` - (Optional) The specification of the Server Load Balancer instance. Default to empty string indicating it is "Shared-Performance" instance.
 Launching "[Performance-guaranteed](https://www.alibabacloud.com/help/doc-detail/27657.htm)" instance, it is must be specified and it valid values are: `slb.s1.small`, `slb.s2.small`, `slb.s2.medium`,
 `slb.s3.small`, `slb.s3.medium`, `slb.s3.large` and `slb.s4.large`.
* `tags` - (Optional) A mapping of tags to assign to the resource. The `tags` can have a maximum of 10 tag for every load balancer instance.
* `payment_type` - (Optional) The billing method of the load balancer. Valid values are `PayAsYouGo` and `Subscription`. Default to `PayAsYouGo`.
* `period` - (Optional) The duration that you will buy the resource, in month. It is valid when `PaymentType` is `Subscription`. Default to 1. Valid values: [1-9, 12, 24, 36]. This attribute is only used to create `Subscription` instance or modify the `PayAsYouGo` instance to `Subscription`. Once effect, it will not be modified that means running `terraform apply` will not affact the resource.
* `master_zone_id` - (Optional, ForceNew) The primary zone ID of the SLB instance. If not specified, the system will be randomly assigned. You can query the primary and standby zones in a region by calling the [DescribeZone](https://help.aliyun.com/document_detail/27585.htm) API.
* `slave_zone_id` - (Optional, ForceNew) The standby zone ID of the SLB instance. If not specified, the system will be randomly assigned. You can query the primary and standby zones in a region by calling the DescribeZone API.
* `delete_protection` - (Optional) Whether enable the deletion protection or not. on: Enable deletion protection. off: Disable deletion protection. Default to off. Only postpaid instance support this function.   
* `address_ip_version` - (Optional) The IP version of the SLB instance to be created, which can be set to `ipv4` or `ipv6` . Default to `ipv4`. Now, only internet instance support `ipv6` address.
* `address` - (Optional) Specify the IP address of the private network for the SLB instance, which must be in the destination CIDR block of the corresponding switch.
* `resource_group_id` - (Optional, ForceNew) The Id of resource group which the SLB belongs.
* `modification_protection_reason` - (Optional) The reason of modification protection. It's effective when `modification_protection_status` is `ConsoleProtection`.
* `modification_protection_status` - (Optional) The status of modification protection. Valid values: `ConsoleProtection` and `NonProtection`. Default value is `NonProtection`.
* `status` - (Optional) The status of slb load balancer. Valid values: `active` and `inactice`. The system default value is `active`.
* `name` - (Optional, Deprecated form v1.123.1) Field `name` has been deprecated from provider version 1.123.1 New field `load_balancer_name` instead.
* `instance_charge_type` - (Optional, Deprecated form v1.124.0) Field `instance_charge_type` has been deprecated from provider version 1.124.0 New field `payment_type` instead.
* `specification` - (Optional, Deprecated form v1.123.1) Field `specification` has been deprecated from provider version 1.123.1 New field `load_balancer_spec` instead.
* `internet` - (Optional, Deprecated form v1.124.0) Field `internet` has been deprecated from provider version 1.124.0 New field `address_type` instead.

-> **NOTE:** A "Shared-Performance" instance can be changed to "Performance-guaranteed", but the change is irreversible.

-> **NOTE:** To change a "Shared-Performance" instance to a "Performance-guaranteed" instance, the SLB will have a short probability of business interruption (10 seconds-30 seconds). Advise to change it during the business downturn, or migrate business to other SLB Instances by using GSLB before changing.

-> **NOTE:** Currently, the alibaba cloud international account does not support creating a `Subscription` SLB instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the load balancer.
* `address` - The IP address of the load balancer.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the SLB load balancer.(until it reaches the initial `active` status). 
* `delete` - (Defaults to 9 mins) Used when terminating the SLB load balancer.

## Import

Load balancer can be imported using the id, e.g.

```
$ terraform import alicloud_slb_load_balancer.example lb-abc123456
```
