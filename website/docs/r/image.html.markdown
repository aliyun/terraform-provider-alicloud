---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_image"
description: |-
  Provides a Alicloud ECS Image resource.
---

# alicloud_image

Provides a ECS Image resource.

-> **NOTE:**  If you want to create a template from an ECS instance, you can specify the instance ID (InstanceId) to create a custom image. You must make sure that the status of the specified instance is Running or Stopped. After a successful invocation, each disk of the specified instance has a new snapshot created.

-> **NOTE:**  If you want to create a custom image based on the system disk of your ECS instance, you can specify one of the system disk snapshots (SnapshotId) to create a custom image. However, the specified snapshot cannot be created on or before July 15, 2013.

-> **NOTE:**  If you want to combine snapshots of multiple disks into an image template, you can specify DiskDeviceMapping to create a custom image.

For information about ECS Image and how to use it, see [What is Image](https://www.alibabacloud.com/help/en/ecs/developer-reference/api-ecs-2014-05-26-createimage).

-> **NOTE:** Available since v1.64.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_image&exampleId=8d8342e8-72da-4f9c-a1aa-e5764f0cddcbad72b16e&activeTab=example&spm=docs.r.image.0.8d8342e872&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_zones" "default" {
  available_resource_creation = "Instance"
}

data "alicloud_instance_types" "default" {
  instance_type_family = "ecs.sn1ne"
}

data "alicloud_images" "default" {
  name_regex = "^ubuntu_18.*64"
  owners     = "system"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name   = "terraform-example"
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  availability_zone          = data.alicloud_zones.default.zones.0.id
  instance_name              = "terraform-example"
  security_groups            = [alicloud_security_group.default.id]
  vswitch_id                 = alicloud_vswitch.default.id
  instance_type              = data.alicloud_instance_types.default.ids[0]
  image_id                   = data.alicloud_images.default.ids[0]
  internet_max_bandwidth_out = 10
}
data "alicloud_resource_manager_resource_groups" "default" {}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_image" "default" {
  instance_id       = alicloud_instance.default.id
  image_name        = "terraform-example-${random_integer.default.result}"
  description       = "terraform-example"
  architecture      = "x86_64"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  tags = {
    FinanceDept = "FinanceDeptJoshua"
  }
}
```

## Argument Reference

The following arguments are supported:
* `architecture` - (Optional, ForceNew, Computed) The system architecture of the system disk. If you specify a data disk snapshot to create the system disk of the custom image, you must use Architecture to specify the system architecture of the system disk. Valid values: `i386`, `x86\_64`, `arm64`. Default value: `x86\_64`.
* `boot_mode` - (Optional, Computed, Available since v1.227.0) The new boot mode of the image. Valid values:

  *   BIOS: Basic Input/Output System (BIOS)

  *   UEFI: Unified Extensible Firmware Interface (UEFI)

  *   UEFI-Preferred: BIOS and UEFI

-> **NOTE:**   Before you change the boot mode, we recommend that you obtain the boot modes supported by the image. If you specify an unsupported boot mode for the image, ECS instances that use the image cannot start as expected. If you do not know which boot modes are supported by the image, we recommend that you use the image check feature to perform a check. For information about the image check feature, see [Overview](https://www.alibabacloud.com/help/en/doc-detail/439819.html).

-> **NOTE:**   For information about the UEFI-Preferred boot mode, see [Best practices for ECS instance boot modes](https://www.alibabacloud.com/help/en/doc-detail/2244655.html).

* `description` - (Optional) The new description of the custom image. The description must be 2 to 256 characters in length It cannot start with `http://` or `https://`. This parameter is empty by default, which specifies that the original description is retained. 
* `detection_strategy` - (Optional, Available since v1.227.0) The mode in which to check the custom image. If you do not specify this parameter, the image is not checked. Only the standard check mode is supported.

