---
subcategory: "Serverless Workflow"
layout: "alicloud"
page_title: "Alicloud: alicloud_fnf_flows"
sidebar_current: "docs-alicloud-datasource-fnf-flows"
description: |-
  Provides a list of Fnf Flows to the user.
---

# alicloud\_fnf\_flows

This data source provides the Fnf Flows of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.105.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_fnf_flows" "example" {
  ids        = ["example_value"]
  name_regex = "the_resource_name"
}

output "first_fnf_flow_id" {
  value = data.alicloud_fnf_flows.example.flows.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Flow IDs.
* `limit` - (Optional, ForceNew, Available in v1.110.0+) The number of resource queries.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Flow name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Flow names.
* `flows` - A list of Fnf Flows. Each element contains the following attributes:
	* `definition` - The definition of the flow. It must comply with the Flow Definition Language (FDL) syntax.
	* `description` - The description of the flow.
	* `flow_id` - The unique ID of the flow.
	* `id` - The ID of the Flow.
	* `last_modified_time` - The time when the flow was last modified.
	* `name` - The name of the flow. The name must be unique in an Alibaba Cloud account.
	* `role_arn` - The ARN of the specified RAM role that Serverless Workflow uses to assume the role when Serverless Workflow executes a flow.
	* `type` - The type of the flow. Set the value to `FDL`.
