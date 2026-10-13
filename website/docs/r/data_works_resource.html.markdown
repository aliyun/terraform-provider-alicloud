---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_resource"
sidebar_current: "docs-alicloud-resource-data-works-resource"
description: |-
  Provides a Alicloud Data Works Resource resource.
---

# alicloud_data_works_resource

Provides a Data Works Resource resource. The resource file is used when computing engine UDFs or scheduling tasks run.

For information about Data Works Resource and how to use it, see [What is Resource](https://www.alibabacloud.com/help/en/dataworks/).

-> **NOTE:** Available since v1.214.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_data_works_project" "example" {
  project_name     = "tf_example"
  display_name     = "tf_example"
  pai_task_enabled = false
}

resource "alicloud_data_works_resource" "example" {
  project_id    = alicloud_data_works_project.example.id
  resource_name = "example_resource"
  spec          = jsonencode({})
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, ForceNew) The ID of the DataWorks workspace to which the resource file belongs.
* `resource_name` - (Optional) The name of the resource.
* `spec` - (Optional, Sensitive) The FlowSpec definition of the resource file in JSON format.
* `resource_file` - (Optional) The resource file, corresponds to the file stream or OSS download URL. This is a write-only field and is not returned on Read.
* `path` - (Optional) The destination path of the resource. This is a write-only field and is not returned on Read.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Data Works Resource. The value formats as `<project_id>:<resource_id>`.
* `resource_id` - The unique identifier of the resource.
* `create_time` - The creation time of the resource.
* `modify_time` - The last modification time of the resource.
* `owner` - The responsible person of the resource file.
* `type` - The resource type. Valid values: `jar`, `python`, `file`, `archive`.
* `source_type` - The file resource source storage type. Valid values: `local`, `oss`.
* `source_path` - The file source path.
* `target_type` - The file destination storage type. Valid values: `gateway`, `oss`, `hdfs`.
* `target_path` - The file destination storage path.
* `data_source` - The data source information. Contains:
  * `name` - The data source name.
  * `type` - The data source type.
* `script` - The script information of the workflow. Contains:
  * `path` - The script path.
  * `runtime_command` - The script type.
  * `script_id` - The script ID.

## Import

Data Works Resource can be imported using the id, e.g.

```shell
$ terraform import alicloud_data_works_resource.example <project_id>:<resource_id>
```
