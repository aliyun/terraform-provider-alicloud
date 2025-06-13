---
subcategory: "Apsara File Storage for HDFS (DFS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dfs_vsc_mount_point"
description: |-
  Provides a Alicloud Apsara File Storage for HDFS (DFS) Vsc Mount Point resource.
---

# alicloud_dfs_vsc_mount_point

Provides a Apsara File Storage for HDFS (DFS) Vsc Mount Point resource.

For information about Apsara File Storage for HDFS (DFS) Vsc Mount Point and how to use it, see [What is Vsc Mount Point](https://www.alibabacloud.com/help/en/aibaba-cloud-storage-services/latest/apsara-file-storage-for-hdfs).

-> **NOTE:** Available since v1.218.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dfs_vsc_mount_point&exampleId=4ebae226-7231-6191-23fa-73660f7cc27f694c4d4b&activeTab=example&spm=docs.r.dfs_vsc_mount_point.0.4ebae22672&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_dfs_file_system" "default" {
  space_capacity       = "1024"
  description          = "for vsc mountpoint RMC test"
  storage_type         = "PERFORMANCE"
  zone_id              = "cn-hangzhou-b"
  protocol_type        = "PANGU"
  data_redundancy_type = "LRS"
  file_system_name     = var.name
}

resource "alicloud_dfs_vsc_mount_point" "DefaultFsForRMCVscMp" {
  file_system_id = alicloud_dfs_file_system.default.id
  alias_prefix   = var.name
  description    = var.name
}
```

## Argument Reference

The following arguments are supported:
* `alias_prefix` - (Optional) Mount point alias prefix, which is used as the prefix for generating VSC mount point aliases.
* `description` - (Optional) The description of the Mount point.  The length is 0 to 100 characters.
* `file_system_id` - (Required, ForceNew) The ID of the HDFS file system resource associated with the VSC mount point.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<file_system_id>:<mount_point_id>`.
* `instances` - The collection of ECS instances on which the HDFS file system is mounted. **The current property is not available**.
  * `status` - The status of the ECS instance on which the HDFS file system is mounted.
  * `vscs` - The VSC list of mounted HDFS file systems.
    * `vsc_id` - VSC Channel primary key representation, used to retrieve the specified VSC Channel.
    * `vsc_status` - VSC Mount status.
    * `vsc_type` - The VSC type.
  * `instance_id` -The ID of the ECS instance to which the HDFS file system is mounted.
* `mount_point_id` - VSC mount point ID, which is the unique identifier of the vsc mount point and is used to access the associated HDFS file system.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vsc Mount Point.
* `delete` - (Defaults to 5 mins) Used when delete the Vsc Mount Point.
* `update` - (Defaults to 5 mins) Used when update the Vsc Mount Point.

## Import

Apsara File Storage for HDFS (DFS) Vsc Mount Point can be imported using the id, e.g.

```shell
$ terraform import alicloud_dfs_vsc_mount_point.example <file_system_id>:<mount_point_id>
```
