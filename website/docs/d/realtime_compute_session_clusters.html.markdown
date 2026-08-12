---
subcategory: "Realtime Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_realtime_compute_session_clusters"
sidebar_current: "docs-alicloud-datasource-realtime-compute-session-clusters"
description: |-
  Provides a list of Realtime Compute Session Cluster owned by an Alibaba Cloud account.
---

# alicloud_realtime_compute_session_clusters

This data source provides Realtime Compute Session Cluster available to the user.[What is Session Cluster](https://next.api.alibabacloud.com/document/ververica/2022-07-18/CreateSessionCluster)

-> **NOTE:** Available since v1.289.0.

## Example Usage

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

data "alicloud_realtime_compute_session_clusters" "default" {
  workspace      = alicloud_realtime_compute_vvp_instance.default.resource_id
  namespace      = "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"
  enable_details = true
}

output "session_cluster_id" {
  value = data.alicloud_realtime_compute_session_clusters.default.clusters.0.id
}
```

## Argument Reference

The following arguments are supported:
* `namespace` - (Required) The name of the namespace
* `workspace` - (Required) The ID of the workspace
* `ids` - (Optional, Computed) A list of Session Cluster IDs. The value is formulated as `<workspace>:<namespace>:<session_cluster_name>`.
* `name_regex` - (Optional) A regex string to filter results by Session Cluster name.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Session Cluster IDs.
* `names` - A list of name of Session Clusters.
* `clusters` - A list of Session Cluster Entries. Each element contains the following attributes:
  * `basic_resource_setting` - Basic resource setting.
    * `jobmanager_resource_setting_spec` - JobManager resource setting.
      * `cpu` - CPU.
      * `memory` - Memory.
    * `parallelism` - Parallelism.
    * `taskmanager_resource_setting_spec` - TaskManager resource setting.
      * `cpu` - CPU.
      * `memory` - Memory.
  * `created_at` - Creation time.
  * `creator` - Creator.
  * `creator_name` - Creator name.
  * `deployment_target_name` - The name of the deployment target.
  * `engine_version` - The engine version of the session cluster.
  * `flink_conf` - Flink configurations.
  * `labels` - The labels of the session cluster.
  * `logging` - Logging config.
    * `log4j2_configuration_template` - Custom log template.
    * `log4j_loggers` - log4j config.
      * `logger_level` - Logger level.
      * `logger_name` - Class name of the output log.
    * `log_reserve_policy` - Log reserve policy.
      * `expiration_days` - Expiration days.
      * `open_history` - Enable log saving.
    * `logging_profile` - System log template.
  * `modified_at` - Modification time.
  * `modifier` - Modifier.
  * `modifier_name` - Modifier name.
  * `namespace` - The name of the namespace.
  * `session_cluster_id` - The ID of the session cluster.
  * `session_cluster_name` - The name of the session cluster.
  * `status` - Status of the session cluster.
    * `current_session_cluster_status` - Current status of the session cluster.
    * `failure` - **NOTE:** This field is only available when `enable_details` is `true`. Failure info of the session cluster.
      * `failed_at` - Failure time.
      * `message` - Failure message.
      * `reason` - Failure reason.
    * `running` - **NOTE:** This field is only available when `enable_details` is `true`. Running info of the session cluster.
      * `last_update_time` - Last update time.
      * `reference_deployment_ids` - IDs of deployments running on the session cluster.
      * `started_at` - Start time.
  * `workspace` - The ID of the workspace.
  * `id` - The ID of the resource supplied above.
