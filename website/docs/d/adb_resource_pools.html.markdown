---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_adb_resource_pools"
sidebar_current: "docs-alicloud-datasource-adb-resource-pools"
description: |-
  Provides a list of Adb Resource Pools to the user.
---

# alicloud\_adb\_resource\_pools

This data source provides the Adb Resource Pools of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.170.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_adb_resource_pools" "ids" {
  db_cluster_id = "example_value"
  ids           = ["example_value-1", "example_value-2"]
}
output "adb_resource_pool_id_1" {
  value = data.alicloud_adb_resource_pools.ids.pools.0.id
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The db cluster id.
* `resource_pool_name` - (Optional, ForceNew) The name of the resource pool.
* `ids` - (Optional, ForceNew, Computed)  A list of Resource Pool IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Pool name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Pool names.
* `pools` - A list of Adb Resource Pools. Each element contains the following attributes:
	* `create_time` - The creation time.
	* `db_cluster_id` - The db cluster id.
	* `id` - The ID of the Resource Pool. The value formats as `<db_cluster_id>:<resource_pool_name>`.
	* `node_num` - The number of nodes. 
	* `resource_pool_name` - The name of the resource pool.
	* `query_type` - The query type.