---
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_access_group"
sidebar_current: "docs-alicloud-resource-nas-access-group"
description: |-
  Provides a Alicloud Nas_Access_Group resource.
---

# alicloud\_nas_access_group

In NAS, the permission group acts as a whitelist that allows you to restrict file system access. You can allow specified IP addresses or CIDR blocks to access the file system, and assign different levels of access permission to different IP addresses or CIDR blocks by adding rules to the permission group.

Provides a Nas_Access_Group resource.

 -> **NOTE:** Available in v1.33.0+.

## Example Usage

Basic Usage

```
resource "alicloud_nas_access_group" "foo" {
    	accessgroup_name = "CreateAccessGroup"
 	accessgroup_type = "Classic"
 	description = "test_AccessG"
  
}
```
## Argument Reference

The following arguments are supported:

* `accessgroup_name` - (Required, ForceNew) The AccessGroupName block for the AccessGroup.
* `accessgroup_type` - (Required, ForceNew) The AccessGroupType block for the AccessGroup
* `description`      - (Optional) The AccessGroup description. Defaults to null.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the AccessGroup.

##Import
Nas Access Group can be imported using the id, e.g.

```
$ terraform import alicloud_nas_access_group.default tf_testAccNasConfig
```
