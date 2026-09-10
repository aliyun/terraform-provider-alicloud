---
subcategory: "Distributed Relational Database Service (DRDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_drds_polardbx_instances"
description: |-
  Provides a list of Alicloud Distributed Relational Database Service (DRDS) Polardbx Instances to the user.
---

# alicloud_drds_polardbx_instances

This data source provides the Distributed Relational Database Service (DRDS) PolarDB-X Instances of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.285.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_drds_polardbx_instances" "ids" {
  ids = ["pxc-xxxxxxxx"]
}

output "first_polardbx_instance_id" {
  value = data.alicloud_drds_polardbx_instances.ids.instances.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of PolarDB-X Instance IDs.
* `description_regex` - (Optional, ForceNew) A regex string to filter results by instance description.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group used to filter instances.
* `status` - (Optional, ForceNew) Filter results by instance status. Valid values: `Creating`, `Running`, `MinorVersionUpgrading`, `ClassChanging`, `NodeCreating`, `NodeDeleting`, `Deleting`.
* `page_number` - (Optional) Number of the returned page. Valid values: 1 to `10000`. Default: `1`.
* `page_size` - (Optional) Number of entries per page. Valid values: 1 to `1000`. Default: `100`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `descriptions` - A list of PolarDB-X instance descriptions.
* `total_count` - Total number of PolarDB-X instances returned by the API.
* `instances` - A list of PolarDB-X instances. Each element contains the following attributes:
  * `id` - The ID of the PolarDB-X instance.
  * `polardbx_instance_id` - The ID of the PolarDB-X instance. Same as `id`.
  * `cn_class` - Compute node specifications of the instance.
  * `cn_node_count` - The number of compute nodes.
  * `create_time` - The creation time of the instance.
  * `description` - Instance remarks.
  * `dn_class` - Storage node specifications of the instance.
  * `dn_node_count` - The number of storage nodes.
  * `engine_version` - Engine version of the instance.
  * `network_type` - The network type of the instance.
  * `payment_type` - The billing method of the instance. Valid values: `Postpaid`, `Prepaid`.
  * `primary_zone` - Primary availability zone.
  * `region_id` - Region ID of the instance.
  * `resource_group_id` - The ID of the resource group the instance belongs to.
  * `secondary_zone` - Secondary availability zone.
  * `status` - Status of the instance.
  * `storage_type` - Storage type of the instance. Valid values: `custom_local_ssd`, `cloud_auto`.
  * `tertiary_zone` - Third availability zone.
  * `topology_type` - Topology type of the instance. Valid values: `1azone`, `3azones`.
  * `vpc_id` - VPC ID of the instance.
  * `zone_id` - Availability zone of the instance.
