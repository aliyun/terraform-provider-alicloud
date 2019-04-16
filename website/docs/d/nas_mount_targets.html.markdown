---
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_mount_targets"
sidebar_current: "docs-alicloud-datasource-nas-mount-targets"
description: |-
  Provides a list of mns topic subscriptions available to the user.
---

# alicloud\_nas_mount_targets

This data source provides MountTargets available to the user.

-> NOTE: Available in 1.35.0+

## Example Usage

```
data "alicloud_nas_mount_targets" "mt" {
  file_system_id = "1a2sc4d"
  access_group_name = "tf-testAccNasConfig"
}

output "alicloud_nas_mount_targets_id" {
  value = "${data.alicloud_nas_mount_targets.alicloud_nas_mount_targets_ds.mount_targets.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `file_system_id` - (Required ForceNew) The ID of the FileSystem that owns the MountTarget.
* `access_group_name` - (Optional) Filter results by a specific AccessGroupName.
* `type` - (Optional) Filter results by a specific NetworkType.
* `mount_target_domain` - (Optional) Filter results by a specific MountTargetDomain.
* `vpc_id` - (Optional) Filter results by a specific VpcId.
* `vswitch_id` - (Optional) Filter results by a specific VSwitchId.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of MountTargetDomain.
* `targets` - A list of MountTargetDomains. Each element contains the following attributes:
   * `id` - ID of the MountTargetDomain.
   * `mount_target_domain` - MountTargetDomain of the MountTarget.
   * `type`- NetworkType of The MountTarget.
   * `vpc_id` - VpcId of The MountTarget.
   * `vswitch_id` - VSwitchId of The MountTarget.
   * `access_group_name` - AccessGroup of The MountTarget.
