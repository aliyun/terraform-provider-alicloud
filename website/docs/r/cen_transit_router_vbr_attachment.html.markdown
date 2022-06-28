---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_vbr_attachment"
sidebar_current: "docs-alicloud-resource-cen-transit_router_vbr_attachment"
description: |-
  Provides a Alicloud CEN transit router VBR attachment resource.
---

# alicloud\_cen_transit_router_vbr_attachment

Provides a CEN transit router VBR attachment resource that associate the VBR with the CEN instance.[What is Cen Transit Router VBR Attachment](https://help.aliyun.com/document_detail/261361.html)

-> **NOTE:** Available in 1.126.0+

## Example Usage

Basic Usage

```terraform
# Create a new instance-attachment and use it to attach one child instance to a new CEN
variable "name" {
  default = "tf-testAccCenTransitRouterVbrAttachment"
}

variable "vbr_id" {
  default = "vbr-xxxxxxxxxx"
}

variable "transit_router_attachment_name" {
  default = "tf-test"
}

variable "transit_router_attachment_description" {
  default = "tf-test"
}

resource "alicloud_cen_instance" "cen" {
  instance_name = var.name
  description   = "terraform01"
}

resource "alicloud_transit_router" "tr" {
  name   = var.name
  cen_id = alicloud_cen_instance.cen.id
}

resource "alicloud_cen_transit_router_vbr_attachment" "foo" {
  vbr_id                                = var.vbr_id
  cen_id                                = alicloud_cen_instance.cen.id
  transit_router_id                     = alicloud_transit_router.tr.transit_router_id
  auto_publish_route_enabled            = true
  transit_router_attachment_name        = var.transit_router_attachment_name
  transit_router_attachment_description = var.transit_router_attachment_description
}
```
## Argument Reference

The following arguments are supported:

* `vbr_id` - (Required, ForceNew) The ID of the VBR.
* `cen_id` - (Required, ForceNew) The ID of the CEN.
* `transit_router_id` - (Optional, ForceNew) The ID of the transit router.
* `auto_publish_route_enabled` - (Optional) Auto publish route enabled.Default value is `false`.
* `transit_router_attachment_name` - (Optional) The name of the transit router vbr attachment.
* `transit_router_attachment_description` - (Optional) The description of the transit router vbr attachment.
* `route_table_association_enabled` - (Optional,ForceNew) Whether to enabled route table association. The system default value is `true`.
* `route_table_propagation_enabled` - (Optional,ForceNew) Whether to enabled route table propagation. The system default value is `true`.  
* `dry_run` - (Optional) The dry run.

->**NOTE:** Ensure that the vbr is not used in Express Connect.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, It is formatted to `<transit_router_id>:<transit_router_attachment_id>`. 
* `status` - The associating status of the network.
* `resource_type` - The resource type of the transit router vbr attachment.  Valid values: `VPC`, `CCN`, `VBR`, `TR`.
* `transit_router_attachment_id` - The id of the transit router vbr attachment.
* `vbr_owner_id` - The owner id of the transit router vbr attachment.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the cen transit router vbr attachment (until it reaches the initial `Attached` status).
* `update` - (Defaults to 10 mins) Used when update the cen transit router vbr attachment.
* `delete` - (Defaults to 10 mins) Used when delete the cen transit router vbr attachment.

## Import

CEN transit router VBR attachment can be imported using the id, e.g.

```
$ terraform import alicloud_cen_transit_router_vbr_attachment.example tr-********:tr-attach-********
```
