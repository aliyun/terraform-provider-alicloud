---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_resource"
sidebar_current: "docs-alicloud-datasource-data-works-resource"
description: |-
  Provides a list of Data Works Resource to the user.
---

# alicloud_data_works_resource

Provides a list of Data Works Resource items owned by the user.

-> **NOTE:** Available since v1.214.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_data_works_project" "example" {
  project_name = "tf_example"
}

data "alicloud_data_works_resource" "example" {
  project_id  = data.alicloud_data_works_project.example.id
  output_file = "resources.txt"
}

output "first_resource_id" {
  value = data.alicloud_data_works_resource.example.resources.0.id
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the DataWorks workspace to query resources for.
* `owner` - (Optional) The owner of the resource, used to filter results.
* `type` - (Optional) The type of the resource. Valid values: `jar`, `python`, `file`, `archive`.
* `ids` - (Optional) A list of resource IDs to filter results.
* `output_file` - (Optional) File path where results will be saved after running `terraform plan`.

## Attributes Reference

The following attributes are exported:

* `id` - The data source ID (hash of the resource IDs).
* `resources` - A list of Data Works Resource items. Each element contains the following attributes:

### resources

* `id` - The resource ID, formats as `<project_id>:<resource_id>`.
* `resource_id` - The unique identifier of the resource.
* `resource_name` - The name of the resource.
* `project_id` - The ID of the workspace.
* `owner` - The responsible person of the resource file.
* `type` - The resource type.
* `create_time` - The creation time of the resource.
* `modify_time` - The last modification time of the resource.
* `source_type` - The file resource source storage type.
* `source_path` - The file source path.
* `target_type` - The file destination storage type.
* `target_path` - The file destination storage path.
* `data_source` - The data source information.
  * `name` - The data source name.
  * `type` - The data source type.
* `script` - The script information of the workflow.
  * `path` - The script path.
  * `runtime_command` - The script type.
  * `script_id` - The script ID.
