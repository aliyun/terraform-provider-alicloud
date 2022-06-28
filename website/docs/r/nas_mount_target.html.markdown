---
subcategory: "Network Attached Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_mount_target"
sidebar_current: "docs-alicloud-resource-nas-mount-target"
description: |-
  Provides a Alicloud NAS MountTarget resource.
---

# alicloud\_nas_mount_target

Provides a NAS Mount Target resource.
For information about NAS Mount Target and how to use it, see [Manage NAS Mount Targets](https://www.alibabacloud.com/help/en/doc-detail/27531.htm).

-> **NOTE**: Available in v1.34.0+.

-> **NOTE**: Currently this resource support create a mount point in a classic network only when current region is China mainland regions.

-> **NOTE**: You must grant NAS with specific RAM permissions when creating a classic mount targets,
and it only can be achieved by creating a classic mount target mannually.
See [Add a mount point](https://www.alibabacloud.com/help/doc-detail/60431.htm) and [Why do I need RAM permissions to create a mount point in a classic network](https://www.alibabacloud.com/help/faq-detail/42176.htm).

## Example Usage

Basic Usage

```terraform
resource "alicloud_nas_file_system" "example" {
  protocol_type = "NFS"
  storage_type  = "Performance"
  description   = "test file system"
}

resource "alicloud_nas_access_group" "example" {
  access_group_name = "test_name"
  access_group_type = "Classic"
  description       = "test access group"
}

resource "alicloud_nas_mount_target" "example" {
  file_system_id    = alicloud_nas_file_system.example.id
  access_group_name = alicloud_nas_access_group.example.access_group_name
}
```

## Argument Reference

The following arguments are supported:

* `file_system_id` - (Required, ForceNew) The ID of the file system.
* `access_group_name` - (Required) The name of the permission group that applies to the mount target.
* `vswitch_id` - (Optional, ForceNew) The ID of the VSwitch in the VPC where the mount target resides.
* `status` - (Optional) Whether the MountTarget is active. The status of the mount target. Valid values: `Active` and `Inactive`, Default value is `Active`. Before you mount a file system, make sure that the mount target is in the Active state.
* `security_group_id` - (Optional, ForceNew, Available in v1.95.0+.) The ID of security group.

## Attributes Reference

The following attributes are exported:

*`id` - This ID of this resource. It is formatted to `<file_system_id>:<mount_target_domain>`. Before version 1.95.0, the value is `<mount_target_domain>`.
* `mount_target_domain` - The IPv4 domain name of the mount target. **NOTE:** Available in v1.161.0+.
### Timeouts

-> **NOTE:** Available in v1.153.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 40 mins) Used when create the mount target (until it reaches the initial `Active` status).
* `delete` - (Defaults to 40 mins) Used when delete the mount target.


## Import

NAS MountTarget  can be imported using the id, e.g.

```
$ terraform import alicloud_nas_mount_target.foo 192094b415:192094b415-luw38.cn-beijing.nas.aliyuncs.com
```
