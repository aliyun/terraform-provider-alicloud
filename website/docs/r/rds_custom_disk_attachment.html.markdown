---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_custom_disk_attachment"
description: |-
  Provides a Alicloud RDS Custom Disk Attachment resource.
---

# alicloud_rds_custom_disk_attachment

Provides a RDS Custom Disk Attachment resource.

Operating cloud disk mount and unmount resources.

For information about RDS Custom Disk Attachment and how to use it, see [What is Custom Disk Attachment](https://next.api.alibabacloud.com/document/Rds/2014-08-15/AttachRCDisk).

-> **NOTE:** Available since v1.275.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-hangzhou-i"
}

data "alicloud_security_groups" "default" {
  vpc_id     = data.alicloud_vpcs.default.ids.0
  name_regex = "default-NODELETING"
}

resource "alicloud_rds_custom" "default" {
  zone_id              = data.alicloud_vswitches.default.zone_id
  instance_charge_type = "PostPaid"
  vswitch_id           = data.alicloud_vswitches.default.ids.0
  amount               = "1"
  security_group_ids   = [data.alicloud_security_groups.default.ids.0]
  system_disk {
    size = "40"
  }
  force         = true
  instance_type = "mysql.x4.xlarge.6cm"
  spot_strategy = "NoSpot"
}

resource "alicloud_rds_custom_disk" "default" {
  zone_id       = data.alicloud_vswitches.default.zone_id
  size          = "40"
  disk_category = "cloud_ssd"
  auto_pay      = true
  disk_name     = "ran_disk_attach"
}

resource "alicloud_rds_custom_disk_attachment" "default" {
  instance_id = alicloud_rds_custom.default.id
  disk_id     = alicloud_rds_custom_disk.default.id
}
```

## Argument Reference

The following arguments are supported:
* `delete_with_instance` - (Optional, ForceNew) Whether the disk is released together with the instance when the instance is released. Value range:
true: Release.
false: Do not release. The disk is converted to a pay-as-you-go data disk and is retained.
When Setting this parameter, you need to pay attention:
After the DeleteWithInstance is set to false, once the instance is under security control, the value "LockReason" : "security" is marked in OperationLocks. When the instance is released, the disk will be ignored and released at the same time.
If the target disk to be attached is an elastic temporary disk, you must set the DeleteWithInstance parameter to true.
This parameter is not supported for cloud disks with the multi-Mount feature enabled.
* `disk_id` - (Required, ForceNew) The ID of the cloud disk to be mounted. The cloud disk ('DiskId') and the instance ('InstanceId') must be in the same zone.
* `instance_id` - (Required, ForceNew) The ID of the target RDS Custom instance.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<disk_id>:<instance_id>`.
* `region_id` - The region ID of the resource.
* `status` - The status of the disk.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Custom Disk Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Custom Disk Attachment.

## Import

RDS Custom Disk Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_rds_custom_disk_attachment.example <disk_id>:<instance_id>
```
