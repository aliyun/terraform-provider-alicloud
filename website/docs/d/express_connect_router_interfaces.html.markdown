---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_router_interfaces"
sidebar_current: "docs-alicloud-datasource-express-connect-router-interfaces"
description: |-
  Provides a list of Router Interface owned by an Alibaba Cloud account.
---

# alicloud\_express\_connect\_router\_interfaces

This data source provides Router Interface available to the user.[What is Router Interface](https://www.alibabacloud.com/help/doc-detail/52412.htm)

-> **NOTE:** Available in 1.199.0+

## Example Usage

```terraform
data "alicloud_express_connect_router_interfaces" "default" {
  ids        = ["${alicloud_router_interface.default.id}"]
  name_regex = alicloud_router_interface.default.name
}

output "alicloud_router_interface_example_id" {
  value = data.alicloud_express_connect_router_interfaces.default.interfaces.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of Router Interface IDs.
* `include_reservation_data` - (Optional, ForceNew) Does it contain renewal data. Valid values: `true`, `false`.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Router Interface IDs.
* `names` - A list of name of Router Interfaces.
* `interfaces` - A list of Router Interface Entries. Each element contains the following attributes:
  * `access_point_id` - The access point ID to which the VBR belongs.
  * `bandwidth` - The bandwidth of the resource.
  * `business_status` - The businessStatus of the resource. Valid Values: `Normal`, `FinancialLocked`, `SecurityLocked`.
  * `connected_time` - The connected time of the resource.
  * `create_time` - The creation time of the resource
  * `cross_border` - The cross border of the resource.
  * `description` - The description of the router interface.
  * `end_time` - The end time of the resource.
  * `has_reservation_data` - The has reservation data of the resource.
  * `hc_rate` - The hc rate of the resource.
  * `hc_threshold` -  The hc threshold of the resource.
  * `health_check_source_ip` - The health check source IP address, must be an unused IP within the local VPC.
  * `health_check_target_ip` - The IP address for health screening purposes.
  * `opposite_access_point_id` -The Access point ID to which the other end belongs.
  * `opposite_bandwidth` -  The opposite bandwidth of the router on the other side.
  * `opposite_interface_business_status` - The opposite interface business status of the router on the other side. Valid Values: `Normal`, `FinancialLocked`, `SecurityLocked`.
  * `opposite_interface_id` - The Interface ID of the router at the other end.
  * `opposite_interface_owner_id` - The AliCloud account ID of the owner of the router interface on the other end.
  * `opposite_interface_spec` - The opposite interface spec of the router on the other side. Valid Values: `Mini.2`, `Mini.5`, `Mini.5`, `Small.2`, `Small.5`, `Middle.1`, `Middle.2`, `Middle.5`, `Large.1`, `Large.2`, `Large.5`, `XLarge.1`, `Negative`.
  * `opposite_interface_status` - The opposite interface status of the router on the other side. Valid Values: `Idle`, `AcceptingConnecting`, `Connecting`, `Activating`, `Active`, `Modifying`, `Deactivating`, `Inactive`, `Deleting`.
  * `opposite_region_id` - The geographical ID of the location of the receiving end of the connection.
  * `opposite_router_id` - The id of the router at the other end.
  * `opposite_router_type` - The opposite router type of the router on the other side. Valid Values: `VRouter`, `VBR`.
  * `opposite_vpc_instance_id` - The opposite vpc instance id of the router on the other side.
  * `payment_type` - The payment methods for router interfaces. Valid Values: `PrePaid`, `PostPaid`.
  * `reservation_active_time` - The reservation active time of the resource.
  * `reservation_bandwidth` - The reservation bandwidth of the resource.
  * `reservation_internet_charge_type` - The reservation internet charge type of the resource.
  * `reservation_order_type` - The reservation order type of the resource.
  * `role` - The role of the router interface. Valid Values: `InitiatingSide`, `AcceptingSide`.
  * `router_id` - The router id associated with the router interface.
  * `router_interface_id` - The first ID of the resource.
  * `router_interface_name` - The name of the resource.
  * `router_type` - The type of router associated with the router interface. Valid Values: `VRouter`, `VBR`.
  * `spec` - The specification of the router interface. Valid Values: `Mini.2`, `Mini.5`, `Mini.5`, `Small.2`, `Small.5`, `Middle.1`, `Middle.2`, `Middle.5`, `Large.1`, `Large.2`, `Large.5`, `XLarge.1`, `Negative`.
  * `status` - The status of the resource. Valid Values: `Idle`, `AcceptingConnecting`, `Connecting`, `Activating`, `Active`, `Modifying`, `Deactivating`, `Inactive`, `Deleting`.
  * `vpc_instance_id` - The vpc instance id of the resource.
