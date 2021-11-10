---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_projects"
sidebar_current: "docs-alicloud-datasource-log-projects"
description: |-
  Provides a list of log projects to the user.
---

# alicloud\_log\_projects

This data source provides the Log Projects of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.126.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_log_projects" "example" {
  ids = ["the_project_name"]
}

output "first_log_project_id" {
  value = data.alicloud_log_projects.example.project.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of project IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by project name.
* `status` - (Optional, ForceNew) The status of log project. Valid values `Normal` and `Disable`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of project names.
* `projects` - A list of Log Project. Each element contains the following attributes:
	* `description` - The description of the project.
	* `id` - The ID of the project.
	* `project_name` - The name of the project. 
	* `last_modify_time` - The last modify time of project.
	* `owner` - The owner of project.
	* `region` - The region of project.
	* `status` - The status of project.
