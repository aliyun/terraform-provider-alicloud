---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_chart_repositories"
sidebar_current: "docs-alicloud-datasource-cr-chart-repositories"
description: |-
  Provides a list of Cr Chart Repositories to the user.
---

# alicloud\_cr\_chart\_repositories

This data source provides the Cr Chart Repositories of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.149.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cr_chart_repositories" "ids" {
  instance_id = "example_value"
  ids         = ["example_value-1", "example_value-2"]
}
output "cr_chart_repository_id_1" {
  value = data.alicloud_cr_chart_repositories.default.ids.0
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Chart Repository IDs.
* `instance_id` - (Required, ForceNew) InstanceId.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by repository name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of matched Container Registry Enterprise Edition repositories.
* `names` - A list of repository names.
* `repositories` - A list of Cr Chart Repositories. Each element contains the following attributes:
	* `chart_repository_id` - The first ID of the resource.
	* `create_time` - The creation time of the resource.
	* `id` - The ID of the Chart Repository.
	* `instance_id` - The ID of the Container Registry instance.
	* `repo_name` - The name of the repository.
	* `repo_namespace_name` - The namespace to which the repository belongs.
	* `repo_type` - The type of the repository. Valid values: `PUBLIC`,`PRIVATE`.
	* `summary` - The summary about the repository.
