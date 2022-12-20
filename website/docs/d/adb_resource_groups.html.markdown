---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_adb_resource_groups"
sidebar_current: "docs-alicloud-datasource-adb-resource-groups"
description: |-
  Provides a list of Alicloud AnalyticDB for MySQL (ADB) Resource Group owned by an Alibaba Cloud account.
---

# alicloud_adb_resource_groups

This data source provides Adb Resource Group available to the user.[What is Resource Group](https://www.alibabacloud.com/help/en/analyticdb-for-mysql/latest/describe-db-resource-group)

-> **NOTE:** Available in 1.195.0+

## Example Usage

```
data "alicloud_adb_resource_groups" "default" {
  db_cluster_id = "am-bp1a16357gty69185"
  group_name    = "TESTOPENAPI"
}

output "alicloud_adb_resource_group_example_id" {
  value = data.alicloud_adb_resource_groups.default.groups.0.id
}
```

## Argument Reference

The following arguments are supported:
* `db_cluster_id` - (Required,ForceNew) DBClusterId
* `group_name` - (ForceNew,Optional) The name of the resource pool, which cannot exceed 64 bytes in length.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `ids` - (Optional, ForceNew, Computed)  A list of AnalyticDB for MySQL (ADB) Resource Group IDs.


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `groups` - A list of Resource Group Entries. Each element contains the following attributes:
    * `id` - The `key` of the resource supplied above.The value is formulated as `<db_cluster_id>:<group_name>`.
    * `create_time` - Creation time.
    * `db_cluster_id` - DB cluster id.
    * `group_name` - The name of the resource pool.
    * `group_type` - Query type, value description:
      * **etl**: Batch query mode.
      * **interactive**: interactive Query mode
      * **default_type**: the default query mode.
    * `node_num` - The number of nodes. The default number of nodes is 0. The number of nodes must be less than or equal to the number of nodes whose resource name is USER_DEFAULT.
    * `update_time` - Update time.
    * `user` - Binding User.
