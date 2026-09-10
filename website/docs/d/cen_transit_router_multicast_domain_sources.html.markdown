---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_multicast_domain_sources"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-multicast-domain-sources"
description: |-
  Provides a list of Cen Transit Router Multicast Domain Source owned by an Alibaba Cloud account.
---

# alicloud_cen_transit_router_multicast_domain_sources

This data source provides Cen Transit Router Multicast Domain Source available to the user.[What is Transit Router Multicast Domain Source](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-registertransitroutermulticastgroupsources)

-> **NOTE:** Available since v1.195.0.

## Example Usage

```terraform
data "alicloud_cen_transit_router_multicast_domain_sources" "default" {
  transit_router_multicast_domain_id = "tr-mcast-domain-xxxxxx"
}

output "alicloud_cen_transit_router_multicast_domain_source_example_id" {
  value = data.alicloud_cen_transit_router_multicast_domain_sources.default.sources.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional)A list of the multicast domain IDs.
* `transit_router_multicast_domain_id` - (Required, ForceNew) The ID of the multicast domain to which the multicast source belongs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `sources` - A list of Transit Router Multicast Domain Source Entries. Each element contains the following attributes:
  * `transit_router_multicast_domain_id` - The ID of the multicast domain to which the multicast source belongs.
  * `group_ip_address` - The IP address of the multicast group to which the multicast source belongs. Value range: **224.0.0.1** to **239.255.255.254**. If the multicast group you specified does not exist in the current multicast domain, the system will automatically create a new multicast group for you.
  * `network_interface_id` - ENI ID of the multicast source.
  * `vpc_id` - The VPC to which the ENI of the multicast source belongs. This field is mandatory for VPCs that is owned by another accounts.
  * `status` - The status of the resource.
  * `id` - The id of the resource.
