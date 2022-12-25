---
subcategory: "Cen"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_multicast_domain_peer_members"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-multicast-domain-peer-members"
description: |-
  Provides a list of Cen Transit Router Multicast Domain Peer Member owned by an Alibaba Cloud account.
---

# alicloud_cen_transit_router_multicast_domain_peer_members

This data source provides Cen Transit Router Multicast Domain Peer Member available to the user.[What is Transit Router Multicast Domain Peer Member](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/api-doc-cbn-2017-09-12-api-doc-registertransitroutermulticastgroupmembers)

-> **NOTE:** Available in 1.195.0+

## Example Usage

```
data "alicloud_cen_transit_router_multicast_domain_peer_members" "default" {
  transit_router_multicast_domain_id = "tr-mcast-domain-2d9oq455uk533zfrxx"
}

output "alicloud_cen_transit_router_multicast_domain_peer_member_example_id" {
  value = data.alicloud_cen_transit_router_multicast_domain_peer_members.default.members.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of Cen Transit Router Multicast Domain Peer Member IDs.
* `transit_router_multicast_domain_id` - (ForceNew, Required) The ID of the multicast domain to which the multicast member belongs.
* `peer_transit_router_multicast_domains` - (ForceNew, Optional) The IDs of the inter-region multicast domains.
* `resource_id` - (ForceNew, Optional) The ID of the resource associated with the multicast resource.
* `resource_type` - (ForceNew, Optional) The type of the multicast resource. Valid values:
  * VPC: queries multicast resources by VPC.
  * TR: queries multicast resources that are also deployed in a different region.
* `transit_router_attachment_id` - (ForceNew, Optional) The ID of the network instance connection.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `members` - A list of Transit Router Multicast Domain Peer Member Entries. Each element contains the following attributes:
    * `id` - The `key` of the resource supplied above.The value is formulated as `<transit_router_multicast_domain_id>:<group_ip_address>:<peer_transit_router_multicast_domain_id>`.
    * `group_ip_address` - The IP address of the multicast group to which the multicast member belongs. Value range: **224.0.0.1** to **239.255.255.254**.If the multicast group you specified does not exist in the current multicast domain, the system will automatically create a new multicast group for you in the current multicast domain.
    * `peer_transit_router_multicast_domain_id` - The multicast domain ID of the peer transit router.
    * `status` - The status of the resource
    * `transit_router_multicast_domain_id` - The ID of the multicast domain to which the multicast member belongs.
