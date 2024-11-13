---
subcategory: "GWLB"
layout: "alicloud"
page_title: "Alicloud: alicloud_gwlb_load_balancer"
description: |-
  Provides a Alicloud GWLB Load Balancer resource.
---

# alicloud_gwlb_load_balancer

Provides a GWLB Load Balancer resource.



For information about GWLB Load Balancer and how to use it, see [What is Load Balancer](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.234.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-wulanchabu"
}

variable "region_id" {
  default = "cn-wulanchabu"
}

variable "zone_id2" {
  default = "cn-wulanchabu-c"
}

variable "zone_id1" {
  default = "cn-wulanchabu-b"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaulti9Axhl" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default9NaKmL" {
  vpc_id       = alicloud_vpc.defaulti9Axhl.id
  zone_id      = var.zone_id1
  cidr_block   = "10.0.0.0/24"
  vswitch_name = format("%s1", var.name)
}

resource "alicloud_vswitch" "defaultH4pKT4" {
  vpc_id       = alicloud_vpc.defaulti9Axhl.id
  zone_id      = var.zone_id2
  cidr_block   = "10.0.1.0/24"
  vswitch_name = format("%s2", var.name)
}


resource "alicloud_gwlb_load_balancer" "default" {
  vpc_id             = alicloud_vpc.defaulti9Axhl.id
  load_balancer_name = var.name
  zone_mappings {
    vswitch_id = alicloud_vswitch.default9NaKmL.id
    zone_id    = var.zone_id1
  }
  address_ip_version = "Ipv4"
}
```

## Argument Reference

The following arguments are supported:
* `address_ip_version` - (Optional, ForceNew, Computed) The protocol version. Value:
  - Ipv4: Ipv4 type
* `dry_run` - (Optional) Specifies whether to perform only a dry run, without performing the actual request. 
* `load_balancer_name` - (Optional) The name of the Gateway Load Balancer instance.

  It must be 2 to 128 English or Chinese characters in length. It must start with a letter or a Chinese character and can contain digits, half-width periods (.), underscores (_), and dashes (-).
* `resource_group_id` - (Optional, Computed) The ID of the resource group
* `tags` - (Optional, Map) The list of tags.
* `vpc_id` - (Required, ForceNew) The ID of the VPC which the Gateway Load Balancer instance belongs.
* `zone_mappings` - (Required, Set) The List of zones and vSwitches mapped. You must add at least one zone and a maximum of 20 zones. If the current region supports two or more zones, we recommend that you add two or more zones. See [`zone_mappings`](#zone_mappings) below.

### `zone_mappings`

The zone_mappings supports the following:
* `vswitch_id` - (Required) The ID of the vSwitch that corresponds to the zone. Each zone can use only one vSwitch and subnet.
* `zone_id` - (Required) The ID of the zone to which the Gateway Load Balancer instance belongs.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The resource creation time, in Greenwich Mean Time, in the format of **yyyy-MM-ddTHH:mm:ssZ**.
* `status` - The status of the Gateway load Balancer instance. Value:, indicating that the instance listener will no longer forward traffic.
* `zone_mappings` - The List of zones and vSwitches mapped. You must add at least one zone and a maximum of 20 zones. If the current region supports two or more zones, we recommend that you add two or more zones.
  * `load_balancer_addresses` - The addresses of the Gateway Load Balancer instance.
    * `eni_id` - The ID of the ENI.
    * `private_ipv4_address` - IPv4 private network address.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Load Balancer.
* `delete` - (Defaults to 5 mins) Used when delete the Load Balancer.
* `update` - (Defaults to 5 mins) Used when update the Load Balancer.

## Import

GWLB Load Balancer can be imported using the id, e.g.

```shell
$ terraform import alicloud_gwlb_load_balancer.example <id>
```