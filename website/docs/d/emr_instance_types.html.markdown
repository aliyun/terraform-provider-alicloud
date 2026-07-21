---
subcategory: "E-MapReduce (EMR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_emr_instance_types"
sidebar_current: "docs-alicloud-datasource-emr-instance-types"
description: |-
    Provides a collection of ecs instance types when create emr cluster according to the specified filters.
---

# alicloud\_emr\_instance\_types

The `alicloud_emr_instance_types` data source provides a collection of ecs
instance types available in Alibaba Cloud account when create a emr cluster.

-> **NOTE:** Available in 1.59.0+

## Example Usage

```
data "alicloud_emr_instance_types" "default" {
  destination_resource  = "InstanceType"
  instance_charge_type  = "PostPaid"
  support_local_storage = false
  cluster_type          = "HADOOP"
  instance_type         = "ecs.g5.2xlarge"
  support_node_type     = ["MASTER", "CORE"]
}

output "first_instance_type" {
  value = "${data.alicloud_emr_instance_types.default.types.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `destination_resource` - (Required) The destination resource of emr cluster instance
* `instance_charge_type` - (Required) Filter the results by charge type. Valid values: `PrePaid` and `PostPaid`. Default to `PostPaid`.
* `cluster_type` - (Required) The cluster type of the emr cluster instance. Possible values: `HADOOP`, `KAFKA`, `ZOOKEEPER`, `DRUID`.
* `instance_type` - (Optional, Available in 1.71.2+) Filter the specific ecs instance type to create emr cluster.
* `support_local_storage` - (Optional, Available in 1.61.0+) Whether the current storage disk is local or not.
* `support_node_type` - (Optional, Available in 1.63.0+) The specific supported node type list.
Possible values may be any one or combination of these: ["MASTER", "CORE", "TASK", "GATEWAY"]
* `zone_id` - (Optional,Available in 1.69.0+) The supported resources of specific zoneId.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of emr instance types IDs. 
* `types` - A list of emr instance types. Each element contains the following attributes:
  * `id` - The ID of the instance type.
  * `zone_id` - The available zone id in Alibaba Cloud account
  * `local_storage_capacity` - Local capacity of the applied ecs instance for emr cluster. Unit: GB.
