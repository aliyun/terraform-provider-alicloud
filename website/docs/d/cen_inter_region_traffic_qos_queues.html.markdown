---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_inter_region_traffic_qos_queues"
sidebar_current: "docs-alicloud-datasource-cen-inter-region-traffic-qos-queues"
description: |-
  Provides a list of Cen Inter Region Traffic Qos Queue owned by an Alibaba Cloud account.
---

# alicloud_cen_inter_region_traffic_qos_queues

This data source provides Cen Inter Region Traffic Qos Queue available to the user.

-> **NOTE:** Available in 1.195.0+

## Example Usage

```
data "alicloud_cen_inter_region_traffic_qos_queues" "default" {
  ids                   = ["${alicloud_cen_inter_region_traffic_qos_queue.default.id}"]
  name_regex            = alicloud_cen_inter_region_traffic_qos_queue.default.name
  traffic_qos_policy_id = "qos-xxxxxxx"
}

output "alicloud_cen_inter_region_traffic_qos_queue_example_id" {
  value = data.alicloud_cen_inter_region_traffic_qos_queues.default.queues.0.id
}
```

## Argument Reference

The following arguments are supported:
* `traffic_qos_policy_id` - (ForceNew,Required) The ID of the traffic scheduling policy.
* `ids` - (Optional, ForceNew, Computed) A list of Inter Region Traffic Qos Queue IDs.
* `names` - (Optional, ForceNew) The name of inter Region Traffic Qos Queue. You can specify at most 10 names.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Inter Region Traffic Qos Queue IDs.
* `names` - A list of name of Inter Region Traffic Qos Queues.
* `queues` - A list of Inter Region Traffic Qos Queue Entries. Each element contains the following attributes:
  * `inter_region_traffic_qos_queue_id` - The ID of the resource.
  * `inter_region_traffic_qos_queue_description` - The description information of the traffic scheduling policy.
  * `inter_region_traffic_qos_queue_name` - The name of the traffic scheduling policy.
  * `traffic_qos_policy_id` - The ID of the traffic scheduling policy.
  * `dscps` - The DSCP value of the traffic packet to be matched in the current queue, ranging from 0 to 63.
  * `remain_bandwidth_percent` - The percentage of cross-region bandwidth that the current queue can use.
  * `status` - The status of the traffic scheduling policy. -**Creating**: The function is being created.-**Active**: available.-**Modifying**: is being modified.-**Deleting**: Deleted.-**Deleted**: Deleted.
