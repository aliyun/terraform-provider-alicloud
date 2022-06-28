---
subcategory: "Brain Industrial"
layout: "alicloud"
page_title: "Alicloud: alicloud_brain_industrial_pid_loops"
sidebar_current: "docs-alicloud-datasource-brain-industrial-pid-loops"
description: |-
  Provides a list of Brain Industrial Pid Loops to the user.
---

# alicloud\_brain\_industrial\_pid\_loops

This data source provides the Brain Industrial Pid Loops of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.117.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_brain_industrial_pid_loops" "example" {
  pid_project_id = "856c6b8f-ca63-40a4-xxxx-xxxx"
  ids            = ["742a3d4e-d8b0-47c8-xxxx-xxxx"]
  name_regex     = "tf-testACC"
}

output "first_brain_industrial_pid_loop_id" {
  value = data.alicloud_brain_industrial_pid_loops.example.loops.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Pid Loop IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Pid Loop name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `pid_loop_name` - (Optional, ForceNew) The name of Pid Loop.
* `pid_project_id` - (Required, ForceNew) The pid project id.
* `status` - (Optional, ForceNew) The status of Pid Loop.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Pid Loop names.
* `loops` - A list of Brain Industrial Pid Loops. Each element contains the following attributes:
	* `id` - The ID of the Pid Loop.
	* `pid_loop_dcs_type` - The dcs type of Pid Loop.
	* `pid_loop_id` - The ID of the Pid Loop.
	* `pid_loop_is_crucial` - Whether is crucial Pid Loop.
	* `pid_loop_name` - The name of Pid Loop.
	* `pid_loop_type` - The type of Pid Loop.
	* `status` - The status of Pid Loop.
