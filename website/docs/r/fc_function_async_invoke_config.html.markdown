---
subcategory: "Function Compute Service (FC)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fc_function_async_invoke_config"
sidebar_current: "docs-alicloud-resource-fc-function-async-invoke-config"
description: |-
  Provides an Alicloud Function Compute Function Async Invoke Config resource. 
---

# alicloud_fc_function_async_invoke_config

Manages an asynchronous invocation configuration for a FC Function or Alias.  
 For the detailed information, please refer to the [developer guide](https://www.alibabacloud.com/help/en/fc/developer-reference/api-fc-open-2021-04-06-putfunctionasyncinvokeconfig).

-> **NOTE:** Available since v1.100.0.

## Example Usage

### Destination Configuration

-> **NOTE** Ensure the FC Function RAM Role has necessary permissions for the destination, such as `mns:SendMessage`, `mns:PublishMessage` or `fc:InvokeFunction`, otherwise the API will return a generic error.

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_fc_function_async_invoke_config&exampleId=c9b2eee5-8400-74ab-91c2-a98ed4b9ece56b26feec&activeTab=example&spm=docs.r.fc_function_async_invoke_config.0.c9b2eee584&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_account" "default" {}
data "alicloud_regions" "default" {
  current = true
}

resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_ram_role" "default" {
  name        = "examplerole${random_integer.default.result}"
  document    = <<DEFINITION
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
	DEFINITION
  description = "this is a example"
  force       = true
}
resource "alicloud_ram_policy" "default" {
  policy_name     = "examplepolicy${random_integer.default.result}"
  policy_document = <<DEFINITION
	{
		"Version": "1",
		"Statement": [
		  {
			"Action": "mns:*",
			"Resource": "*",
			"Effect": "Allow"
		  }
		]
	  }
	DEFINITION
}
resource "alicloud_ram_role_policy_attachment" "default" {
  role_name   = alicloud_ram_role.default.name
  policy_name = alicloud_ram_policy.default.policy_name
  policy_type = "Custom"
}

resource "alicloud_fc_service" "default" {
  name            = "example-value-${random_integer.default.result}"
  description     = "example-value"
  role            = alicloud_ram_role.default.arn
  internet_access = false
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

resource "alicloud_fc_function" "default" {
  service     = alicloud_fc_service.default.name
  name        = "terraform-example-${random_integer.default.result}"
  description = "example"
  oss_bucket  = alicloud_oss_bucket.default.id
  oss_key     = alicloud_oss_bucket_object.default.key
  memory_size = "512"
  runtime     = "python3.10"
  handler     = "hello.handler"
}

resource "alicloud_mns_queue" "default" {
  name = "terraform-example-${random_integer.default.result}"
}
resource "alicloud_mns_topic" "default" {
  name = "terraform-example-${random_integer.default.result}"
}

resource "alicloud_fc_function_async_invoke_config" "default" {
  service_name  = alicloud_fc_service.default.name
  function_name = alicloud_fc_function.default.name

  destination_config {
    on_failure {
      destination = "acs:mns:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:/queues/${alicloud_mns_queue.default.name}/messages"
    }

    on_success {
      destination = "acs:mns:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:/topics/${alicloud_mns_topic.default.name}/messages"
    }
  }


  # Error Handling Configuration
  maximum_event_age_in_seconds = 60
  maximum_retry_attempts       = 0

  # Async Job Configuration
  stateful_invocation = true

  # Configuration for Function Latest Unpublished Version
  qualifier = "LATEST"

}
```

## Argument Reference

* `service_name` - (Required, ForceNew) Name of the Function Compute Function, omitting any version or alias qualifier.
* `function_name` - (Required, ForceNew) Name of the Function Compute Function.
* `destination_config` - (Optional) Configuration block with destination configuration. See [`destination_config`](#destination_config) below.
* `maximum_event_age_in_seconds` - (Optional) Maximum age of a request that Function Compute sends to a function for processing in seconds. Valid values between 1 and 2592000 (between 60 and 21600 before v1.167.0).
* `maximum_retry_attempts` - (Optional) Maximum number of times to retry when the function returns an error. Valid values between 0 and 8 (between 0 and 2 before v1.167.0). Defaults to 2.
* `stateful_invocation` - (Optional) Function Compute async job configuration(also known as Task Mode). valid values true or false, default `false`
* `qualifier` - (Optional, ForceNew) Function Compute Function published version, `LATEST`, or Function Compute Alias name. The default value is `LATEST`.

### `destination_config`

-> **NOTE:** At least one of `on_failure` or `on_success` must be configured when using this configuration block, otherwise remove it completely to prevent perpetual differences in Terraform runs.

The following arguments are optional:

* `on_failure` - (Optional) Configuration block with destination configuration for failed asynchronous invocations. See [`on_failure`](#destination_config-on_failure) below.
* `on_success` - (Optional) Configuration block with destination configuration for successful asynchronous invocations. See [`on_success`](#destination_config-on_success) below.

### `destination_config-on_failure`

The following arguments are required:

* `destination` - (Required) Alicloud Resource Name (ARN) of the destination resource. See the [Developer Guide](https://www.alibabacloud.com/help/doc-detail/181866.htm) for acceptable resource types and associated RAM permissions.

### `destination_config-on_success`

The following arguments are required:

* `destination` - (Required) Alicloud Resource Name (ARN) of the destination resource. See the [Developer Guide](https://www.alibabacloud.com/help/doc-detail/181866.htm) for acceptable resource types and associated RAM permissions.

## Attributes Reference

In addition to all arguments above, the following arguments are exported:

* `id` - Fully qualified Function Compute Function name (`service_name:function_name:qualifier`) or Alicloud Resource Name (ARN).
* `created_time` - The date this resource was created.
* `last_modified_time` - The date this resource was last modified.

## Import

Function Compute Function Async Invoke Configs can be imported using the id, e.g.

```shell
$ terraform import alicloud_fc_function_async_invoke_config.example my_function
```
