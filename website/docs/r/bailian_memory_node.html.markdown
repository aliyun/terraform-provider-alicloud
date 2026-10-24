---
subcategory: "Bailian (SFM)"
layout: "alicloud"
page_title: "Alicloud: alicloud_bailian_memory_node"
sidebar_current: "docs-alicloud-resource-bailian-memory-node"
description: |-
  Provides a Bailian Memory Node resource.
---

# alicloud_bailian_memory_node

Provides a Bailian Memory Node resource.

For information about Bailian Memory Node and how to use it, see [What is Memory Node](https://www.alibabacloud.com/help/en/model-studio/developer-api/api-reference-37).

-> **NOTE:** Available since v1.243.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_bailian_memory" "default" {
  workspace_id = "llm-xxx"
  description  = "tf-example-memory"
}

resource "alicloud_bailian_memory_node" "default" {
  workspace_id = alicloud_bailian_memory.default.workspace_id
  memory_id    = alicloud_bailian_memory.default.memory_id
  content      = "tf-example-memory-node-content"
}
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required, ForceNew) The workspace ID of Bailian.
* `memory_id` - (Required, ForceNew) The memory ID the node belongs to.
* `content` - (Required) The content of the memory node.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, format as `<workspace_id>:<memory_id>:<memory_node_id>`.
* `memory_node_id` - The ID of the memory node.

## Import

Bailian Memory Node can be imported using the id, e.g.

```shell
$ terraform import alicloud_bailian_memory_node.example <workspace_id>:<memory_id>:<memory_node_id>
```
