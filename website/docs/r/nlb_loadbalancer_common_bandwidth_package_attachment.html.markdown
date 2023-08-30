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
  default = "tf-example"
}
data "alicloud_resource_manager_resource_groups" "default" {}
data "alicloud_nlb_zones" "default" {}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_nlb_zones.default.zones.0.id
}
resource "alicloud_vswitch" "default1" {
  vswitch_name = var.name
  cidr_block   = "10.4.1.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_nlb_zones.default.zones.1.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_nlb_load_balancer" "default" {
  load_balancer_name = var.name
  resource_group_id  = data.alicloud_resource_manager_resource_groups.default.ids.0
  load_balancer_type = "Network"
  address_type       = "Internet"
  address_ip_version = "Ipv4"
  vpc_id             = alicloud_vpc.default.id
  tags = {
    Created = "TF",
    For     = "example",
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default.id
    zone_id    = data.alicloud_nlb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default1.id
    zone_id    = data.alicloud_nlb_zones.default.zones.1.id
  }
}

resource "alicloud_common_bandwidth_package" "default" {
  bandwidth              = 2
  internet_charge_type   = "PayByBandwidth"
  bandwidth_package_name = var.name
  description            = var.name
}

resource "alicloud_nlb_loadbalancer_common_bandwidth_package_attachment" "default" {
  bandwidth_package_id = alicloud_common_bandwidth_package.default.id
  load_balancer_id     = alicloud_nlb_load_balancer.default.id
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