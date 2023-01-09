---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_multicast_domain_source"
sidebar_current: "docs-alicloud-resource-cen-transit-router-multicast-domain-source"
description: |-
  Provides a Alicloud Cen Transit Router Multicast Domain Source resource.
---

# alicloud_cen_transit_router_multicast_domain_source

Provides a Cen Transit Router Multicast Domain Source resource.

For information about Cen Transit Router Multicast Domain Source and how to use it, see [What is Transit Router Multicast Domain Source](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/api-doc-cbn-2017-09-12-api-doc-registertransitroutermulticastgroupsources).

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cen_transit_router_multicast_domain_source" "default" {
  vpc_id                             = "vpc-xxxx"
  transit_router_multicast_domain_id = "tr-mcast-domain-xxxx"
  network_interface_id               = "eni-xxxx"
  group_ip_address                   = "239.1.1.1"
}
```

## Argument Reference

The following arguments are supported:
* `transit_router_multicast_domain_id` - (Required,ForceNew) The ID of the multicast domain to which the multicast source belongs.
* `group_ip_address` - (Required,ForceNew) The IP address of the multicast group to which the multicast source belongs. Value range: **224.0.0.1** to **239.255.255.254**. If the multicast group you specified does not exist in the current multicast domain, the system will automatically create a new multicast group for you.
* `network_interface_id` - (Required,ForceNew) ENI ID of the multicast source.
* `vpc_id` - (ForceNew,Optional) The VPC to which the ENI of the multicast source belongs. This field is mandatory for VPCs that is owned by another accounts.


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.The value is formulated as `&lt;transit_router_multicast_domain_id&gt;:&lt;group_ip_address&gt;:&lt;network_interface_id&gt;`.
* `status` - The status of the resource

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Transit Router Multicast Domain Source.
* `delete` - (Defaults to 5 mins) Used when delete the Transit Router Multicast Domain Source.

## Import

Cen Transit Router Multicast Domain Source can be imported using the id, e.g.

```shell
$terraform import alicloud_cen_transit_router_multicast_domain_source.example <transit_router_multicast_domain_id>:<group_ip_address>:<network_interface_id>
```