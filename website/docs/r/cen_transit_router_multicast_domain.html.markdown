---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_multicast_domain"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Transit Router Multicast Domain resource.
---

# alicloud_cen_transit_router_multicast_domain

Provides a Cloud Enterprise Network (CEN) Transit Router Multicast Domain resource.



For information about Cloud Enterprise Network (CEN) Transit Router Multicast Domain and how to use it, see [What is Transit Router Multicast Domain](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.242.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hongkong"
}

resource "alicloud_cen_instance" "defaultPNnbnI" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "defaultSwwLm7" {
  support_multicast   = true
  cen_id              = alicloud_cen_instance.defaultPNnbnI.id
  transit_router_name = format("%s1", var.name)
}


resource "alicloud_cen_transit_router_multicast_domain" "default" {
  transit_router_multicast_domain_name        = var.name
  transit_router_multicast_domain_description = "description"
  transit_router_id                           = alicloud_cen_transit_router.defaultSwwLm7.transit_router_id
  options {
    igmpv2_support = "disable"
  }
}
```

## Argument Reference

The following arguments are supported:
* `options` - (Optional, List) Options See [`options`](#options) below.
* `tags` - (Optional, Map) The tag of the resource
* `transit_router_id` - (Optional, ForceNew) The ID of the forwarding router instance.
* `transit_router_multicast_domain_description` - (Optional) The description of the multicast domain.
* `transit_router_multicast_domain_name` - (Optional) The name of the multicast domain.

### `options`

The options supports the following:
* `igmpv2_support` - (Optional) Igmpv2Support

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `region_id` - The ID of the region to which the forwarding router instance belongs.

  You can call the [DescribeChildInstanceRegions](~~ 132080 ~~) operation to obtain the region ID.
* `status` - The status of the multicast domain.

  Only value: `Active`, which indicates that the multicast domain is currently available.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Transit Router Multicast Domain.
* `delete` - (Defaults to 5 mins) Used when delete the Transit Router Multicast Domain.
* `update` - (Defaults to 5 mins) Used when update the Transit Router Multicast Domain.

## Import

Cloud Enterprise Network (CEN) Transit Router Multicast Domain can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_multicast_domain.example <id>
```