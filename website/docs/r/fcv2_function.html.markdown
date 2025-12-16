---
subcategory: "Function Compute Service (FC)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fcv2_function"
description: |-
  Provides a Alicloud FCV2 Function resource.
---

# alicloud_fcv2_function

Provides a FCV2 Function resource. Function is the unit of system scheduling and operation. Functions must be subordinate to services. All functions under the same service share some identical settings, such as service authorization and log configuration.

For information about FCV2 Function and how to use it, see [What is Function](https://www.alibabacloud.com/help/en/resource-orchestration-service/latest/aliyun-fc-function).

-> **NOTE:** Available since v1.208.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_fcv2_function&exampleId=0a7b499c-14c5-4650-98aa-ecd41560f0b72abab691&activeTab=example&spm=docs.r.fcv2_function.0.0a7b499c14&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_log_project" "default" {
  name        = var.name
  description = var.name
}

resource "alicloud_log_store" "default" {
  project          = alicloud_log_project.default.name
  name             = var.name
  retention_period = "3000"
  shard_count      = 1
}

# add index for logstore, which is used to query logs
locals {
  sls_default_token = ", '\";=()[]{}?@&<>/:\n\t\r"
}

resource "alicloud_log_store_index" "example" {
  project  = alicloud_log_project.default.name
  logstore = alicloud_log_store.default.name
  full_text {
    case_sensitive = false
    token          = local.sls_default_token
  }
  field_search {
    name             = "aggPeriodSeconds"
    enable_analytics = true
    type             = "long"
    token            = local.sls_default_token
  }
  field_search {
    name             = "concurrentRequests"
    enable_analytics = true
    type             = "long"
    token            = local.sls_default_token
  }
  field_search {
    name             = "cpuPercent"
    enable_analytics = true
    type             = "double"
    token            = local.sls_default_token
  }
  field_search {
    name             = "cpuQuotaPercent"
    enable_analytics = true
    type             = "double"
    token            = local.sls_default_token
  }
  field_search {
    name             = "functionName"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
    case_sensitive   = true
  }
  field_search {
    name             = "hostname"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
  }
  field_search {
    name             = "instanceID"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
  }
  field_search {
    name             = "ipAddress"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
  }
  field_search {
    name             = "memoryLimitMB"
    enable_analytics = true
    type             = "double"
    token            = local.sls_default_token
  }
  field_search {
    name             = "memoryUsageMB"
    enable_analytics = true
    type             = "double"
    token            = local.sls_default_token
  }
  field_search {
    name             = "memoryUsagePercent"
    enable_analytics = true
    type             = "double"
    token            = local.sls_default_token
  }
  field_search {
    name             = "operation"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
  }
  field_search {
    name             = "qualifier"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
    case_sensitive   = true
  }
  field_search {
    name             = "rxBytes"
    enable_analytics = true
    type             = "long"
    token            = local.sls_default_token
  }
  field_search {
    name             = "rxTotalBytes"
    enable_analytics = true
    type             = "long"
    token            = local.sls_default_token
  }
  field_search {
    name             = "serviceName"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
    case_sensitive   = true
  }
  field_search {
    name             = "txBytes"
    enable_analytics = true
    type             = "long"
    token            = local.sls_default_token
  }
  field_search {
    name             = "txTotalBytes"
    enable_analytics = true
    type             = "long"
    token            = local.sls_default_token
  }
  field_search {
    name             = "versionId"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
  }
  field_search {
    name             = "events"
    enable_analytics = true
    type             = "json"
    token            = local.sls_default_token
  }
  field_search {
    name             = "isColdStart"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
  }
  field_search {
    name             = "hasFunctionError"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
  }
  field_search {
    name             = "errorType"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
  }
  field_search {
    name             = "triggerType"
    enable_analytics = true
    type             = "text"
    token            = local.sls_default_token
  }
  field_search {
    name             = "durationMs"
    enable_analytics = true
    type             = "double"
    token            = local.sls_default_token
  }
  field_search {
    name             = "statusCode"
    enable_analytics = true
    type             = "long"
    token            = local.sls_default_token
  }
}

