---
subcategory: "Cen"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_multicast_domain_members"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-multicast-domain-members"
description: |-
  Provides a list of Cen Transit Router Multicast Domain Member owned by an Alibaba Cloud account.
---

# alicloud_cen_transit_router_multicast_domain_members

This data source provides Cen Transit Router Multicast Domain Member available to the user.[What is Transit Router Multicast Domain Member](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/api-doc-cbn-2017-09-12-api-doc-registertransitroutermulticastgroupmembers)

-> **NOTE:** Available in 1.195.0+

## Example Usage

```
data "alicloud_cen_transit_router_multicast_domain_members" "default" {
  transit_router_multicast_domain_id = "tr-mcast-domain-2d9oq455uk533zfr29"
}

output "alicloud_cen_transit_router_multicast_domain_member_example_id" {
  value = data.alicloud_cen_transit_router_multicast_domain_members.default.members.0.id
}
```

## Argument Reference

The following arguments are supported:
* `transit_router_multicast_domain_id` - (ForceNew,Required) The ID of the multicast domain to which the multicast member belongs.
* `network_interface_id` - (ForceNew,Optional) The ID of the ENI.
* `ids` - (ForceNew,Optional,Computed) A list of Transit Router Multicast Domain Member IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `members` - A list of Transit Router Multicast Domain Member Entries. Each element contains the following attributes:
    * `id` - The `key` of the resource supplied above.The value is formulated as `<transit_router_multicast_domain_id>:<group_ip_address>:<network_interface_id>`.
    * `group_ip_address` - The IP address of the multicast group to which the multicast member belongs. If the multicast group you specified does not exist in the current multicast domain, the system will automatically create a new multicast group for you in the current multicast domain.
    * `network_interface_id` - ENI ID of multicast member.
    * `status` - The status of the resource
    * `transit_router_multicast_domain_id` - The ID of the multicast domain to which the multicast member belongs.
    * `vpc_id` - The VPC to which the ENI of the multicast member belongs. This field is mandatory for VPCs owned by another accounts.
