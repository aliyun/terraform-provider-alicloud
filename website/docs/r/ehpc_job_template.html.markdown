---
subcategory: "Elastic High Performance Computing(ehpc)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ehpc_job_template"
sidebar_current: "docs-alicloud-resource-ehpc-job-template"
description: |-
  Provides a Alicloud Ehpc Job Template resource.
---

# alicloud\_ehpc\_job\_template

Provides a Ehpc Job Template resource.

For information about Ehpc Job Template and how to use it, see [What is Job Template](https://www.alibabacloud.com/help/product/57664.html).

-> **NOTE:** Available in v1.133.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ehpc_job_template" "default" {
  job_template_name = "example_value"
  command_line      = "./LammpsTest/lammps.pbs"
}

```

## Argument Reference

The following arguments are supported:

* `array_request` - (Optional) Queue Jobs, Is of the Form: 1-10:2.
* `clock_time` - (Optional) Job Maximum Run Time.
* `command_line` - (Required) Job Commands.
* `gpu` - (Optional) A Single Compute Node Using the GPU Number.Possible Values: 1~20000.
* `job_template_name` - (Required) A Job Template Name.
* `mem` - (Optional) A Single Compute Node Maximum Memory.
* `node` - (Optional) Submit a Task Is Required for Computing the Number of Data Nodes to Be. Possible Values: 1~5000 .
* `package_path` - (Optional) Job Commands the Directory.
* `priority` - (Optional) The Job Priority.
* `queue` - (Optional) The Job Queue.
* `re_runable` - (Optional) If the Job Is Support for the Re-Run.
* `runas_user` - (Optional) The name of the user who performed the job.
* `stderr_redirect_path` - (Optional) Error Output Path.
* `stdout_redirect_path` - (Optional) Standard Output Path and.
* `task` - (Optional) A Single Compute Node Required Number of Tasks. Possible Values: 1~20000 .
* `thread` - (Optional) A Single Task and the Number of Required Threads.
* `variables` - (Optional) The Job of the Environment Variable.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Job Template.

## Import

Ehpc Job Template can be imported using the id, e.g.

```
$ terraform import alicloud_ehpc_job_template.example <id>
```
