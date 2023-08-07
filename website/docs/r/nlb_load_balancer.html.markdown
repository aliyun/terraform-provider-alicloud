---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_load_balancer"
sidebar_current: "docs-alicloud-resource-nlb-load-balancer"
description: |-
  Provides a Alicloud NLB Load Balancer resource.
---

# alicloud\_nlb\_load\_balancer

Provides a NLB Load Balancer resource.

For information about NLB Load Balancer and how to use it, see [What is Load Balancer](https://www.alibabacloud.com/help/en/server-load-balancer/latest/createloadbalancer).

-> **NOTE:** Available since v1.191.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_nlb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_vswitches" "default_1" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_nlb_zones.default.zones.0.id
}
data "alicloud_vswitches" "default_2" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_nlb_zones.default.zones.1.id
}
locals {
  zone_id_1    = data.alicloud_nlb_zones.default.zones.0.id
  vswitch_id_1 = data.alicloud_vswitches.default_1.ids[0]
  zone_id_2    = data.alicloud_nlb_zones.default.zones.1.id
  vswitch_id_2 = data.alicloud_vswitches.default_2.ids[0]
}

resource "alicloud_nlb_load_balancer" "default" {
  load_balancer_name = var.name
  resource_group_id  = data.alicloud_resource_manager_resource_groups.default.ids.0
  load_balancer_type = "Network"
  address_type       = "Internet"
  address_ip_version = "Ipv4"
  tags = {
    Created = "tfTestAcc0"
    For     = "Tftestacc 0"
  }
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_mappings {
    vswitch_id = local.vswitch_id_1
    zone_id    = local.zone_id_1
  }
  zone_mappings {
    vswitch_id = local.vswitch_id_2
    zone_id    = local.zone_id_2
  }
}
```

## Argument Reference

The following arguments are supported:

* `address_ip_version` - (Optional, Computed, ForceNew) The protocol version. Valid values:
  - ipv4 (default): IPv4
  - DualStack: dual stack
* `address_type` - (Required) The type of IPv4 address used by the NLB instance. Valid values:
  - Internet: The NLB instance uses a public IP address. The domain name of the NLB instance is resolved to the public IP address. Therefore, the NLB instance can be accessed over the Internet.
  - Intranet: The NLB instance uses a private IP address. The domain name of the NLB instance is resolved to the private IP address. Therefore, the NLB instance can be accessed over the virtual private cloud (VPC) where the NLB instance is deployed.
* `cross_zone_enabled` - (Optional, Computed) Specifies whether to enable cross-zone load balancing for the NLB instance.
* `load_balancer_name` - (Optional, Computed) The name of the NLB instance. The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). The name must start with a letter.
* `load_balancer_type` - (Optional, Computed, ForceNew) The type of the instance. Set the value to `Network`, which specifies an NLB instance.
* `resource_group_id` - (Optional, Computed, ForceNew) The ID of the resource group.
* `vpc_id` - (Required, ForceNew) The ID of the VPC where the NLB instance is deployed.
* `zone_mappings` - (Required) Available Area Configuration List. You must add at least two zones. You can add a maximum of 10 zones. See [`zone_mappings`](#zone_mappings) below.
* `bandwidth_package_id` - (Optional) The ID of the EIP bandwidth plan that is associated with the NLB instance if the NLB instance uses a public IP address.
* `deletion_protection_enabled` - (Optional, Computed, Available in 1.206.0+) Specifies whether to enable deletion protection. Default value: `false`. Valid values:
  - `true`: Enable deletion protection.
  - `false`: Disable deletion protection. You cannot set the `deletion_protection_reason`. If the `deletion_protection_reason` is set, the value is cleared.
* `deletion_protection_reason` - (Optional, Available in 1.206.0+) The reason why the deletion protection feature is enabled or disabled. The `deletion_protection_reason` takes effect only when `deletion_protection_enabled` is set to `true`.
* `modification_protection_status` - (Optional, Computed, Available in 1.206.0+) Specifies whether to enable the configuration read-only mode. Default value: `NonProtection`. Valid values:
  - `NonProtection`: Does not enable the configuration read-only mode. You cannot set the `modification_protection_reason`. If the `modification_protection_reason` is set, the value is cleared.
  - `ConsoleProtection`: Enables the configuration read-only mode. You can set the `modification_protection_reason`.
* `modification_protection_reason` - (Optional, Available in 1.206.0+) The reason why the configuration read-only mode is enabled. The `modification_protection_reason` takes effect only when `modification_protection_status` is set to `ConsoleProtection`.
* `tags` - (Optional) A mapping of tags to assign to the resource.

### `zone_mappings`

The zone_mappings supports the following: 

* `allocation_id` - (Optional, Computed) The ID of the EIP associated with the Internet-facing NLB instance.
* `public_ipv4_address` - (Computed) The public IPv4 address of the NLB instance.
* `vswitch_id` - (Required) The vSwitch in the zone. You can specify only one vSwitch (subnet) in each zone of an NLB instance.
* `zone_id` - (Required) The ID of the zone of the NLB instance.
* `eni_id` - (Computed) The ID of the elastic network interface (ENI).
* `ipv6_address` - (Computed) The IPv6 address of the NLB instance.
* `private_ipv4_address` - (Optional, Computed) The private IPv4 address of the NLB instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Load Balancer.
* `status` - The status of the NLB instance.
* `load_balancer_business_status` - The business status of the NLB instance.
* `ipv6_address_type` - The type of IPv6 address used by the NLB instance.
* `create_time` - The time when the resource was created. The time is displayed in UTC in `yyyy-MM-ddTHH:mm:ssZ` format.
* `dns_name` - The domain name of the NLB instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Load Balancer.
* `delete` - (Defaults to 1 mins) Used when delete the Load Balancer.
* `update` - (Defaults to 5 mins) Used when update the Load Balancer.

## Import

NLB Load Balancer can be imported using the id, e.g.

```shell
$ terraform import alicloud_nlb_load_balancer.example <id>
```
