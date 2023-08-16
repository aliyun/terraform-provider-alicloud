---
subcategory: "SLB"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_load_balancer"
description: |-
  Provides a Alicloud SLB Load Balancer resource.
---

# alicloud_slb_load_balancer

Provides a SLB Load Balancer resource. The load balancing service that distributes traffic to multiple cloud servers can expand the external service capability of the application system through traffic distribution and improve the availability of the application system by eliminating single points of failure.

For information about SLB Load Balancer and how to use it, see [What is Load Balancer](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.210.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "vpc" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "vsw1" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = data.alicloud_zones.default.zones.0.id
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "${var.name}1"
}

resource "alicloud_resource_manager_resource_group" "rg1" {
  display_name        = "slb01"
  resource_group_name = "${var.name}2"
}


resource "alicloud_slb_load_balancer" "default" {
  status                         = "active"
  address_ip_version             = "ipv4"
  address                        = "10.0.10.1"
  instance_charge_type           = "PayByCLCU"
  vswitch_id                     = "vsw-bp1fpj92chcwmdla73oxg"
  slave_zone_id                  = "cn-hangzhou-h"
  modification_protection_status = "ConsoleProtection"
  load_balancer_name             = var.name
  delete_protection              = "off"
  vpc_id                         = "vpc-bp18uccoyc62e4gk6033e"
  payment_type                   = "PayAsYouGo"
  modification_protection_reason = "test"
  address_type                   = "intranet"
  master_zone_id                 = "cn-hangzhou-j"
  tags {
    tag_key   = "k1"
    tag_value = "v1"
  }
  resource_group_id = alicloud_resource_manager_resource_group.rg1.id
}
```

### Deleting `alicloud_slb_load_balancer` or removing it from your configuration

The `alicloud_slb_load_balancer` resource allows you to manage  `payment_type = "Subscription"`  instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Instance.
You can resume managing the subscription instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:
* `address` - (Optional, ForceNew, Available since v1.123.1) The service address of the load balancing instance.
* `address_ip_version` - (Optional, ForceNew, Available since v1.123.1) IP version.
* `address_type` - (Optional, ForceNew, Available since v1.123.1) The address type of the load balancing instance.
* `auto_pay` - (Optional) Whether the bill for a prepaid public network instance is automatically paid. Value:
  - **true:** automatic payment. Immediately after calling the API, an SLB instance is generated.
  - **false** (default): The SLB order was created after the API was called, but no payment was made. You can see unpaid orders in the console. The SLB instance will not be created because the order is not paid.
-> **NOTE:**  This parameter is only applicable to China site and is only valid for subscription instances, that is, it is valid when the **PayType** parameter value is **PrePay.


.
* `bandwidth` - (Optional, Available since v1.123.1) Peak bandwidth of a public networked instance billed by bandwidth.
* `delete_protection` - (Optional, Available since v1.123.1) Whether to enable instance deletion protection.
* `duration` - (Optional) The purchase duration of the prepaid public network instance. Valid values:

* If **PricingCycle** is **month**, the value is **1 to 9 * *.

* If **PricingCycle** is **year**, the value is **1 to 3 * *.
-> **NOTE:**  This parameter is only applicable to China site and is only valid for subscription instances, that is, it is valid when the **PayType** parameter value is **PrePay.
* `instance_charge_type` - (Optional, Available since v1.123.1) Instance billing method. Value:
  - **PayBySpec (default)**: billed by specification.
  - **PayByCLCU**: Billed by usage.
-> **NOTE:**  This parameter takes effect when the value of **PayType** (instance payment mode) is **PayOnDemand** (pay-as-you-go).
.
* `internet_charge_type` - (Optional, Available since v1.123.1) Public network type instance payment method.
* `load_balancer_name` - (Optional, Available since v1.123.1) Name of the load balancing instance.
* `load_balancer_spec` - (Optional, Available since v1.123.1) Specifications for instances.
* `master_zone_id` - (Optional, ForceNew, Available since v1.123.1) ID of the primary available area of the load balancing instance.
* `modification_protection_reason` - (Optional, Available since v1.123.1) Sets the reason for modifying the protected state.
* `modification_protection_status` - (Optional, Available since v1.123.1) Load balancing modifies the protection state.
* `payment_type` - (Optional, Computed, Available since v1.123.1) Load balancing instance payment type.
* `pricing_cycle` - (Optional) Billing cycle.
* `resource_group_id` - (Optional, Computed, Available since v1.123.1) Resource group id.
* `slave_zone_id` - (Optional, ForceNew, Available since v1.123.1) ID of the ready-to-use zone of the load balancing instance.
* `status` - (Optional, Computed, Available since v1.123.1) The status of SLB.
* `tags` - (Optional, Map, Available since v1.123.1) The tags of VSwitch.
* `vswitch_id` - (Optional, ForceNew, Available since v1.123.1) The ID of the switch.
* `vpc_id` - (Optional, ForceNew) VPC ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the load balancing instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Load Balancer.
* `delete` - (Defaults to 5 mins) Used when delete the Load Balancer.
* `update` - (Defaults to 5 mins) Used when update the Load Balancer.

## Import

SLB Load Balancer can be imported using the id, e.g.

```shell
$ terraform import alicloud_slb_load_balancer.example <id>
```