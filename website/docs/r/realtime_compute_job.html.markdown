---
subcategory: "Realtime Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_realtime_compute_job"
description: |-
  Provides a Alicloud Realtime Compute Job resource.
---

# alicloud_realtime_compute_job

Provides a Realtime Compute Job resource.



For information about Realtime Compute Job and how to use it, see [What is Job](https://next.api.alibabacloud.com/document/ververica/2022-07-18/StartJobWithParams).

-> **NOTE:** Available since v1.264.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vpc" "create_Vpc5" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = "example-tf-vpc-deployment"
}

resource "alicloud_vswitch" "create_Vswitch5" {
  is_default   = false
  vpc_id       = alicloud_vpc.create_Vpc5.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "example-tf-vSwitch-deployment"
}

resource "alicloud_oss_bucket" "create_bucket5" {
}

resource "alicloud_realtime_compute_vvp_instance" "create_VvpInstance5" {
  vvp_instance_name = "code-example-tf-deployment"
  storage {
    oss {
      bucket = alicloud_oss_bucket.create_bucket5.id
    }
  }
  vpc_id      = alicloud_vpc.create_Vpc5.id
  vswitch_ids = ["{{$.create_Vswitch5.VSwitchId}}"]
  resource_spec {
    cpu       = "4"
    memory_gb = "16"
  }
  payment_type = "PayAsYouGo"
}

resource "alicloud_realtime_compute_deployment" "create_Deployment5" {
  deployment_name = "tf-example-deployment-sql-24"
  engine_version  = "vvr-8.0.10-flink-1.17"
  resource_id     = alicloud_realtime_compute_vvp_instance.create_VvpInstance5.resource_id
  execution_mode  = "STREAMING"
  deployment_target {
    mode = "PER_JOB"
    name = "default-queue"
  }
  namespace = "${alicloud_realtime_compute_vvp_instance.create_VvpInstance5.vvp_instance_name}-default"
  artifact {
    kind = "SQLSCRIPT"
    sql_artifact {
      sql_script = "create temporary table `datagen` ( id varchar, name varchar ) with ( \\'connector\\' = \\'datagen\\' );  create temporary table `blackhole` ( id varchar, name varchar ) with ( \\'connector\\' = \\'blackhole\\' );  insert into blackhole select * from datagen;"
    }
  }
}


resource "alicloud_realtime_compute_job" "default" {
  deployment_id       = alicloud_realtime_compute_deployment.create_Deployment5.deployment_id
  resource_queue_name = "default-queue"
  resource_id         = alicloud_realtime_compute_vvp_instance.create_VvpInstance5.resource_id
  local_variables {
    value = "qq"
    name  = "tt"
  }
  restore_strategy {
    kind                 = "NONE"
    job_start_time_in_ms = "1763694521254"
  }
  namespace = "${alicloud_realtime_compute_vvp_instance.create_VvpInstance5.vvp_instance_name}-default"
}
```

## Argument Reference

The following arguments are supported:
* `resource_id` - (Required, ForceNew) resourceId
* `deployment_id` - (Optional, ForceNew) deploymentId
* `local_variables` - (Optional, ForceNew, List) Local variables See [`local_variables`](#local_variables) below.
* `namespace` - (Required, ForceNew) namespace
* `resource_queue_name` - (Optional) Resource Queue for Job Run

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `restore_strategy` - (Optional, ForceNew, List) Restore strategy See [`restore_strategy`](#restore_strategy) below.
* `status` - (Optional, Computed, List) job status See [`status`](#status) below.
* `stop_strategy` - (Optional) Job Stop Policy

-> **NOTE:** This parameter only applies during resource update. If modified in isolation without other property changes, Terraform will not trigger any action.


### `local_variables`

The local_variables supports the following:
* `name` - (Optional, ForceNew) Local variables name
* `value` - (Optional, ForceNew) Local variables value

### `restore_strategy`

The restore_strategy supports the following:
* `allow_non_restored_state` - (Optional, ForceNew) Stateless startup
* `job_start_time_in_ms` - (Optional, ForceNew, Int) Stateless start time. When stateless start is selected, you can set this parameter to enable all source tables that support startTime to read data from this time.
* `kind` - (Optional, ForceNew) Restore type
* `savepoint_id` - (Optional, ForceNew) SavepointId

### `status`

The status supports the following:
* `current_job_status` - (Optional) Job current status

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<resource_id>:<namespace>:<job_id>`.
* `job_id` - The first ID of the resource
* `status` - job status
  * `failure` - Job failure information
    * `failed_at` - Job failure time
    * `message` - Failure Information Details
    * `reason` - Failure Reason
  * `health_score` - Job Run Health Score
  * `risk_level` - Risk level, which indicates the risk level of the operation status of the job.
  * `running` - job running status, which has value when the job is Running.
    * `observed_flink_job_restarts` - Number of job restarts
    * `observed_flink_job_status` - Flink job status

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Job.
* `delete` - (Defaults to 5 mins) Used when delete the Job.
* `update` - (Defaults to 5 mins) Used when update the Job.

## Import

Realtime Compute Job can be imported using the id, e.g.

```shell
$ terraform import alicloud_realtime_compute_job.example <resource_id>:<namespace>:<job_id>
```