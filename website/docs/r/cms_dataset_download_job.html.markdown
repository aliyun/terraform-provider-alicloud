---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_dataset_download_job"
description: |-
  Provides a Alicloud Cms Dataset Download Job resource.
---

# alicloud_cms_dataset_download_job

Provides a Cms Dataset Download Job resource.

For information about Cms Dataset Download Job and how to use it, see [What is DatasetDownloadJob](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateDatasetDownloadJob).

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
```

## Argument Reference

The following arguments are supported:

* `compression_type` - (Optional, ForceNew) The compression format of the download result.
* `content_type` - (Optional, ForceNew) The content type of the download result.
* `dataset_name` - (Required, ForceNew) The name of the dataset to which the download job belongs.
* `job_name` - (Required, ForceNew) The name of the download job.
* `query` - (Required, ForceNew) The query statement of the download job.
* `workspace` - (Required, ForceNew) The name of the workspace to which the dataset belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource supplied above. It is formatted as `<workspace>:<dataset_name>:<job_name>`.
* `create_time` - The creation time of the resource.
* `region_id` - The region ID of the resource.
* `update_time` - The last modified time of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Dataset Download Job.
* `delete` - (Defaults to 5 mins) Used when delete the Dataset Download Job.

## Import

Cms Dataset Download Job can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_dataset_download_job.example <workspace>:<dataset_name>:<job_name>
```
