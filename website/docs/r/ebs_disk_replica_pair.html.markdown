---
subcategory: "Elastic Block Storage(EBS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ebs_disk_replica_pair"
sidebar_current: "docs-alicloud-resource-ebs-disk-replica-pair"
description: |-
  Provides a Alicloud Ebs Disk Replica Pair resource.
---

# alicloud_ebs_disk_replica_pair

Provides a Ebs Disk Replica Pair resource.

For information about Ebs Disk Replica Pair and how to use it, see [What is Disk Replica Pair](https://www.alibabacloud.com/help/en/ecs/developer-reference/api-ebs-2021-07-30-creatediskreplicapair).

-> **NOTE:** Available since v1.196.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ebs_disk_replica_pair&exampleId=e7b6f5df-04a7-24cb-0d99-60580ff8c4a41b721f93&activeTab=example&spm=docs.r.ebs_disk_replica_pair.0.e7b6f5df04&intl_lang=EN_US" target="_blank">
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

resource "alicloud_ecs_disk" "default" {
  zone_id              = data.alicloud_ebs_regions.default.regions[0].zones[0].zone_id
  category             = "cloud_essd"
  delete_auto_snapshot = "true"
  delete_with_instance = "true"
  description          = var.name
  disk_name            = var.name
  enable_auto_snapshot = "true"
  encrypted            = "true"
  size                 = "500"
  tags = {
    Created      = "TF",
    For          = "example",
    controlledBy = "ear"
  }
}

resource "alicloud_ecs_disk" "destination" {
  zone_id              = data.alicloud_ebs_regions.default.regions[0].zones[1].zone_id
  category             = "cloud_essd"
  delete_auto_snapshot = "true"
  delete_with_instance = "true"
  description          = format("%s-destination", var.name)
  disk_name            = var.name
  enable_auto_snapshot = "true"
  encrypted            = "true"
  size                 = "500"
  tags = {
    Created      = "TF",
    For          = "example",
    controlledBy = "ear"
  }
}

resource "alicloud_ebs_disk_replica_pair" "default" {
  destination_disk_id   = alicloud_ecs_disk.destination.id
  destination_region_id = data.alicloud_regions.default.regions.0.id
  payment_type          = "POSTPAY"
  destination_zone_id   = alicloud_ecs_disk.destination.zone_id
  source_zone_id        = alicloud_ecs_disk.default.zone_id
  disk_id               = alicloud_ecs_disk.default.id
  description           = var.name
}
```

## Argument Reference

The following arguments are supported:
* `bandwidth` - (Optional, ForceNew) The bandwidth for asynchronous data replication between cloud disks. The unit is Kbps. Value range:-10240 Kbps: equal to 10 Mbps.-20480 Kbps: equal to 20 Mbps.-51200 Kbps: equal to 50 Mbps.-102400 Kbps: equal to 100 Mbps.Default value: 10240.This parameter cannot be specified when the ChargeType value is POSTPAY. The system value is 0, which indicates that the disk is dynamically allocated according to data write changes during asynchronous replication.
* `description` - (Optional) The description of the asynchronous replication relationship. 2 to 256 English or Chinese characters in length and cannot start with' http:// 'or' https.
* `destination_disk_id` - (Required, ForceNew) The ID of the standby disk.
* `destination_region_id` - (Required, ForceNew) The ID of the region to which the disaster recovery site belongs.
* `destination_zone_id` - (Required, ForceNew) The ID of the zone to which the disaster recovery site belongs.
* `disk_id` - (Required, ForceNew) The ID of the primary disk.
* `pair_name` - (Optional) The name of the asynchronous replication relationship. The length must be 2 to 128 characters in length and must start with a letter or Chinese name. It cannot start with http:// or https. It can contain Chinese, English, numbers, half-width colons (:), underscores (_), half-width periods (.), or dashes (-).
* `payment_type` - (Optional, ForceNew) The payment type of the resource
* `rpo` - (Optional, ForceNew) The RPO value set by the consistency group in seconds. Currently only 900 seconds are supported.
* `source_zone_id` - (Required, ForceNew) The ID of the zone to which the production site belongs.
* `period_unit` - (Optional) The units of asynchronous replication relationship purchase length. Valid values: `Week` and `Month`. Default value: `Month`.
* `period` - (Optional) The length of the purchase for the asynchronous replication relationship. When ChargeType=PrePay, this parameter is mandatory. The unit of duration is specified by PeriodUnit and takes on a range of values. When PeriodUnit=Week, this parameter takes values in the range `1`, `2`, `3` and `4`. When PeriodUnit=Month, the parameter takes on the values `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `12`, `24`, `36`, `48`, `60`.
* `replica_pair_id` - (Optional) The first ID of the resource.


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `create_time` - The creation time of the resource
* `status` - The status of the resource
* `resource_group_id` - The ID of the resource group


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the Disk Replica Pair.
* `delete` - (Defaults to 5 mins) Used when delete the Disk Replica Pair.
* `update` - (Defaults to 5 mins) Used when update the Disk Replica Pair.

## Import

Ebs Disk Replica Pair can be imported using the id, e.g.

```shell
$ terraform import alicloud_ebs_disk_replica_pair.example <id>
```