---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_multicast_domain_peer_member"
sidebar_current: "docs-alicloud-resource-cen-transit-router-multicast-domain-peer-member"
description: |-
  Provides a Alicloud Cen Transit Router Multicast Domain Peer Member resource.
---

# alicloud_cen_transit_router_multicast_domain_peer_member

Provides a Cen Transit Router Multicast Domain Peer Member resource.

For information about Cen Transit Router Multicast Domain Peer Member and how to use it, see [What is Transit Router Multicast Domain Peer Member](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/api-doc-cbn-2017-09-12-api-doc-deregistertransitroutermulticastgroupmembers).

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cen_transit_router_multicast_domain_peer_member" "default" {
  peer_transit_router_multicast_domain_id = "tr-mcast-domain-itc67v79yk4xrkr9f3"
  transit_router_multicast_domain_id      = "tr-mcast-domain-2d9oq455uk533zfr29"
  group_ip_address                        = "239.1.1.1"
}
```

## Argument Reference

The following arguments are supported:
* `group_ip_address` - (Required, ForceNew) The IP address of the multicast group to which the multicast member belongs. Value range: **224.0.0.1** to **239.255.255.254**.If the multicast group you specified does not exist in the current multicast domain, the system will automatically create a new multicast group for you in the current multicast domain.
* `transit_router_multicast_domain_id` - (Required, ForceNew) The ID of the multicast domain to which the multicast member belongs.
* `peer_transit_router_multicast_domain_id` - (Required, ForceNew) The IDs of the inter-region multicast domains.
* `dry_run` - (Optional) Specifies whether only to precheck the request.

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.The value is formulated as `<transit_router_multicast_domain_id>:<group_ip_address>:<peer_transit_router_multicast_domain_id>`.
* `peer_transit_router_multicast_domain_id` - The multicast domain ID of the peer transit router.
* `status` - The status of the multicast resource. Valid values:
  - Registering: being created
  - Registered: available
  - Deregistering: being deleted

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the Transit Router Multicast Domain Peer Member.
* `delete` - (Defaults to 10 mins) Used when delete the Transit Router Multicast Domain Peer Member.

## Import

Cen Transit Router Multicast Domain Peer Member can be imported using the id, e.g.

```shell
$terraform import alicloud_cen_transit_router_multicast_domain_peer_member.example <transit_router_multicast_domain_id>:<group_ip_address>:<peer_transit_router_multicast_domain_id>
```