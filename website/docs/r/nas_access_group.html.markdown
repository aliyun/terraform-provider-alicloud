---
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_access_group"
sidebar_current: "docs-alicloud-resource-nas-access-group"
description: |-
  Provides a Alicloud Nas Access Group resource.
---

# alicloud\_nas_access_group

Provides a Nas Access Group resource.

In NAS, the permission group acts as a whitelist that allows you to restrict file system access. You can allow specified IP addresses or CIDR blocks to access the file system, and assign different levels of access permission to different IP addresses or CIDR blocks by adding rules to the permission group.

-> **NOTE:** Available in v1.33.0+.

## Example Usage

Basic Usage

```
resource "alicloud_nas_access_group" "foo" {
    name = "CreateAccessGroup"
 	type = "Classic"
 	description = "test_AccessG"
  
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) A Name of one Access Group.
* `type` - (Required, ForceNew) A Type of one Access Group. Valid values: `Vpc` and `Classic`.
* `description` - (Optional) The Access Group description.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Access Group.

## Import

Nas Access Group can be imported using the id, e.g.

```
$ terraform import alicloud_nas_access_group.default tf_testAccNasConfig
```
