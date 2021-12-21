---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_sharding_nodes"
sidebar_current: "docs-alicloud-datasource-mongodb-sharding-nodes"
description: |-
  Provides a list of Mongodb Sharding Nodes to the user.
---

# alicloud\_mongodb\_sharding\_nodes

This data source provides the Mongodb Sharding Nodes of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.149.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_mongodb_sharding_nodes" "ids" {
  db_instance_id = "example_value"
  ids            = ["example_value-1", "example_value-2"]
  status         = "Running"
}
output "mongodb_sharding_node_id_1" {
  value = data.alicloud_mongodb_sharding_nodes.ids.nodes.0.id
}

```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, ForceNew) The ID of the sharded cluster instance.
* `ids` - (Optional, ForceNew, Computed)  A list of Sharding Node IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `Creating`, `Deleting` or `Running`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `nodes` - A list of Mongodb Sharding Nodes. Each element contains the following attributes:
  * `db_instance_id` - The ID of the sharded cluster instance.	
  * `id` - The ID of the Sharding Node.	The value formats as `<db_instance_id>:<node_id>`.
  * `node_class` - The specifications of Shard nodes.
  * `node_id` - The ID of the node.	
  * `node_storage` - The disk space of the Node.	
  * `status` - The status of the resource.