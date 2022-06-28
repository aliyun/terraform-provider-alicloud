---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_image"
sidebar_current: "docs-alicloud-resource-image"
description: |-
  Provides an ECS image resource.
---

# alicloud\_image

Creates a custom image. You can then use a custom image to create ECS instances (RunInstances) or change the system disk for an existing instance (ReplaceSystemDisk).

-> **NOTE:**  If you want to create a template from an ECS instance, you can specify the instance ID (InstanceId) to create a custom image. You must make sure that the status of the specified instance is Running or Stopped. After a successful invocation, each disk of the specified instance has a new snapshot created.

-> **NOTE:**  If you want to create a custom image based on the system disk of your ECS instance, you can specify one of the system disk snapshots (SnapshotId) to create a custom image. However, the specified snapshot cannot be created on or before July 15, 2013.

-> **NOTE:**  If you want to combine snapshots of multiple disks into an image template, you can specify DiskDeviceMapping to create a custom image.

-> **NOTE:**  Available in 1.64.0+

## Example Usage

```
resource "alicloud_image" "default" {
  instance_id        = "i-bp1g6zv0ce8oghu7k***"
  image_name         = "test-image"
  description        = "test-image"
  architecture       = "x86_64"
  platform           = "CentOS"
  resource_group_id  = "rg-bp67acfmxazb4ph***"
  tags = {
    FinanceDept = "FinanceDeptJoshua"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Optional, ForceNew, Conflict with `snapshot_id ` and `disk_device_mapping `) The instance ID.
* `image_name` - (Optional) The image name. It must be 2 to 128 characters in length, and must begin with a letter or Chinese character (beginning with http:// or https:// is not allowed). It can contain digits, colons (:), underscores (_), or hyphens (-). Default value: null.
* `description` - (Optional) The description of the image. It must be 2 to 256 characters in length and must not start with http:// or https://. Default value: null.
* `snapshot_id` - (Optional, ForceNew, Conflict with `instance_id ` and `disk_device_mapping `) Specifies a snapshot that is used to create a custom image.
* `architecture` - (Optional, ForceNew) Specifies the architecture of the system disk after you specify a data disk snapshot as the data source of the system disk for creating an image. Valid values: `i386` , Default is `x86_64`.
* `platform` - (Optional, ForceNew) Specifies the operating system platform of the system disk after you specify a data disk snapshot as the data source of the system disk for creating an image. Valid values: `CentOS`, `Ubuntu`, `SUSE`, `OpenSUSE`, `RedHat`, `Debian`, `CoreOS`, `Aliyun Linux`, `Windows Server 2003`, `Windows Server 2008`, `Windows Server 2012`, `Windows 7`, Default is `Others Linux`.
* `tags` - (Optional) The tag value of an image. The value of N ranges from 1 to 20.
* `resource_group_id` - (Optional, Available in 1.115.0+) The ID of the enterprise resource group to which a custom image belongs
* `disk_device_mapping` - (Optional, ForceNew, Conflict with `snapshot_id ` and `instance_id `) Description of the system with disks and snapshots under the image.
  * `disk_type` - (Optional, ForceNew) Specifies the type of a disk in the combined custom image. If you specify this parameter, you can use a data disk snapshot as the data source of a system disk for creating an image. If it is not specified, the disk type is determined by the corresponding snapshot. Valid values: `system`, `data`,
  * `size` - (Optional, ForceNew) Specifies the size of a disk in the combined custom image, in GiB. Value range: 5 to 2000.
  * `snapshot_id` - (Optional, ForceNew) Specifies a snapshot that is used to create a combined custom image.
  * `device` - (Optional, ForceNew)Specifies the name of a disk in the combined custom image. Value range: /dev/xvda to /dev/xvdz.
* `force` - (Optional) Indicates whether to force delete the custom image, Default is `false`. 
  - true：Force deletes the custom image, regardless of whether the image is currently being used by other instances.
  - false：Verifies that the image is not currently in use by any other instances before deleting the image.
   
## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the image (until it reaches the initial `Available` status). 
* `delete` - (Defaults to 10 mins) Used when terminating the image.
   
   
## Attributes Reference
 
 The following attributes are exported:
 
* `id` - ID of the image.

## Import
 
 image can be imported using the id, e.g.

```
$ terraform import alicloud_image.default m-uf66871ape***yg1q***
```
