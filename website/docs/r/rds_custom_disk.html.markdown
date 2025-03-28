---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_custom_disk"
description: |-
  Provides a Alicloud RDS Custom Disk resource.
---

# alicloud_rds_custom_disk

Provides a RDS Custom Disk resource.

RDS User dedicated host disk.

For information about RDS Custom Disk and how to use it, see [What is Custom Disk](https://next.api.alibabacloud.com/document/Rds/2014-08-15/CreateRCDisk).

-> **NOTE:** Available since v1.247.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

variable "region_id" {
  default = "cn-beijing"
}


resource "alicloud_rds_custom_disk" "default" {
  description          = "zcc测试用例"
  zone_id              = "cn-beijing-i"
  size                 = "40"
  performance_level    = "PL1"
  instance_charge_type = "Postpaid"
  disk_category        = "cloud_essd"
  disk_name            = "custom_disk_001"
  auto_renew           = false
  period               = "1"
  auto_pay             = true
  period_unit          = "1"
}
```

## Argument Reference

The following arguments are supported:
* `auto_pay` - (Optional) Whether to pay automatically. Value range:
  - `true`: automatic payment. You need to ensure that your account balance is sufficient.
  - `false`: only orders are generated without deduction.



-> **NOTE:**  The default value is true. If the balance of your payment method is insufficient, you can set the AutoPay parameter to false. In this case, unpaid orders will be generated. You can log on to the RDS management console to pay by yourself.

-> **NOTE:** >

* `auto_renew` - (Optional) Whether to automatically renew. This parameter is passed in only when you create a data disk. Valid values:
  - `true`: Yes
  - `false`: No

-> **NOTE:**  When purchasing by month, the automatic renewal period is 1 month.
When purchasing by year, the automatic renewal period is 1 year.
* `description` - (Optional, ForceNew) The disk description. It must be 2 to 256 characters in length and cannot start with 'http:// 'or 'https.
Default value: empty.
* `disk_category` - (Required) The type of the data disk. Value range:
  - `cloud` (default): a normal cloud disk.
  - `cloud_efficiency`: The ultra cloud disk.
  - `cloud_ssd`:SSD cloud disk.
  - `cloud_essd`: the ESSD cloud disk.
  - `cloud_auto`:ESSD AutoPL cloud disk.
  - `Cloud_essd_entry`: the ESSD Entry disk.
  - `Elastic_ephemeral_disk_standard`: Elastic temporary disk-standard version.
  - `Elastic_ephemeral_disk_premium`: Elastic temporary disk-Pro version.
* `disk_name` - (Optional, ForceNew) The disk name. It can be 2 to 128 characters in length. It supports letters in Unicode (including English, Chinese, and numbers). Can contain a colon (:), an underscore (_), a period (.), or a dash (-).
Default value: empty.
* `dry_run` - (Optional) Whether to pre-check the instance creation operation. Valid values:
  - `true`: The PreCheck operation is performed without creating an instance. Check items include request parameters, request formats, business restrictions, and inventory.
  - `false` (default): Sends a normal request and directly creates an instance after the check is passed.
* `instance_charge_type` - (Optional) The Payment type. Only `Postpaid`: Pay-As-You-Go is supported.
* `performance_level` - (Optional, ForceNew) When creating an ESSD cloud disk, set the performance level of the disk. Value range:
  - `PL0`: The maximum random read/write IOPS 10000 for a single disk.
  - `PL1` (default): The maximum number of random read/write IOPS 50000 for a single disk.
  - `PL2`: maximum random read/write IOPS 100000 for a single disk.
  - `PL3`: The maximum random read/write IOPS 1 million for a single disk.

For more information about how to select an ESSD performance level, see [ESSD cloud disk](~~ 122389 ~~).
* `period` - (Optional, Int) Reserved parameters, no need to fill in.
* `period_unit` - (Optional) Reserved parameters, no need to fill in.
* `size` - (Required, Int) Capacity size. Unit: GiB. You must pass in a parameter value for this parameter. Value range:
  - `cloud`:5~2,000.
  - `cloud_efficiency`:20 to 32,768.
  - `cloud_ssd`:20 to 32,768.
  - `cloud_essd`: The specific value range is related to the value of PerformanceLevel.
  - PL0:1~65,536.
  - PL1:20~65,536.
  - PL2:461~65,536.
  - PL3:1,261~65,536.
  - `cloud_auto`:1~65,536.
  - `Cloud_essd_entry`:10 to 32,768.
  - `Elastic_ephemeral_disk_standard`:64-8,192.
  - `Elastic_ephemeral_disk_premium`:64 to 8,192.

If you specify the 'SnapshotId' parameter, the 'SnapshotId' parameter and the 'Size' parameter have the following limitations:
  - If the snapshot capacity corresponding to the 'SnapshotId' parameter is greater than the set 'Size' parameter value, the actual size of the cloud disk created is the size of the specified snapshot.
  - If the snapshot capacity corresponding to the 'SnapshotId' parameter is less than the set 'Size' parameter value, the size of the cloud disk created is the specified 'Size' parameter value.
* `snapshot_id` - (Optional) The snapshot used to create the cloud disk. Snapshots made on or before July 15, 2013 cannot be used to create cloud disks. The 'SnapshotId' parameter and the 'Size' parameter have the following limitations:
  - If the snapshot capacity corresponding to the 'SnapshotId' parameter is greater than the set 'Size' parameter value, the actual size of the cloud disk created is the size of the specified snapshot.
  - If the snapshot capacity corresponding to the 'SnapshotId' parameter is less than the set 'Size' parameter value, the size of the cloud disk created is the specified 'Size' parameter value.
  - Snapshots are not supported for creating elastic temporary disks.
* `type` - (Optional) The method of expanding the disk. Value range:
offline (default): offline expansion. After the expansion, the instance must be restarted to take effect.
online: online expansion, which can be completed without restarting the instance.
* `zone_id` - (Required, ForceNew) The zone ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Creation time.
* `region_id` - The region ID. You can view the region ID through the DescribeRegions interface.
* `resource_group_id` - The ID of the resource group to which the disk belongs.
* `status` - Disk status. Value Description:_use: In use.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Custom Disk.
* `delete` - (Defaults to 5 mins) Used when delete the Custom Disk.
* `update` - (Defaults to 5 mins) Used when update the Custom Disk.

## Import

RDS Custom Disk can be imported using the id, e.g.

```shell
$ terraform import alicloud_rds_custom_disk.example <id>
```