---
subcategory: "PAI"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_node_output"
sidebar_current: "docs-alicloud-datasource-pai-node-output"
description: |-
  Provides a single output of a node in a PAI Studio experiment.
---

# alicloud_pai_node_output

This data source provides a single output of a node in a PAI Studio experiment, identified by the experiment id, node id and output id.

-> **NOTE:** Available since v1.296.0.

## Example Usage

```terraform
data "alicloud_pai_node_output" "default" {
  experiment_id = "ex-xxx"
  node_id       = "node-xxx"
  output_id     = "outputModel"
}
```

## Argument Reference

The following arguments are supported:

* `experiment_id` - (Required) The ID of the PAI Studio experiment.
* `node_id` - (Required) The ID of the node in the experiment.
* `output_id` - (Required) The ID of the node output.
* `output_index` - (Optional) The index of the node output. When the node emits multiple outputs sharing the same OutputId, set this to select one.

## Attributes Reference

The following attributes are exported in addition to the arguments above:

* `id` - The ID of the node output, formatted as `<experiment_id>:<node_id>:<output_id>`.
* `node_name` - The name of the node.
* `algo_name` - The algorithm name of the node.
* `display_name` - The display name of the node output.
* `type` - The type of the node output, e.g. `Model`, `MaxComputeTable`.
* `location_type` - The location type of the node output, e.g. `MaxComputeTable`.
* `value` - The value payload of the node output, serialized as a JSON string. Its structure varies by `type`.
