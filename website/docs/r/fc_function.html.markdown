---
subcategory: "Function Compute Service (FC)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fc_function"
sidebar_current: "docs-alicloud-resource-fc"
description: |-
  Provides a Alicloud Function Compute Function resource. Function allows you to trigger execution of code in response to events in Alibaba Cloud. The Function itself includes source code and runtime configuration.
---

# alicloud_fc_function

Provides a Alicloud Function Compute Function resource. Function allows you to trigger execution of code in response to events in Alibaba Cloud. The Function itself includes source code and runtime configuration.
 For information about Service and how to use it, see [What is Function Compute](https://www.alibabacloud.com/help/en/fc/developer-reference/api-createfunction).

-> **NOTE:** The resource requires a provider field 'account_id'. [See account_id](https://www.terraform.io/docs/providers/alicloud/index.html#account_id).

-> **NOTE:** Available since v1.10.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_fc_function&exampleId=845c1d41-b528-e598-35fd-2d8154569f4db35523af&activeTab=example&spm=docs.r.fc_function.0.845c1d41b5&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_log_project" "default" {
  project_name = "example-value-${random_integer.default.result}"
}

resource "alicloud_log_store" "default" {
  project_name  = alicloud_log_project.default.project_name
  logstore_name = "example-value"
}

resource "alicloud_ram_role" "default" {
  name        = "fcservicerole-${random_integer.default.result}"
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
  description = "this is a example"
  force       = true
}

resource "alicloud_ram_role_policy_attachment" "default" {
  role_name   = alicloud_ram_role.default.name
  policy_name = "AliyunLogFullAccess"
  policy_type = "System"
}

resource "alicloud_fc_service" "default" {
  name        = "example-value-${random_integer.default.result}"
  description = "example-value"
  role        = alicloud_ram_role.default.arn
  log_config {
    project                 = alicloud_log_project.default.project_name
    logstore                = alicloud_log_store.default.logstore_name
    enable_instance_metrics = true
    enable_request_metrics  = true
  }
}

resource "alicloud_oss_bucket" "default" {
  bucket = "terraform-example-${random_integer.default.result}"
}
# If you upload the function by OSS Bucket, you need to specify path can't upload by content.
resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.id
  key     = "index.py"
  content = "import logging \ndef handler(event, context): \nlogger = logging.getLogger() \nlogger.info('hello world') \nreturn 'hello world'"
}

resource "alicloud_fc_function" "foo" {
  service     = alicloud_fc_service.default.name
  name        = "terraform-example"
  description = "example"
  oss_bucket  = alicloud_oss_bucket.default.id
  oss_key     = alicloud_oss_bucket_object.default.key
  memory_size = "512"
  runtime     = "python3.10"
  handler     = "hello.handler"
  environment_variables = {
    prefix = "terraform"
  }
}
```

## Module Support

You can use to the existing [fc module](https://registry.terraform.io/modules/terraform-alicloud-modules/fc/alicloud) 
to create a function quickly and set several triggers for it.

## Argument Reference

The following arguments are supported:

* `service` - (Required, ForceNew) The Function Compute service name.
* `name` - (Optional, ForceNew) The Function Compute function name. It is the only in one service and is conflict with "name_prefix".
* `name_prefix` - (Optional, ForceNew) Setting a prefix to get a only function name. It is conflict with "name".
* `description` - (Optional) The Function Compute function description.
* `filename` - (Optional) The path to the function's deployment package within the local filesystem. It is conflict with the `oss_`-prefixed options.
* `oss_bucket` - (Optional) The OSS bucket location containing the function's deployment package. Conflicts with `filename`. This bucket must reside in the same Alibaba Cloud region where you are creating the function.
* `oss_key` - (Optional) The OSS key of an object containing the function's deployment package. Conflicts with `filename`.
* `handler` - (Required) The function [entry point](https://www.alibabacloud.com/help/doc-detail/157704.htm) in your code.
* `memory_size` - (Optional) Amount of memory in MB your function can use at runtime. Defaults to `128`. Limits to [128, 32768].
* `runtime` - (Required) See [Runtimes][https://www.alibabacloud.com/help/zh/function-compute/latest/manage-functions#multiTask3514] for valid values.
* `timeout` - (Optional) The amount of time your function has to run in seconds.
* `environment_variables` - (Optional, available in 1.36.0+) A map that defines environment variables for the function.
* `code_checksum` - (Optional, available in 1.59.0+) The checksum (crc64) of the function code.Used to trigger updates.The value can be generated by data source [alicloud_file_crc64_checksum](https://www.terraform.io/docs/providers/alicloud/d/file_crc64_checksum).
-> **NOTE:** For more information, see [Limits](https://www.alibabacloud.com/help/doc-detail/51907.htm).
* `initializer` - (Optional, available in 1.96.0+) The entry point of the function's [initialization](https://www.alibabacloud.com/help/doc-detail/157704.htm).
* `initialization_timeout` - (Optional, available in 1.96.0+) The maximum length of time, in seconds, that the function's initialization should be run for.
* `instance_type` - (Optional, available in 1.96.0+) The instance type of the function.
* `instance_concurrency` - (Optional, available in 1.96.0+) The maximum number of requests can be executed concurrently within the single function instance.
* `ca_port` - (Optional, available in 1.96.0+) The port that the function listen to, only valid for [custom runtime](https://www.alibabacloud.com/help/doc-detail/132044.htm) and [custom container runtime](https://www.alibabacloud.com/help/doc-detail/179368.htm).
* `custom_container_config` - (Optional, available in 1.96.0+) The configuration for custom container runtime.See [`custom_container_config`](#custom_container_config) below.
* `layers` - (Optional, available in 1.187.0+) The configuration for layers.


### `custom_container_config`

The custom_container_config following arguments:

* `image` - (Required) The container image address.
* `command` - (Optional) The entry point of the container, which specifies the actual command run by the container.
* `args` - (Optional) The args field specifies the arguments passed to the command.

## Attributes Reference

The following arguments are exported:

* `id` - The ID of the function. It formats as `<service>:<name>`.
* `last_modified` - The date this resource was last modified.
* `function_id` - The Function Compute service function ID.
* `function_arn` - The Function Compute service function arn. It formats as `acs:fc:<region>:<uid>:services/<serviceName>.LATEST/functions/<functionName>`.

## Import

Function Compute function can be imported using the id, e.g.

```shell
$ terraform import alicloud_fc_function.foo my-fc-service:hello-world
```