resource "alicloud_ram_role" "default" {
  name        = var.name
  document    = <<EOF
  {
      "Statement": [
        {
          "Action": "sts:AssumeRole",
          "Effect": "Allow",
          "Principal": {
            "Service": [
              "fc.aliyuncs.com"
            ]
          }
        }
      ],
      "Version": "1"
  }
  EOF
  description = var.name
  force       = true
}

resource "alicloud_ram_role_policy_attachment" "default" {
  role_name   = alicloud_ram_role.default.name
  policy_name = "AliyunLogFullAccess"
  policy_type = "System"
}

resource "alicloud_fc_service" "default" {
  name        = var.name
  description = var.name
  log_config {
    project  = alicloud_log_project.default.name
    logstore = alicloud_log_store.default.name
  }
  role = alicloud_ram_role.default.arn
}

resource "alicloud_fcv2_function" "default" {
  function_name          = var.name
  memory_size            = 1024
  runtime                = "custom.debian10"
  description            = var.name
  service_name           = alicloud_fc_service.default.name
  initializer            = "index.initializer"
  initialization_timeout = 10
  timeout                = 60
  handler                = "index.handler"
  instance_type          = "e1"
  instance_lifecycle_config {
    pre_freeze {
      handler = "index.prefreeze"
      timeout = 30
    }
    pre_stop {
      handler = "index.prestop"
      timeout = 30
    }
  }
  code {
    oss_bucket_name = "code-sample-cn-hangzhou"
    oss_object_name = "quick-start-sample-codes/quick-start-sample-codes-nodejs/RocketMQ-producer-nodejs14-event/code.zip"
  }
  custom_dns {
    name_servers = ["223.5.5.5"]
    searches     = ["mydomain.com"]
    dns_options {
      name  = var.name
      value = "1"
    }
  }
  disk_size            = 512
  instance_concurrency = 10
  layers               = ["d3fc5de8d120687be2bfab761518d5de#Nodejs-Aliyun-SDK#2", "d3fc5de8d120687be2bfab761518d5de#Python39#2"]
  cpu                  = 1
  custom_health_check_config {
    http_get_url          = "/healthcheck"
    initial_delay_seconds = 3
    period_seconds        = 3
    timeout_seconds       = 3
    failure_threshold     = 1
    success_threshold     = 1
  }
  ca_port = 9000
  custom_runtime_config {
    command = ["npm"]
    args    = ["run", "start"]
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_fcv2_function&spm=docs.r.fcv2_function.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `ca_port` - (Optional, Computed) The listening port of the HTTP Server when the Custom Runtime or Custom Container is running.
* `code` - (Optional) Function Code ZIP package. code and customContainerConfig choose one. See [`code`](#code) below.
* `code_checksum` - (Optional, Computed) crc64 of function code.
* `cpu` - (Optional) The CPU specification of the function. The unit is vCPU, which is a multiple of the 0.05 vCPU.
* `custom_container_config` - (Optional) Custom-container runtime related function configuration. See [`custom_container_config`](#custom_container_config) below.
* `custom_dns` - (Optional) Function custom DNS configuration. See [`custom_dns`](#custom_dns) below.
* `custom_health_check_config` - (Optional) Custom runtime/container Custom health check configuration. See [`custom_health_check_config`](#custom_health_check_config) below.
* `custom_runtime_config` - (Optional) Detailed configuration of Custom Runtime function. See [`custom_runtime_config`](#custom_runtime_config) below.
* `description` - (Optional) description of function.
* `disk_size` - (Optional) The disk specification of the function. The unit is MB. The optional value is 512 MB or 10240MB.
* `environment_variables` - (Optional, Map) The environment variable set for the function can get the value of the environment variable in the function. For more information, see [Environment Variables](~~ 69777 ~~).
* `function_name` - (Required, ForceNew) function name.
* `gpu_memory_size` - (Optional) The GPU memory specification of the function, in MB, is a multiple of 1024MB.
* `handler` - (Required) entry point of function.
* `initialization_timeout` - (Optional, Computed) max running time of initializer.
* `initializer` - (Optional) initializer entry point of function.
* `instance_concurrency` - (Optional, Computed) The maximum concurrency allowed for a single function instance.
* `instance_lifecycle_config` - (Optional) Instance lifecycle configuration. See [`instance_lifecycle_config`](#instance_lifecycle_config) below.
* `instance_type` - (Optional, Computed) The instance type of the function. Valid values:
  - **e1**: Elastic instance.
  - **c1**: performance instance.
  - **fc.gpu.tesla.1**: the T4 card type of the Tesla series of GPU instances.
  - **fc.gpu.ampere.1**: The Ampere series A10 card type of the GPU instance.
  - **g1**: Same as **fc.gpu.tesla.1**.
* `layers` - (Optional) List of layers.
-> **NOTE:**  Multiple layers will be merged in the order of array subscripts from large to small, and the contents of layers with small subscripts will overwrite the files with the same name of layers with large subscripts.
* `memory_size` - (Optional, Computed) memory size needed by function.
* `runtime` - (Required) runtime of function code.
* `service_name` - (Required, ForceNew) The name of the function Service.
* `timeout` - (Optional, Computed) max running time of function.

### `code`

The code supports the following:
* `oss_bucket_name` - (Optional) The OSS bucket name of the function code package.
* `oss_object_name` - (Optional) The OSS object name of the function code package.
* `zip_file` - (Optional) Upload the base64 encoding of the code zip package directly in the request body.

### `custom_container_config`

The custom_container_config supports the following:
* `acceleration_type` - (Optional) Image acceleration type. The value Default is to enable acceleration and None is to disable acceleration.
* `args` - (Optional) Container startup parameters.
* `command` - (Optional) Container start command, equivalent to Docker ENTRYPOINT.
* `image` - (Optional) Container Image address. Example value: registry-vpc.cn-hangzhou.aliyuncs.com/fc-demo/helloworld:v1beta1.
* `web_server_mode` - (Optional) Whether the image is run in Web Server mode. The value of true needs to implement the Web Server in the container image to listen to the port and process the request. The value of false needs to actively exit the process after the container runs, and the ExitCode needs to be 0. Default true.

### `custom_dns`

The custom_dns supports the following:
* `dns_options` - (Optional) DNS resolver configuration parameter list. See [`dns_options`](#custom_dns-dns_options) below.
* `name_servers` - (Optional) List of IP addresses of DNS servers.
* `searches` - (Optional) List of DNS search domains.

### `custom_dns-dns_options`

The dns_options supports the following:
* `name` - (Optional) DNS option name.
* `value` - (Optional) DNS option value.

### `custom_health_check_config`

The custom_health_check_config supports the following:
* `failure_threshold` - (Optional) The threshold for the number of health check failures. The system considers the check failed after the health check fails.
* `http_get_url` - (Optional) Container custom health check URL address.
* `initial_delay_seconds` - (Optional) Delay from container startup to initiation of health check.
* `period_seconds` - (Optional) Health check cycle.
* `success_threshold` - (Optional) The threshold for the number of successful health checks. After the health check is reached, the system considers the check successful.
* `timeout_seconds` - (Optional) Health check timeout.

### `custom_runtime_config`

The custom_runtime_config supports the following:
* `args` - (Optional) Parameters received by the start entry command.
* `command` - (Optional) List of Custom entry commands started by Custom Runtime. When there are multiple commands in the list, they are spliced in sequence.

### `instance_lifecycle_config`

The instance_lifecycle_config supports the following:
* `pre_freeze` - (Optional) PreFreeze function configuration. See [`pre_freeze`](#instance_lifecycle_config-pre_freeze) below.
* `pre_stop` - (Optional) PreStop function configuration. See [`pre_stop`](#instance_lifecycle_config-pre_stop) below.

### `instance_lifecycle_config-pre_freeze`

The pre_freeze supports the following:
* `handler` - (Optional) Entry for function execution.
* `timeout` - (Optional) The timeout of the run, in seconds.

### `instance_lifecycle_config-pre_stop`

The pre_stop supports the following:
* `handler` - (Optional) Entry for function execution.
* `timeout` - (Optional) Timeout of run.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<service_name>:<function_name>`.
* `create_time` - create time of function.
* `function_arn` - The Function Compute service function arn. It formats as `acs:fc:<region>:<uid>:services/<serviceName>.LATEST/functions/<functionName>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Function.
* `delete` - (Defaults to 5 mins) Used when delete the Function.
* `update` - (Defaults to 5 mins) Used when update the Function.

## Import

FCV2 Function can be imported using the id, e.g.

```shell
$ terraform import alicloud_fcv2_function.example <service_name>:<function_name>
```