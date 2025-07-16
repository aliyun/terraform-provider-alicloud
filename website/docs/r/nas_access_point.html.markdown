---
subcategory: "File Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_access_point"
description: |-
  Provides a Alicloud File Storage (NAS) Access Point resource.
---

# alicloud_nas_access_point

Provides a File Storage (NAS) Access Point resource.



For information about NAS Access Point and how to use it, see [What is Access Point](https://www.alibabacloud.com/help/zh/nas/developer-reference/api-nas-2017-06-26-createaccesspoint).

-> **NOTE:** Available since v1.224.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nas_access_point&exampleId=1b819cd0-58c8-ea04-c8f2-4408a7538987a0b2b99c&activeTab=example&spm=docs.r.nas_access_point.0.1b819cd058&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "azone" {
  default = "cn-hangzhou-g"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_vpc" "defaultkyVC70" {
  cidr_block  = "172.16.0.0/12"
  description = "接入点测试noRootDirectory"
}

resource "alicloud_vswitch" "defaultoZAPmO" {
  vpc_id     = alicloud_vpc.defaultkyVC70.id
  zone_id    = data.alicloud_zones.default.zones.0.id
  cidr_block = "172.16.0.0/24"
}

resource "alicloud_nas_access_group" "defaultBbc7ev" {
  access_group_type = "Vpc"
  access_group_name = "${var.name}-${random_integer.default.result}"
  file_system_type  = "standard"
}

resource "alicloud_nas_file_system" "defaultVtUpDh" {
  storage_type     = "Performance"
  zone_id          = var.azone
  encrypt_type     = "0"
  protocol_type    = "NFS"
  file_system_type = "standard"
  description      = "AccessPointnoRootDirectory"
}

resource "alicloud_nas_access_point" "default" {
  vpc_id            = alicloud_vpc.defaultkyVC70.id
  access_group      = alicloud_nas_access_group.defaultBbc7ev.access_group_name
  vswitch_id        = alicloud_vswitch.defaultoZAPmO.id
  file_system_id    = alicloud_nas_file_system.defaultVtUpDh.id
  access_point_name = var.name
  posix_user {
    posix_group_id = "123"
    posix_user_id  = "123"
  }
  root_path_permission {
    owner_group_id = "1"
    owner_user_id  = "1"
    permission     = "0777"
  }
}
```

## Argument Reference

The following arguments are supported:
* `access_group` - (Required) The name of the permission group.
* `access_point_name` - (Optional) The name of the access point.
* `enabled_ram` - (Optional, Bool) Specifies whether to enable the RAM policy. Default value: `false`. Valid values:
  - `true`: The RAM policy is enabled.
  - `false`: The RAM policy is disabled.
* `file_system_id` - (Required, ForceNew) The ID of the file system.
* `posix_user` - (Optional, ForceNew, Set) The Posix user. See [`posix_user`](#posix_user) below.
* `root_path` - (Optional, ForceNew) The root directory of the access point.
* `root_path_permission` - (Optional, ForceNew, Set) Root permissions. See [`root_path_permission`](#root_path_permission) below.
* `vswitch_id` - (Required, ForceNew) The vSwitch ID.
* `vpc_id` - (Required, ForceNew) The ID of the VPC.

### `posix_user`

The posix_user supports the following:
* `posix_group_id` - (Optional, ForceNew, Int) The ID of the Posix user group.
* `posix_user_id` - (Optional, ForceNew, Int) The Posix user ID.

### `root_path_permission`

The root_path_permission supports the following:
* `owner_group_id` - (Optional, ForceNew, Int) The ID of the primary user group.
* `owner_user_id` - (Optional, ForceNew, Int) The owner user ID.
* `permission` - (Optional, ForceNew) The Portable Operating System Interface for UNIX (POSIX) permission.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<file_system_id>:<access_point_id>`.
* `access_point_id` - The ID of the access point.
* `create_time` - The time when the access point was created.
* `posix_user` - The Posix user.
  * `posix_secondary_group_ids` - The ID of the second user group.
* `region_id` - (Available since v1.254.0) The region ID.
* `status` - The status of the access point.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Access Point.
* `delete` - (Defaults to 5 mins) Used when delete the Access Point.
* `update` - (Defaults to 5 mins) Used when update the Access Point.

## Import

File Storage (NAS) Access Point can be imported using the id, e.g.

```shell
$ terraform import alicloud_nas_access_point.example <file_system_id>:<access_point_id>
```
