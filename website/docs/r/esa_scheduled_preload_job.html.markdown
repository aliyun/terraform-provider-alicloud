---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_scheduled_preload_job"
description: |-
  Provides a Alicloud ESA Scheduled Preload Job resource.
---

# alicloud_esa_scheduled_preload_job

Provides a ESA Scheduled Preload Job resource.



For information about ESA Scheduled Preload Job and how to use it, see [What is Scheduled Preload Job](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateScheduledPreloadJob).

-> **NOTE:** Available since v1.248.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "default" {
  site_name   = "terraform.cn"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_scheduled_preload_job" "default" {
  insert_way                 = "textBox"
  site_id                    = alicloud_esa_site.default.id
  scheduled_preload_job_name = "example_scheduledpreloadexecution_job"
  url_list                   = "http://example.gositecdn.cn/example/example.txt"
}
```

## Argument Reference

The following arguments are supported:
* `insert_way` - (Required, ForceNew) The method to submit the URLs to be prefetched.
* `oss_url` - (Optional) Preheat OSS files regularly and fill in the OSS file address. Note: The OSS file contains the URL that you need to warm up.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `scheduled_preload_job_name` - (Required, ForceNew) The task name.
* `site_id` - (Required, ForceNew) The site ID.
* `url_list` - (Optional) A list of URLs to be preheated, which is used when uploading a preheated file in the text box mode.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<scheduled_preload_job_id>`.
* `create_time` - The time when the task was created.
* `scheduled_preload_job_id` - The ID of the prefetch task.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Scheduled Preload Job.
* `delete` - (Defaults to 5 mins) Used when delete the Scheduled Preload Job.

## Import

ESA Scheduled Preload Job can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_scheduled_preload_job.example <site_id>:<scheduled_preload_job_id>
```