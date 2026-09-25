---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_node_on_baseline"
sidebar_current: "docs-alicloud-datasource-data-works-node-on-baseline"
description: |-
  Provides a Data Works Node On Baseline resource of the current Alibaba Cloud user.
---

# alicloud\_data\_works\_node\_on\_baseline

This data source provides a Data Works Node On Baseline of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.235.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_data_works_node_on_baseline" "example" {
  baseline_id = "123"
}

output "node_owner" {
  value = data.alicloud_data_works_node_on_baseline.example.owner
}
```

Filter by node name pattern

```terraform
data "alicloud_data_works_node_on_baseline" "example" {
  baseline_id = "123"
  name_regex  = "node_.*"
}
```

## Argument Reference

The following arguments are supported:

* `baseline_id` - (Required, ForceNew) The ID of the DataWorks Baseline that the node belongs to.
* `name_regex` - (Optional, ForceNew) A regex string used to filter the result by node name. The first matching node is returned.
* `node_id` - (Optional, ForceNew) Filter the result by the specified node ID. The first matching node is returned.
* `owner` - (Optional, ForceNew) Filter the result by the specified node owner. The first matching node is returned.
* `project_id` - (Optional, ForceNew) Filter the result by the specified DataWorks project ID. The first matching node is returned.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The resource ID in the form of `<baseline_id>:<node_id>`.
* `node_id` - The ID of the node mounted on the baseline.
* `node_name` - The name of the node mounted on the baseline.
* `owner` - The owner of the node mounted on the baseline.
* `project_id` - The ID of the DataWorks workspace that the node belongs to.
* `region_id` - The region ID of the node.
