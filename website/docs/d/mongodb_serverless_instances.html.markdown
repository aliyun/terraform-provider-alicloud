---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_serverless_instances"
sidebar_current: "docs-alicloud-datasource-mongodb-serverless-instances"
description: |-
  Provides a list of Mongodb Serverless Instances to the user.
---

# alicloud\_mongodb\_serverless\_instances

This data source provides the Mongodb Serverless Instances of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.148.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_mongodb_serverless_instances" "example" {
  ids                     = ["example_value"]
  db_instance_class       = "example_value"
  db_instance_description = "example_value"
  network_type            = "VPC"
  resource_group_id       = "example_value"
  status                  = "Running"
  vpc_id                  = "example_value"
  vswitch_id              = "example_value"
  zone_id                 = "example_value"
  tags = {
    Created = "MongodbServerlessInstance"
    For     = "TF"
  }
}
output "mongodb_serverless_instance_id_1" {
  value = data.alicloud_mongodb_serverless_instances.example.instances.0.id
}

```

## Argument Reference

The following arguments are supported:

* `db_instance_class` - (Optional, ForceNew) The db instance class.
* `db_instance_description` - (Optional, ForceNew) The db instance description.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Serverless Instance IDs.
* `network_type` - (Optional, ForceNew) The network type of the instance. Valid values: `Classic` or `VPC`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `status` - (Optional, ForceNew) The instance status. Valid values: `Creating`, `DBInstanceClassChanging`, `DBInstanceNetTypeChanging`, `Deleting`, `EngineVersionUpgrading`, `GuardSwitching`, `HASwitching`, `Importing`, `ImportingFromOthers`, `LinkSwitching`, `MinorVersionUpgrading`, `NodeCreating`, `NodeDeleting`, `Rebooting`, `Restoring`, `Running`, `SSLModifying`, `TDEModifying`, `TempDBInstanceCreating`, `Transing`, `TransingToOthers`, `released`.
* `vpc_id` - (Optional, ForceNew) The ID of the VPC network.
* `vswitch_id` - (Optional, ForceNew) The id of the vswitch.
* `zone_id` - (Optional, ForceNew) The ID of the zone.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `instances` - A list of Mongodb Serverless Instances. Each element contains the following attributes:
  * `capacity_unit` - The read/write throughput consumed by the instance.
  * `payment_type` - The Payment type of the instance.
  * `db_instance_class` - The db instance class.
  * `db_instance_description` - The db instance description.
  * `db_instance_id` - The db instance id.
  * `db_instance_release_protection` - The db instance release protection.
  * `db_instance_storage` - The db instance storage.
  * `engine` - The database engine of the instance.
  * `engine_version` - The database version number. Valid values: `4.2`.
  * `expire_time` - The time when the subscription instance expires. The time is in the `yyyy-MM-ddTHH:mmZ` format. The time is displayed in UTC.
  * `id` - The ID of the Serverless Instance.
  * `kind_code` - Indicates the type of the instance. Valid values: `0`: physical machine. `1`: ECS. `2`: DOCKER. `18`: k8s new architecture instance.
  * `lock_mode` - The locked status of the instance.
  * `maintain_end_time` - The start time of the maintenance window. The time is in the `HH:mmZ` format. The time is displayed in UTC.
  * `maintain_start_time` - The end time of the maintenance window. The time is in the `HH:mmZ` format. The time is displayed in UTC.
  * `max_connections` - Instance maximum connections.
  * `max_iops` - The maximum IOPS of the instance.
  * `network_type` - The network type of the instance.
  * `protocol_type` - The access protocol type of the instance. Valid values: `mongodb`, `dynamodb`.
  * `resource_group_id` - The ID of the resource group.
  * `security_ip_groups` - The security ip list.
      * `security_ip_group_attribute` - The attribute of the IP whitelist. This parameter is empty by default.
      * `security_ip_group_name` - The name of the IP whitelist.
      * `security_ip_list` - The IP addresses in the whitelist.
  * `status` - The status of the instance.
  * `storage_engine` - The storage engine used by the instance.
  * `tags` - The tag of the resource.
  * `vpc_auth_mode` - Intranet secret free access mode.
  * `vpc_id` - The ID of the VPC network.
  * `vswitch_id` - The id of the vswitch.
  * `zone_id` - The ID of the zone.