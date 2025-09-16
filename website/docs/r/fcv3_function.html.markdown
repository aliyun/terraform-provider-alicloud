---
subcategory: "Function Compute Service V3 (FCV3)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fcv3_function"
description: |-
  Provides a Alicloud Function Compute Service V3 (FCV3) Function resource.
---

# alicloud_fcv3_function

Provides a Function Compute Service V3 (FCV3) Function resource.

The resource scheduling and running of Function Compute is based on functions. The FC function consists of function code and function configuration.

For information about Function Compute Service V3 (FCV3) Function and how to use it, see [What is Function](https://www.alibabacloud.com/help/en/functioncompute/developer-reference/api-fc-2023-03-30-getfunction).

-> **NOTE:** Available since v1.228.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_fcv3_function&exampleId=9c93b29f-29a3-2a4a-410e-7ff176f24eb199f2f587&activeTab=example&spm=docs.r.fcv3_function.0.9c93b29f29&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shanghai"
}

resource "random_uuid" "default" {
}

resource "alicloud_oss_bucket" "default" {
  bucket = "${var.name}-${random_uuid.default.result}"
}

resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.bucket
  key     = "FCV3Py39.zip"
  content = "print('hello')"
}

resource "alicloud_fcv3_function" "default" {
  description = "Create"
  memory_size = "512"
  layers = [
    "acs:fc:cn-shanghai:official:layers/Python39-Aliyun-SDK/versions/3"
  ]
  timeout   = "3"
  runtime   = "custom.debian10"
  handler   = "index.handler"
  disk_size = "512"
  custom_runtime_config {
    command = [
      "python",
      "-c",
      "example"
    ]
    args = [
      "app.py",
      "xx",
      "x"
    ]
    port = "9000"
    health_check_config {
      http_get_url          = "/ready"
      initial_delay_seconds = "1"
      period_seconds        = "10"
      success_threshold     = "1"
      timeout_seconds       = "1"
      failure_threshold     = "3"
    }

  }

  log_config {
    log_begin_rule = "None"
  }

  code {
    oss_bucket_name = alicloud_oss_bucket.default.bucket
    oss_object_name = alicloud_oss_bucket_object.default.key
    checksum        = "4270285996107335518"
  }

  instance_lifecycle_config {
    initializer {
      timeout = "1"
      handler = "index.init"
    }

    pre_stop {
      timeout = "1"
      handler = "index.stop"
    }

  }

  cpu                  = "0.5"
  instance_concurrency = "2"
  function_name        = "${var.name}-${random_uuid.default.result}"
  environment_variables = {
    "EnvKey" = "EnvVal"
  }
  internet_access = "true"
}
```

## Argument Reference

The following arguments are supported:
* `code` - (Optional, List) Function code ZIP package. code and customContainerConfig. See [`code`](#code) below.
* `cpu` - (Optional, Computed, Float) The CPU specification of the function. The unit is vCPU, which is a multiple of the 0.05 vCPU.
* `custom_container_config` - (Optional, List) The configuration of the custom container runtime. After the configuration is successful, the function can use the custom container image to execute the function. code and customContainerConfig. See [`custom_container_config`](#custom_container_config) below.
* `custom_dns` - (Optional, List) Function custom DNS configuration See [`custom_dns`](#custom_dns) below.
* `custom_runtime_config` - (Optional, List) Customize the runtime configuration. See [`custom_runtime_config`](#custom_runtime_config) below.
* `description` - (Optional) The description of the function. The function compute system does not use this attribute value, but we recommend that you set a concise and clear description for the function.
* `disk_size` - (Optional, Computed, Int) The disk specification of the function, in MB. The optional value is 512 MB or 10240MB.
* `environment_variables` - (Optional, Map) The environment variable set for the function, you can get the value of the environment variable in the function.
* `function_name` - (Optional, ForceNew, Computed) The function name. Consists of uppercase and lowercase letters, digits (0 to 9), underscores (_), and dashes (-). It must begin with an English letter (a ~ z), (A ~ Z), or an underscore (_). Case sensitive. The length is 1~128 characters.
* `gpu_config` - (Optional, List) Function GPU configuration. See [`gpu_config`](#gpu_config) below.
* `handler` - (Required) Function Handler: the call entry for the function compute system to run your function.
* `instance_concurrency` - (Optional, Computed, Int) Maximum instance concurrency.
* `instance_isolation_mode` - (Optional, Available since v1.256.0) Instance isolation mode
* `instance_lifecycle_config` - (Optional, List) Instance lifecycle callback method configuration. See [`instance_lifecycle_config`](#instance_lifecycle_config) below.
* `internet_access` - (Optional, Computed) Allow function to access public network
* `invocation_restriction` - (Optional, List, Available since v1.255.0) Invocation Restriction Detail See [`invocation_restriction`](#invocation_restriction) below.
* `layers` - (Optional, List) The list of layers.
* `log_config` - (Optional, List) The logs generated by the function are written to the configured Logstore. See [`log_config`](#log_config) below.
* `memory_size` - (Optional, Computed, Int) The memory specification of the function. The unit is MB. The memory size is a multiple of 64MB. The minimum value is 128MB and the maximum value is 32GB. At the same time, the ratio of cpu to memorySize (calculated by GB) should be between 1:1 and 1:4.
* `nas_config` - (Optional, Computed, List) NAS configuration. After this parameter is configured, the function can access the specified NAS resource. See [`nas_config`](#nas_config) below.
* `oss_mount_config` - (Optional, Computed, List) OSS mount configuration See [`oss_mount_config`](#oss_mount_config) below.
* `resource_group_id` - (Optional, Computed, Available since v1.260.0) Resource Group ID.
* `role` - (Optional) The user is authorized to the RAM role of function compute. After the configuration, function compute will assume this role to generate temporary access credentials. In the function, you can use the temporary access credentials of the role to access the specified Alibaba cloud service, such as OSS and OTS
* `runtime` - (Required) Function runtime type.
* `session_affinity` - (Optional, Available since v1.256.0) The affinity policy of the function compute call request. To implement the request affinity of the MCP SSE protocol, set it to MCP_SSE. If Cookie affinity is used, it can be set to GENERATED_COOKIE. If Header affinity is used, it can be set to HEADER_FIELD. If it is not set or set to NONE, the affinity effect is not set, and the request is routed according to the default scheduling policy of the function calculation system.
* `session_affinity_config` - (Optional, Available since v1.256.0) When you set the sessionAffinity affinity type, you need to set the relevant affinity configuration. For example, the MCP_SSE affinity needs to fill in the mcpssessionaffinityconfig configuration. The Cookie affinity needs to be filled with the CookieSessionAffinityConfig configuration, and the Header Field affinity needs to be filled with the HeaderFieldSessionAffinityConfig configuration.
* `tags` - (Optional, Map, Available since v1.242.0) The tag of the resource
* `timeout` - (Optional, Computed, Int) The maximum running time of the function, in seconds.
* `vpc_config` - (Optional, Computed, List) VPC configuration. After this parameter is configured, the function can access the specified VPC resources. See [`vpc_config`](#vpc_config) below.

### `code`

The code supports the following:
* `checksum` - (Optional) The CRC-64 value of the function code package.
* `oss_bucket_name` - (Optional) The name of the OSS Bucket that stores the function code ZIP package.
* `oss_object_name` - (Optional) The name of the OSS Object that stores the function code ZIP package.
* `zip_file` - (Optional) The Base 64 encoding of the function code ZIP package.

### `custom_container_config`

The custom_container_config supports the following:
* `acceleration_type` - (Optional, Deprecated since v1.228.0) Whether to enable Image acceleration. Default: The Default value, indicating that image acceleration is enabled. None: indicates that image acceleration is disabled. (Obsolete)
* `acr_instance_id` - (Optional, Deprecated since v1.228.0) ACR Enterprise version Image Repository ID, which must be entered when using ACR Enterprise version image. (Obsolete)
* `command` - (Optional, List) Container startup parameters.
* `entrypoint` - (Optional, List) Container start command.
* `health_check_config` - (Optional, List) Function custom health check configuration See [`health_check_config`](#custom_container_config-health_check_config) below.
* `image` - (Optional) The container Image address.
* `port` - (Optional, Int) The listening port of the HTTP Server when the custom container runs.

### `custom_container_config-health_check_config`

The custom_container_config-health_check_config supports the following:
* `failure_threshold` - (Optional, Int) The health check failure threshold. The system considers the health check failure when the health check fails. The value range is 1~120. The default value is 3.
* `http_get_url` - (Optional) The URL of the container's custom health check.
* `initial_delay_seconds` - (Optional, Int) The delay between the start of the container and the initiation of the health check. Value range 0~120. The default value is 0.
* `period_seconds` - (Optional, Int) Health check cycle. The value range is 1~120. The default value is 3.
* `success_threshold` - (Optional, Int) The threshold for the number of successful health checks. When the threshold is reached, the system considers that the health check is successful. The value range is 1~120. The default value is 1.
* `timeout_seconds` - (Optional, Int) Health check timeout. Value range 1~3. The default value is 1.

### `custom_dns`

The custom_dns supports the following:
* `dns_options` - (Optional, List) List of configuration items in the resolv.conf file. Each item corresponds to a key-value pair in the format of key:value, where the key is required. See [`dns_options`](#custom_dns-dns_options) below.
* `name_servers` - (Optional, List) IP Address List of DNS servers
* `searches` - (Optional, List) DNS search domain list

### `custom_dns-dns_options`

The custom_dns-dns_options supports the following:
* `name` - (Optional) Configuration Item Name
* `value` - (Optional) Configuration Item Value

### `custom_runtime_config`

The custom_runtime_config supports the following:
* `args` - (Optional, List) Instance startup parameters.
* `command` - (Optional, List) Instance start command.
* `health_check_config` - (Optional, List) Function custom health check configuration. See [`health_check_config`](#custom_runtime_config-health_check_config) below.
* `port` - (Optional, Computed, Int) The listening port of the HTTP Server.

### `custom_runtime_config-health_check_config`

The custom_runtime_config-health_check_config supports the following:
* `failure_threshold` - (Optional, Computed, Int) The health check failure threshold. The system considers the health check failure when the health check fails. The value range is 1~120. The default value is 3.
* `http_get_url` - (Optional) The URL of the container's custom health check. No more than 2048 characters in length.
* `initial_delay_seconds` - (Optional, Int) The delay between the start of the container and the initiation of the health check. Value range 0~120. The default value is 0.
* `period_seconds` - (Optional, Int) Health check cycle. The value range is 1~120. The default value is 3.
* `success_threshold` - (Optional, Int) The threshold for the number of successful health checks. When the threshold is reached, the system considers that the health check is successful. The value range is 1~120. The default value is 1.
* `timeout_seconds` - (Optional, Int) Health check timeout. Value range 1~3. The default value is 1.

### `gpu_config`

The gpu_config supports the following:
* `gpu_memory_size` - (Optional, Int) GPU memory specification, unit: MB, multiple of 1024MB
* `gpu_type` - (Optional) GPU card architecture.
  - fc.gpu.tesla.1 indicates the type of the Tesla Architecture Series card of the GPU instance (the same as the NVIDIA T4 card type).
  - fc.gpu.ampere.1 indicates the GPU instance type of Ampere Architecture Series card (same as NVIDIA A10 card type).
  - fc.gpu.ada.1 Indicates the GPU instance Ada Lovelace architecture family card type.

### `instance_lifecycle_config`

The instance_lifecycle_config supports the following:
* `initializer` - (Optional, List) Initializer handler method configuration See [`initializer`](#instance_lifecycle_config-initializer) below.
* `pre_stop` - (Optional, List) PreStop handler method configuration See [`pre_stop`](#instance_lifecycle_config-pre_stop) below.

### `instance_lifecycle_config-initializer`

The instance_lifecycle_config-initializer supports the following:
* `command` - (Optional, List, Available since v1.260.0) Lifecycle Initialization Phase Callback Instructions.
* `handler` - (Optional) The execution entry of the callback method, which is similar to the request handler.
* `timeout` - (Optional, Int) The timeout time of the callback method, in seconds.

### `instance_lifecycle_config-pre_stop`

The instance_lifecycle_config-pre_stop supports the following:
* `handler` - (Optional) The execution entry of the callback method, which is similar to the request handler.
* `timeout` - (Optional, Int) The timeout time of the callback method, in seconds.

### `invocation_restriction`

The invocation_restriction supports the following:
* `disable` - (Optional, Available since v1.255.0) Whether invocation is disabled
* `reason` - (Optional, Available since v1.255.0) Disable Reason

### `log_config`

The log_config supports the following:
* `enable_instance_metrics` - (Optional, Computed) After this feature is enabled, you can view core metrics such as instance-level CPU usage, memory usage, instance network status, and the number of requests within an instance. false: The default value, which means that instance-level metrics are turned off. true: indicates that instance-level metrics are enabled.
* `enable_request_metrics` - (Optional, Computed) After this function is enabled, you can view the time and memory consumed by a call to all functions under this service. false: indicates that request-level metrics are turned off. true: The default value, indicating that request-level metrics are enabled.
* `log_begin_rule` - (Optional, Computed) Log Line First Matching Rules
* `logstore` - (Optional) The Logstore name of log service.
* `project` - (Optional) The name of the log service Project.

### `nas_config`

The nas_config supports the following:
* `group_id` - (Optional, Computed, Int) Group ID
* `mount_points` - (Optional, List) Mount point list See [`mount_points`](#nas_config-mount_points) below.
* `user_id` - (Optional, Computed, Int) Account ID

### `nas_config-mount_points`

The nas_config-mount_points supports the following:
* `enable_tls` - (Optional) Use transport encryption to mount. Note: only general-purpose NAS supports transmission encryption.
* `mount_dir` - (Optional) Local Mount Directory
* `server_addr` - (Optional) NAS server address

### `oss_mount_config`

The oss_mount_config supports the following:
* `mount_points` - (Optional, List) OSS mount point list See [`mount_points`](#oss_mount_config-mount_points) below.

### `oss_mount_config-mount_points`

The oss_mount_config-mount_points supports the following:
* `bucket_name` - (Optional) OSS Bucket name
* `bucket_path` - (Optional) Path of the mounted OSS Bucket
* `endpoint` - (Optional) OSS access endpoint
* `mount_dir` - (Optional) Mount Directory
* `read_only` - (Optional) Read-only

### `vpc_config`

The vpc_config supports the following:
* `security_group_id` - (Optional) Security group ID
* `vswitch_ids` - (Optional, List) Switch List
* `vpc_id` - (Optional) VPC network ID

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `code_size` - The code package size of the function returned by the system, in byte Example : 1024
* `create_time` - The creation time of the function.
* `custom_container_config` - The configuration of the custom container runtime. After the configuration is successful, the function can use the custom container image to execute the function. code and customContainerConfig.
  * `acceleration_info` - (Deprecated since v1.242.0) Image Acceleration Information (Obsolete)
    * `status` - Image Acceleration Status (Deprecated)
  * `resolved_image_uri` - The actual digest version of the deployed Image. The code version specified by this digest is used when the function starts.
* `function_arn` - ARN of function
* `function_id` - The first ID of the resource
* `invocation_restriction` - Invocation Restriction Detail
  * `last_modified_time` - Last modified time of invocation restriction
* `last_modified_time` - Last time the function was Updated
* `last_update_status` - The status of the last function update operation. When the function is created successfully, the value is Successful. Optional values are Successful, Failed, and InProgress.
* `last_update_status_reason` - The reason that caused the last function to update the Operation State to the current value
* `last_update_status_reason_code` - Status code of the reason that caused the last function update operation status to the current value
* `state` - Function Status
* `state_reason` - The reason why the function is in the current state
* `state_reason_code` - The status code of the reason the function is in the current state.
* `tracing_config` - Tracing configuration
  * `params` - Tracing parameters
  * `type` - The tracing protocol type. Currently, only Jaeger is supported.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Function.
* `delete` - (Defaults to 5 mins) Used when delete the Function.
* `update` - (Defaults to 5 mins) Used when update the Function.

## Import

Function Compute Service V3 (FCV3) Function can be imported using the id, e.g.

```shell
$ terraform import alicloud_fcv3_function.example <id>
```