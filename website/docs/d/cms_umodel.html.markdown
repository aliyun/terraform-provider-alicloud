---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_umodel"
description: |-
  Provides a Alicloud Cms Umodel data source.
---

# alicloud_cms_umodel

Provides a Cms Umodel data source that fetches the umodel configuration of a CMS workspace.

For information about Cms Umodel and how to use it, see [GetUmodel](https://next.api.alibabacloud.com/document/Cms/2024-03-30/GetUmodel).

-> **NOTE:** Available since v1.292.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.name
  sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_umodel" "default" {
  workspace   = alicloud_cms_workspace.default.workspace_name
  description = "terraform-example"
}

data "alicloud_cms_umodel" "default" {
  workspace = alicloud_cms_umodel.default.workspace
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) The name of the workspace whose umodel is to be fetched.
* `output_file` - (Optional) File path where to save the result after running `terraform plan`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the data source. It is the workspace name.
* `description` - The description of the umodel.
* `region_id` - The region ID of the umodel.
* `common_schema_ref` - The common schema references of the umodel.
  * `group` - The group of the common schema reference.
  * `items` - The item list of the common schema reference.
  * `version` - The version of the common schema reference.
