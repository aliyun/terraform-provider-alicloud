---
subcategory: "Function Compute Service V3 (FCV3)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fcv3_functions"
sidebar_current: "docs-alicloud-datasource-fcv3-functions"
description: |-
  Provides a list of Fcv3 Function owned by an Alibaba Cloud account.
---

# alicloud_fcv3_functions

This data source provides Fcv3 Function available to the user.[What is Function](https://next.api.alibabacloud.com/document/FC/2023-03-30/CreateFunction)

-> **NOTE:** Available since v1.264.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_fcv3_functions" "default" {
  prefix = "terraform-example-for-function-alias"
}

output "alicloud_fcv3_function_example_id" {
  value = data.alicloud_fcv3_functions.default.functions.0.function_name
}
```

## Argument Reference

The following arguments are supported:
* `resource_group_id` - (ForceNew, Optional) Resource Group ID
* `ids` - (Optional, ForceNew, Computed) A list of Function IDs. 
* `name_regex` - (Optional, ForceNew) A regex string to filter results by function name.
* `prefix` - (Optional, ForceNew) A prefix string to filter results by function name.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Function IDs.
* `names` - A list of name of Functions.
* `functions` - A list of Function Entries. Each element contains the following attributes:
  * `code_size` - The code package size of the function returned by the system, in byte Example : 1024
  * `cpu` - The CPU specification of the function. The unit is vCPU, which is a multiple of the 0.05 vCPU.
  * `create_time` - The creation time of the function.
  * `custom_container_config` - The configuration of the custom container runtime. After the configuration is successful, the function can use the custom container image to execute the function. code and customContainerConfig.
    * `acceleration_info` - Image Acceleration Information (Obsolete).
      * `status` - Image Acceleration Status (Deprecated).
    * `acceleration_type` - Whether to enable Image acceleration. Default: The Default value, indicating that image acceleration is enabled. None: indicates that image acceleration is disabled. (Obsolete).
    * `acr_instance_id` - ACR Enterprise version Image Repository ID, which must be entered when using ACR Enterprise version image. (Obsolete).
    * `command` - Container startup parameters.
    * `entrypoint` - Container start command.
    * `health_check_config` - Function custom health check configuration.
      * `failure_threshold` - The health check failure threshold. The system considers the health check failure when the health check fails. The value range is 1~120. The default value is 3.
      * `http_get_url` - The URL of the container's custom health check.
      * `initial_delay_seconds` - The delay between the start of the container and the initiation of the health check. Value range 0~120. The default value is 0.
      * `period_seconds` - Health check cycle. The value range is 1~120. The default value is 3.
      * `success_threshold` - The threshold for the number of successful health checks. When the threshold is reached, the system considers that the health check is successful. The value range is 1~120. The default value is 1.
      * `timeout_seconds` - Health check timeout. Value range 1~3. The default value is 1.
    * `image` - The container Image address.
    * `port` - The listening port of the HTTP Server when the custom container runs.
    * `resolved_image_uri` - The actual digest version of the deployed Image. The code version specified by this digest is used when the function starts.
  * `custom_dns` - Function custom DNS configuration
    * `dns_options` - List of configuration items in the resolv.conf file. Each item corresponds to a key-value pair in the format of key:value, where the key is required.
      * `name` - Configuration Item Name.
      * `value` - Configuration Item Value.
    * `name_servers` - IP Address List of DNS servers.
    * `searches` - DNS search domain list.
  * `custom_runtime_config` - Customize the runtime configuration.
    * `args` - Instance startup parameters.
    * `command` - Instance start command.
    * `health_check_config` - Function custom health check configuration.
      * `failure_threshold` - The health check failure threshold. The system considers the health check failure when the health check fails. The value range is 1~120. The default value is 3.
      * `http_get_url` - The URL of the container's custom health check. No more than 2048 characters in length.
      * `initial_delay_seconds` - The delay between the start of the container and the initiation of the health check. Value range 0~120. The default value is 0.
      * `period_seconds` - Health check cycle. The value range is 1~120. The default value is 3.
      * `success_threshold` - The threshold for the number of successful health checks. When the threshold is reached, the system considers that the health check is successful. The value range is 1~120. The default value is 1.
      * `timeout_seconds` - Health check timeout. Value range 1~3. The default value is 1.
    * `port` - The listening port of the HTTP Server.
  * `description` - The description of the function. The function compute system does not use this attribute value, but we recommend that you set a concise and clear description for the function.
  * `disk_size` - The disk specification of the function, in MB. The optional value is 512 MB or 10240MB.
  * `environment_variables` - The environment variable set for the function, you can get the value of the environment variable in the function.
  * `function_arn` - ARN of function
  * `function_id` - The first ID of the resource
  * `function_name` - The function name. Consists of uppercase and lowercase letters, digits (0 to 9), underscores (_), and dashes (-). It must begin with an English letter (a ~ z), (A ~ Z), or an underscore (_). Case sensitive. The length is 1~128 characters.
  * `gpu_config` - Function GPU configuration.
    * `gpu_memory_size` - GPU memory specification, unit: MB, multiple of 1024MB.
    * `gpu_type` - GPU card architecture.-fc.gpu.tesla.1 indicates the type of the Tesla Architecture Series card of the GPU instance (the same as the NVIDIA T4 card type).-fc.gpu.ampere.1 indicates the GPU instance type of Ampere Architecture Series card (same as NVIDIA A10 card type).-fc.gpu.ada.1 Indicates the GPU instance Ada Lovelace architecture family card type.
  * `handler` - Function Handler: the call entry for the function compute system to run your function.
  * `idle_timeout` - Destroy an instance when the instance no-request duration exceeds this attribute. -1 means that the threshold is cleared and the system default behavior is used.
  * `instance_concurrency` - Maximum instance concurrency.
  * `instance_isolation_mode` - Instance isolation mode
  * `instance_lifecycle_config` - Instance lifecycle callback method configuration.
    * `initializer` - Initializer handler method configuration.
      * `command` - Lifecycle Initialization Phase Callback Instructions.
      * `handler` - The execution entry of the callback method, which is similar to the request handler.
      * `timeout` - The timeout time of the callback method, in seconds.
    * `pre_stop` - PreStop handler method configuration.
      * `handler` - The execution entry of the callback method, which is similar to the request handler.
      * `timeout` - The timeout time of the callback method, in seconds.
  * `internet_access` - Allow function to access public network
  * `invocation_restriction` - Invocation Restriction Detail
    * `disable` - Whether invocation is disabled.
    * `last_modified_time` - Last modified time of invocation restriction.
    * `reason` - Disable Reason.
  * `last_modified_time` - Last time the function was Updated
  * `last_update_status` - The status of the last function update operation. When the function is created successfully, the value is Successful. Optional values are Successful, Failed, and InProgress.
  * `last_update_status_reason` - The reason that caused the last function to update the Operation State to the current value
  * `last_update_status_reason_code` - Status code of the reason that caused the last function update operation status to the current value
  * `layers` - The list of layers.
  * `log_config` - The logs generated by the function are written to the configured Logstore.
    * `enable_instance_metrics` - After this feature is enabled, you can view core metrics such as instance-level CPU usage, memory usage, instance network status, and the number of requests within an instance. false: The default value, which means that instance-level metrics are turned off. true: indicates that instance-level metrics are enabled.
    * `enable_request_metrics` - After this function is enabled, you can view the time and memory consumed by a call to all functions under this service. false: indicates that request-level metrics are turned off. true: The default value, indicating that request-level metrics are enabled.
    * `log_begin_rule` - Log Line First Matching Rules.
    * `logstore` - The Logstore name of log service.
    * `project` - The name of the log service Project.
  * `memory_size` - The memory specification of the function. The unit is MB. The memory size is a multiple of 64MB. The minimum value is 128MB and the maximum value is 32GB. At the same time, the ratio of cpu to memorySize (calculated by GB) should be between 1:1 and 1:4.
  * `nas_config` - NAS configuration. After this parameter is configured, the function can access the specified NAS resource.
    * `group_id` - Group ID.
    * `mount_points` - Mount point list.
      * `enable_tls` - Use transport encryption to mount. Note: only general-purpose NAS supports transmission encryption.
      * `mount_dir` - Local Mount Directory.
      * `server_addr` - NAS server address.
    * `user_id` - Account ID.
  * `oss_mount_config` - OSS mount configuration
    * `mount_points` - OSS mount point list.
      * `bucket_name` - OSS Bucket name.
      * `bucket_path` - Path of the mounted OSS Bucket.
      * `endpoint` - OSS access endpoint.
      * `mount_dir` - Mount Directory.
      * `read_only` - Read-only.
  * `resource_group_id` - Resource Group ID
  * `role` - The user is authorized to the RAM role of function compute. After the configuration, function compute will assume this role to generate temporary access credentials. In the function, you can use the temporary access credentials of the role to access the specified Alibaba cloud service, such as OSS and OTS
  * `runtime` - Function runtime type
  * `session_affinity` - The affinity policy of the function compute call request. To implement the request affinity of the MCP SSE protocol, set it to MCP_SSE. If Cookie affinity is used, it can be set to GENERATED_COOKIE. If Header affinity is used, it can be set to HEADER_FIELD. If it is not set or set to NONE, the affinity effect is not set, and the request is routed according to the default scheduling policy of the function calculation system.
  * `session_affinity_config` - When you set the sessionAffinity affinity type, you need to set the relevant affinity configuration. For example, the MCP_SSE affinity needs to fill in the mcpssessionaffinityconfig configuration. The Cookie affinity needs to be filled with the CookieSessionAffinityConfig configuration, and the Header Field affinity needs to be filled with the HeaderFieldSessionAffinityConfig configuration.
  * `state` - Function Status
  * `state_reason` - The reason why the function is in the current state
  * `state_reason_code` - The status code of the reason the function is in the current state.
  * `tags` - The tag of the resource
  * `timeout` - The maximum running time of the function, in seconds.
  * `tracing_config` - Tracing configuration
    * `params` - Tracing parameters.
    * `type` - The tracing protocol type. Currently, only Jaeger is supported.
  * `vpc_config` - VPC configuration. After this parameter is configured, the function can access the specified VPC resources.
    * `security_group_id` - Security group ID.
    * `vswitch_ids` - Switch List.
    * `vpc_id` - VPC network ID.
  * `id` - The ID of the resource supplied above.
