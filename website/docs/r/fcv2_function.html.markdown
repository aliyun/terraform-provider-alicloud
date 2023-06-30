---
subcategory: "FCV2"
layout: "alicloud"
page_title: "Alicloud: alicloud_fcv2_function"
description: |-
  Provides a Alicloud FCV2 Function resource.
---

# alicloud_fcv2_function

Provides a FCV2 Function resource. Function is the unit of system scheduling and operation. Functions must be subordinate to services. All functions under the same service share some identical settings, such as service authorization and log configuration.

For information about FCV2 Function and how to use it, see [What is Function](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.207.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}


resource "alicloud_fcv2_function" "default" {
  function_name          = var.name
  memory_size            = 1024
  runtime                = "custom.debian10"
  description            = "terraform测试case"
  service_name           = "terraform-e2e-service-base"
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
    zip_file = "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"
  }
  environment_variables {
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
  - **g1**: Same as **fc.gpu.tesla.1 * *.
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
* `instance_id` - (Optional) ACR enterprise image Warehouse ID, which must be passed in when using ACR enterprise image.
* `web_server_mode` - (Optional) Whether the image is run in Web Server mode. The value of true needs to implement the Web Server in the container image to listen to the port and process the request. The value of false needs to actively exit the process after the container runs, and the ExitCode needs to be 0. Default true.

### `dns_options`

The dns_options supports the following:
* `name` - (Optional) DNS option name.
* `value` - (Optional) DNS option value.

### `custom_dns`

The custom_dns supports the following:
* `dns_options` - (Optional) DNS resolver configuration parameter list. See [`dns_options`](#dns_options) below.
* `name_servers` - (Optional) List of IP addresses of DNS servers.
* `searches` - (Optional) List of DNS search domains.

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

### `pre_freeze`

The pre_freeze supports the following:
* `handler` - (Optional) Entry for function execution.
* `timeout` - (Optional) The timeout of the run, in seconds.

### `pre_stop`

The pre_stop supports the following:
* `handler` - (Optional) Entry for function execution.
* `timeout` - (Optional) Timeout of run.

### `instance_lifecycle_config`

The instance_lifecycle_config supports the following:
* `pre_freeze` - (Optional) PreFreeze function configuration. See [`pre_freeze`](#pre_freeze) below.
* `pre_stop` - (Optional) PreStop function configuration. See [`pre_stop`](#pre_stop) below.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<service_name>:<function_name>`.
* `create_time` - create time of function.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Function.
* `delete` - (Defaults to 5 mins) Used when delete the Function.
* `update` - (Defaults to 5 mins) Used when update the Function.

## Import

FCV2 Function can be imported using the id, e.g.

```shell
$ terraform import alicloud_fcv2_function.example <service_name>:<function_name>
```