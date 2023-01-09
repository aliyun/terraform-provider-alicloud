---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_inter_region_traffic_qos_queue"
sidebar_current: "docs-alicloud-resource-cen-inter-region-traffic-qos-queue"
description: |-
  Provides a Alicloud Cen Inter Region Traffic Qos Queue resource.
---

# alicloud_cen_inter_region_traffic_qos_queue

Provides a Cen Inter Region Traffic Qos Queue resource.

For information about Cen Inter Region Traffic Qos Queue and how to use it, see [What is Inter Region Traffic Qos Queue](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/api-doc-cbn-2017-09-12-api-doc-createceninterregiontrafficqosqueue).

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cen_inter_region_traffic_qos_queue" "default" {
  remain_bandwidth_percent                   = 20
  traffic_qos_policy_id                      = "qos-xxxxxx"
  dscps                                      = [1, 2]
  inter_region_traffic_qos_queue_description = "test"
}
```

## Argument Reference

The following arguments are supported:
* `traffic_qos_policy_id` - (Required,ForceNew) The ID of the traffic scheduling policy.
* `remain_bandwidth_percent` - (Required) The percentage of cross-region bandwidth that the current queue can use.
* `dscps` - (Required) The DSCP value of the traffic packet to be matched in the current queue, ranging from 0 to 63.
* `inter_region_traffic_qos_queue_name` - (Optional) The name of the traffic scheduling policy.
* `inter_region_traffic_qos_queue_description` - (Optional) The description information of the traffic scheduling policy.



## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `status` - The status of the traffic scheduling policy. -**Creating**: The function is being created.-**Active**: available.-**Modifying**: is being modified.-**Deleting**: Deleted.-**Deleted**: Deleted.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Inter Region Traffic Qos Queue.
* `delete` - (Defaults to 5 mins) Used when delete the Inter Region Traffic Qos Queue.
* `update` - (Defaults to 5 mins) Used when update the Inter Region Traffic Qos Queue.

## Import

Cen Inter Region Traffic Qos Queue can be imported using the id, e.g.

```shell
$terraform import alicloud_cen_inter_region_traffic_qos_queue.example <id>
```