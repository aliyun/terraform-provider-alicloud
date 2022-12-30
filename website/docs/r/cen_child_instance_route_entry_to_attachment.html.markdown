---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_child_instance_route_entry_to_attachment"
sidebar_current: "docs-alicloud-resource-cen-child-instance-route-entry-to-attachment"
description: |-
  Provides a Alicloud Cen Child Instance Route Entry To Attachment resource.
---

# alicloud_cen_child_instance_route_entry_to_attachment

Provides a Cen Child Instance Route Entry To Attachment resource.

For information about Cen Child Instance Route Entry To Attachment and how to use it, see [What is Child Instance Route Entry To Attachment](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/api-doc-cbn-2017-09-12-api-doc-createcenchildinstancerouteentrytoattachment).

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cen_child_instance_route_entry_to_attachment" "default" {
  transit_router_attachment_id  = "tr-attach-f1fd1y50rql00emvej"
  cen_id                        = "cen-3sgjn0u745c3i0o3dk"
  destination_cidr_block        = "10.0.0.0/24"
  child_instance_route_table_id = "vtb-t4nt0z5xxbti85c78nkzy"
}
```

## Argument Reference

The following arguments are supported:
* `cen_id` - (Required,ForceNew) The ID of the CEN instance.
* `child_instance_route_table_id` - (Required,ForceNew) The first ID of the resource
* `destination_cidr_block` - (Required,ForceNew) DestinationCidrBlock
* `transit_router_attachment_id` - (Required,ForceNew) TransitRouterAttachmentId
* `dry_run` - (Optional) Whether to perform pre-check on this request, including permission and instance status verification.

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.The value is formulated as `<cen_id>:<child_instance_route_table_id>:<transit_router_attachment_id>:<destination_cidr_block>`.
* `service_type` - ServiceType
* `status` - The status of the resource

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Child Instance Route Entry To Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Child Instance Route Entry To Attachment.

## Import

Cen Child Instance Route Entry To Attachment can be imported using the id, e.g.

```shell
$terraform import alicloud_cen_child_instance_route_entry_to_attachment.example <cen_id>:<child_instance_route_table_id>:<transit_router_attachment_id>:<destination_cidr_block>
```