---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_vbr_attachments"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-vbr-attachments"
description: |-
  Provides a list of CEN Transit Router VBR Attachments owned by an Alibaba Cloud account.
---

# alicloud_cen_transit_router_vbr_attachments

This data source provides CEN Transit Router VBR Attachments available to the user.[What is Cen Transit Router VBR Attachments](https://www.alibabacloud.com/help/doc-detail/261226.html)

-> **NOTE:** Available since v1.126.0+

## Example Usage

```terraform
data "alicloud_cen_transit_router_vbr_attachments" "default" {
  cen_id = "cen-id1"
}

output "first_transit_router_vbr_attachments_vbr_id" {
  value = data.alicloud_cen_transit_router_vbr_attachments.default.transit_router_attachments.0.vbr_id
}
```

## Argument Reference

The following arguments are supported:

* `cen_id` - (Required) ID of the CEN instance.
* `ids` - (Optional) A list of resource id. The element value is same as `transit_router_id`.
* `status` - (Optional) The status of the resource. Valid values `Attached`, `Attaching` and `Detaching`.  
* `transit_router_id` - (Optional) ID of the transit router.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of CEN Transit Router VBR attachment IDs.
* `attachments` - A list of CEN Transit Router VBR Attachments. Each element contains the following attributes:
  * `id` - The ID of the CEN Transit Router VBR Attachment.
  * `transit_router_attachment_id` - ID of the transit router attachment.
  * `transit_router_attachment_name` - Name of the transit router attachment.
  * `transit_router_attachment_description` - The description of the transit router attachment.
  * `resource_type` - Type of the resource.
  * `status` - The status of the transit router attachment.
  * `vbr_id` - ID of the VBR.
  * `vbr_owner_id` - The Owner ID of the VBR.
  * `auto_publish_route_enabled` - Auto publish route enabled.
