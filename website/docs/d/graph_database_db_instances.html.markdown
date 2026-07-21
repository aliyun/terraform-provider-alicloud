---
subcategory: "Graph Database"
layout: "alicloud"
page_title: "Alicloud: alicloud_graph_database_db_instances"
sidebar_current: "docs-alicloud-datasource-graph-database-db-instances"
description: |-
  Provides a list of Graph Database Db Instances to the user.
---

# alicloud\_graph\_database\_db\_instances

This data source provides the Graph Database Db Instances of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.136.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_graph_database_db_instances" "ids" {
  ids = ["example_id"]
}
output "graph_database_db_instance_id_1" {
  value = data.alicloud_graph_database_db_instances.ids.instances.0.id
}

data "alicloud_graph_database_db_instances" "status" {
  ids    = ["example_id"]
  status = "Running"
}
output "graph_database_db_instance_id_2" {
  value = data.alicloud_graph_database_db_instances.status.instances.0.id
}

data "alicloud_graph_database_db_instances" "description" {
  ids                     = ["example_id"]
  db_instance_description = "example_value"
}
output "graph_database_db_instance_id_3" {
  value = data.alicloud_graph_database_db_instances.description.instances.0.id
}

```

## Argument Reference

The following arguments are supported:

* `db_instance_description` - (Optional, ForceNew) According to the practical example or notes.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Db Instance IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) Instance status. Value range: `Creating`, `Running`, `Deleting`, `DBInstanceClassChanging`, `NetAddressCreating` and `NetAddressDeleting`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `instances` - A list of Graph Database Db Instances. Each element contains the following attributes:
    * `connection_string` - Virtual Private Cloud (vpc connection such as a VPN connection or leased line domain name).
    * `create_time` - Creation time, which follows the format of `YYYY-MM-DD 'T'hh:mm:ssZ`, such as `2011-05-30 T12:11:4Z`.
    * `current_minor_version` - The current kernel image version.
    * `db_instance_category` - The category of the db instance.
    * `db_instance_cpu` - For example, instances can be grouped according to Cpu core count.
    * `db_instance_description` - According to the practical example or notes.
    * `db_instance_id` - The ID of the instance.
    * `db_instance_memory` - Instance memory, which is measured in MB.
    * `db_instance_network_type` - The network type of the db instance.
    * `db_instance_storage_type` - Disk storage type. Valid values: `cloud_essd`, `cloud_ssd`.
    * `db_instance_type` - The type of the db instance.
    * `db_node_class` - The class of the db node.
    * `db_node_count` - The count of the db node.
    * `db_node_storage` - Instance storage space, which is measured in GB.
    * `db_version` - Kernel Version. Value range: `1.0` or `1.0-OpenCypher`. `1.0`: represented as gremlin, `1.0-OpenCypher`: said opencypher.
    * `expire_time` - The instance after it expires time for subscription instance.
    * `expired` - The expire status of the db instance.
    * `id` - The ID of the Db Instance.
    * `latest_minor_version` - The latest kernel image version.
    * `lock_mode` - Instance lock state. Value range: `Unlock`, `ManualLock`, `LockByExpiration`, `LockByRestoration` and `LockByDiskQuota`. `Unlock`: normal. `ManualLock`: the manual trigger lock. `LockByExpiration`: that represents the instance expires automatically lock. `LockByRestoration`: indicates that the instance rollback before auto-lock. `LockByDiskQuota`: that represents the instance space full automatic lock.
    * `lock_reason` - An instance is locked the reason.
    * `maintain_time` - Instance maintenance time such as `00:00Z-02:00Z`, 0 to 2 points to carry out routine maintenance.
    * `master_db_instance_id` - The master instance ID of the db instance.
    * `payment_type` - The paymen type of the resource.
    * `port` - Application Port.
    * `public_connection_string` - The public connection string ID of the resource.
    * `public_port` - The public port ID of the resource.
    * `read_only_db_instance_ids` - The array of the readonly db instances.
    * `status` - Instance status. Value range: `Creating`, `Running`, `Deleting`, `Rebooting`, `DBInstanceClassChanging`, `NetAddressCreating` and `NetAddressDeleting`.
    * `vpc_id` - The vpc id of the db instance.
    * `vswitch_id` - The vswitch id.
    * `zone_id` - The zone ID of the resource.
    * `db_instance_ip_array` - IP ADDRESS whitelist for the instance group list.
        * `db_instance_ip_array_attribute` - The default is empty. To distinguish between the different property console does not display a `hidden` label grouping.
        * `db_instance_ip_array_name` - IP ADDRESS whitelist group name.
        * `security_ips` - IP ADDRESS whitelist addresses in the IP ADDRESS list, and a maximum of 1000 comma-separated format is as follows: `0.0.0.0/0` and `10.23.12.24`(IP) or `10.23.12.24/24`(CIDR mode, CIDR (Classless Inter-Domain Routing)/24 represents the address prefixes in the length of the range [1,32]).
