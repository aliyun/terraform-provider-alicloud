---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_peer_attachment"
sidebar_current: "docs-alicloud-resource-cen-transit_router_peer_attachment"
description: |-
  Provides a Alicloud CEN transit router peer attachment resource.
---

# alicloud_cen_transit_router_peer_attachment

Provides a CEN transit router peer attachment resource that associate the transit router with the CEN instance. [What is CEN transit router peer attachment](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/api-doc-cbn-2017-09-12-api-doc-createtransitrouterpeerattachment)

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

* `dry_run` - (Optional) Whether to perform pre-check for this request, including permission, instance status verification, etc.
* `cen_id` - (Required, ForceNew) The ID of the CEN.
* `transit_router_id` - (Optional, ForceNew) The ID of the transit router to attach.
* `peer_transit_router_id` - (Required, ForceNew) The ID of the peer transit router.
* `peer_transit_router_region_id` - (Required, ForceNew) The region ID of peer transit router.
* `resource_type` - (Optional, ForceNew) The resource type to attachment. Only support `VR` and default value is `VR`.
* `cen_bandwidth_package_id` - (Optional) The ID of the bandwidth package. If you do not enter the ID of the package, it means you are using the test. The system default test is 1bps, demonstrating that you test network connectivity
* `bandwidth` - (Optional) The bandwidth of the bandwidth package.
* `auto_publish_route_enabled` - (Optional) Auto publish route enabled. The system default value is `false`.
* `route_table_association_enabled` - (Optional, ForceNew) Whether to association route table. System default is `false`.
* `route_table_propagation_enabled` - (Optional, ForceNew) Whether to propagation route table. System default is `false`.
* `transit_router_attachment_description` - (Optional) The description of transit router attachment. The description is 2~256 characters long and must start with a letter or Chinese, but cannot start with `http://` or `https://`.
* `transit_router_attachment_name` - (Optional) The name of transit router attachment. The name is 2~128 characters in length, starts with uppercase and lowercase letters or Chinese, and can contain numbers, underscores (_) and dashes (-)
* `bandwidth_type` - (Optional, Available in v1.157.0+) The method that is used to allocate bandwidth to the cross-region connection. Valid values: `BandwidthPackage` and `DataTransfer`.
  * `DataTransfer` - uses pay-by-data-transfer bandwidth.
  * `BandwidthPackage` - allocates bandwidth from a bandwidth plan.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, It is formatted to `<transit_router_id>:<transit_router_attachment_id>`. 
* `transit_router_attachment_id` - The ID of transit router attachment id.
* `status` - The associating status of the network.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when creating the cen transit router peer attachment (until it reaches the initial `Attached` status).
* `update` - (Defaults to 3 mins) Used when update the cen transit router peer attachment.
* `delete` - (Defaults to 3 mins) Used when delete the cen transit router peer attachment.

## Import

CEN instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_peer_attachment.example tr-********:tr-attach-*******
```
