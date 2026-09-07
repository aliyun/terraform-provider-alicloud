---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_instances"
sidebar_current: "docs-alicloud-datasource-gpdb-instances"
description: |-
  Provides a collection of AnalyticDB for PostgreSQL instances according to the specified filters.
---

# alicloud\_gpdb\_instances

This data source provides the AnalyticDB for PostgreSQL instances of the current Alibaba Cloud user.

-> **NOTE:**  Available in 1.47.0+

## Example Usage

Basic Usage

```terraform
data "alicloud_gpdb_instances" "ids" {}
output "gpdb_db_instance_id_1" {
  value = data.alicloud_gpdb_instances.ids.instances.0.id
}

```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the instance name.
* `availability_zone` - (Optional) Instance availability zone.
* `vswitch_id` - (Optional) Used to retrieve instances belong to specified `vswitch` resources.
* `db_instance_categories` - (Optional) The db instance categories.
* `db_instance_modes` - (Optional) The db instance modes.
* `description` - (Optional) The description of the instance.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional)  A list of DB Instance IDs.
* `instance_network_type` - (Optional) The network type of the instance.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional) The ID of the enterprise resource group to which the instance belongs.
* `status` - (Optional) The status of the instance. Valid values: `Creating`, `DBInstanceClassChanging`, `DBInstanceNetTypeChanging`, `Deleting`, `EngineVersionUpgrading`, `GuardDBInstanceCreating`, `GuardSwitching`, `Importing`, `ImportingFromOtherInstance`, `Rebooting`, `Restoring`, `Running`, `Transfering`, `TransferingToOtherInstance`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - The ids list of AnalyticDB for PostgreSQL instances.
* `names` - The names list of AnalyticDB for PostgreSQL instance.
* `instances` - A list of Gpdb Db Instances. Each element contains the following attributes:
	* `connection_string` - The endpoint of the instance.
	* `cpu_cores` - The number of CPU cores of the computing node. Unit: Core.
	* `create_time` - The time when the instance was created. The time is in the YYYY-MM-DDThh:mm:ssZ format, such as 2011-05-30T12:11:4Z.
  	* `db_instance_category` - The db instance category. Valid values: `HighAvailability`, `Basic`.
  	* `db_instance_class` - The db instance class.
	* `db_instance_id` - The db instance id.
	* `db_instance_mode` - The db instance mode. Valid values: `StorageElastic`, `Serverless`, `Classic`.
	* `description` - The description of the instance.
	* `engine` - The database engine used by the instance.
	* `engine_version` - The version of the database engine used by the instance.
	* `id` - The ID of the db Instance.
	* `instance_network_type` - The network type of the instance.
	* `ip_whitelist` - The ip whitelist.
    * `ip_group_attribute` - The value of this parameter is empty by default. The attribute of the whitelist group. The console does not display the whitelist group whose value of this parameter is hidden.
    * `ip_group_name` - IP whitelist group name
    * `security_ip_list` - List of IP addresses allowed to access all databases of an instance. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]). System default to `["127.0.0.1"]`.
	* `maintain_end_time` - The end time of the maintenance window for the instance.
	* `maintain_start_time` - The start time of the maintenance window for the instance.
	* `master_node_num` - The number of Master nodes. Valid values: 1 to 2. if it is not filled in, the default value is 1 Master node.
	* `memory_size` - The memory size of the compute node.
	* `payment_type` - The billing method of the instance. Valid values: `Subscription`, `PayAsYouGo`.
	* `seg_node_num` - Calculate the number of nodes. The value range of the high-availability version of the storage elastic mode is 4 to 512, and the value must be a multiple of 4. The value range of the basic version of the storage elastic mode is 2 to 512, and the value must be a multiple of 2. The-Serverless version has a value range of 2 to 512. The value must be a multiple of 2.
	* `status` - The status of the instance. Valid values: `Creating`, `DBInstanceClassChanging`, `DBInstanceNetTypeChanging`, `Deleting`, `EngineVersionUpgrading`, `GuardDBInstanceCreating`, `GuardSwitching`, `Importing`, `ImportingFromOtherInstance`, `Rebooting`, `Restoring`, `Running`, `Transfering`, `TransferingToOtherInstance`.
	* `storage_size` - The storage capacity. Unit: GB. Value: `50` to `4000`.
	* `storage_type` - The type of disks. Valid values: `cloud_essd`, `cloud_efficiency`.
	* `tags` - The tags of the instance.
	* `vpc_id` - The ID of the VPC。.
	* `vswitch_id` - The vswitch id.
	* `zone_id` - The zone ID of the instance.
  	* `region_id` - Region ID the instance belongs to.
	* `connection_string` - (Available in 1.196.0+) The connection string of the instance.
	* `port` - (Available in 1.196.0+) The connection port of the instance.