-> **NOTE:**   This parameter is supported for most Linux and Windows operating system versions. For information about image check items and operating system limits for image check, see [Overview of image check](https://www.alibabacloud.com/help/en/doc-detail/439819.html) and [Operating system limits for image check](https://www.alibabacloud.com/help/en/doc-detail/475800.html).

* `disk_device_mapping` - (Optional, ForceNew, Computed) Snapshot information for the image See [`disk_device_mapping`](#disk_device_mapping) below.
* `features` - (Optional, Computed, Available since v1.227.0) Features See [`features`](#features) below.
* `force` - (Optional) Whether to perform forced deletion. Value range:
  - true: forcibly deletes the custom image, ignoring whether the current image is used by other instances.
  - false: The custom image is deleted normally. Before deleting the custom image, check whether the current image is used by other instances.

  Default value: false
* `delete_auto_snapshot` - (Optional, Available since 1.136.0) Not the public attribute and it used to automatically delete dependence snapshots while deleting the image.
* `image_family` - (Optional, Available since v1.227.0) The name of the image family. The name must be 2 to 128 characters in length. It must start with a letter and cannot start with acs: or aliyun. It cannot contain http:// or https://. It can contain letters, digits, periods (.), colons (:), underscores (\_), and hyphens (-). By default, this parameter is empty. 
* `image_name` - (Optional) The name of the custom image. The name must be 2 to 128 characters in length. It must start with a letter and cannot start with acs: or aliyun. It cannot contain http:// or https://. It can contain letters, digits, periods (.), colons (:), underscores (\_), and hyphens (-). By default, this parameter is empty. In this case, the original name is retained. 
* `image_version` - (Optional, ForceNew, Available since v1.227.0) The image version.

-> **NOTE:**  If you specify an instance by configuring `InstanceId`, and the instance uses an Alibaba Cloud Marketplace image or a custom image that is created from an Alibaba Cloud Marketplace image, you must leave this parameter empty or set this parameter to the value of ImageVersion of the instance.

* `instance_id` - (Optional) The instance ID. 
* `license_type` - (Optional, Available since v1.227.0) The type of the license that is used to activate the operating system after the image is imported. Set the value to BYOL. BYOL: The license that comes with the source operating system is used. When you use the BYOL license, make sure that your license key is supported by Alibaba Cloud. 
* `platform` - (Optional, ForceNew, Computed) The operating system distribution for the system disk in the custom image. If you specify a data disk snapshot to create the system disk of the custom image, use Platform to specify the operating system distribution for the system disk. Valid values: `Aliyun`, `Anolis`, `CentOS`, `Ubuntu`, `CoreOS`, `SUSE`, `Debian`, `OpenSUSE`, `FreeBSD`, `RedHat`, `Kylin`, `UOS`, `Fedora`, `Fedora CoreOS`, `CentOS Stream`, `AlmaLinux`, `Rocky Linux`, `Gentoo`, `Customized Linux`, `Others Linux`, `Windows Server 2022`, `Windows Server 2019`, `Windows Server 2016`, `Windows Server 2012`, `Windows Server 2008`, `Windows Server 2003`. Default value: `Others Linux`. 
* `resource_group_id` - (Optional, Computed, Available since 1.115.0) The ID of the resource group to which to assign the custom image. If you do not specify this parameter, the image is assigned to the default resource group.

-> **NOTE:**   If you call the CreateImage operation as a Resource Access Management (RAM) user who does not have the permissions to manage the default resource group and do not specify `ResourceGroupId`, the `Forbbiden: User not authorized to operate on the specified resource` error message is returned. You must specify the ID of a resource group that the RAM user has the permissions to manage or grant the RAM user the permissions to manage the default resource group before you call the CreateImage operation again.

* `snapshot_id` - (Optional) The ID of the snapshot that you want to use to create the custom image. 
* `tags` - (Optional, Map) The tag

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated since v1.227.0). Field 'name' has been deprecated from provider version 1.227.0. New field 'image_name' instead.

### `disk_device_mapping`

The disk_device_mapping supports the following:
* `device` - (Optional, ForceNew, Computed) The device name of disk N in the custom image. Valid values:
  - For disks other than basic disks, such as standard SSDs, ultra disks, and enhanced SSDs (ESSDs), the valid values range from /dev/vda to /dev/vdz in alphabetical order.
  - For basic disks, the valid values range from /dev/xvda to /dev/xvdz in alphabetical order.
* `disk_type` - (Optional, ForceNew, Computed) The type of disk N in the custom image. You can specify this parameter to create the system disk of the custom image from a data disk snapshot. If you do not specify this parameter, the disk type is determined by the corresponding snapshot. Valid values:
  - system: system disk. You can specify only one snapshot to use to create the system disk in the custom image.
  - data: data disk. You can specify up to 16 snapshots to use to create data disks in the custom image.
* `size` - (Optional, ForceNew, Computed) The size of disk N in the custom image. Unit: GiB. The valid values and default value of DiskDeviceMapping.N.Size vary based on the value of DiskDeviceMapping.N.SnapshotId.
  - If no corresponding snapshot IDs are specified in the value of DiskDeviceMapping.N.SnapshotId, DiskDeviceMapping.N.Size has the following valid values and default values:
    *   For basic disks, the valid values range from 5 to 2000, and the default value is 5.
    *   For other disks, the valid values range from 20 to 32768, and the default value is 20.
  - If a corresponding snapshot ID is specified in the value of DiskDeviceMapping.N.SnapshotId, the value of DiskDeviceMapping.N.Size must be greater than or equal to the size of the specified snapshot. The default value of DiskDeviceMapping.N.Size is the size of the specified snapshot.
* `snapshot_id` - (Optional, ForceNew, Computed) The ID of snapshot N to use to create the custom image. .

### `features`

The features supports the following:
* `nvme_support` - (Optional, ForceNew, Computed) Specifies whether to support the Non-Volatile Memory Express (NVMe) protocol. Valid values:
  - supported: The image supports NVMe. Instances created from this image also support NVMe.
  - unsupported: The image does not support NVMe. Instances created from this image do not support NVMe.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The create time
* `disk_device_mapping` - Snapshot information for the image
  * `format` - Image format.
  * `import_oss_object` - Import the object of the OSS to which the image file belongs.
  * `import_oss_bucket` - Import the bucket of the OSS to which the image belongs.
  * `progress` - Copy the progress of the task.
  * `remain_time` - For an image being replicated, return the remaining time of the replication task, in seconds.
* `status` - The status of the image. By default, if you do not specify this parameter, only images in the Available state are returned. 

  Default value: Available. You can specify multiple values for this parameter. Separate the values with commas (,).


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the Image.
* `delete` - (Defaults to 10 mins) Used when delete the Image.
* `update` - (Defaults to 10 mins) Used when update the Image.

## Import

ECS Image can be imported using the id, e.g.

```shell
$ terraform import alicloud_image.example <id>
```