---
subcategory: "Elastic Block Storage(EBS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ebs_replica_group_drill"
description: |-
  Provides a Alicloud EBS Replica Group Drill resource.
---

# alicloud_ebs_replica_group_drill

Provides a EBS Replica Group Drill resource. 

For information about Elastic Block Storage(EBS) Replica Group Drill and how to use it, see [What is Replica Group Drill](https://next.api.alibabacloud.com/document/ebs/2021-07-30/StartReplicaGroupDrill).

-> **NOTE:** Available since v1.215.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ebs_replica_group_drill&exampleId=b4f5d619-a475-5796-efa2-0281eaacf6cee43c0313&activeTab=example&spm=docs.r.ebs_replica_group_drill.0.b4f5d619a4&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_ebs_replica_group_drill" "default" {
  group_id = "pg-m1H9aaOUIGsDUwgZ"
}
```

## Argument Reference

The following arguments are supported:
* `group_id` - (Required, ForceNew) The ID of the replication group. You can use the [describediskreplicaggroups](~~ 426614 ~~) interface to query the asynchronous replication group list to obtain the value of the replication group ID input parameter.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<group_id>:<replica_group_drill_id>`.
* `replica_group_drill_id` - The first ID of the resource.
* `status` - Walkthrough status. _failed: Execution failed._failed: Cleanup failed.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Replica Group Drill.
* `delete` - (Defaults to 5 mins) Used when delete the Replica Group Drill.

## Import

EBS Replica Group Drill can be imported using the id, e.g.

```shell
$ terraform import alicloud_ebs_replica_group_drill.example <group_id>:<replica_group_drill_id>
```