---
subcategory: "Realtime Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_realtime_compute_deployment"
description: |-
  Provides a Alicloud Realtime Compute Deployment resource.
---

# alicloud_realtime_compute_deployment

Provides a Realtime Compute Deployment resource.

Deployment in the Realtime Compute console.

For information about Realtime Compute Deployment and how to use it, see [What is Deployment](https://next.api.alibabacloud.com/document/ververica/2022-07-18/CreateDeployment).

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

resource "alicloud_vpc" "create_Vpc" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = "example-tf-vpc-deployment"
}

resource "alicloud_vswitch" "create_Vswitch" {
  is_default   = false
  vpc_id       = alicloud_vpc.create_Vpc.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "example-tf-vSwitch-deployment"
}

resource "alicloud_oss_bucket" "create_bucket" {
}

resource "alicloud_realtime_compute_vvp_instance" "create_VvpInstance" {
  vvp_instance_name = "code-example-tf-deployment"
  storage {
    oss {
      bucket = alicloud_oss_bucket.create_bucket.id
    }
  }
  vpc_id      = alicloud_vpc.create_Vpc.id
  vswitch_ids = ["${alicloud_vswitch.create_Vswitch.id}"]
  resource_spec {
    cpu       = "4"
    memory_gb = "16"
  }
  payment_type = "PayAsYouGo"
}


