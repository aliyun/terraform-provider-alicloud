---
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_accessgroup"
sidebar_current: "docs-alicloud-resource-nas_accessgroup"
description: |-
  Provides a Alicloud NAS_AccessGroup resource.
---

# alicloud\_nas_accessgroup

Provides a NAS_AccessGroup resource.

## Example Usage

Basic Usage

```
resource "alicloud_nas_accessgroup" "foo" {
    accessgroup_name = "CreateAccessGroup"
 	accessgroup_type = "Classic"
 	description = "test_AccessG"
  
}
```
## Argument Reference

The following arguments are supported:

* `accessgroup_name` - (Required, Forces new resource) The AccessGroupName block for the AccessGroup.
* `accessgroup_type` - (Required, Forces new resource) The AccessGroupType block for the AccessGroup
* `description`      - (Optional) The AccessGroup description. Defaults to null.

## Attributes Reference

The following attributes are exported:


* `accessgroup_name`    - The ID of the AccessGroup.
* `accessgroup_type`    - The AccessGroupType block for the AccessGroup.
* `RuleCount`           - The RuleCount block for the AccessGroup.
* `MountTargetCount`    - The MountTargetCount of the AccessGroup.
* `description`         - The Description of the AccessGroup.


