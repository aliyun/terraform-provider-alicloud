---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_instances"
description: |-
  Provides a list of MongoDB Instances to the user.
---

# alicloud_mongodb_instances

This data source provides the MongoDB Instances of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.13.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_mongodb_zones" "default" {
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_mongodb_zones.default.zones.0.id
}

data "alicloud_security_groups" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_mongodb_instance" "default" {
  engine_version      = "4.4"
  db_instance_class   = "mdb.shard.2x.xlarge.d"
  db_instance_storage = 20
  vswitch_id          = data.alicloud_vswitches.default.ids.0
  name                = var.name
  tags = {
    Created = "TF"
    For     = "Instance"
  }
}

data "alicloud_mongodb_instances" "ids" {
  ids = [alicloud_mongodb_instance.default.id]
}

output "mongodb_instances_id_0" {
  value = data.alicloud_mongodb_instances.ids.instances.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, List, Available since v1.53.0) A list of Instance IDs.
* `name_regex` - (Optional) A regex string to filter results by Instance name.
* `instance_type` - (Optional) The instance architecture. Default value: `replicate`. Valid values: `replicate`, `sharding`.
* `instance_class` - (Optional) The instance type.
* `availability_zone` - (Optional) The zone ID.
* `status` - (Optional, Available since v1.271.0) The instance status.
* `tags` - (Optional, Available since v1.66.0) A mapping of tags to assign to the resource.
* `enable_details` - (Optional, Bool, Available since v1.271.0) Whether to query the detailed list of resource attributes. Default value: `false`.
* `output_file` - (Optional) The name of file that can save the collection of instances after running `terraform plan`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - (Available since v1.42.0) A list of Instance names.
* `instances` -  A list of Instances. Each element contains the following attributes:
  * `id` - The instance ID.
  * `engine` - The database engine.
  * `engine_version` - The database engine version.
  * `instance_type` - The instance architecture.
  * `instance_class` - The instance type.
  * `storage` - The storage space of the instance.
  * `network_type` - The network type of the instance.
  * `availability_zone` - The zone ID of the instance.
  * `name` - The name of the instance.
  * `charge_type` - The billing method of the instance.
  * `replication` - The number of nodes in the instance.
  * `lock_mode` - The lock status of the instance.
  * `region_id` - The region ID of the instance.
  * `creation_time` - The time when the instance was created.
  * `expiration_time` - The time when the instance expires.
  * `status` - The instance status.
  * `tags` - (Available since v1.66.0) The details of the resource tags.
  * `mongos` - The mongo nodes of the instance. **Note:** `mongos` takes effect only if `instance_type` is set to `sharding`.
      * `node_id` - The ID of the Mongos node.
      * `class` - The instance type of the Mongos node.
      * `description` - The description of the Mongos node.
  * `shards` - The information of the shard node. **Note:** `shards` takes effect only if `instance_type` is set to `sharding`.
      * `node_id` - The ID of the shard node.
      * `class` - The instance type of the shard node.
      * `storage` - The storage space of the shard node.
      * `description` - The description of the shard node.
  * `restore_ranges` - (Available since v1.271.0) A list of time ranges available for point-in-time recovery. **Note:** `restore_ranges` takes effect only if `enable_details` is set to `true`.
      * `restore_type` - The restoration method.
      * `restore_begin_time` - The beginning of the recoverable time range.
      * `restore_end_time` - The end of the recoverable time range.
