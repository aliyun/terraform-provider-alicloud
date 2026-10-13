---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_node_on_baselines"
sidebar_current: "docs-alicloud-datasource-data-works-node-on-baselines"
description: |-
  Provides a list of Data Works Node On Baseline resources of the current Alibaba Cloud user.
---

# alicloud\_data\_works\_node\_on\_baselines

This data source provides a list of Data Works Node On Baseline resources of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.235.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_data_works_node_on_baselines" "example" {
  baseline_id = "123"
}

output "first_node_id" {
  value = data.alicloud_data_works_node_on_baselines.example.nodes.0.node_id
}
```

Filter by node name pattern

```terraform
data "alicloud_data_works_node_on_baselines" "example" {
  baseline_id = "123"
  name_regex  = "node_.*"
}
```

## Argument Reference

The following arguments are supported:

* `baseline_id` - (Required, ForceNew) The ID of the DataWorks Baseline that the nodes belong to.
* `ids` - (Optional, ForceNew) A list of node IDs. Only nodes whose IDs match this list are returned.
* `name_regex` - (Optional, ForceNew) A regex string used to filter the result by node name.
* `node_id` - (Optional, ForceNew) Filter the result by the specified node ID.
* `owner` - (Optional, ForceNew) Filter the result by the specified node owner.
* `project_id` - (Optional, ForceNew) Filter the result by the specified DataWorks project ID.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The resource ID of the data source, computed from the matched node IDs.
* `nodes` - A list of Data Works Node On Baseline resources. Each element contains the following attributes:

### nodes

* `id` - The resource ID in the form of `<baseline_id>:<node_id>`.
* `baseline_id` - The ID of the DataWorks Baseline that the node belongs to.
* `node_id` - The ID of the node mounted on the baseline.
* `node_name` - The name of the node mounted on the baseline.
* `owner` - The owner of the node mounted on the baseline.
* `project_id` - The ID of the DataWorks workspace that the node belongs to.
* `region_id` - The region ID of the node.
