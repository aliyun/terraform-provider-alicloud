---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_data_lake_storage_transfer_job"
description: |-
  Provides an Alicloud OSS Data Lake Storage Transfer Job resource.
---

# alicloud_oss_data_lake_storage_transfer_job

Provides an OSS Data Lake Storage Transfer Job resource. A data lake storage transfer job configures the rule that OSS-HDFS uses to replicate objects between a source bucket path and a destination, with optional prefix filtering, verification and a tag.

For information about OSS Data Lake Storage Transfer Job and how to use it, see [Manage data lake storage transfer jobs](https://www.alibabacloud.com/help/en/oss/user-guide/data-lake-storage-transfer).

-> **NOTE:** Available since v2.1.0.

-> **NOTE:** The DataLakeStorageTransferJob plane is served from `cn-hangzhou`. The OSS-HDFS data lake transfer feature must be enabled for your account before the resource can be created; otherwise the create returns `AccessDenied`.

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

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = var.name
}

data "alicloud_caller_identity" "current" {}

resource "alicloud_ram_role" "executor" {
  name     = "${var.name}-role"
  document = <<EOF
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": ["oss.aliyuncs.com"]
      }
    }
  ],
  "Version": "1"
}
EOF
  force    = true
}

resource "alicloud_oss_data_lake_storage_transfer_job" "default" {
  bucket                 = alicloud_oss_bucket.CreateBucket.bucket
  executor_role_id       = alicloud_ram_role.executor.arn
  log_base_dir           = "logs/transfer"
  need_verify            = false
  prefix_filter_includes = ["logs/", "data/"]
  tag                    = "tf-example"
  enable                 = false
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, ForceNew) The name of the bucket to which the data lake storage transfer job belongs.
* `executor_role_id` - (Required) The ARN of the RAM role that OSS-HDFS assumes to run the transfer.
* `log_base_dir` - (Required) The base directory under which the transfer logs are written.
* `prefix_filter_includes` - (Optional) The list of object key prefixes that the job transfers. Objects whose key starts with one of these prefixes are included.
* `need_verify` - (Optional) Whether to verify transferred objects. Defaults to `false`.
* `tag` - (Optional) A tag attached to the job for identification.
* `enable` - (Optional) Whether the job is running. Setting `enable` to `true` starts the job (idle -> running); setting it to `false` stops a running job. It maps to the `StartDataLakeStorageTransferJob` / `StopDataLakeStorageTransferJob` operational actions and is not part of the Rule payload.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. The value is formulated as `<bucket>:<job_id>`.
* `job_id` - The server-generated id of the data lake storage transfer job.
* `status` - The current status of the job, e.g. `REPLICATION_JOB_IDLE` or `RUNNING`.
* `type` - The job type.
* `progress_percent` - The transfer progress percentage reported by the job.
* `history_id` - The id of the latest transfer history record.
* `create_time` - The creation time of the job.
* `last_modify_time` - The last modification time of the job.

## Import

OSS Data Lake Storage Transfer Job can be imported using the id, which consists of bucket and job_id, e.g.

```shell
$ terraform import alicloud_oss_data_lake_storage_transfer_job.example <bucket>:<job_id>
```
