---
subcategory: "DFS"
layout: "alicloud"
page_title: "Alicloud: alicloud_dfs_access_rule"
description: |-
  Provides a Alicloud DFS Access Rule resource.
---

# alicloud_dfs_access_rule

Provides a DFS Access Rule resource. 

For information about DFS Access Rule and how to use it, see [What is Access Rule](https://www.alibabacloud.com/help/en/aibaba-cloud-storage-services/latest/apsara-file-storage-for-hdfs).

-> **NOTE:** Available since v1.140.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_dfs_access_group" "default" {
  description       = "example"
  network_type      = "VPC"
  access_group_name = var.name
}

resource "alicloud_dfs_access_rule" "default" {
  description     = "example"
  rw_access_type  = "RDWR"
  priority        = "1"
  network_segment = "192.168.81.1"
  access_group_id = alicloud_dfs_access_group.default.id
}
```

## Argument Reference

The following arguments are supported:
* `access_group_id` - (Required, ForceNew) Permission group resource ID. You must specify the permission group ID when creating a permission rule.
* `description` - (Optional) Permission rule description.  No more than 32 characters in length.
* `network_segment` - (Required, ForceNew) The IP address or network segment of the authorized object.
* `priority` - (Required) Permission rule priority. When the same authorization object matches multiple rules, the high-priority rule takes effect. Value range: 1~100,1 is the highest priority.
* `rw_access_type` - (Required) The read and write permissions of the authorized object on the file system. Value: RDWR: readable and writable RDONLY: Read only.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<access_group_id>:<access_rule_id>`.
* `access_rule_id` - The unique identity of the permission rule, which is used to retrieve the permission rule for a specific day in the permission group.
* `create_time` - Permission rule resource creation time.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Access Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Access Rule.
* `update` - (Defaults to 5 mins) Used when update the Access Rule.

## Import

DFS Access Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_dfs_access_rule.example <access_group_id>:<access_rule_id>
```