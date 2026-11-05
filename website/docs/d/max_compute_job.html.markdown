---
subcategory: "Max Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_max_compute_job"
sidebar_current: "docs-alicloud-datasource-max-compute-job"
description: |-
  Provides information about an Alibaba Cloud MaxCompute job.
---

# alicloud_max_compute_job

The `alicloud_max_compute_job` data source provides information about a MaxCompute job instance. MaxCompute jobs are read-only/queryable entities: there is no Create/Update/Delete API for a job, so this data source reads an existing job by its instance ID.

For information on the API and its fields, see [GetJobInfo](https://www.alibabacloud.com/help/en/maxcompute/developer-reference/api-maxcompute-2022-01-04-getjobinfo).

-> **NOTE:** Available since v2.3.0.

## Example Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "instance_id" {
  description = "The instance ID of the MaxCompute job to query."
  type        = string
}

data "alicloud_max_compute_job" "example" {
  instance_id = var.instance_id
}

output "job_status" {
  value = data.alicloud_max_compute_job.example.status
}

output "job_type" {
  value = data.alicloud_max_compute_job.example.job_type
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) The instance ID of the MaxCompute job. This is the `instanceId` path parameter of the GetJobInfo API.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - (Read-Only, String) The ID of the data source, set to the job instance ID.
* `tenant_id` - (Read-Only, String) The ID of the tenant to which the job belongs.
* `region` - (Read-Only, String) The region where the job runs.
* `job_owner` - (Read-Only, String) The owner of the job.
* `project` - (Read-Only, String) The name of the MaxCompute project to which the job belongs.
* `job_type` - (Read-Only, String) The type of the job, such as `SQL`, `MAPREDUCE`, `GRAPH`, `STREAM` and `CUPID`.
* `priority` - (Read-Only, Int) The priority of the job.
* `status` - (Read-Only, String) The status of the job, such as `running`, `terminated` and `failed`.
* `quota_nickname` - (Read-Only, String) The nickname of the quota used by the job.
* `quota_type` - (Read-Only, String) The type of the quota used by the job, such as `subscription` and `payAsYouGo`.
* `cu_usage` - (Read-Only, Int) The number of compute units (CUs) used by the job.
* `memory_usage` - (Read-Only, Int) The memory usage of the job.
* `submitted_at_time` - (Read-Only, Int) The time at which the job was submitted, as a Unix timestamp in seconds.
* `waiting_time` - (Read-Only, Int) The time that the job spent waiting, in seconds.
* `running_at_time` - (Read-Only, Int) The time at which the job started running, as a Unix timestamp in seconds.
* `running_time` - (Read-Only, Int) The time that the job spent running, in seconds.
* `end_at_time` - (Read-Only, Int) The time at which the job ended, as a Unix timestamp in seconds.
* `total_time` - (Read-Only, Int) The total elapsed time of the job, in seconds.
* `signature` - (Read-Only, String) The signature of the job.
* `ext_plant_from` - (Read-Only, String) The external platform from which the job was submitted.
* `ext_node_id` - (Read-Only, String) The external node ID associated with the job.
* `ext_node_on_duty` - (Read-Only, String) The external on-duty node associated with the job.
* `input_bytes` - (Read-Only, Int) The input data size of the job, in bytes.
* `scene_results` - (Read-Only, List) The scene analysis results of the job. Each element contains the following attributes:
  * `scene_tag` - (Read-Only, String) The scene tag.
  * `scene` - (Read-Only, String) The scene name.
  * `type` - (Read-Only, String) The type of the scene result, such as `warning`.
  * `summary` - (Read-Only, String) A short summary of the scene result.
  * `description` - (Read-Only, String) The detailed description of the scene result.
  * `params` - (Read-Only, Map) The key-value parameters of the scene result.
* `job_sub_status_list` - (Read-Only, List) The sub-status list of the job. Each element contains the following attributes:
  * `code` - (Read-Only, Int) The code of the sub-status.
  * `start_time` - (Read-Only, String) The start time of the sub-status.
  * `description` - (Read-Only, String) The description of the sub-status.
