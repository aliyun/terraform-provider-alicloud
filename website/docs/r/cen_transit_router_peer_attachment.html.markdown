---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_peer_attachment"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Transit Router Peer Attachment resource.
---

# alicloud_cen_transit_router_peer_attachment

Provides a Cloud Enterprise Network (CEN) Transit Router Peer Attachment resource.



For information about Cloud Enterprise Network (CEN) Transit Router Peer Attachment and how to use it, see [What is Transit Router Peer Attachment](https://next.api.alibabacloud.com/document/Cbn/2017-09-12/CreateTransitRouterPeerAttachment).

-> **NOTE:** Available since v1.128.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}
variable "region" {
  default = "cn-hangzhou"
}
variable "peer_region" {
  default = "cn-beijing"
}
provider "alicloud" {
  alias  = "hz"
  region = var.region
}
provider "alicloud" {
  alias  = "bj"
  region = var.peer_region
}

resource "alicloud_cen_instance" "example" {
  provider          = alicloud.bj
  cen_instance_name = var.name
  protection_level  = "REDUCED"
}

resource "alicloud_cen_bandwidth_package" "example" {
  provider                   = alicloud.bj
  bandwidth                  = 5
  cen_bandwidth_package_name = "tf_example"
  geographic_region_a_id     = "China"
  geographic_region_b_id     = "China"
}

resource "alicloud_cen_bandwidth_package_attachment" "example" {
  provider             = alicloud.bj
  instance_id          = alicloud_cen_instance.example.id
  bandwidth_package_id = alicloud_cen_bandwidth_package.example.id
}

resource "alicloud_cen_transit_router" "example" {
  provider = alicloud.hz
  cen_id   = alicloud_cen_bandwidth_package_attachment.example.instance_id
}

resource "alicloud_cen_transit_router" "peer" {
  provider = alicloud.bj
  cen_id   = alicloud_cen_transit_router.example.cen_id
}

resource "alicloud_cen_transit_router_peer_attachment" "example" {
  provider                              = alicloud.hz
  cen_id                                = alicloud_cen_instance.example.id
  transit_router_id                     = alicloud_cen_transit_router.example.transit_router_id
  peer_transit_router_region_id         = var.peer_region
  peer_transit_router_id                = alicloud_cen_transit_router.peer.transit_router_id
  cen_bandwidth_package_id              = alicloud_cen_bandwidth_package_attachment.example.bandwidth_package_id
  bandwidth                             = 5
  transit_router_attachment_description = var.name
  transit_router_attachment_name        = var.name
}
```

## Argument Reference

The following arguments are supported:
* `auto_publish_route_enabled` - (Optional) Specifies whether to enable the local Enterprise Edition transit router to automatically advertise the routes of the inter-region connection to the peer transit router. Valid values:

  - `false` (default): no
  - `true`: yes
* `bandwidth` - (Optional, Int) The bandwidth value of the inter-region connection. Unit: Mbit/s.

  - This parameter specifies the maximum bandwidth value for the inter-region connection if you set `BandwidthType` to `BandwidthPackage`.
  - This parameter specifies the bandwidth throttling threshold for the inter-region connection if you set `BandwidthType` to `DataTransfer`.
* `bandwidth_type` - (Optional, Available since v1.157.0) The method that is used to allocate bandwidth to the inter-region connection. Valid values:

  - `BandwidthPackage`: allocates bandwidth from a bandwidth plan.
  - `DataTransfer`: bandwidth is billed based on the pay-by-data-transfer metering method.
* `cen_bandwidth_package_id` - (Optional) The ID of the bandwidth plan that is used to allocate bandwidth to the inter-region connection.

-> **NOTE:**   If you set `BandwidthType` to `DataTransfer`, you do not need to set this parameter.

* `cen_id` - (Optional, ForceNew) The ID of the Cloud Enterprise Network (CEN) instance.
* `default_link_type` - (Optional, Available since v1.223.1) The default line type.
Valid values: Platinum and Gold.
Platinum is supported only when BandwidthType is set to DataTransfer.
* `dry_run` - (Optional) Whether to perform PreCheck on this request, including permissions and instance status verification. Value:
  - `false` (default): A normal request is sent, and a cross-region connection is created after the check is passed.
  - `true`: The check request is sent only for verification, and no cross-region connection is created. Check items include whether required parameters and request format are filled in. If the check does not pass, the corresponding error is returned. If the check is passed, the corresponding request ID is returned.
* `peer_transit_router_id` - (Required, ForceNew) The ID of the peer transit router.
* `peer_transit_router_region_id` - (Optional, ForceNew) The ID of the region where the peer transit router is deployed.
* `tags` - (Optional, Map, Available since v1.247.0) The tag of the resource
* `transit_router_attachment_description` - (Optional) The new description of the inter-region connection.
This parameter is optional. If you enter a description, it must be 1 to 256 characters in length, and cannot start with http:// or https://.
* `transit_router_id` - (Optional, ForceNew) The ID of the local Enterprise Edition transit router.
* `transit_router_peer_attachment_name` - (Optional, Available since v1.247.0) The new name of the inter-region connection.
The name can be empty or 1 to 128 characters in length, and cannot start with http:// or https://.
* `resource_type` - (Optional, ForceNew) The resource type to attachment. Only support `VR` and default value is `VR`.

The following arguments will be discarded. Please use new fields as soon as possible:
* `transit_router_attachment_name` - (Deprecated since v1.247.0). Field 'transit_router_attachment_name' has been deprecated from provider version 1.247.0. New field 'transit_router_peer_attachment_name' instead.
* `route_table_association_enabled` - (Deprecated since v1.230.0) Field `route_table_association_enabled` has been deprecated from provider version 1.230.0.
* `route_table_propagation_enabled` - (Deprecated since v1.230.0) Field `route_table_propagation_enabled` has been deprecated from provider version 1.230.0.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource
* `region_id` - The ID of the region where the local Enterprise Edition transit router is deployed.
* `status` - The status of the resource
* `transit_router_attachment_id` - The ID of the inter-region connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Transit Router Peer Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Transit Router Peer Attachment.
* `update` - (Defaults to 5 mins) Used when update the Transit Router Peer Attachment.

## Import

Cloud Enterprise Network (CEN) Transit Router Peer Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_peer_attachment.example <cen_id>:<transit_router_attachment_id>
```