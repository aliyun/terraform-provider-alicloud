---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_loadbalancer_common_bandwidth_package_attachment"
description: |-
  Provides a Alicloud NLB Loadbalancer Common Bandwidth Package Attachment resource.
---

# alicloud_nlb_loadbalancer_common_bandwidth_package_attachment

Provides a NLB Loadbalancer Common Bandwidth Package Attachment resource. Bandwidth Package Operation.

For information about NLB Loadbalancer Common Bandwidth Package Attachment and how to use it, see [What is Loadbalancer Common Bandwidth Package Attachment](https://www.alibabacloud.com/help/en/server-load-balancer/latest/nlb-instances-change).

-> **NOTE:** Available since v1.209.0.

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
  cidr_block = "10.0.0.0/16"
  vpc_name   = var.name

}

resource "alicloud_vswitch" "vswtich" {
  vpc_id     = alicloud_vpc.vpc.id
  zone_id    = data.alicloud_zones.default.zones.0.id
  cidr_block = "10.0.1.0/24"
}

resource "alicloud_vswitch" "vswtich2" {
  vpc_id     = alicloud_vpc.vpc.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "10.0.2.0/24"
}

resource "alicloud_nlb_load_balancer" "nlb" {
  zone_mappings {
    vswitch_id = alicloud_vswitch.vswtich.id
    zone_id    = alicloud_vswitch.vswtich.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vswtich2.id
    zone_id    = alicloud_vswitch.vswtich2.zone_id
  }
  load_balancer_type = "Network"
  vpc_id             = alicloud_vpc.vpc.id
  address_type       = "Internet"
  address_ip_version = "Ipv4"
}

resource "alicloud_common_bandwidth_package" "cbwp" {
  description          = "nlb-tf-test"
  bandwidth            = "1000"
  internet_charge_type = "PayByBandwidth"
}


resource "alicloud_nlb_loadbalancer_common_bandwidth_package_attachment" "default" {
  load_balancer_id     = alicloud_nlb_load_balancer.nlb.id
  bandwidth_package_id = alicloud_common_bandwidth_package.cbwp.id
}
```

## Argument Reference

The following arguments are supported:
* `bandwidth_package_id` - (Required, ForceNew) The ID of the bound shared bandwidth package.
* `load_balancer_id` - (Required, ForceNew) The ID of the network-based server load balancer instance.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<load_balancer_id>:<bandwidth_package_id>`.
* `status` - Network-based load balancing instance status. Value:, indicating that the instance listener will no longer forward traffic.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Loadbalancer Common Bandwidth Package Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Loadbalancer Common Bandwidth Package Attachment.

## Import

NLB Loadbalancer Common Bandwidth Package Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_nlb_loadbalancer_common_bandwidth_package_attachment.example <load_balancer_id>:<bandwidth_package_id>
```