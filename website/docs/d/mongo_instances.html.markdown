---
layout: "alicloud"
page_title: "Alicloud: alicloud_mongo_instances"
sidebar_current: "docs-alicloud-datasource-mongo-instances"
description: |-
    Provides a collection of MongoDB instances according to the specified filters.
---

# alicloud\_mongo\_instances

The `alicloud_mongo_instances` data source provides a collection of MongoDB instances available in Alicloud account.
Filters support regular expression for the instance name, engine or instance type.

## Example Usage

```
data "alicloud_mongo_instances" "mongo" {
  name_regex        = "dds-.+\\d+"
  instance_type     = "replicate"
  instance_class    = "dds.mongo.mid"
  availability_zone = "eu-central-1a"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the instance name.
* `instance_type` - (Optional) Type of the instance to be queried. If it is set to `sharding`, the sharded cluster instances are listed. If it is set to `replicate`, replica set instances are listed. Default value `replicate`.
* `instance_class` - (Optional) Sizing of the instance to be queried.
* `availability_zone` - (Optional) Instance availability zone.
* `output_file` - (Optional) The name of file that can save the collection of instances after running `terraform plan`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `instances` - A list of MongoDB instances. Its every element contains the following attributes:
  * `id` - The ID of the MongoDB instance.
  * `name` - The name of the MongoDB instance.
  * `charge_type` - Billing method. Value options are `PostPaid` for  Pay-As-You-Go and `PrePaid` for yearly or monthly subscription.
  * `instance_type` - Instance type. Optional values `sharding` or `replicate`.
  * `region_id` - Region ID the instance belongs to.
  * `creation_time` - Creation time of the instance in RFC3339 format.
  * `expiration_time` - Expiration time in RFC3339 format. Pay-As-You-Go instances are never expire.
  * `status` - Status of the instance.
  * `replication` - Replication factor corresponds to number of nodes. Optional values are `1` for single node and `3` for three nodes replica set.
  * `engine` - Database engine type. Supported option is `MongoDB`.
  * `engine_version` - Database engine version.
  * `network_type` - Classic network or VPC.
  * `instance_class` - Sizing of the MongoDB instance.
  * `lock_mode` - Lock status of the instance.
  * `storage` - Storage size.
  * `mongos` - Array composed of Mongos.
    * `node_id` - Mongos instance ID.
    * `description` - Mongos instance description.
    * `class` - Mongos instance specification.
  * `shards` - Array composed of shards.
    * `node_id` - Shard instance ID.
    * `description` - Shard instance description.
    * `class` - Shard instance specification.
    * `storage` - Shard disk.
  * `availability_zone` - Instance availability zone.
