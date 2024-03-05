---
subcategory: "Elastic Block Storage(EBS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ebs_replica_pair_drill"
description: |-
  Provides a Alicloud EBS Replica Pair Drill resource.
---

# alicloud_ebs_replica_pair_drill

Provides a EBS Replica Pair Drill resource. 

For information about EBS Replica Pair Drill and how to use it, see [What is Replica Pair Drill](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.215.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_ebs_replica_pair_drill" "default" {
  pair_id = "pair-cn-wwo3kjfq5001"
}
```

## Argument Reference

The following arguments are supported:
* `pair_id` - (Required, ForceNew) Copy the ID of the pair. You can call [DescribeDiskReplicaPairs](~~ 354206 ~~) to query the list of asynchronous replication pairs to obtain the replication pair ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<pair_id>:<replica_pair_drill_id>`.
* `replica_pair_drill_id` - The first ID of the resource.
* `status` - Walkthrough status. _failed: Execution failed._failed: Cleanup failed.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Replica Pair Drill.
* `delete` - (Defaults to 5 mins) Used when delete the Replica Pair Drill.

## Import

EBS Replica Pair Drill can be imported using the id, e.g.

```shell
$ terraform import alicloud_ebs_replica_pair_drill.example <pair_id>:<replica_pair_drill_id>
```