---
subcategory: "File Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_protocol_mount_target"
description: |-
  Provides a Alicloud File Storage (NAS) Protocol Mount Target resource.
---

# alicloud_nas_protocol_mount_target

Provides a File Storage (NAS) Protocol Mount Target resource.



For information about File Storage (NAS) Protocol Mount Target and how to use it, see [What is Protocol Mount Target](https://next.api.alibabacloud.com/document/NAS/2017-06-26/CreateProtocolMountTarget).

-> **NOTE:** Available since v1.267.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nas_protocol_mount_target&exampleId=2a808cbc-fe91-9557-53c4-917fa18144f423b24a31&activeTab=example&spm=docs.r.nas_protocol_mount_target.0.2a808cbcfe&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

resource "alicloud_vpc" "example" {
  is_default  = false
  cidr_block  = "192.168.0.0/16"
  vpc_name    = "nas-examplee1223-vpc"
  enable_ipv6 = true
}

resource "alicloud_vswitch" "example" {
  is_default   = false
  vpc_id       = alicloud_vpc.example.id
  zone_id      = "cn-beijing-i"
  cidr_block   = "192.168.3.0/24"
  vswitch_name = "nas-examplee1223-vsw2sdw-C"
}

resource "alicloud_nas_file_system" "example" {
  description      = var.name
  storage_type     = "advance_100"
  zone_id          = "cn-beijing-i"
  vpc_id           = alicloud_vpc.example.id
  capacity         = "3600"
  protocol_type    = "cpfs"
  vswitch_id       = alicloud_vswitch.example.id
  file_system_type = "cpfs"
}

resource "alicloud_nas_protocol_service" "example" {
  vpc_id         = alicloud_vpc.example.id
  protocol_type  = "NFS"
  protocol_spec  = "General"
  vswitch_id     = alicloud_vswitch.example.id
  dry_run        = false
  file_system_id = alicloud_nas_file_system.example.id
}

resource "alicloud_nas_fileset" "example" {
  file_system_path = "/examplefileset/"
  description      = "cpfs-LRS-filesetexample-wyf"
  file_system_id   = alicloud_nas_file_system.example.id
}


resource "alicloud_nas_protocol_mount_target" "default" {
  fset_id             = alicloud_nas_fileset.example.fileset_id
  description         = var.name
  vpc_id              = alicloud_vpc.example.id
  vswitch_id          = alicloud_vswitch.example.id
  access_group_name   = "DEFAULT_VPC_GROUP_NAME"
  dry_run             = false
  file_system_id      = alicloud_nas_file_system.example.id
  protocol_service_id = alicloud_nas_protocol_service.example.protocol_service_id
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_nas_protocol_mount_target&spm=docs.r.nas_protocol_mount_target.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `access_group_name` - (Optional, ForceNew, Computed) The permission group name.
Default value: DEFAULT_VPC_GROUP_NAME
* `description` - (Optional) Description of the protocol service mount target. Display as the export directory name in the console.

Limitations:
  - Length is 2~128 English or Chinese characters.
  - It must start with an uppercase or lowercase letter or a Chinese character. It cannot start with http:// or https.
  - Can contain numbers, colons (:), underscores (_), or dashes (-).
* `dry_run` - (Optional) DryRun

-> **NOTE:** This parameter only applies during resource creation, update or deletion. If modified in isolation without other property changes, Terraform will not trigger any action.

* `file_system_id` - (Required, ForceNew) The ID of the file system.
* `fset_id` - (Optional, ForceNew) The ID of the Fileset to be mounted.

Limitations:
  - The Fileset must already exist.
  - A Fileset allows only one export directory to be created.
  - Fileset and Path can and must specify only one.
* `path` - (Optional, ForceNew) The path of the CPFS directory to be mounted.

Limitations:
  - The directory must be an existing directory on the CPFS.
  - Only one export is allowed for the same directory.
  - Fileset and Path can and must specify only one.

Format:
  - 1~1024 characters in length.
  - Use UTF-8 encoding.
  - Must start and end with a forward slash (/) and root directory is/.
* `protocol_service_id` - (Required, ForceNew) Protocol Service ID
* `vswitch_id` - (Optional, ForceNew) The vSwitch ID of the protocol service mount target.
* `vswitch_ids` - (Optional, ForceNew, List) The vSwitch IDs of the protocol service mount target.
When the storage redundancy type of the file system is ZRS, if VpcId is set, the vSwitch ID of three different zones under the Vpc must be set in this field.
* `vpc_id` - (Optional, ForceNew) The VPC ID of the protocol service mount point.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<file_system_id>:<protocol_service_id>:<export_id>`.
* `create_time` - The creation time of the resource
* `export_id` - Protocol Service Mount Target ID
* `status` - Status

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Protocol Mount Target.
* `delete` - (Defaults to 5 mins) Used when delete the Protocol Mount Target.
* `update` - (Defaults to 5 mins) Used when update the Protocol Mount Target.

## Import

File Storage (NAS) Protocol Mount Target can be imported using the id, e.g.

```shell
$ terraform import alicloud_nas_protocol_mount_target.example <file_system_id>:<protocol_service_id>:<export_id>
```