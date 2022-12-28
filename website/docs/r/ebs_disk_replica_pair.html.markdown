---
subcategory: "Elastic Block Storage(EBS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ebs_disk_replica_pair"
sidebar_current: "docs-alicloud-resource-ebs-disk-replica-pair"
description: |-
  Provides a Alicloud Ebs Disk Replica Pair resource.
---

# alicloud\_ebs\_disk\_replica\_pair

Provides a Ebs Disk Replica Pair resource.

For information about Ebs Disk Replica Pair and how to use it, see [What is Disk Replica Pair](https://www.alibabacloud.com/help/en/elastic-compute-service/latest/CreateDiskReplicaPair).

-> **NOTE:** Available in v1.196.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecs_disk" "default" {
  zone_id              = "cn-hangzhou-onebox-nebula"
  category             = "cloud_essd"
  delete_auto_snapshot = "true"
  delete_with_instance = "true"
  description          = "Test For Terraform"
  disk_name            = var.name
  enable_auto_snapshot = "true"
  encrypted            = "true"
  size                 = "500"
  tags = {
    Created     = "TF"
    Environment = "Acceptance-test"
  }
}

resource "alicloud_ecs_disk" "defaultone" {
  zone_id              = "cn-hangzhou-onebox-nebula-b"
  category             = "cloud_essd"
  delete_auto_snapshot = "true"
  delete_with_instance = "true"
  description          = "Test For Terraform"
  disk_name            = var.name
  enable_auto_snapshot = "true"
  encrypted            = "true"
  size                 = "500"
  tags = {
    Created     = "TF"
    Environment = "Acceptance-test"
  }
}

resource "alicloud_ebs_disk_replica_pair" "default" {
  destination_disk_id   = alicloud_ecs_disk.default.id
  destination_region_id = "cn-hangzhou-onebox-nebula"
  bandwidth             = 10240
  destination_zone_id   = "cn-hangzhou-onebox-nebula-e"
  source_zone_id        = "cn-hangzhou-onebox-nebula-b"
  disk_id               = alicloud_ecs_disk.defaultone.id
  description           = "abc"
}
```

## Argument Reference

The following arguments are supported:
* `bandwidth` - (ForceNew,Optional) The bandwidth for asynchronous data replication between cloud disks. The unit is Kbps. Value range:-10240 Kbps: equal to 10 Mbps.-20480 Kbps: equal to 20 Mbps.-51200 Kbps: equal to 50 Mbps.-102400 Kbps: equal to 100 Mbps.Default value: 10240.This parameter cannot be specified when the ChargeType value is POSTPAY. The system value is 0, which indicates that the disk is dynamically allocated according to data write changes during asynchronous replication.
* `description` - (Optional) The description of the asynchronous replication relationship. 2 to 256 English or Chinese characters in length and cannot start with' http:// 'or' https.
* `destination_disk_id` - (Required,ForceNew) The ID of the standby disk.
* `destination_region_id` - (Required,ForceNew) The ID of the region to which the disaster recovery site belongs.
* `destination_zone_id` - (Required,ForceNew) The ID of the zone to which the disaster recovery site belongs.
* `disk_id` - (Required,ForceNew) The ID of the primary disk.
* `pair_name` - (Optional) The name of the asynchronous replication relationship. The length must be 2 to 128 characters in length and must start with a letter or Chinese name. It cannot start with http:// or https. It can contain Chinese, English, numbers, half-width colons (:), underscores (_), half-width periods (.), or dashes (-).
* `payment_type` - (ForceNew,Optional) The payment type of the resource
* `rpo` - (ForceNew,Optional) The RPO value set by the consistency group in seconds. Currently only 900 seconds are supported.
* `source_zone_id` - (Required,ForceNew) The ID of the zone to which the production site belongs.
* `period_unit` - (Optional) The units of asynchronous replication relationship purchase length. Valid values: `Week` and `Month`. Default value: `Month`.
* `period` - (Optional) The length of the purchase for the asynchronous replication relationship. When ChargeType=PrePay, this parameter is mandatory. The unit of duration is specified by PeriodUnit and takes on a range of values. When PeriodUnit=Week, this parameter takes values in the range `1`, `2`, `3` and `4`. When PeriodUnit=Month, the parameter takes on the values `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `12`, `24`, `36`, `48`, `60`.


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `create_time` - The creation time of the resource
* `replica_pair_id` - The first ID of the resource
* `status` - The status of the resource
* `resource_group_id` - The ID of the resource group


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the Disk Replica Pair.
* `delete` - (Defaults to 5 mins) Used when delete the Disk Replica Pair.
* `update` - (Defaults to 5 mins) Used when update the Disk Replica Pair.

## Import

Ebs Disk Replica Pair can be imported using the id, e.g.

```shell
$ terraform import alicloud_ebs_disk_replica_pair.example <id>
```