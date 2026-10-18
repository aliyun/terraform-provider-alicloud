---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_vcc_route_entry"
description: |-
  Provides a Alicloud Eflo Vcc Route Entry resource.
---

# alicloud_eflo_vcc_route_entry

Provides a Eflo Vcc Route Entry resource.

VccRouteEntry is immutable: there is no UpdateVccRouteEntry API, so changing any
settable field (`vcc_id`, `destination_cidr_block`, `region_id`) forces replacement.

For information about Eflo VccRouteEntry and how to use it, see
[What is VccRouteEntry](https://next.api.alibabacloud.com/document/eflo/2022-05-30/CreateVccRouteEntry).

-> **NOTE:** Available since v1.258.0.

## Example Usage

```terraform
resource "alicloud_eflo_vcc_route_entry" "default" {
  vcc_id                 = "vcc-cn-zvp2w222001"
  destination_cidr_block = "192.168.10.0/24"
}
```

## Argument Reference

The following arguments are supported:

* `vcc_id` - (Required, ForceNew) The ID of the Eflo VCC (灵骏连接) instance.
* `destination_cidr_block` - (Required, ForceNew) The destination CIDR block of the route entry.
* `region_id` - (Optional, Computed, ForceNew) The region ID. Defaults to the provider region.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, formatted as `<vcc_id>:<vcc_route_entry_id>`.
* `vcc_route_entry_id` - The ID of the route entry.
* `status` - The status of the route entry.
* `route_type` - The route type, e.g. `System`, `BGP`.
* `next_hop_id` - The next hop instance ID.
* `next_hop_type` - The next hop type, e.g. `VCC`.
* `message` - The message returned with the route entry.
* `gmt_modified` - The last modification time.
* `tenant_id` - The tenant ID.
* `resource_group_id` - The resource group ID.

## Import

Eflo Vcc Route Entry can be imported using the composite id `<vcc_id>:<vcc_route_entry_id>`, e.g.

```shell
$ terraform import alicloud_eflo_vcc_route_entry.example vcc-cn-zvp2w222001:vcc-rte-maysfadg
```
