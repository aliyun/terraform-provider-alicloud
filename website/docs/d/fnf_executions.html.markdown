---
subcategory: "Serverless Workflow"
layout: "alicloud"
page_title: "Alicloud: alicloud_fnf_executions"
sidebar_current: "docs-alicloud-datasource-fnf-executions"
description: |-
  Provides a list of FnF Executions to the user.
---

# alicloud\_fnf\_executions

This data source provides the FnF Executions of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.149.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_fnf_executions" "ids" {
  flow_name = "example_value"
  ids       = ["my-Execution-1", "my-Execution-2"]
}
output "fnf_execution_id_1" {
  value = data.alicloud_fn_f_executions.ids.executions.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `flow_name` - (Required, ForceNew) The name of the flow.
* `ids` - (Optional, ForceNew, Computed)  A list of Execution IDs. The value formats as `<flow_name>:<execution_name>`.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Execution name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `Running`, `Stopped`, `Succeeded`, `Failed`, `TimedOut`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Execution names.
* `executions` - A list of Fn F Executions. Each element contains the following attributes:
	* `execution_name` - The name of the execution.
	* `flow_name` - The name of the flow.
	* `id` - The ID of the Execution. The value formats as `<flow_name>:<execution_name>`.
	* `input` - The Input information for this execution.
	* `output` - The output of the execution.
	* `started_time` - The started time of the execution.
	* `status` - The status of the resource.
	* `stopped_time` - The stopped time of the execution.