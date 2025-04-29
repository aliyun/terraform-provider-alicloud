---
subcategory: "Elastic Block Storage(EBS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ebs_disk_replica_group"
description: |-
  Provides a Alicloud Elastic Block Storage(EBS) Disk Replica Group resource.
---

# alicloud_ebs_disk_replica_group

Provides a Elastic Block Storage(EBS) Disk Replica Group resource.

consistent replica group.

For information about Elastic Block Storage(EBS) Disk Replica Group and how to use it, see [What is Disk Replica Group](https://www.alibabacloud.com/help/en/elastic-compute-service/latest/creatediskreplicagroup).

-> **NOTE:** Available since v1.187.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ebs_disk_replica_group&exampleId=6d26c356-67cb-e060-9ac4-cf20ccf54b04caba09dd&activeTab=example&spm=docs.r.ebs_disk_replica_group.0.6d26c35667&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}
data "alicloud_regions" "default" {
  current = true
}
data "alicloud_ebs_regions" "default" {
  region_id = data.alicloud_regions.default.regions.0.id
}

resource "alicloud_ebs_disk_replica_group" "default" {
  source_region_id      = data.alicloud_regions.default.regions.0.id
  source_zone_id        = data.alicloud_ebs_regions.default.regions[0].zones[0].zone_id
  destination_region_id = data.alicloud_regions.default.regions.0.id
  destination_zone_id   = data.alicloud_ebs_regions.default.regions[0].zones[1].zone_id
  group_name            = var.name
  description           = var.name
  rpo                   = 900
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) The description of the consistent replication group.
* `destination_region_id` - (Required, ForceNew) The ID of the region to which the disaster recovery site belongs.
* `destination_zone_id` - (Required, ForceNew) The ID of the zone to which the disaster recovery site belongs.
* `disk_replica_group_name` - (Optional, Available since v1.245.0) Consistent replication group name.
* `one_shot` - (Optional, Available since v1.245.0) Whether to synchronize immediately. Value range:
  - true: Start data synchronization immediately.
  - false: Data Synchronization starts after the RPO time period.

Default value: false.
* `pair_ids` - (Optional, Set, Available since v1.245.0) List of replication pair IDs contained in a consistent replication group.
* `rpo` - (Optional, ForceNew, Int) The RPO value set by the consistency group in seconds. Currently only 900 seconds are supported.
* `resource_group_id` - (Optional, Computed, Available since v1.245.0) resource group ID of enterprise
* `reverse_replicate` - (Optional, Available since v1.245.0) Specifies whether to enable the reverse replication sub-feature. Valid values: true and false. Default value: true.
* `source_region_id` - (Required, ForceNew) The ID of the region to which the production site belongs.
* `source_zone_id` - (Required, ForceNew) The ID of the zone to which the production site belongs.
* `status` - (Optional, Computed) The status of the consistent replication group. Possible values:
  - invalid: invalid. This state indicates that there is an exception to the replication pair in the consistent replication group.
  - creating: creating.
  - created: created.
  - create_failed: creation failed.
  - manual_syncing: in a single synchronization. If it is the first single synchronization, this status is also displayed in the synchronization.
  - syncing: synchronization. This state is the first time data is copied asynchronously between the master and slave disks.
  - normal: normal. When data replication is completed within the current cycle of asynchronous replication, it will be in this state.
  - stopping: stopping.
  - stopped: stopped.
  - stop_failed: Stop failed.
  - Failover: failover.
  - Failed: failover completed.
  - failover_failed: failover failed.
  - Reprotection: In reverse copy operation.
  - reprotect_failed: reverse replication failed.
  - deleting: deleting.
  - delete_failed: delete failed.
  - deleted: deleted.
* `tags` - (Optional, Map, Available since v1.245.0) The tag of the resource

The following arguments will be discarded. Please use new fields as soon as possible:
* `group_name` - (Deprecated since v1.245.0). Field 'group_name' has been deprecated from provider version 1.245.0. New field 'disk_replica_group_name' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Disk Replica Group.
* `delete` - (Defaults to 5 mins) Used when delete the Disk Replica Group.
* `update` - (Defaults to 20 mins) Used when update the Disk Replica Group.

## Import

Elastic Block Storage(EBS) Disk Replica Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_ebs_disk_replica_group.example <id>
```