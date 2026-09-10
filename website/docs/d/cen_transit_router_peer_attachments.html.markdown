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

-> **NOTE:** Available in 1.128.0+

## Example Usage

```terraform
data "alicloud_cen_transit_router_peer_attachments" "default" {
  cen_id = "cen-id1"
}

output "first_transit_router_peer_attachments_transit_router_attachment_resource_type" {
  value = "${data.alicloud_cen_transit_router_peer_attachments.default.transit_router_attachments.0.resource_type}"
}
```

## Argument Reference

The following arguments are supported:

* `cen_id` - (Required, ForceNew) ID of the CEN instance.
* `name_regex` - (Optional, ForceNew) A regex string to filter CEN Transit Router peer attachments by name.
* `route_table_id` - (Optional, ForceNew) ID of the route table of the VPC or VBR.
* `resource_type` - (Optional, ForceNew) Type of the resource.
* `status` - (Optional, ForceNew) The status of CEN Transit Router peer attachment. Valid values `Attached`, `Attaching` and `Detaching`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `transit_router_attachment_id` - (Optional, ForceNew) The ID of CEN Transit Router peer attachments.
* `transit_router_id` - (Optional, ForceNew) The ID of transit router.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:


* `ids` - A list of CEN Transit Router peer attachments IDs.
* `names` - A list of CEN Transit Router peer attachments names. 
* `transit_router_attachments` - A list of CEN Transit Router peer Attachments. Each element contains the following attributes:
    * `auto_publish_route_enabled` - Auto publish route enabled.
    * `id` - The ID of CEN Transit Router peer attachments.
    * `transit_router_attachment_description` - The description of CEN Transit Router peer attachments.
    * `transit_router_attachment_id` - ID of the transit router attachment.
    * `peer_transit_router_region_id` - Region ID of the peer transit router.
    * `peer_transit_router_owner_id` - Owner ID of the peer transit router.
    * `peer_transit_router_id` - ID of the peer transit router.
    * `transit_router_attachment_name` - Name of the transit router attachment.
    * `resource_type` - Type of the resource.
    * `status` - The status of the transit router attachment.
    * `transit_router_id` - ID of the transit router.
    * `bandwidth` - The bandwidth of the bandwidth package.
    * `bandwidth_package_id` - ID of the bandwidth package.
    * `geographic_span_id` - ID of the geographic.
    * `cen_bandwidth_package_id` - ID of the CEN bandwidth package.
