---
subcategory: "Brain Industrial"
layout: "alicloud"
page_title: "Alicloud: alicloud_brain_industrial_pid_projects"
sidebar_current: "docs-alicloud-datasource-brain-industrial-pid-projects"
description: |-
  Provides a list of Brain Industrial Pid Projects to the user.
---

# alicloud\_brain\_industrial\_pid\_projects

This data source provides the Brain Industrial Pid Projects of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.113.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_brain_industrial_pid_projects" "example" {
  ids        = ["3e74e684-cbb5-xxxx"]
  name_regex = "tf-testAcc"
}

output "first_brain_industrial_pid_project_id" {
  value = data.alicloud_brain_industrial_pid_projects.example.projects.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Pid Project IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Pid Project name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `pid_organization_id` - (Optional, ForceNew) The ID of Pid Organization.
* `pid_project_name` - (Optional, ForceNew) The name of Pid Project.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Pid Project names.
* `projects` - A list of Brain Industrial Pid Projects. Each element contains the following attributes:
	* `id` - The ID of the Pid Project.
	* `pid_organization_id` - The ID of Pid Organization.
	* `pid_project_desc` - The description of Pid Project.
	* `pid_project_id` - The ID of Pid Project.
	* `pid_project_name` - The name of Pid Project.
