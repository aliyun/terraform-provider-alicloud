---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_child_instance_route_entry_to_attachments"
sidebar_current: "docs-alicloud-datasource-cen-child-instance-route-entry-to-attachments"
description: |-
  Provides a list of Cen Child Instance Route Entry To Attachment owned by an Alibaba Cloud account.
---

# alicloud_cen_child_instance_route_entry_to_attachments

This data source provides Cen Child Instance Route Entry To Attachment available to the user.[What is Child Instance Route Entry To Attachment](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-createcenchildinstancerouteentrytoattachment)

-> **NOTE:** Available in 1.195.0+

## Example Usage

```
data "alicloud_cen_child_instance_route_entry_to_attachments" "default" {
  child_instance_route_table_id = "vtb-t4nt0z5xxbti85c78nkzy"
  transit_router_attachment_id  = "tr-attach-f1fd1y50rql00emvej"
}

output "alicloud_cen_child_instance_route_entry_to_attachment_example_id" {
  value = data.alicloud_cen_child_instance_route_entry_to_attachments.default.attachments.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional) Limit search to a list of specific IDs.The value is formulated as `<cen_id>:<child_instance_route_table_id>:<transit_router_attachment_id>:<destination_cidr_block>`.
* `child_instance_route_table_id` - (Required,ForceNew) The first ID of the resource
* `service_type` - (ForceNew,Optional) ServiceType
* `cen_id` - (ForceNew,Optional) The ID of the CEN instance.
* `transit_router_attachment_id` - (Required,ForceNew) TransitRouterAttachmentId
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - Limit search to a list of specific IDs.The value is formulated as `<cen_id>:<child_instance_route_table_id>:<transit_router_attachment_id>:<destination_cidr_block>`.
* `attachments` - A list of Child Instance Route Entry To Attachment Entries. Each element contains the following attributes:
    * `id` - The ID of the resource. The value is formulated as `<cen_id>:<child_instance_route_table_id>:<transit_router_attachment_id>:<destination_cidr_block>`.
    * `cen_id` - The ID of the CEN instance.
    * `child_instance_route_table_id` - The first ID of the resource
    * `destination_cidr_block` - DestinationCidrBlock
    * `service_type` - ServiceType
    * `status` - The status of the resource
    * `transit_router_attachment_id` - TransitRouterAttachmentId
