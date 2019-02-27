---
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_mounttarget"
sidebar_current: "docs-alicloud-resource-nas_mounttarget"
description: |-
  Provides a Alicloud NAS_MountTarget resource.
---

# alicloud\_nas_mounttarget

Provides a NAS_MountTarget resource.


## Example Usage

Basic Usage

```
resource "alicloud_nas_mounttarget" "foo" {
  filesystem_id = "123467"
  accessgroup_name = "CreateAG"
  networktype = "Classic"
}
```
## Argument Reference

The following arguments are supported:

* `filesystem_id` - (Required, Forces new resource) The FileSystemId block for the MountTarget.
* `accessgroup_name` - (Required, Forces new resource) The AccessGroupName block for the MountTarget
* `networktype` - (Required, Forces new resource) The NetworkType block for the MountTarget.
* `vpc_id` - (Optional) The MountTarget VpcId. Defaults to null.
* `vsw_id` - (Optional) The MountTarget VSwitchId. Defaults to null.


## Attributes Reference

The following attributes are exported:

* `mounttarget_domain`  - The MountTargetDomain of the MountTarget.
* `status`              - The Status block for the MountTarget.
* `network_type`        - The NetworkType block for the MountTarget.
* `access_group`        - The AccessGroup of the MountTarget.
* `vpc_id`              - The VpcId of the MountTarget.
* `vsw_id`              - The VswId of the MountTarget.


