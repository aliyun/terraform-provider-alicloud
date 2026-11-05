---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_dataworks_resource_file"
sidebar_current: "docs-alicloud-resource-dataworks-resource-file"
description: |-
  Provides a Alicloud Data Works Resource File resource.
---

# alicloud\_dataworks\_resource\_file

Provides a Data Works Resource File resource.

For information about Data Works Resource File and how to use it, see [CreateResourceFile](https://next.api.alibabacloud.com/document/dataworks-public/2020-05-18/CreateResourceFile).

-> **NOTE:** Available since v2.1.0.

-> **NOTE:** `DeleteFile` is asynchronous: when the file has been committed to the scheduler it returns a `DeploymentId`, and Terraform polls `GetDeployment` until the deletion deployment succeeds (status `1`) or fails (status `2`). An empty `DeploymentId` means the file was already removed and no polling is performed.

## Example Usage

Basic Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "terraform-example"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_data_works_project" "default" {
  project_name      = var.name
  display_name      = var.name
  pai_task_enabled  = false
  description       = "workspace for resource file"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  status            = "Available"
}

resource "alicloud_dataworks_resource_file" "default" {
  project_id       = alicloud_data_works_project.default.id
  file_name        = var.name
  file_type        = 10
  content          = "SHOW TABLES;"
  file_description = "an example resource file"
  file_folder_path = "Business Flow/tfTestAcc/folderDi"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, ForceNew) The ID of the DataWorks workspace. The resource file is created within this workspace. Changing this creates a new resource.
* `file_name` - (Required) The name of the file.
* `file_type` - (Required, ForceNew) The code type of the file. Different file types use different codes; see the [DataWorks node collection](https://www.alibabacloud.com/help/doc-detail/600169.html). You can also query the code type by calling the `ListFileType` API. Changing this creates a new resource.
* `content` - (Optional) The code (content) of the file.
* `file_description` - (Optional) The description of the file.
* `owner` - (Optional) The owner of the file.
* `file_folder_path` - (Optional) The folder path of the file. This field is writable on create and update but is not returned by `GetFile`, so Terraform keeps the configured value on refresh.
* `origin_resource_name` - (Optional, ForceNew) The name of the original resource. This field is accepted only by `CreateResourceFile` and is not returned by `GetFile`, so Terraform keeps the configured value on refresh. Changing this creates a new resource.
* `register_to_calc_engine` - (Optional, ForceNew) Whether to register the file to the compute engine. This field is accepted only by `CreateResourceFile` and is not returned by `GetFile`, so Terraform keeps the configured value on refresh. Changing this creates a new resource.
* `storage_url` - (Optional, ForceNew) The storage URL of the file. This field is accepted only by `CreateResourceFile` and is not returned by `GetFile`, so Terraform keeps the configured value on refresh. Changing this creates a new resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of the resource file. The value formats as `<file_id>:<project_id>`.
* `file_id` - The ID of the file. It is the server-assigned primary key returned by `CreateResourceFile` and `GetFile`.
* `region_id` - The region ID of the DataWorks workspace.

## Import

Data Works Resource File can be imported using the id, e.g.

```shell
$ terraform import alicloud_dataworks_resource_file.example <file_id>:<project_id>
```
