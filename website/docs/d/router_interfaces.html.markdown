---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_router_interfaces"
sidebar_current: "docs-alicloud-datasource-router-interfaces"
description: |-
    Provides a list of router interfaces to the user.
---

# alicloud\_router\_interfaces

This data source provides information about [router interfaces](https://www.alibabacloud.com/help/doc-detail/52412.htm)
that connect VPCs together.

## Example Usage

```
data "alicloud_router_interfaces" "router_interfaces_ds" {
  name_regex = "^testenv"
  status     = "Active"
}

output "first_router_interface_id" {
  value = "${data.alicloud_router_interfaces.router_interfaces_ds.interfaces.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string used to filter by router interface name.
* `status` - (Optional) Expected status. Valid values are `Active`, `Inactive` and `Idle`.
* `specification` - (Optional) Specification of the link, such as `Small.1` (10Mb), `Middle.1` (100Mb), `Large.2` (2Gb), ...etc.
* `router_id` - (Optional) ID of the VRouter located in the local region.
* `router_type` - (Optional) Router type in the local region. Valid values are `VRouter` and `VBR` (physical connection).
* `role` - (Optional) Role of the router interface. Valid values are `InitiatingSide` (connection initiator) and 
  `AcceptingSide` (connection receiver). The value of this parameter must be `InitiatingSide` if the `router_type` is set to `VBR`.
* `opposite_interface_id` - (Optional) ID of the peer router interface.
* `opposite_interface_owner_id` - (Optional) Account ID of the owner of the peer router interface.
* `ids` - (Optional, Available in 1.44.0+) A list of router interface IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of router interface IDs.
* `names` - A list of router interface names.
* `interfaces` - A list of router interfaces. Each element contains the following attributes:
  * `id` - Router interface ID.
  * `status` - Router interface status. Possible values: `Active`, `Inactive` and `Idle`.
  * `name` - Router interface name.
  * `description` - Router interface description.
  * `role` - Router interface role. Possible values: `InitiatingSide` and `AcceptingSide`.
  * `specification` - Router interface specification. Possible values: `Small.1`, `Middle.1`, `Large.2`, ...etc.
  * `router_id` - ID of the VRouter located in the local region.
  * `router_type` - Router type in the local region. Possible values: `VRouter` and `VBR`.
  * `vpc_id` - ID of the VPC that owns the router in the local region.
  * `access_point_id` - ID of the access point used by the VBR.
  * `creation_time` - Router interface creation time.
  * `opposite_region_id` - Peer router region ID.
  * `opposite_interface_id` - Peer router interface ID.
  * `opposite_router_id` - Peer router ID.
  * `opposite_router_type` - Router type in the peer region. Possible values: `VRouter` and `VBR`.
  * `opposite_interface_owner_id` - Account ID of the owner of the peer router interface.
  * `health_check_source_ip` - Source IP address used to perform health check on the physical connection.
  * `health_check_target_ip` - Destination IP address used to perform health check on the physical connection.
