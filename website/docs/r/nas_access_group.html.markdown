---
subcategory: "Network Attached Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_access_group"
sidebar_current: "docs-alicloud-resource-nas-access-group"
description: |-
  Provides a Alicloud NAS Access Group resource.
---

# alicloud\_nas\_access\_group

Provides a NAS Access Group resource.

In NAS, the permission group acts as a whitelist that allows you to restrict file system access. You can allow specified IP addresses or CIDR blocks to access the file system, and assign different levels of access permission to different IP addresses or CIDR blocks by adding rules to the permission group.
For information about NAS Access Group and how to use it, see [What is NAS Access Group](https://www.alibabacloud.com/help/en/doc-detail/27534)

-> **NOTE:** Available in v1.33.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_nas_access_group" "foo" {
  name        = "CreateAccessGroup"
  type        = "Classic"
  description = "test_AccessG"

}
```

Example after v1.92.0

```terraform
resource "alicloud_nas_access_group" "foo" {
  access_group_name = "CreateAccessGroup"
  access_group_type = "Vpc"
  description       = "test_AccessG"
  file_system_type  = "extreme"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, ForceNew, Deprecated in v1.92.0+) Replaced by `access_group_name` after version 1.92.0.
* `type` - (Optional, ForceNew, Deprecated in v1.92.0+) Replaced by `access_group_type` after version 1.92.0.
* `access_group_name` - (Optional, ForceNew, Available in v1.92.0+) A Name of one Access Group.
* `access_group_type` - (Optional, ForceNew, Available in v1.92.0+) A Type of one Access Group. Valid values: `Vpc` and `Classic`.
* `description` - (Optional) The Access Group description.
* `file_system_type` - (Optional, ForceNew, Available in v1.92.0+) The type of file system. Valid values: `standard` and `extreme`. Default to `standard`. Note that the extreme only support Vpc Network.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Access Group. The value as `<access_group_name>`. 

-> **NOTE:** The ID value as `<access_group_name>`:`<file_system_type>` after version 1.92.0.

## Import

NAS Access Group can be imported using the id, e.g.

```
$ terraform import alicloud_nas_access_group.foo tf_testAccNasConfig:standard
```
