---
subcategory: "Realtime Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_realtime_compute_variables"
description: |-
  Provides a list of Realtime Compute Variables to the user.
---

# alicloud_realtime_compute_variables

This data source provides the Realtime Compute Variables of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.266.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_realtime_compute_variables" "default" {
  workspace  = alicloud_realtime_compute_vvp_instance.default.resource_id
  namespace  = "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"
  name_regex = "terraform-variable-example"
}

output "first_variable_name" {
  value = data.alicloud_realtime_compute_variables.default.variables.0.name
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) The ID of the workspace.
* `namespace` - (Required) The name of the namespace (project space).
* `name_regex` - (Optional) A regex string to filter variables by name.
* `ids` - (Optional) A list of variable IDs. The ID format is `<workspace>:<namespace>:<name>`.
* `output_file` - (Optional) The file name of the output results, if not set, the results will be printed to stdout.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of variable IDs.
* `variables` - A list of Realtime Compute Variables. Each element contains the following attributes:
  * `id` - The ID of the variable, formatted as `<workspace>:<namespace>:<name>`.
  * `workspace` - The ID of the workspace.
  * `namespace` - The name of the namespace.
  * `name` - The name of the variable.
  * `kind` - The kind of the variable, currently `Plain`.
  * `value` - The value of the variable.
  * `description` - The description of the variable.
  * `region_id` - The region ID of the resource.
