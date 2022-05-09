---
subcategory: "Function Compute Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_fc_function_async_invoke_config"
sidebar_current: "docs-alicloud-resource-fc-function-async-invoke-config"
description: |-
  Provides an Alicloud Function Compute Function Async Invoke Config resource. 
---

# alicloud\_fc\_function\_async\_invoke\_config

Manages an asynchronous invocation configuration for a FC Function or Alias.  
 For the detailed information, please refer to the [developer guide](https://www.alibabacloud.com/help/doc-detail/181866.htm).

-> **NOTE:** Available in 1.100.0+

## Example Usage

### Destination Configuration

-> **NOTE** Ensure the FC Function RAM Role has necessary permissions for the destination, such as `mns:SendMessage`, `mns:PublishMessage` or `fc:InvokeFunction`, otherwise the API will return a generic error.

```terraform
resource "alicloud_fc_function_async_invoke_config" "example" {
  service_name  = alicloud_fc_service.example.name
  function_name = alicloud_fc_function.example.name

  destination_config {
    on_failure {
      destination = the_example_mns_queue_arn
    }

    on_success {
      destination = the_example_mns_topic_arn
    }
  }
}
```

### Error Handling Configuration

```terraform
resource "alicloud_fc_function_async_invoke_config" "example" {
  service_name                 = alicloud_fc_service.example.name
  function_name                = alicloud_fc_function.example.name
  maximum_event_age_in_seconds = 60
  maximum_retry_attempts       = 0
}
```

### Async Job Configuration

```terraform
resource "alicloud_fc_function_async_invoke_config" "example" {
  service_name        = alicloud_fc_service.example.name
  function_name       = alicloud_fc_function.example.name
  stateful_invocation = true
}
```

### Configuration for Function Latest Unpublished Version

```terraform
resource "alicloud_fc_function_async_invoke_config" "example" {
  service_name  = alicloud_fc_service.example.name
  function_name = alicloud_fc_function.example.name
  qualifier     = "LATEST"

  # ... other configuration ...
}
```

## Argument Reference

* `service_name` - (Required, ForceNew) Name of the Function Compute Function, omitting any version or alias qualifier.
* `function_name` - (Required, ForceNew) Name of the Function Compute Function.
* `destination_config` - (Optional) Configuration block with destination configuration. See below for details.
* `maximum_event_age_in_seconds` - (Optional) Maximum age of a request that Function Compute sends to a function for processing in seconds. Valid values between 1 and 2592000 (between 60 and 21600 before v1.167.0).
* `maximum_retry_attempts` - (Optional) Maximum number of times to retry when the function returns an error. Valid values between 0 and 8 (between 0 and 2 before v1.167.0). Defaults to 2.
* `stateful_invocation` - (Optional) Function Compute async job configuration. valid values true or false, default `false`
* `qualifier` - (Optional, ForceNew) Function Compute Function published version, `LATEST`, or Function Compute Alias name. The default value is `LATEST`.

### destination_config Configuration Block

-> **NOTE:** At least one of `on_failure` or `on_success` must be configured when using this configuration block, otherwise remove it completely to prevent perpetual differences in Terraform runs.

The following arguments are optional:

* `on_failure` - (Optional) Configuration block with destination configuration for failed asynchronous invocations. See below for details.
* `on_success` - (Optional) Configuration block with destination configuration for successful asynchronous invocations. See below for details.

#### destination_config on_failure Configuration Block

The following arguments are required:

* `destination` - (Required) Alicloud Resource Name (ARN) of the destination resource. See the [Developer Guide](https://www.alibabacloud.com/help/doc-detail/181866.htm) for acceptable resource types and associated RAM permissions.

#### destination_config on_success Configuration Block

The following arguments are required:

* `destination` - (Required) Alicloud Resource Name (ARN) of the destination resource. See the [Developer Guide](https://www.alibabacloud.com/help/doc-detail/181866.htm) for acceptable resource types and associated RAM permissions.

## Attributes Reference

In addition to all arguments above, the following arguments are exported:

* `id` - Fully qualified Function Compute Function name (`service_name:function_name:qualifier`) or Alicloud Resource Name (ARN).
* `created_time` - The date this resource was created.
* `last_modified_time` - The date this resource was last modified.

## Import

Function Compute Function Async Invoke Configs can be imported using the id, e.g.

```
$ terraform import alicloud_fc_function_async_invoke_config.example my_function
```
