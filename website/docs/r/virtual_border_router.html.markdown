---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_virtual_border_router"
sidebar_current: "docs-alicloud-resource-virtual-border-router"
description: |-
  Provides a Alicloud Virtual Border Router (VBR) resource.
---

# alicloud\_virtual_border_router

Provides a Alicloud Virtual Border Router (VBR) resource.

-> **NOTE:**  Available in 1.83.0+

For information about VBR and how to use it, see [What is a Virtual Border Router](https://www.alibabacloud.com/help/doc-detail/44854.htm).

## Example Usage

Basic Usage

```
resource "alicloud_virtual_border_router" "default" {
  physical_connection_id = "pc-fakeid"
  vlan_id                = 1000
  local_gateway_ip       = "10.0.0.1"
  peer_gateway_ip        = "10.0.0.2"
  peering_subnet_mask    = "255.255.255.0"
}
```
## Argument Reference

The following arguments are supported:

* `physical_connection_id` - (Required, ForceNew) The ID of the physical connection.
* `vlan_id` - (Required) The VLAN ID of the VBR. Value range: 1 to 2999.
* `local_gateway_ip` - (Required) The Alibaba Cloud-side IP address used by the VBR.
* `peer_gateway_ip` - (Required) The customer-side IP address used by the VBR.
* `peering_subnet_mask` - (Required) The subnet mask for the Alibaba Cloud-side IP address and the customer-side IP address. The two IP addresses must be in the same subnet.
* `name` - (Optional) The name of the VBR. Defaults to null.
* `description` - (Optional) The description of the VBR. Defaults to null.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the VBR.
* `vlan_interface_id` - ID of the VRouter located in the local region.
* `route_table_id` - The ID of the route table of the VBR.

## Import

VBR can be imported using the id, e.g.

```
$ terraform import alicloud_virtual_border_router.example vbr-abc123456
```
