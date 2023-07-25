---
subcategory: "Apsara File Storage for HDFS (DFS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dfs_access_group"
sidebar_current: "docs-alicloud-resource-dfs-access-group"
description: |-
  Provides a Alicloud DFS Access Group resource.
---

# alicloud_dfs_access_group

Provides a DFS Access Group resource.

For information about DFS Access Group and how to use it, see [What is Access Group](https://www.alibabacloud.com/help/doc-detail/207144.htm).

-> **NOTE:** Available since v1.133.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

resource "alicloud_dfs_access_group" "default" {
  access_group_name = var.name
  network_type      = "VPC"
}
```

## Argument Reference

The following arguments are supported:

* `access_group_name` - (Required) The Name of Access Group.The length of `access_group_name` does not exceed 100 bytes.
* `description` - (Optional) The Description of Access Group. The length of `description` does not exceed 100 bytes.
* `network_type` - (Required, ForceNew) The NetworkType of Access Group. Valid values: `VPC`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Access Group.

## Import

DFS Access Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_dfs_access_group.example <id>
```
