---
subcategory: "File Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_mount_targets"
sidebar_current: "docs-alicloud-datasource-nas-mount-targets"
description: |-
  Provides a list of mns topic subscriptions available to the user.
---

# alicloud\_nas_mount_targets

This data source provides MountTargets available to the user.

-> **NOTE**: Available in 1.35.0+

## Example Usage

```terraform
data "alicloud_nas_mount_targets" "example" {
  file_system_id    = "1a2sc4d"
  access_group_name = "tf-testAccNasConfig"
}

output "the_first_mount_target_domain" {
  value = data.alicloud_nas_mount_targets.example.targets.0.id
}
```

## Argument Reference

The following arguments are supported:

* `file_system_id` - (Required, ForceNew) The ID of the FileSystem that owns the MountTarget.
* `access_group_name` - (Optional, ForceNew) Filter results by a specific AccessGroupName.
* `type` - (Optional, Deprecated in 1.95.0+) Field `type` has been deprecated from provider version 1.95.0. New field `network_type` replaces it.
* `network_type` - (Optional, ForceNew, Available 1.95.0+) Filter results by a specific NetworkType.
* `mount_target_domain` - (Optional, Deprecated in 1.53.+) Field `mount_target_domain` has been deprecated from provider version 1.53.0. New field `ids` replaces it.
* `vpc_id` - (Optional, ForceNew) Filter results by a specific VpcId.
* `vswitch_id` - (Optional, ForceNew) Filter results by a specific VSwitchId.
* `ids` - (Optional, ForceNew, Available 1.53.0+) A list of MountTargetDomain.
* `status` - (Optional, ForceNew, Available 1.95.0+) Filter results by the status of mount target. Valid values: `Active`, `Inactive` and `Pending`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of MountTargetDomain.
* `targets` - A list of MountTargetDomains. Each element contains the following attributes:
   * `id` - ID of the MountTargetDomain.
   * `mount_target_domain` - MountTargetDomain of the MountTarget.
   * `type`- Field `type` has been deprecated from provider version 1.95.0. New field `network_type` replaces it. 
   * `network_type`- (Available 1.95.0+) NetworkType of The MountTarget.
   * `status`- (Available 1.95.0+) The status of the mount target. 
   * `vpc_id` - VpcId of The MountTarget.
   * `vswitch_id` - VSwitchId of The MountTarget.
   * `access_group_name` - AccessGroup of The MountTarget.
