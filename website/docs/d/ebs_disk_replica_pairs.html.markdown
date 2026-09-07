---
subcategory: "Elastic Block Storage(EBS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ebs_disk_replica_pairs"
sidebar_current: "docs-alicloud-datasource-ebs-disk-replica-pairs"
description: |-
  Provides a list of Ebs Disk Replica Pair owned by an Alibaba Cloud account.
---

# alicloud_ebs_disk_replica_pairs

This data source provides Ebs Disk Replica Pair available to the user.

-> **NOTE:** Available in 1.196.0+

## Example Usage

```terraform
data "alicloud_ebs_disk_replica_pairs" "default" {
  ids = ["${alicloud_ebs_disk_replica_pair.default.id}"]
}

output "alicloud_ebs_disk_replica_pair_example_id" {
  value = data.alicloud_ebs_disk_replica_pairs.default.pairs.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew) A list of Disk Replica Pair IDs.
* `replica_group_id` - (Optional, ForceNew) Consistent Replication Group ID, you can specify a consistent replication group ID to query the replication pairs within the group.
* `site` - (Optional, ForceNew) Get data for replication pairs where this Region is the production site or the disaster recovery site.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Disk Replica Pair IDs.
* `pairs` - A list of Disk Replica Pair Entries. Each element contains the following attributes:
  * `bandwidth` - The bandwidth for asynchronous data replication between cloud disks. The unit is Kbps. Value range:-10240 Kbps: equal to 10 Mbps.-20480 Kbps: equal to 20 Mbps.-51200 Kbps: equal to 50 Mbps.-102400 Kbps: equal to 100 Mbps.Default value: 10240.This parameter cannot be specified when the ChargeType value is POSTPAY. The system value is 0, which indicates that the disk is dynamically allocated according to data write changes during asynchronous replication.
  * `description` - The description of the asynchronous replication relationship. 2 to 256 English or Chinese characters in length and cannot start with' http:// 'or' https.
  * `rpo` - The RPO value set by the consistency group in seconds. Currently only 900 seconds are supported.
  * `replica_pair_id` - The first ID of the resource
  * `resource_group_id` - The ID of the resource group
  * `destination_disk_id` - The ID of the standby disk.
  * `destination_region_id` - The ID of the region to which the disaster recovery site belongs.
  * `destination_zone_id` - The ID of the zone to which the disaster recovery site belongs.
  * `disk_id` - The ID of the primary disk.
  * `pair_name` - The name of the asynchronous replication relationship. The length must be 2 to 128 characters in length and must start with a letter or Chinese name. It cannot start with http:// or https. It can contain Chinese, English, numbers, half-width colons (:), underscores (_), half-width periods (.), or dashes (-).
  * `payment_type` - The payment type of the resource.
  * `source_zone_id` - The ID of the zone to which the production site belongs.
  * `status` -  The status of the resource.
