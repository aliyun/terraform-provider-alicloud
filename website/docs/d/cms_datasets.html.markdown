---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_datasets"
sidebar_current: "docs-alicloud-datasource-cms-datasets"
description: |-
  Provides a list of Cms Datasets to the user.
---

# alicloud\_cms\_datasets

This data source provides the Cms Datasets of the current Alibaba Cloud user.

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

resource "alicloud_cms_dataset" "default" {
  workspace    = alicloud_cms_workspace.default.workspace_name
  dataset_name = var.name
  description  = "terraform-example"
  schema       = jsonencode({ type = "record", name = "example", fields = [{ name = "metric", type = "string" }] })
}

data "alicloud_cms_datasets" "default" {
  workspace = alicloud_cms_workspace.default.workspace_name
  ids       = [alicloud_cms_dataset.default.id]
}
output "cms_dataset_id_1" {
  value = data.alicloud_cms_datasets.default.datasets.0.id
}
```

## Argument Reference

The following arguments are supported:

* `dataset_name_regex` - (Optional, ForceNew) A regex string to filter results by Dataset name.
* `ids` - (Optional, ForceNew, Computed) A list of Dataset IDs. Its element value is formatted as `<workspace>:<dataset_name>`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `workspace` - (Required, ForceNew) The name of the workspace to which the datasets belong.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `datasets` - A list of Cms Datasets. Each element contains the following attributes:
  * `create_time` - The creation time of the resource.
  * `dataset_name` - The name of the resource.
  * `description` - The description of the dataset.
  * `id` - The ID of the resource. It is formatted as `<workspace>:<dataset_name>`.
  * `region_id` - The region ID of the resource.
  * `update_time` - The last modified time of the resource.
  * `workspace` - The name of the workspace to which the dataset belongs.
