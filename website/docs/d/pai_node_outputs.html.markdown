---
subcategory: "PAI"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_node_outputs"
sidebar_current: "docs-alicloud-datasource-pai-node-outputs"
description: |-
  Provides the list of outputs of a node in a PAI Studio experiment.
---

# alicloud_pai_node_outputs

This data source provides the list of outputs of a node in a PAI Studio experiment.

-> **NOTE:** Available since v1.296.0.

## Example Usage

```terraform
data "alicloud_pai_node_outputs" "default" {
  experiment_id = "ex-xxx"
  node_id       = "node-xxx"
}

output "first_output_id" {
  value = data.alicloud_pai_node_outputs.default.outputs.0.output_id
}
```

## Argument Reference

The following arguments are supported:

* `experiment_id` - (Required) The ID of the PAI Studio experiment.
* `node_id` - (Required) The ID of the node in the experiment.
* `output_id_regex` - (Optional) A regex string to filter results by OutputId. The filter is applied client-side.
* `type` - (Optional) Filter results by node output Type, e.g. `Model`. The filter is applied client-side.

## Attributes Reference

The following attributes are exported in addition to the arguments above:

* `id` - The ID of the data source, formatted as `<experiment_id>:<node_id>`.
* `ids` - A list of OutputIds of the matched node outputs.
* `outputs` - The list of node outputs. Each element contains the following attributes:
  * `output_id` - The ID of the node output.
  * `output_index` - The index of the node output.
  * `node_name` - The name of the node.
  * `algo_name` - The algorithm name of the node.
  * `display_name` - The display name of the node output.
  * `type` - The type of the node output, e.g. `Model`, `MaxComputeTable`.
  * `location_type` - The location type of the node output, e.g. `MaxComputeTable`.
  * `value` - The value payload of the node output, serialized as a JSON string. Its structure varies by `type`.
