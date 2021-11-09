---
subcategory: "Apsara File Storage for HDFS"
layout: "alicloud"
page_title: "Alicloud: alicloud_dfs_access_rule"
sidebar_current: "docs-alicloud-resource-dfs-access-rule"
description: |-
  Provides a Alicloud DFS Access Rule resource.
---

# alicloud\_dfs\_access\_rule

Provides a DFS Access Rule resource.

For information about DFS Access Rule and how to use it, see [What is Access Rule](https://www.alibabacloud.com/help/doc-detail/207144.htm).

-> **NOTE:** Available in v1.140.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "example_name"
}

resource "alicloud_dfs_access_group" "default" {
  network_type      = "VPC"
  access_group_name = var.name
  description       = var.name
}

resource "alicloud_dfs_access_rule" "default" {
  network_segment = "192.0.2.0/24"
  access_group_id = alicloud_dfs_access_group.default.id
  description     = var.name
  rw_access_type  = "RDWR"
  priority        = "10"
}

```

## Argument Reference

The following arguments are supported:

* `access_group_id` - (Required, ForceNew) The resource ID of Access Group.
* `description` - (Optional) The Description of the Access Rule.
* `network_segment` - (Required, ForceNew) The NetworkSegment of the Access Rule.
* `priority` - (Required) The Priority of the Access Rule. Valid values: `1` to `100`. **NOTE:** When multiple rules are matched by the same authorized object, the high-priority rule takes effect. `1` is the highest priority.
* `rw_access_type` - (Required) The RWAccessType of the Access Rule. Valid values: `RDONLY`, `RDWR`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of the Access Rule. The value formats as `<access_group_id>:<access_rule_id>`.
* `access_rule_id` - The ID of the Access Rule.

## Import

DFS Access Rule can be imported using the id, e.g.

```
$ terraform import alicloud_dfs_access_rule.example <access_group_id>:<access_rule_id>
```
