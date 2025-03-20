---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_inter_region_traffic_qos_queue"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Inter Region Traffic Qos Queue resource.
---

# alicloud_cen_inter_region_traffic_qos_queue

Provides a Cloud Enterprise Network (CEN) Inter Region Traffic Qos Queue resource.



For information about Cloud Enterprise Network (CEN) Inter Region Traffic Qos Queue and how to use it, see [What is Inter Region Traffic Qos Queue](https://next.api.alibabacloud.com/document/Cbn/2017-09-12/CreateCenInterRegionTrafficQosQueue).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}
variable "default_region" {
  default = "cn-hangzhou"
}
variable "peer_region" {
  default = "cn-beijing"
}
provider "alicloud" {
  alias  = "hz"
  region = var.default_region
}
provider "alicloud" {
  alias  = "bj"
  region = var.peer_region
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  protection_level  = "REDUCED"
}

resource "alicloud_cen_bandwidth_package" "default" {
  provider                   = alicloud.hz
  bandwidth                  = 5
  cen_bandwidth_package_name = "tf_example"
  geographic_region_a_id     = "China"
  geographic_region_b_id     = "China"
}

resource "alicloud_cen_bandwidth_package_attachment" "default" {
  provider             = alicloud.hz
  instance_id          = alicloud_cen_instance.default.id
  bandwidth_package_id = alicloud_cen_bandwidth_package.default.id
}

resource "alicloud_cen_transit_router" "default" {
  provider          = alicloud.hz
  cen_id            = alicloud_cen_instance.default.id
  support_multicast = true
}

resource "alicloud_cen_transit_router" "peer" {
  provider          = alicloud.bj
  cen_id            = alicloud_cen_transit_router.default.cen_id
  support_multicast = true
}

resource "alicloud_cen_transit_router_peer_attachment" "default" {
  provider                              = alicloud.hz
  cen_id                                = alicloud_cen_instance.default.id
  transit_router_id                     = alicloud_cen_transit_router.default.transit_router_id
  peer_transit_router_region_id         = var.peer_region
  peer_transit_router_id                = alicloud_cen_transit_router.peer.transit_router_id
  cen_bandwidth_package_id              = alicloud_cen_bandwidth_package_attachment.default.bandwidth_package_id
  bandwidth                             = 5
  transit_router_attachment_description = var.name
  transit_router_attachment_name        = var.name
}

resource "alicloud_cen_inter_region_traffic_qos_policy" "default" {
  provider                                    = alicloud.hz
  transit_router_id                           = alicloud_cen_transit_router.default.transit_router_id
  transit_router_attachment_id                = alicloud_cen_transit_router_peer_attachment.default.transit_router_attachment_id
  inter_region_traffic_qos_policy_name        = var.name
  inter_region_traffic_qos_policy_description = var.name
}

resource "alicloud_cen_inter_region_traffic_qos_queue" "default" {
  remain_bandwidth_percent                   = 20
  traffic_qos_policy_id                      = alicloud_cen_inter_region_traffic_qos_policy.default.id
  dscps                                      = [1, 2]
  inter_region_traffic_qos_queue_description = var.name
}
```

## Argument Reference

The following arguments are supported:
* `bandwidth` - (Optional, Available since v1.246.0) The guaranteed bandwidth value. If guaranteed by bandwidth is selected for TrafficQosPolicy, this value is valid.
* `dscps` - (Required, List) The DSCP value of the traffic packet to be matched in the current queue, ranging from 0 to 63.
* `inter_region_traffic_qos_queue_description` - (Optional) The description information of the traffic scheduling policy.
* `inter_region_traffic_qos_queue_name` - (Optional) The name of the traffic scheduling policy.
* `remain_bandwidth_percent` - (Optional, Int) The percentage of cross-region bandwidth that the current queue can use.
* `traffic_qos_policy_id` - (Required, ForceNew) The ID of the traffic scheduling policy.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - The status of the traffic scheduling policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Inter Region Traffic Qos Queue.
* `delete` - (Defaults to 5 mins) Used when delete the Inter Region Traffic Qos Queue.
* `update` - (Defaults to 5 mins) Used when update the Inter Region Traffic Qos Queue.

## Import

Cloud Enterprise Network (CEN) Inter Region Traffic Qos Queue can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_inter_region_traffic_qos_queue.example <id>
```