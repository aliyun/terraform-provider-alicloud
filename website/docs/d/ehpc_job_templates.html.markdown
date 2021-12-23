---
subcategory: "Elastic High Performance Computing(ehpc)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ehpc_job_templates"
sidebar_current: "docs-alicloud-datasource-ehpc-job-templates"
description: |-
  Provides a list of Ehpc Job Templates to the user.
---

# alicloud\_ehpc\_job\_templates

This data source provides the Ehpc Job Templates of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.133.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ehpc_job_template" "default" {
  job_template_name = "example_value"
  command_line      = "./LammpsTest/lammps.pbs"
}
data "alicloud_ehpc_job_templates" "ids" {
  ids = [alicloud_ehpc_job_template.default.id]
}
output "ehpc_job_template_id_1" {
  value = data.alicloud_ehpc_job_templates.ids.id
}


```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Job Template IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `templates` - A list of Ehpc Job Templates. Each element contains the following attributes:
	* `array_request` - Queue Jobs, Is of the Form: 1-10:2.
	* `clock_time` - Job Maximum Run Time.
	* `command_line` - Job Commands.
	* `gpu` - A Single Compute Node Using the GPU Number.Possible Values: 1~20000.
	* `id` - The ID of the Job Template.
	* `job_template_id` - The first ID of the resource.
	*  `job_template_name` - A Job Template Name.
	* `mem` - A Single Compute Node Maximum Memory.
	* `node` - Submit a Task Is Required for Computing the Number of Data Nodes to Be. Possible Values: 1~5000 .
	* `package_path` - Job Commands the Directory.
	* `priority` - The Job Priority.Possible Values: 0~9.
	* `queue` - The Job Queue.
	* `re_runable` - If the Job Is Support for the Re-Run.
	* `runas_user` - The name of the user who performed the job.
	* `stderr_redirect_path` - Error Output Path.
	* `stdout_redirect_path` - Standard Output Path and.
	* `task` - A Single Compute Node Required Number of Tasks. Possible Values: 1~20000 .
	* `thread` - A Single Task and the Number of Required Threads.Possible Values: 1~20000.
	* `variables` - The Job of the Environment Variable.
