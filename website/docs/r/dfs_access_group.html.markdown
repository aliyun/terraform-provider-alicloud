---
subcategory: "Apsara File Storage for HDFS (DFS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dfs_access_group"
description: |-
  Provides a Alicloud DFS Access Group resource.
---

# alicloud_dfs_access_group

Provides a DFS Access Group resource. 

For information about DFS Access Group and how to use it, see [What is Access Group](https://www.alibabacloud.com/help/en/aibaba-cloud-storage-services/latest/apsara-file-storage-for-hdfs).

-> **NOTE:** Available since v1.133.0.

## Example Usage

Basic Usage

```terraform
resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_dfs_access_group" "default" {
  access_group_name = "tf-example-${random_integer.default.result}"
  network_type      = "VPC"
}
```

## Argument Reference

The following arguments are supported:
* `access_group_name` - (Required) The permission group name. The naming rules are as follows: The length is 6~64 characters. Globally unique and cannot be an empty string. English letters are supported and can contain numbers, underscores (_), and dashes (-).
* `description` - (Optional) The permission group description.  No more than 32 characters in length.
* `network_type` - (Required, ForceNew) The permission group type. Only VPC (VPC) is supported.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the permission group resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Access Group.
* `delete` - (Defaults to 5 mins) Used when delete the Access Group.
* `update` - (Defaults to 5 mins) Used when update the Access Group.

## Import

DFS Access Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_dfs_access_group.example <id>
```