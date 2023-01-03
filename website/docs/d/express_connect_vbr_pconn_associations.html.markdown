---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_vbr_pconn_associations"
sidebar_current: "docs-alicloud-datasource-express-connect-vbr-pconn-associations"
description: |-
  Provides a list of Express Connect Vbr Pconn Association owned by an Alibaba Cloud account.
---

# alicloud_express_connect_vbr_pconn_associations

This data source provides Express Connect Vbr Pconn Association available to the user.

-> **NOTE:** Available in 1.196.0+

## Example Usage

```terraform
data "alicloud_express_connect_vbr_pconn_associations" "default" {
  ids    = ["example_id"]
  vbr_id = alicloud_express_connect_vbr_pconn_association.default.vbr_id
}

output "alicloud_express_connect_vbr_pconn_association_example_id" {
  value = data.alicloud_express_connect_vbr_pconn_associations.default.associations.0.id
}
```

## Argument Reference

The following arguments are supported:
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `ids` - (Optional, ForceNew, Computed) A list of Vbr Pconn Association IDs.
* `vbr_id` - (Optional, ForceNew) The ID of the VBR instance.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `associations` - A list of Vbr Pconn Association Entries. Each element contains the following attributes:
  * `circuit_code` - The circuit code provided by the operator for the physical connection.
  * `enable_ipv6` - Whether IPv6 is enabled. 
  * `local_gateway_ip` - The Alibaba cloud IP address of the VBR instance.
  * `local_ipv6_gateway_ip` - The IPv6 address on the Alibaba Cloud side of the VBR instance.
  * `peer_gateway_ip` - The client IP address of the VBR instance.
  * `peer_ipv6_gateway_ip` - The IPv6 address of the client side of the VBR instance.
  * `peering_ipv6_subnet_mask` - The subnet mask of the Alibaba Cloud side and the client side of the VBR instance.Two IPv6 addresses must be in the same subnet.
  * `peering_subnet_mask` - The subnet mask of the Alibaba Cloud side and the client side of the VBR instance.The two IP addresses must be in the same subnet.
  * `physical_connection_id` - The ID of the leased line instance.
  * `status` - The status of the resource
  * `vbr_id` - The ID of the VBR instance.
  * `vlan_id` - VLAN ID of the VBR.
  * `id` - The ID of the Vbr Pconn Association.
