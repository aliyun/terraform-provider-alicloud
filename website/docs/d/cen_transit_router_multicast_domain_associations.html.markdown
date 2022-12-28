---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_multicast_domain_associations"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-multicast-domain-associations"
description: |-
  Provides a list of Cen Transit Router Multicast Domain Associations to the user.
---

# alicloud\_cen\_transit\_router\_multicast\_domain\_associations

This data source provides the Cen Transit Router Multicast Domain Associations of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cen_transit_router_multicast_domain_associations" "ids" {
  ids                                = ["example_id"]
  transit_router_multicast_domain_id = "your_transit_router_multicast_domain_id"
}

output "cen_transit_router_multicast_domain_id_0" {
  value = data.alicloud_cen_transit_router_multicast_domain_associations.ids.associations.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Transit Router Multicast Domain Association IDs.
* `transit_router_multicast_domain_id` - (Required, ForceNew) The ID of the multicast domain.
* `transit_router_attachment_id` - (Optional, ForceNew) The ID of the network instance connection.
* `vswitch_id` - (Optional, ForceNew) The ID of the vSwitch.
* `resource_id` - (Optional, ForceNew) The ID of the resource associated with the multicast domain.
* `resource_type` - (Optional, ForceNew) The type of resource associated with the multicast domain. Valid Value: `VPC`.
* `status` - (Optional, ForceNew) The status of the associated resource. Valid Value: `Associated`, `Associating`, `Dissociating`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `associations` - A list of Cen Transit Router Multicast Domain Associations. Each element contains the following attributes:
	* `id` - The ID of the Transit Router Multicast Domain Association. It formats as `<transit_router_multicast_domain_id>:<transit_router_attachment_id>:<vswitch_id>`.
	* `transit_router_multicast_domain_id` - The ID of the multicast domain.
	* `transit_router_attachment_id` - The ID of the network instance connection.
	* `vswitch_id` - The ID of the vSwitch.
	* `resource_id` - The ID of the resource associated with the multicast domain.
	* `resource_owner_id` - The ID of the Alibaba Cloud account to which the resource associated with the multicast domain belongs.
	* `resource_type` - The type of resource associated with the multicast domain.
	* `status` - The status of the associated resource.
	