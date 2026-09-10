---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_virtual_physical_connections"
sidebar_current: "docs-alicloud-datasource-express-connect-virtual-physical-connections"
description: |-
  Provides a list of Express Connect Virtual Physical Connection owned by an Alibaba Cloud account.
---

# alicloud_express_connect_virtual_physical_connections

This data source provides Express Connect Virtual Physical Connection available to the user.

-> **NOTE:** Available in 1.196.0+

## Example Usage

```
data "alicloud_express_connect_virtual_physical_connections" "default" {
  ids                           = ["${alicloud_express_connect_virtual_physical_connection.default.id}"]
  name_regex                    = alicloud_express_connect_virtual_physical_connection.default.name
  parent_physical_connection_id = alicloud_express_connect_virtual_physical_connection.default.parent_physical_connection_id
  vlan_id                       = 789
  vpconn_ali_uid                = 1234567890
}

output "alicloud_express_connect_virtual_physical_connection_example_id" {
  value = data.alicloud_express_connect_virtual_physical_connections.default.connections.0.id
}
```

## Argument Reference

The following arguments are supported:
* `business_status` - (ForceNew,Optional) The commercial status of the physical line. Value:
  - **Normal**: activated.
  - **Financialized**: Arrears locked.
  - **SecurityLocked**: locked for security reasons.
* `parent_physical_connection_id` - (ForceNew,Optional) The ID of the instance of the physical connection.
* `virtual_physical_connection_ids` - (ForceNew,Optional) The ID of the hosted connection. You can specify multiple hosted connection IDs.
* `virtual_physical_connection_status` - (ForceNew,Optional) The business status of the shared line. Value:
  - **Confirmed**: The shared line has been Confirmed to receive.
  - **UnConfirmed**: The shared line has not been confirmed to be received.
  - **Deleted**: The shared line has been Deleted.
* `vlan_ids` - (ForceNew,Optional) The VLAN ID of the hosted connection. You can specify multiple VLAN IDs.
* `vpconn_ali_uid` - (ForceNew,Optional) The ID of the Alibaba Cloud account (primary account) of the owner of the shared line.
* `ids` - (Optional, ForceNew, Computed) A list of Virtual Physical Connection IDs.
* `virtual_physical_connection_names` - (Optional, ForceNew) The name of the Virtual Physical Connection. You can specify at most 10 names.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Virtual Physical Connection IDs.
* `names` - A list of name of Virtual Physical Connections.
* `connections` - A list of Virtual Physical Connection Entries. Each element contains the following attributes:
  * `access_point_id` - The ID of the access point of the physical connection.
  * `ad_location` - The physical location where the physical connection access device is located.
  * `bandwidth` - The bandwidth of the physical connection. Unit: Mbps.
  * `business_status` - The commercial status of the physical line. Value:-**Normal**: activated.-**Financialized**: Arrears locked.-**SecurityLocked**: locked for security reasons.
  * `circuit_code` - The circuit code provided by the operator for the physical connection.
  * `create_time` - The creation time of the resource
  * `description` - The description of the physical connection.
  * `enabled_time` - The opening time of the physical connection.
  * `end_time` - The expiration time of the shared line.Time is expressed according to ISO8601 standard and UTC time is used. The format is: YYYY-MM-DDThh:mm:ssZ.
  * `expect_spec` - The estimated bandwidth value of the shared line. The expected bandwidth value will not take effect until the payment is completed.Unit: **M** indicates Mbps,**G** indicates Gbps.
  * `line_operator` - Operators that provide access to physical lines. Value:-**CT**: China Telecom.-**CU**: China Unicom.-**CM**: China Mobile.-**CO**: China Other.-**Equinix**:Equinix.-**Other**: Other abroad.
  * `loa_status` - The state of LOA. Value:-**Applying**:LOA application.-**Accept**:LOA application passed.-**Available**:LOA is Available.-**Rejected**:LOA application Rejected.-**Completing**: The dedicated line is under construction.-**Complete**: The construction of the dedicated line is completed.-**Deleted**:LOA has been Deleted.
  * `order_mode` - The payment method of shared dedicated line. Value:-**PayByPhysicalConnectionOwner**: indicates that the owner of the physical line associated with the shared line pays.-**PayByVirtualPhysicalConnectionOwner**: indicates that the owner of the shared line pays.
  * `parent_physical_connection_ali_uid` - The ID of the Alibaba Cloud account (primary account) to which the physical connection belongs.
  * `parent_physical_connection_id` - The ID of the instance of the physical connection.
  * `peer_location` - The geographic location of the local data center.
  * `port_number` - The port number of the physical connection device.
  * `port_type` - Physical connection port type. Value:-**100Base-T**: 100 megabytes port.-**1000Base-T**: Gigabit port.-**1000Base-LX**: Gigabit single mode optical port (10km).-**10GBase-T**: 10 Gigabit port.-**10GBase-LR**: 10 Gigabit single mode optical port (10km).-**40GBase-LR**: 40 megabytes single-mode optical port.-**100GBase-LR**: 100,000 megabytes single-mode optical port.
  * `redundant_physical_connection_id` - The ID of the redundant physical connection.
  * `resource_group_id` - The resource group id
  * `spec` - The bandwidth value of the shared line.Unit: **M** indicates Mbps,**G** indicates Gbps.
  * `status` - The status of the resource
  * `virtual_physical_connection_id` - The ID of the hosted connection
  * `virtual_physical_connection_name` - The name of the physical connection.
  * `virtual_physical_connection_status` - The business status of the shared line. Value:-**Confirmed**: The shared line has been Confirmed to receive.-**UnConfirmed**: The shared line has not been confirmed to be received.-**Deleted**: The shared line has been Deleted.
  * `vlan_id` - The VLAN ID of the shared leased line.
  * `vpconn_ali_uid` - The ID of the Alibaba Cloud account (primary account) of the owner of the shared line.
  * `id` - The ID of the Virtual Physical Connection.
