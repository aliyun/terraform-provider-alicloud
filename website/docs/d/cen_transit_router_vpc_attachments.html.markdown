---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_vpc_attachments"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-vpc-attachments"
description: |-
Provides a list of CEN Transit Router VPC Attachments owned by an Alibaba Cloud account.
---

# alicloud\_cen\_transit\_router\_vpc\_attachments

This data source provides CEN Transit Router VPC Attachments available to the user.

## Example Usage

```
data "alicloud_cen_transit_router_vpc_attachments" "default" {
  cen_id    = "cen-id1"
  region_id = "cn-****"
}

output "first_transit_router_vpc_attachments_transit_router_attachment_vpc_id" {
  value = "${data.alicloud_cen_transit_router_vpc_attachments.default.transit_router_attachments.0.vpc_id}"
}
```

## Argument Reference

The following arguments are supported:

* `cen_id` - (Required) ID of the CEN instance.
* `region_id` - (Optional) Region ID of the VPC.
* `resource_type` - (Optional) Type of the resource.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `transit_router_attachments` - A list of CEN Transit Router VPC Attachments. Each element contains the following attributes:
    * `transit_router_attachment_id` - ID of the transit router attachment.
    * `transit_router_attachment_name` - Name of the transit router attachment.
    * `resource_type` - Type of the resource.
    * `status` - The status of the transit router attachment.
    * `creation_time` - The time when it's created.
    * `vpc_id` - ID of the VPC.      
    * `vpc_owner_id` - The Owner ID of the VPC.     
    * `transit_router_id` - ID of the transit router.
    * `vpc_region_id` - ID of the region where the conflicted VPC is located.
    * `auto_publish_route_enabled` - ID of the region where the conflicted VBR is located.