resource "alicloud_realtime_compute_deployment" "default" {
  logging {
    logging_profile = "default"
    log4j_loggers {
      logger_level = "INFO"
    }
    log_reserve_policy {
      open_history    = true
      expiration_days = "7"
    }
  }
  deployment_name = "tf-example-deployment-sql-74"
  description     = "This is a example deployment."
  engine_version  = "vvr-8.0.6-flink-1.17"
  local_variables {
    value = "value"
    name  = "name"
  }
  execution_mode = "STREAMING"
  labels {
  }
  deployment_target {
    mode = "PER_JOB"
    name = "default-queue"
  }
  streaming_resource_setting {
    basic_resource_setting {
      taskmanager_resource_setting_spec {
        memory = "1Gi"
        cpu    = 1
      }
      parallelism = "1"
      jobmanager_resource_setting_spec {
        memory = "1Gi"
        cpu    = 1
      }
    }
    resource_setting_mode = "BASIC"
  }
  namespace = "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.vvp_instance_name}-default"
  artifact {
    kind = "SQLSCRIPT"
    sql_artifact {
      sql_script              = "create temporary table `datagen` ( id varchar, name varchar ) with ( \\'connector\\' = \\'datagen\\' );  create temporary table `blackhole` ( id varchar, name varchar ) with ( \\'connector\\' = \\'blackhole\\' );  insert into blackhole select * from datagen;"
      additional_dependencies = ["oss://bucket-name/a.jar"]
    }
  }
  resource_id = alicloud_realtime_compute_vvp_instance.create_VvpInstance.resource_id
  flink_conf {
  }
}
```

## Argument Reference

The following arguments are supported:
* `artifact` - (Required, List) The content of deployment See [`artifact`](#artifact) below.
* `batch_resource_setting` - (Optional, List) batch resource setting See [`batch_resource_setting`](#batch_resource_setting) below.
* `deployment_name` - (Required) Name of the deployment
* `deployment_target` - (Required, List) Deployment target See [`deployment_target`](#deployment_target) below.
* `description` - (Optional) The description of deployment
* `engine_version` - (Optional, Computed) The engine version of the deployment
* `execution_mode` - (Required, ForceNew) Execution mode，STREAMING/BATCH
* `flink_conf` - (Optional, Computed, Map) flink configurations
* `labels` - (Optional, Map) deployment label
* `local_variables` - (Optional, Set) Local variables See [`local_variables`](#local_variables) below.
* `logging` - (Optional, Computed, List) logging config See [`logging`](#logging) below.
* `namespace` - (Required, ForceNew) The name of vvpnamespace.
* `resource_id` - (Optional, ForceNew, Computed) ResourceId of vvpinstance
* `streaming_resource_setting` - (Optional, Computed, List) streaming resource setting See [`streaming_resource_setting`](#streaming_resource_setting) below.

### `artifact`

The artifact supports the following:
* `jar_artifact` - (Optional, List) Jar artifact See [`jar_artifact`](#artifact-jar_artifact) below.
* `kind` - (Required, ForceNew) Artifact kind
* `python_artifact` - (Optional, List) Python artifact See [`python_artifact`](#artifact-python_artifact) below.
* `sql_artifact` - (Optional, List) Sql artifact See [`sql_artifact`](#artifact-sql_artifact) below.

### `artifact-jar_artifact`

The artifact-jar_artifact supports the following:
* `additional_dependencies` - (Optional, List) The additional dependencies of jar
* `entry_class` - (Optional) entry class of jar
* `jar_uri` - (Optional) The url of JAR
* `main_args` - (Optional) Main args of jar

### `artifact-python_artifact`

The artifact-python_artifact supports the following:
* `additional_dependencies` - (Optional, List) The url of  additional dependencies 
* `additional_python_archives` - (Optional, List) The url of python archives
* `additional_python_libraries` - (Optional, List) The url of python lib
* `entry_module` - (Optional) Entry module
* `main_args` - (Optional) Main args
* `python_artifact_uri` - (Optional) The url of python artifact

### `artifact-sql_artifact`

The artifact-sql_artifact supports the following:
* `additional_dependencies` - (Optional, List) The url of additional dependencies
* `sql_script` - (Optional) Sql script

### `batch_resource_setting`

The batch_resource_setting supports the following:
* `basic_resource_setting` - (Optional, List) basic resource setting See [`basic_resource_setting`](#batch_resource_setting-basic_resource_setting) below.
* `max_slot` - (Optional, Int) max slot

### `batch_resource_setting-basic_resource_setting`

The batch_resource_setting-basic_resource_setting supports the following:
* `jobmanager_resource_setting_spec` - (Optional, List) JobManager  resource setting See [`jobmanager_resource_setting_spec`](#batch_resource_setting-basic_resource_setting-jobmanager_resource_setting_spec) below.
* `parallelism` - (Optional, Int) Parallelism
* `taskmanager_resource_setting_spec` - (Optional, List) TaskManager resource setting See [`taskmanager_resource_setting_spec`](#batch_resource_setting-basic_resource_setting-taskmanager_resource_setting_spec) below.

### `batch_resource_setting-basic_resource_setting-jobmanager_resource_setting_spec`

The batch_resource_setting-basic_resource_setting-jobmanager_resource_setting_spec supports the following:
* `cpu` - (Optional, Float) CPU
* `memory` - (Optional) Memory

### `batch_resource_setting-basic_resource_setting-taskmanager_resource_setting_spec`

The batch_resource_setting-basic_resource_setting-taskmanager_resource_setting_spec supports the following:
* `cpu` - (Optional, Float) CPU
* `memory` - (Optional) Memory

### `deployment_target`

The deployment_target supports the following:
* `mode` - (Required) deployment mode
* `name` - (Required) target name

### `local_variables`

The local_variables supports the following:
* `name` - (Optional) Local variable name
* `value` - (Optional) Local variable value

### `logging`

The logging supports the following:
* `log4j2_configuration_template` - (Optional, Computed) Custom Log Template
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

### `streaming_resource_setting`

The streaming_resource_setting supports the following:
* `basic_resource_setting` - (Optional, Computed, List) resource setting in basic mode See [`basic_resource_setting`](#streaming_resource_setting-basic_resource_setting) below.
* `expert_resource_setting` - (Optional, Computed, List) expert resource setting See [`expert_resource_setting`](#streaming_resource_setting-expert_resource_setting) below.
* `resource_setting_mode` - (Optional, Computed) Resource setting mode

### `streaming_resource_setting-basic_resource_setting`

The streaming_resource_setting-basic_resource_setting supports the following:
* `jobmanager_resource_setting_spec` - (Optional, Computed, List) JobManager resource setting See [`jobmanager_resource_setting_spec`](#streaming_resource_setting-basic_resource_setting-jobmanager_resource_setting_spec) below.
* `parallelism` - (Optional, Computed, Int) Parallelism
* `taskmanager_resource_setting_spec` - (Optional, Computed, List) TaskManager basic resource setting See [`taskmanager_resource_setting_spec`](#streaming_resource_setting-basic_resource_setting-taskmanager_resource_setting_spec) below.

### `streaming_resource_setting-expert_resource_setting`

The streaming_resource_setting-expert_resource_setting supports the following:
* `jobmanager_resource_setting_spec` - (Optional, List) JobManager resource setting See [`jobmanager_resource_setting_spec`](#streaming_resource_setting-expert_resource_setting-jobmanager_resource_setting_spec) below.
* `resource_plan` - (Optional) resource plan in expert mode

### `streaming_resource_setting-expert_resource_setting-jobmanager_resource_setting_spec`

The streaming_resource_setting-expert_resource_setting-jobmanager_resource_setting_spec supports the following:
* `cpu` - (Optional, Float) CPU
* `memory` - (Optional) 内存

### `streaming_resource_setting-basic_resource_setting-jobmanager_resource_setting_spec`

The streaming_resource_setting-basic_resource_setting-jobmanager_resource_setting_spec supports the following:
* `cpu` - (Optional, Computed, Float) CPU
* `memory` - (Optional, Computed) Memory

### `streaming_resource_setting-basic_resource_setting-taskmanager_resource_setting_spec`

The streaming_resource_setting-basic_resource_setting-taskmanager_resource_setting_spec supports the following:
* `cpu` - (Optional, Computed, Float) CPU
* `memory` - (Optional, Computed) Memory

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<resource_id>:<namespace>:<deployment_id>`.
* `deployment_id` - The first ID of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Deployment.
* `delete` - (Defaults to 5 mins) Used when delete the Deployment.
* `update` - (Defaults to 5 mins) Used when update the Deployment.

## Import

Realtime Compute Deployment can be imported using the id, e.g.

```shell
$ terraform import alicloud_realtime_compute_deployment.example <resource_id>:<namespace>:<deployment_id>
```