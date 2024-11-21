---
subcategory: "Apsara File Storage for HDFS (DFS)"
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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dfs_access_rule&exampleId=17bc700c-1e01-fecd-da0d-6fc671258fa64f5e9431&activeTab=example&spm=docs.r.dfs_access_rule.0.17bc700c1e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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