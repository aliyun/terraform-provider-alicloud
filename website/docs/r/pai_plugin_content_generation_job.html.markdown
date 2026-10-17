---
subcategory: "PAIPlugin"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_plugin_content_generation_job"
description: |-
  Provides a Alicloud PAIPlugin Content Generation Job resource.
---

# alicloud_pai_plugin_content_generation_job

Provides a PAIPlugin Content Generation Job resource. A content generation job
runs against an existing content generation job config.

For information about PAIPlugin Content Generation Job and how to use it, see
[What is Content Generation Job](https://www.alibabacloud.com/help/en/pai/developer-reference/api-paiplugin-2022-01-12-createcontentgenerationjob).

-> **NOTE:** Available since v1.249.0.

## Example Usage

```terraform
resource "alicloud_pai_plugin_content_generation_job" "default" {
  content_generation_job_config_id = alicloud_pai_plugin_content_generation_job_config.default.id
  contents                         = ["content-1"]
}
```

## Argument Reference

The following arguments are supported:

* `content_generation_job_config_id` - (Required, ForceNew) The config id of the content generation job, obtained from the job config.
* `contents` - (Optional) The contents of the content generation job. Updatable via UpdateContentGenerationJob.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. Same as `content_generation_job_id`.
* `content_generation_job_id` - The job id.
* `content_generation_job_config_name` - The config name.
* `status` - The status. Valid values: `Running`, `Failure`, `Success`.
* `results` - The results of the job.
* `signature_id` - The signature id.
* `signature_name` - The signature name.
* `gmt_create_time` - The create time.
* `gmt_modified_time` - The modified time.

-> **NOTE:** PAIPlugin exposes no delete API for content generation jobs; destroying the resource only removes it from Terraform state.

## Import

PAIPlugin Content Generation Job can be imported using the id, e.g.

```shell
terraform import alicloud_pai_plugin_content_generation_job.default <content_generation_job_id>
```
