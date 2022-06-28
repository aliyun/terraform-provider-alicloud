---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_vpc_attachments"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-vpc-attachments"
description: |-
  Provides a list of CEN Transit Router VPC Attachments owned by an Alibaba Cloud account.
---

# alicloud\_cen\_transit\_router\_vpc\_attachments

This data source provides CEN Transit Router VPC Attachments available to the user.[What is Cen Transit Router VPC Attachments](https://help.aliyun.com/document_detail/261222.html)

-> **NOTE:** Available in 1.126.0+

## Example Usage

```
data "alicloud_cen_transit_router_vpc_attachments" "default" {
  cen_id    = "cen-id1"
}

output "first_transit_router_vpc_attachments_transit_router_attachment_vpc_id" {
  value = data.alicloud_cen_transit_router_vpc_attachments.default.transit_router_attachments.0.vpc_id
}
```

## Argument Reference

The following arguments are supported:

* `cen_id` - (Required) ID of the CEN instance.
* `ids` - (Optional) A list of resource id. The element value is same as `transit_router_id`.
* `status` - (Optional) The status of the resource. Valid values `Attached`, `Attaching` and `Detaching`.
* `transit_router_id` - (Optional) The transit router ID.
* `transit_router_attachment_id` - (Optional) ID of the transit router VBR attachment.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `attachments` - A list of CEN Transit Router VPC Attachments. Each element contains the following attributes:
    * `transit_router_attachment_id` - ID of the transit router attachment.
    * `transit_router_attachment_name` - Name of the transit router attachment.
    * `resource_type` - Type of the resource.
    * `transit_router_attachment_description` - The description of transit router attachment.
    * `status` - The status of the transit router attachment.
    * `vpc_id` - ID of the VPC.
    * `id` -  The ID of the transit router.
    * `vpc_owner_id` - The Owner ID of the VPC.     
    * `transit_router_id` - ID of the transit router.
    * `payment_type` - The payment type of the resource.
    * `zone_mappings` - The mappings of zone
        * `vswitch_id` - The VSwitch ID.
        * `zone_id` - The zone ID.
