---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_vbr_attachments"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-vbr-attachments"
description: |-
  Provides a list of CEN Transit Router VBR Attachments owned by an Alibaba Cloud account.
---

# alicloud\_cen\_transit\_router\_vbr\_attachments

This data source provides CEN Transit Router VBR Attachments available to the user.[What is Cen Transit Router VBR Attachments](https://help.aliyun.com/document_detail/261226.html)

-> **NOTE:** Available in 1.126.0+

## Example Usage

```
data "alicloud_cen_transit_router_vbr_attachments" "default" {
  cen_id    = "cen-id1"
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
* `transit_router_attachment_id` - (Optional) ID of the transit router VBR attachment.
* `resource_type` - (Optional) Type of the resource.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of CEN Transit Router VBR attachment IDs.
* `names` - A list of name of CEN Transit VBR attachment Tables.
* `attachments` - A list of CEN Transit Router VBR Attachments. Each element contains the following attributes:
    * `transit_router_attachment_id` - ID of the transit router attachment.
    * `transit_router_attachment_name` - Name of the transit router attachment.
    * `resource_type` - Type of the resource.
    * `status` - The status of the transit router attachment.
    * `creation_time` - The time when it's created.
    * `vbr_id` - ID of the VBR.
    * `vbr_owner_id` - The Owner ID of the VBR.
    * `vbr_region_id` - ID of the region where the conflicted VBR is located.
    * `auto_publish_route_enabled` - ID of the region where the conflicted VBR is located.
