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

-> **NOTE:** Available since v1.265.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_realtime_compute_deployment&exampleId=766c53f8-d02a-7fff-35f3-c1cbd11deba50ba16ccd&activeTab=example&spm=docs.r.realtime_compute_deployment.0.766c53f8d0&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

resource "alicloud_vpc" "default" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = "example-tf-vpc-deployment"
}

resource "alicloud_vswitch" "default" {
  is_default   = false
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "example-tf-vSwitch-deployment"
}

resource "alicloud_oss_bucket" "default" {
}

resource "alicloud_realtime_compute_vvp_instance" "default" {
  vvp_instance_name = "code-example-tf-deployment"
  storage {
    oss {
      bucket = alicloud_oss_bucket.default.id
    }
  }
  vpc_id      = alicloud_vpc.default.id
  vswitch_ids = ["${alicloud_vswitch.default.id}"]
  resource_spec {
    cpu       = "8"
    memory_gb = "32"
  }
  payment_type = "PayAsYouGo"
  zone_id      = alicloud_vswitch.default.zone_id
}

resource "alicloud_realtime_compute_deployment" "create_Deployment9" {
  deployment_name = "tf-example-deployment-sql-56"
  engine_version  = "vvr-8.0.10-flink-1.17"
  resource_id     = alicloud_realtime_compute_vvp_instance.default.resource_id
  execution_mode  = "STREAMING"
  deployment_target {
    mode = "PER_JOB"
    name = "default-queue"
  }
  namespace = "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"
  artifact {
    kind = "SQLSCRIPT"
    sql_artifact {
      sql_script = "create temporary table `datagen` ( id varchar, name varchar ) with ( 'connector' = 'datagen' );  create temporary table `blackhole` ( id varchar, name varchar ) with ( 'connector' = 'blackhole' );  insert into blackhole select * from datagen;"
    }
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_realtime_compute_deployment&spm=docs.r.realtime_compute_deployment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `artifact` - (Required, List) Content of the deployment job See [`artifact`](#artifact) below.
* `batch_resource_setting` - (Optional, List) Batch job resource settings See [`batch_resource_setting`](#batch_resource_setting) below.
* `deployment_name` - (Required) Name of the deployment
* `deployment_target` - (Required, List) Deployment target See [`deployment_target`](#deployment_target) below.
* `description` - (Optional) Description
* `engine_version` - (Optional, Computed) Engine version of the deployment instance
* `execution_mode` - (Required, ForceNew) Execution mode. Valid values: STREAMING or BATCH.
* `flink_conf` - (Optional, Computed, Map) Flink configuration
* `labels` - (Optional, Map) Deployment labels
* `local_variables` - (Optional, Set) Job variables See [`local_variables`](#local_variables) below.
* `logging` - (Optional, Computed, List) Job log configuration   See [`logging`](#logging) below.
* `namespace` - (Required, ForceNew) Namespace name
* `resource_id` - (Optional, ForceNew, Computed) Workspace resource ID
* `streaming_resource_setting` - (Optional, Computed, List) Resource settings for streaming mode See [`streaming_resource_setting`](#streaming_resource_setting) below.

### `artifact`

The artifact supports the following:
* `jar_artifact` - (Optional, List) JarArtifact See [`jar_artifact`](#artifact-jar_artifact) below.
* `kind` - (Required, ForceNew) Artifact type
* `python_artifact` - (Optional, List) PythonArtifact See [`python_artifact`](#artifact-python_artifact) below.
* `sql_artifact` - (Optional, List) SqlArtifact See [`sql_artifact`](#artifact-sql_artifact) below.

### `artifact-jar_artifact`

The artifact-jar_artifact supports the following:
* `additional_dependencies` - (Optional, List) Full URL paths of additional dependencies; you can specify other dependencies required by the JAR here
* `entry_class` - (Optional) Main class; you must specify the fully qualified class name
* `jar_uri` - (Optional) Full URL path of the JAR job
* `main_args` - (Optional) Arguments required by the main class

### `artifact-python_artifact`

The artifact-python_artifact supports the following:
* `additional_dependencies` - (Optional, List) Full URL paths of additional dependencies  
* `additional_python_archives` - (Optional, List) URL paths of dependent Python archive files  
* `additional_python_libraries` - (Optional, List) URL paths of dependent Python library files  
* `entry_module` - (Optional) Entry module for Python
* `main_args` - (Optional) Startup arguments
* `python_artifact_uri` - (Optional) Full URL path of the Python job

### `artifact-sql_artifact`

The artifact-sql_artifact supports the following:
* `additional_dependencies` - (Optional, List) Full URL path of additional files. If you need to use dependencies such as UDFs, connectors, or formats that are not registered on the VVP platform, you must add them using this method. Dependencies already registered on the platform do not require this approach.
* `sql_script` - (Optional) Text content of the SQL job

### `batch_resource_setting`

The batch_resource_setting supports the following:
* `basic_resource_setting` - (Optional, List) Resource settings for basic mode See [`basic_resource_setting`](#batch_resource_setting-basic_resource_setting) below.
* `max_slot` - (Optional, Int) Maximum number of slots

### `batch_resource_setting-basic_resource_setting`

The batch_resource_setting-basic_resource_setting supports the following:
* `jobmanager_resource_setting_spec` - (Optional, List) JobManager resource settings See [`jobmanager_resource_setting_spec`](#batch_resource_setting-basic_resource_setting-jobmanager_resource_setting_spec) below.
* `parallelism` - (Optional, Int) Parallelism.
* `taskmanager_resource_setting_spec` - (Optional, List) TaskManager resource settings See [`taskmanager_resource_setting_spec`](#batch_resource_setting-basic_resource_setting-taskmanager_resource_setting_spec) below.

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
* `mode` - (Required) Deployment mode, valid values: PER_JOB or SESSION
* `name` - (Required) Deployment target name

### `local_variables`

The local_variables supports the following:
* `name` - (Optional) Job variable name
* `value` - (Optional) Job variable value

### `logging`

The logging supports the following:
* `log4j2_configuration_template` - (Optional, Computed) Custom log template  
* `log4j_loggers` - (Optional, Computed, Set) log4j configuration   See [`log4j_loggers`](#logging-log4j_loggers) below.
* `log_reserve_policy` - (Optional, Computed, List) Log retention policy   See [`log_reserve_policy`](#logging-log_reserve_policy) below.
* `logging_profile` - (Optional, Computed) Default system log template  

### `logging-log4j_loggers`

The logging-log4j_loggers supports the following:
* `logger_level` - (Optional, Computed) Log output level  
* `logger_name` - (Optional, Computed) Class name for log output  

### `logging-log_reserve_policy`

The logging-log_reserve_policy supports the following:
* `expiration_days` - (Optional, Computed, Int) Number of days to retain logs after log retention is enabled
* `open_history` - (Optional, Computed) Whether to enable log retention

### `streaming_resource_setting`

The streaming_resource_setting supports the following:
* `basic_resource_setting` - (Optional, Computed, List) Resource settings for basic mode See [`basic_resource_setting`](#streaming_resource_setting-basic_resource_setting) below.
* `expert_resource_setting` - (Optional, Computed, List) Expert mode resource settings See [`expert_resource_setting`](#streaming_resource_setting-expert_resource_setting) below.
* `resource_setting_mode` - (Optional, Computed) Resource mode used in streaming mode, valid values: BASIC or EXPERT

### `streaming_resource_setting-basic_resource_setting`

The streaming_resource_setting-basic_resource_setting supports the following:
* `jobmanager_resource_setting_spec` - (Optional, Computed, List) JobManager resource settings See [`jobmanager_resource_setting_spec`](#streaming_resource_setting-basic_resource_setting-jobmanager_resource_setting_spec) below.
* `parallelism` - (Optional, Computed, Int) Parallelism
* `taskmanager_resource_setting_spec` - (Optional, Computed, List) TaskManager resource settings See [`taskmanager_resource_setting_spec`](#streaming_resource_setting-basic_resource_setting-taskmanager_resource_setting_spec) below.

### `streaming_resource_setting-expert_resource_setting`

The streaming_resource_setting-expert_resource_setting supports the following:
* `jobmanager_resource_setting_spec` - (Optional, List) Basic resource settings for JobManager See [`jobmanager_resource_setting_spec`](#streaming_resource_setting-expert_resource_setting-jobmanager_resource_setting_spec) below.
* `resource_plan` - (Optional) Resource plan for expert mode

### `streaming_resource_setting-expert_resource_setting-jobmanager_resource_setting_spec`

The streaming_resource_setting-expert_resource_setting-jobmanager_resource_setting_spec supports the following:
* `cpu` - (Optional, Float) CPU
* `memory` - (Optional) Memory

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
* `deployment_id` - Resource property field representing the primary resource ID

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