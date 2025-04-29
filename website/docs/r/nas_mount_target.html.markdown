---
subcategory: "File Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_mount_target"
description: |-
  Provides a Alicloud File Storage (NAS) Mount Target resource.
---

# alicloud_nas_mount_target

Provides a File Storage (NAS) Mount Target resource.

File system mount point.

For information about File Storage (NAS) Mount Target and how to use it, see [What is Mount Target](https://www.alibabacloud.com/help/en/doc-detail/27531.htm).

-> **NOTE:** Available since v1.34.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nas_mount_target&exampleId=22f8ab4c-6826-906f-09ad-67827e96eaae2128860b&activeTab=example&spm=docs.r.nas_mount_target.0.22f8ab4c68&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_nas_zones" "default" {
  file_system_type = "extreme"
}

locals {
  count_size = length(data.alicloud_nas_zones.default.zones)
  zone_id    = data.alicloud_nas_zones.default.zones[local.count_size - 1].zone_id
}

resource "alicloud_vpc" "example" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = alicloud_vpc.example.vpc_name
  cidr_block   = alicloud_vpc.example.cidr_block
  vpc_id       = alicloud_vpc.example.id
  zone_id      = local.zone_id
}

resource "alicloud_nas_file_system" "example" {
  protocol_type    = "NFS"
  storage_type     = "advance"
  file_system_type = "extreme"
  capacity         = "100"
  zone_id          = local.zone_id
}

resource "alicloud_nas_access_group" "example" {
  access_group_name = "access_group_xxx"
  access_group_type = "Vpc"
  description       = "test_access_group"
  file_system_type  = "extreme"
}

resource "alicloud_nas_mount_target" "example" {
  file_system_id    = alicloud_nas_file_system.example.id
  access_group_name = alicloud_nas_access_group.example.access_group_name
  vswitch_id        = alicloud_vswitch.example.id
  vpc_id            = alicloud_vpc.example.id
  network_type      = alicloud_nas_access_group.example.access_group_type
}
```

## Argument Reference

The following arguments are supported:
* `access_group_name` - (Optional) The name of the permission group.
* `dual_stack` - (Optional, Available since v1.247.0) Whether to create an IPv6 mount point.

Value:
  - true: create
  - false (default): not created

-> **NOTE:**  currently, only extreme NAS supports IPv6 function in various regions in mainland China, and IPv6 function needs to be turned on for this file system.

* `file_system_id` - (Required, ForceNew) The ID of the file system.
* `network_type` - (Optional, ForceNew, Available since v1.208.1) Network type.
* `security_group_id` - (Optional) The ID of the security group.
* `status` - (Optional, Computed) The current status of the Mount point, including Active and Inactive, can be used to mount the file system only when the status is Active.
* `vswitch_id` - (Optional, ForceNew) The ID of the switch.
* `vpc_id` - (Optional, ForceNew, Available since v1.208.1) VPC ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<file_system_id>:<mount_target_domain>`.
* `mount_target_domain` - The domain name of the Mount point.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Mount Target.
* `delete` - (Defaults to 5 mins) Used when delete the Mount Target.
* `update` - (Defaults to 5 mins) Used when update the Mount Target.

## Import

File Storage (NAS) Mount Target can be imported using the id, e.g.

```shell
$ terraform import alicloud_nas_mount_target.example <file_system_id>:<mount_target_domain>
```