---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_vbr_attachment"
sidebar_current: "docs-alicloud-resource-cen-transit_router_vbr_attachment"
description: |-
  Provides a Alicloud CEN transit router VBR attachment resource.
---

# alicloud_cen_transit_router_vbr_attachment

Provides a CEN transit router VBR attachment resource that associate the VBR with the CEN instance.[What is Cen Transit Router VBR Attachment](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-createtransitroutervbrattachment)

-> **NOTE:** Available since v1.126.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform-example"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  protection_level  = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
}

data "alicloud_express_connect_physical_connections" "nameRegex" {
  name_regex = "^preserved-NODELETING"
}

resource "alicloud_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.nameRegex.connections.0.id
  virtual_border_router_name = var.name
  vlan_id                    = 2420
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

resource "alicloud_cen_transit_router_vbr_attachment" "default" {
  transit_router_id                     = alicloud_cen_transit_router.default.transit_router_id
  transit_router_attachment_name        = "example"
  transit_router_attachment_description = "example"
  vbr_id                                = alicloud_express_connect_virtual_border_router.default.id
  cen_id                                = alicloud_cen_instance.default.id
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
* `route_table_association_enabled` - (Optional, ForceNew) Whether to enabled route table association. The system default value is `true`.
* `route_table_propagation_enabled` - (Optional, ForceNew) Whether to enabled route table propagation. The system default value is `true`.  
* `dry_run` - (Optional) The dry run.
* `tags` - (Optional, Available in v1.193.1+) A mapping of tags to assign to the resource.
* `vbr_owner_id` - (Optional, ForceNew) The owner id of the transit router vbr attachment.
* `resource_type` - (Optional) The resource type of the transit router vbr attachment.  Valid values: `VPC`, `CCN`, `VBR`, `TR`.

->**NOTE:** Ensure that the vbr is not used in Express Connect.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, It is formatted to `<transit_router_id>:<transit_router_attachment_id>`. 
* `status` - The associating status of the network.
* `transit_router_attachment_id` - The id of the transit router vbr attachment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the cen transit router vbr attachment (until it reaches the initial `Attached` status).
* `update` - (Defaults to 10 mins) Used when update the cen transit router vbr attachment.
* `delete` - (Defaults to 10 mins) Used when delete the cen transit router vbr attachment.

## Import

CEN transit router VBR attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_vbr_attachment.example tr-********:tr-attach-********
```
