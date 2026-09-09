---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_dataset_download_jobs"
sidebar_current: "docs-alicloud-datasource-cms-dataset-download-jobs"
description: |-
  Provides a list of Cms Dataset Download Jobs to the user.
---

# alicloud\_cms\_dataset\_download\_jobs

This data source provides the Cms Dataset Download Jobs of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.296.0.

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

resource "alicloud_cms_dataset_download_job" "default" {
  workspace        = alicloud_cms_workspace.default.workspace_name
  dataset_name     = alicloud_cms_dataset.default.dataset_name
  job_name         = var.name
  query            = "* | select count(*) as total"
  compression_type = "gzip"
  content_type     = "csv"
}

data "alicloud_cms_dataset_download_jobs" "default" {
  workspace    = alicloud_cms_workspace.default.workspace_name
  dataset_name = alicloud_cms_dataset.default.dataset_name
  ids          = [alicloud_cms_dataset_download_job.default.id]
}
output "cms_dataset_download_job_id_1" {
  value = data.alicloud_cms_dataset_download_jobs.default.download_jobs.0.id
}
```

## Argument Reference

The following arguments are supported:

* `dataset_name` - (Required, ForceNew) The name of the dataset to which the download jobs belong.
* `ids` - (Optional, ForceNew, Computed) A list of Dataset Download Job IDs. Its element value is formatted as `<workspace>:<dataset_name>:<job_name>`.
* `job_name_regex` - (Optional, ForceNew) A regex string to filter results by Job name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `workspace` - (Required, ForceNew) The name of the workspace to which the dataset belongs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `download_jobs` - A list of Cms Dataset Download Jobs. Each element contains the following attributes:
  * `compression_type` - The compression format of the download result.
  * `content_type` - The content type of the download result.
  * `create_time` - The creation time of the resource.
  * `dataset_name` - The name of the dataset to which the download job belongs.
  * `id` - The ID of the resource. It is formatted as `<workspace>:<dataset_name>:<job_name>`.
  * `job_name` - The name of the download job.
  * `query` - The query statement of the download job.
  * `update_time` - The last modified time of the resource.
  * `workspace` - The name of the workspace to which the dataset belongs.
