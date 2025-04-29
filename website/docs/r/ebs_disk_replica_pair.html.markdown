---
subcategory: "Elastic Block Storage(EBS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ebs_disk_replica_pair"
description: |-
  Provides a Alicloud Elastic Block Storage(EBS) Disk Replica Pair resource.
---

# alicloud_ebs_disk_replica_pair

Provides a Elastic Block Storage(EBS) Disk Replica Pair resource.



For information about Elastic Block Storage(EBS) Disk Replica Pair and how to use it, see [What is Disk Replica Pair](https://www.alibabacloud.com/help/en/ecs/developer-reference/api-ebs-2021-07-30-creatediskreplicapair).

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
* `bandwidth` - (Optional, ForceNew, Int) The bandwidth for asynchronous data replication between cloud disks. The unit is Kbps. Value range:
  - 10240 Kbps: equal to 10 Mbps.
  - 20480 Kbps: equal to 20 Mbps.
  - 51200 Kbps: equal to 50 Mbps.
  - 102400 Kbps: equal to 100 Mbps.

Default value: 10240.
This parameter cannot be specified when the ChargeType value is PayAsYouGo The system value is 0, which indicates that the disk is dynamically allocated according to data write changes during asynchronous replication.
* `description` - (Optional) The description of the asynchronous replication relationship. 2 to 256 English or Chinese characters in length and cannot start with' http:// 'or' https.
* `destination_disk_id` - (Required, ForceNew) The ID of the standby disk.
* `destination_region_id` - (Required, ForceNew) The ID of the region to which the disaster recovery site belongs.
* `destination_zone_id` - (Required, ForceNew) The ID of the zone to which the disaster recovery site belongs.
* `disk_id` - (Required, ForceNew) The ID of the primary disk.
* `disk_replica_pair_name` - (Optional, Available since v1.245.0) The name of the asynchronous replication relationship. The length must be 2 to 128 characters in length and must start with a letter or Chinese name. It cannot start with http:// or https. It can contain Chinese, English, numbers, half-width colons (:), underscores (_), half-width periods (.), or dashes (-).
* `one_shot` - (Optional, Available since v1.245.0) Whether to synchronize immediately. Value range:
  - true: Start data synchronization immediately.
  - false: Data Synchronization starts after the RPO time period.

Default value: false.
* `payment_type` - (Optional, ForceNew, Computed) The payment type of the resource
* `period` - (Optional, Int) The purchase duration of the asynchronous replication relationship. This parameter is required when 'ChargeType = PrePay. The duration unit is specified by'periodunit', and the value range is:
  - When 'PeriodUnit = Week', the value range of this parameter is 1, 2, 3, and 4.
  - When 'PeriodUnit = Month', the value range of this parameter is 1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36, 48, 60.
* `period_unit` - (Optional) The unit of the purchase time of the asynchronous replication relationship. Value range:
  - Week: Week.
  - Month: Month.

Default value: Month.
* `rpo` - (Optional, ForceNew, Int) The RPO value set by the consistency group in seconds. Currently only 900 seconds are supported.
* `resource_group_id` - (Optional, Computed) The ID of the resource group
* `reverse_replicate` - (Optional, Available since v1.245.0) Specifies whether to enable the reverse replication sub-feature. Valid values: true and false. Default value: true.
* `source_zone_id` - (Required, ForceNew) The ID of the zone to which the production site belongs.
* `tags` - (Optional, Map, Available since v1.245.0) The tag of the resource

The following arguments will be discarded. Please use new fields as soon as possible:
* `pair_name` - (Deprecated since v1.245.0). Field 'pair_name' has been deprecated from provider version 1.245.0. New field 'disk_replica_pair_name' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource
* `region_id` - The region ID  of the resource
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Disk Replica Pair.
* `delete` - (Defaults to 5 mins) Used when delete the Disk Replica Pair.
* `update` - (Defaults to 20 mins) Used when update the Disk Replica Pair.

## Import

Elastic Block Storage(EBS) Disk Replica Pair can be imported using the id, e.g.

```shell
$ terraform import alicloud_ebs_disk_replica_pair.example <id>
```