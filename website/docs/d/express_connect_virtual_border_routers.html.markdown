---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_virtual_border_routers"
sidebar_current: "docs-alicloud-datasource-express-connect-virtual-border-routers"
description: |-
  Provides a list of Express Connect Virtual Border Routers to the user.
---

# alicloud\_express\_connect\_virtual\_border\_routers

This data source provides the Express Connect Virtual Border Routers of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_express_connect_virtual_border_routers" "ids" {}
output "express_connect_virtual_border_router_id_1" {
  value = data.alicloud_express_connect_virtual_border_routers.ids.routers.0.id
}

data "alicloud_express_connect_virtual_border_routers" "nameRegex" {
  name_regex = "^my-VirtualBorderRouter"
}
output "express_connect_virtual_border_router_id_2" {
  value = data.alicloud_express_connect_virtual_border_routers.nameRegex.routers.0.id
}

data "alicloud_express_connect_virtual_border_routers" "filter" {
  filter {
    key    = "PhysicalConnectionId"
    values = ["pc-xxxx1"]
  }
  filter {
    key    = "VbrId"
    values = ["vbr-xxxx1", "vbr-xxxx2"]
  }
}
output "express_connect_virtual_border_router_id_3" {
  value = data.alicloud_express_connect_virtual_border_routers.filter.routers.0.id
}

```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) Custom filter block as described below.
* `ids` - (Optional, ForceNew, Computed)  A list of Virtual Border Router IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Virtual Border Router name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The instance state with. Valid values: `active`, `deleting`, `recovering`, `terminated`, `terminating`, `unconfirmed`.

### Block filter

More complex filters can be expressed using one or more `filter` sub-blocks,
which take the following arguments:
* `key` - (Required) The key of the field to filter by, as defined by
  [Alibaba Cloud API](https://www.alibabacloud.com/help/en/doc-detail/124791.htm).
* `values` - (Required) Set of values that are accepted for the given field.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Virtual Border Router names.
* `routers` - A list of Express Connect Virtual Border Routers. Each element contains the following attributes:
	* `access_point_id` - The physical leased line access point ID.
	* `activation_time` - The first activation time of VBR.
	* `circuit_code` - Operators for physical connection circuit provided coding.
	* `cloud_box_instance_id` - Box Instance Id.
	* `create_time` - The representative of the creation time resources attribute field.
	* `description` - The description of VBR. Length is from 2 to 256 characters, must start with a letter or the Chinese at the beginning, but not at the http:// Or https:// at the beginning.
	* `detect_multiplier` - Detection time multiplier that recipient allows the sender to send a message of the maximum allowable connections for the number of packets, used to detect whether the link normal. Value: 3~10.
	* `ecc_id` - High Speed Migration Service Instance Id.
	* `enable_ipv6` - Whether to Enable IPv6.
	* `id` - The ID of the Virtual Border Router.
	* `local_gateway_ip` - Alibaba Cloud-Connected IPv4 address.
	* `local_ipv6_gateway_ip` - Alibaba Cloud-Connected IPv6 Address.
	* `min_rx_interval` - Configure BFD packet reception interval of values include: 200~1000, unit: ms.
	* `min_tx_interval` - Configure BFD packet transmission interval maximum value: 200~1000, unit: ms.
	* `payment_vbr_expire_time` - The Billing of the Extended Time.
	* `peer_gateway_ip` - The Client-Side Interconnection IPv4 Address.
	* `peer_ipv6_gateway_ip` - The Client-Side Interconnection IPv6 Address.
	* `peering_ipv6_subnet_mask` - Alibaba Cloud-Connected IPv6 with Client-Side Interconnection IPv6 of Subnet Mask.
	* `peering_subnet_mask` - Alibaba Cloud-Connected IPv4 and Client-Side Interconnection IPv4 of Subnet Mask.
	* `physical_connection_business_status` - Physical Private Line Service Status Value Normal: Normal, financiallocked: If You Lock.
	* `physical_connection_id` - The ID of the Physical Connection to Which the ID.
	* `physical_connection_owner_uid` - Physical Private Line Where the Account ID.
	* `physical_connection_status` - Physical Private Line State.
	* `recovery_time` - The Last from a Terminated State to the Active State of the Time.
	* `route_table_id` - Route Table ID.
	* `status` - The VBR state.
	* `termination_time` - The Most Recent Was Aborted by the Time.
	* `type` - VBR Type.
	* `virtual_border_router_id` - The VBR ID.
	* `virtual_border_router_name` - The name of VBR. Length is from 2 to 128 characters, must start with a letter or the Chinese at the beginning can contain numbers, the underscore character (_) and dash (-). But do not start with http:// or https:// at the beginning.
	* `vlan_id` - The VLAN ID of the VBR. Value range: 0~2999.
	* `vlan_interface_id` - The ID of the Router Interface.
