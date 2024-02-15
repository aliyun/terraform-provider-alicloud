---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_load_balancer"
description: |-
  Provides a Alicloud NLB Load Balancer resource.
---

# alicloud_nlb_load_balancer

Provides a NLB Load Balancer resource. 

For information about NLB Load Balancer and how to use it, see [What is Load Balancer](https://www.alibabacloud.com/help/en/server-load-balancer/latest/api-nlb-2022-04-30-createloadbalancer).

-> **NOTE:** Available since v1.191.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "vpc" {
  vpc_name = var.name

  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "vsj" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = data.alicloud_zones.default.zones.0.id
  cidr_block   = "192.168.10.0/24"
  vswitch_name = var.name

}

resource "alicloud_vswitch" "vsk" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = data.alicloud_zones.default.zones.1.id
  cidr_block   = "192.168.20.0/24"
  vswitch_name = var.name

}

resource "alicloud_security_group" "defaultLkkjal" {
  vpc_id              = alicloud_vpc.vpc.id
  security_group_name = var.name

}

resource "alicloud_security_group" "defaultmlAdy7" {
  vpc_id              = alicloud_vpc.vpc.id
  security_group_name = var.name

}

resource "alicloud_security_group" "defaultCr6BU3" {
  vpc_id              = alicloud_vpc.vpc.id
  security_group_name = var.name

}

resource "alicloud_resource_manager_resource_group" "rg1" {
  display_name        = "rsg1"
  resource_group_name = var.name

}

resource "alicloud_resource_manager_resource_group" "rg2" {
  display_name        = "rsg2"
  resource_group_name = var.name

}

resource "alicloud_vswitch" "vsg" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = data.alicloud_zones.default.zones.2.id
  cidr_block   = "192.168.30.0/24"
  vswitch_name = var.name

}


resource "alicloud_nlb_load_balancer" "default" {
  load_balancer_name = var.name

  zone_mappings {
    vswitch_id = alicloud_vswitch.vsj.id
    zone_id    = alicloud_vswitch.vsj.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vsk.id
    zone_id    = alicloud_vswitch.vsk.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vsg.id
    zone_id    = alicloud_vswitch.vsg.zone_id
  }
  address_type       = "Intranet"
  address_ip_version = "Ipv4"
  load_balancer_type = "Network"
  vpc_id             = alicloud_vpc.vpc.id
  resource_group_id  = alicloud_resource_manager_resource_group.rg1.id
  security_group_ids = []
  deletion_protection_config {
  }
  modification_protection_config {
    status = "NonProtection"
  }
  cross_zone_enabled = true
}
```

## Argument Reference

The following arguments are supported:
* `address_ip_version` - (Optional, ForceNew) Protocol version. Value:
  - **ipv4**:IPv4 type.
  - **DualStack**: Double Stack type.
* `address_type` - (Required) The network address type of IPv4 for network load balancing. Value:
  - **Internet**: public network. Load balancer has a public network IP address, and the DNS domain name is resolved to a public network IP address, so it can be accessed in a public network environment.
  - **Intranet**: private network. The server load balancer only has a private IP address, and the DNS domain name is resolved to the private IP address, so it can only be accessed by the intranet environment of the VPC where the server load balancer is located.
* `bandwidth_package_id` - (Optional, ForceNew) The ID of the shared bandwidth package associated with the public network instance.
* `cross_zone_enabled` - (Optional) Whether cross-zone is enabled for a network-based load balancing instance. Value:
  - **true**: on.
  - **false**: closed.
* `deletion_protection_config` - (Optional, Available since v1.214.0) Delete protection. See [`deletion_protection_config`](#deletion_protection_config) below.
* `ipv6_address_type` - (Optional) The IPv6 address type of network load balancing. Value:
  - **Internet**: Server Load Balancer has a public IP address, and the DNS domain name is resolved to a public IP address, so it can be accessed in a public network environment.
  - **Intranet**: SLB only has the private IP address, and the DNS domain name is resolved to the private IP address, so it can only be accessed by the Intranet environment of the VPC where SLB is located.
* `load_balancer_name` - (Optional) The name of the network-based load balancing instance.  2 to 128 English or Chinese characters in length, which must start with a letter or Chinese, and can contain numbers, half-width periods (.), underscores (_), and dashes (-).
* `load_balancer_type` - (Optional, ForceNew, Computed) Load balancing type. Only value: **network**, which indicates network-based load balancing.
* `modification_protection_config` - (Optional, Available since v1.214.0) Modify protection. See [`modification_protection_config`](#modification_protection_config) below.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `security_group_ids` - (Optional, Available since v1.214.0) The security group to which the network-based SLB instance belongs.
* `tags` - (Optional, Map) List of labels.
* `vpc_id` - (Required, ForceNew) The ID of the network-based SLB instance.
* `zone_mappings` - (Required) The list of zones and vSwitch mappings. You must add at least two zones and a maximum of 10 zones. See [`zone_mappings`](#zone_mappings) below.

### `deletion_protection_config`

The deletion_protection_config supports the following:
* `enabled` - (Optional) Delete protection enable.
* `reason` - (Optional) Reason for opening.

### `modification_protection_config`

The modification_protection_config supports the following:
* `reason` - (Optional) Reason for opening.
* `status` - (Optional) ON.

### `zone_mappings`

The zone_mappings supports the following:
* `allocation_id` - (Optional) The ID of the elastic IP address.
* `private_ipv4_address` - (Optional) The private IPv4 address of a network-based server load balancer instance.
* `status` - (Optional) Zone Status.
* `vswitch_id` - (Required) The switch corresponding to the zone. Each zone uses one switch and one subnet by default.
* `zone_id` - (Required) The name of the zone. You can call the [DescribeZones](~~ 443890 ~~) operation to obtain the name of the zone.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Resource creation time, using Greenwich Mean Time, formating' yyyy-MM-ddTHH:mm:ssZ '.
* `deletion_protection_config` - Delete protection.
  * `enabled_time` - Opening time.
* `modification_protection_config` - Modify protection.
  * `enabled_time` - Opening time.
* `status` - The status of the resource.
* `zone_mappings` - The list of zones and vSwitch mappings. You must add at least two zones and a maximum of 10 zones.
  * `eni_id` - The ID of ENI.
  * `ipv6_address` - The IPv6 address of a network-based server load balancer instance.
  * `public_ipv4_address` - Public IPv4 address of a network-based server load balancer instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Load Balancer.
* `delete` - (Defaults to 5 mins) Used when delete the Load Balancer.
* `update` - (Defaults to 5 mins) Used when update the Load Balancer.

## Import

NLB Load Balancer can be imported using the id, e.g.

```shell
$ terraform import alicloud_nlb_load_balancer.example <id>
```