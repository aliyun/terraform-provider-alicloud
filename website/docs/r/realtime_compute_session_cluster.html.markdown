---
subcategory: "Realtime Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_realtime_compute_session_cluster"
description: |-
  Provides a Alicloud Realtime Compute Session Cluster resource.
---

# alicloud_realtime_compute_session_cluster

Provides a Realtime Compute Session Cluster resource.

Session cluster of Realtime Compute for Apache Flink.

For information about Realtime Compute Session Cluster and how to use it, see [What is Session Cluster](https://next.api.alibabacloud.com/document/ververica/2022-07-18/CreateSessionCluster).

-> **NOTE:** Available since v1.289.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tfexample-session-cluster"
}

resource "alicloud_vpc" "default" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default" {
  is_default   = false
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
}

resource "alicloud_oss_bucket" "default" {
}

resource "alicloud_realtime_compute_vvp_instance" "default" {
  vvp_instance_name = var.name
  storage {
    oss {
      bucket = alicloud_oss_bucket.default.id
    }
  }
  vpc_id      = alicloud_vpc.default.id
  vswitch_ids = [alicloud_vswitch.default.id]
  resource_spec {
    cpu       = "4"
    memory_gb = "16"
  }
  payment_type = "PayAsYouGo"
  zone_id      = alicloud_vswitch.default.zone_id
}

resource "alicloud_realtime_compute_session_cluster" "default" {
  workspace              = alicloud_realtime_compute_vvp_instance.default.resource_id
  namespace              = "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"
  session_cluster_name   = var.name
  deployment_target_name = "default-queue"
  engine_version         = "vvr-8.0.11-flink-1.17"
  basic_resource_setting {
    parallelism = 1
    jobmanager_resource_setting_spec {
      cpu    = 1
      memory = "1Gi"
    }
    taskmanager_resource_setting_spec {
      cpu    = 1
      memory = "1Gi"
    }
  }
}
```

## Argument Reference

The following arguments are supported:
* `basic_resource_setting` - (Required, List) Basic resource setting See [`basic_resource_setting`](#basic_resource_setting) below.
* `deployment_target_name` - (Required, ForceNew) The name of the deployment target
* `engine_version` - (Required, ForceNew) The engine version of the session cluster
* `flink_conf` - (Optional, Computed, Map) Flink configurations
* `labels` - (Optional, Map) The labels of the session cluster
* `logging` - (Optional, Computed, List) Logging config See [`logging`](#logging) below.
* `namespace` - (Required, ForceNew) The name of the namespace
* `session_cluster_name` - (Required, ForceNew) The name of the session cluster
* `workspace` - (Optional, ForceNew, Computed) The ID of the workspace

### `basic_resource_setting`

The basic_resource_setting supports the following:
* `jobmanager_resource_setting_spec` - (Required, List) JobManager resource setting See [`jobmanager_resource_setting_spec`](#basic_resource_setting-jobmanager_resource_setting_spec) below.
* `parallelism` - (Required, Int) Parallelism
* `taskmanager_resource_setting_spec` - (Required, List) TaskManager resource setting See [`taskmanager_resource_setting_spec`](#basic_resource_setting-taskmanager_resource_setting_spec) below.

### `basic_resource_setting-jobmanager_resource_setting_spec`

The basic_resource_setting-jobmanager_resource_setting_spec supports the following:
* `cpu` - (Required, Float) CPU
* `memory` - (Required) Memory

### `basic_resource_setting-taskmanager_resource_setting_spec`

The basic_resource_setting-taskmanager_resource_setting_spec supports the following:
* `cpu` - (Required, Float) CPU
* `memory` - (Required) Memory

### `logging`

The logging supports the following:
* `log4j2_configuration_template` - (Optional, Computed) Custom log template
* `log4j_loggers` - (Optional, Computed, Set) log4j config See [`log4j_loggers`](#logging-log4j_loggers) below.
* `log_reserve_policy` - (Optional, Computed, List) Log reserve policy See [`log_reserve_policy`](#logging-log_reserve_policy) below.
* `logging_profile` - (Optional, Computed) System log template

### `logging-log4j_loggers`

The logging-log4j_loggers supports the following:
* `logger_level` - (Optional, Computed) Logger level
* `logger_name` - (Optional, Computed) Class name of the output log

### `logging-log_reserve_policy`

The logging-log_reserve_policy supports the following:
* `expiration_days` - (Optional, Computed, Int) Expiration days
* `open_history` - (Optional, Computed) Enable log saving

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<workspace>:<namespace>:<session_cluster_name>`.
* `created_at` - Creation time.
* `creator` - Creator.
* `creator_name` - Creator name.
* `modified_at` - Modification time.
* `modifier` - Modifier.
* `modifier_name` - Modifier name.
* `session_cluster_id` - The ID of the session cluster.
* `status` - Status of the session cluster.
  * `current_session_cluster_status` - Current status of the session cluster.
  * `failure` - Failure info of the session cluster.
    * `failed_at` - Failure time.
    * `message` - Failure message.
    * `reason` - Failure reason.
  * `running` - Running info of the session cluster.
    * `last_update_time` - Last update time.
    * `reference_deployment_ids` - IDs of deployments running on the session cluster.
    * `started_at` - Start time.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Session Cluster.
* `delete` - (Defaults to 5 mins) Used when delete the Session Cluster.
* `update` - (Defaults to 5 mins) Used when update the Session Cluster.

## Import

Realtime Compute Session Cluster can be imported using the id, e.g.

```shell
$ terraform import alicloud_realtime_compute_session_cluster.example <workspace>:<namespace>:<session_cluster_name>
```
