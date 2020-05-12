---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_virtual_border_routers"
sidebar_current: "docs-alicloud-datasource-virtual-border-routers"
description: |-
    Provides a list of virtual border routers owned by an Alibaba Cloud account.
---

# alicloud\_virtual_border_routers

This data source provides a list of virtual border routers owned by an Alibaba Cloud account.

-> **NOTE:**  Available in 1.83.0+

## Example Usage

```
data "alicloud_virtual_border_routers" "default" {
}

output "first_virtual_border_router_id" {
  value = "${data.alicloud_virtual_border_routers.default.vbrs.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `status` - (Optional) Expected status. Valid values are `active` and `terminated`.
* `name_regex` - (Optional) A regex string used to filter by virtual border router name.
* `physical_connection_id` - (Optional) The ID of the physical connection.
* `physical_connection_owner_uid` - (Optional) The UID of the owner of the physical connection.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of virtual border router IDs.
* `names` - A list of virtual border router names.
* `vbrs` - A list of virtual border routers. Each element contains the following attributes:
  * `id` - Virtual border router ID.
  * `status` - Status of the VBR. Valid values are `active` and `terminated`.
  * `name` - The name of the VBR.
  * `description` - The description of the VBR.
  * `vlan_id` - The VLAN ID of the VBR.
  * `route_table_id` - The ID of the route table of the VBR.
  * `vlan_interface_id` - ID of the VRouter located in the local region.
  * `local_gateway_ip` - The Alibaba Cloud-side IP address used by the VBR.
  * `peer_gateway_ip` - The customer-side IP address used by the VBR.
  * `peering_subnet_mask` - The subnet mask for the Alibaba Cloud-side IP address and the customer-side IP address.
  * `physical_connection_id` - The ID of the physical connection.
  * `physical_connection_owner_uid` - The UID of the owner of the physical connection.
  * `access_point_id` - The ID of the physical connection access point.
  * `creation_time` - Virtual border router creation time.
  * `circuit_code` - The leased line code provided by the service provider for the physical connection.
