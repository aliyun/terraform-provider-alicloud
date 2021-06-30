---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_peer_attachments"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-peer-attachments"
description: |-
Provides a list of CEN Transit Router peer attachments owned by an Alibaba Cloud account.
---

# alicloud\_cen\_transit\_router\_peer\_attachments

This data source provides CEN Transit Router peer attachments available to the user.

-> **NOTE:** Available in 1.125.0+

## Example Usage

```
data "alicloud_cen_transit_router_peer_attachments" "default" {
  cen_id    = "cen-id1"
}

output "first_transit_router_peer_attachments_transit_router_attachment_resource_type" {
  value = "${data.alicloud_cen_transit_router_peer_attachments.default.transit_router_attachments.0.resource_type}"
}
```

## Argument Reference

The following arguments are supported:

* `cen_id` - (Required) ID of the CEN instance.
* `route_table_id` - (Optional) ID of the route table of the VPC or VBR.
* `resource_type` - (Optional) Type of the resource.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `transit_router_attachments` - A list of CEN Transit Router peer Attachments. Each element contains the following attributes:
    * `transit_router_attachment_id` - ID of the transit router attachment.
    * `peer_transit_router_region_id` - Region ID of the peer transit router.
    * `peer_transit_router_owner_id` - Owner ID of the peer transit router.
    * `peer_transit_router_id` - ID of the peer transit router.
    * `transit_router_attachment_name` - Name of the transit router attachment.
    * `resource_type` - Type of the resource.
    * `status` - The status of the transit router attachment.
    * `creation_time` - The time when it's created.
    * `transit_router_id` - ID of the transit router.
    * `bandwidth` - The bandwidth of the bandwidth package.
    * `bandwidth_package_id` - ID of the bandwidth package.
    * `region_id` - ID of the region where the conflicted route entry is located.
    * `geographic_span_id` - ID of the geographic.
    * `cen_bandwidth_package_id` - ID of the CEN bandwidth package.
