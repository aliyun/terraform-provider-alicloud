---
subcategory: "Elastic Block Storage(EBS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ebs_replica_pair_drill"
description: |-
  Provides a Alicloud EBS Replica Pair Drill resource.
---

# alicloud_ebs_replica_pair_drill

Provides a EBS Replica Pair Drill resource. 

For information about Elastic Block Storage(EBS) Replica Pair Drill and how to use it, see [What is Replica Pair Drill](https://next.api.alibabacloud.com/document/ebs/2021-07-30/StartPairDrill).

-> **NOTE:** Available since v1.215.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ebs_replica_pair_drill&exampleId=0eaa5d39-60ab-458f-69c8-863f7a8e424fdcfb6b61&activeTab=example&spm=docs.r.ebs_replica_pair_drill.0.0eaa5d3960&intl_lang=EN_US" target="_blank">
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

resource "alicloud_ebs_replica_pair_drill" "default" {
  pair_id = "pair-cn-wwo3kjfq5001"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ebs_replica_pair_drill&spm=docs.r.ebs_replica_pair_drill.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `pair_id` - (Required, ForceNew) Copy the ID of the pair. You can call [DescribeDiskReplicaPairs](~~ 354206 ~~) to query the list of asynchronous replication pairs to obtain the replication pair ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<pair_id>:<replica_pair_drill_id>`.
* `replica_pair_drill_id` - The first ID of the resource.
* `status` - Walkthrough status. _failed: Execution failed._failed: Cleanup failed.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Replica Pair Drill.
* `delete` - (Defaults to 5 mins) Used when delete the Replica Pair Drill.

## Import

EBS Replica Pair Drill can be imported using the id, e.g.

```shell
$ terraform import alicloud_ebs_replica_pair_drill.example <pair_id>:<replica_pair_drill_id>
```