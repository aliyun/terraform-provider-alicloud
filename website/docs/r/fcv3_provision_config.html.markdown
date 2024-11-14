---
subcategory: "Function Compute Service V3 (FCV3)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fcv3_provision_config"
description: |-
  Provides a Alicloud FCV3 Provision Config resource.
---

# alicloud_fcv3_provision_config

Provides a FCV3 Provision Config resource.

Function Reservation Configuration.

For information about FCV3 Provision Config and how to use it, see [What is Provision Config](https://www.alibabacloud.com/help/en/functioncompute/fc-3-0/developer-reference/api-fc-2023-03-30-putprovisionconfig).

-> **NOTE:** Available since v1.230.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-shanghai"
}

variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_log_project" "default" {
  project_name = "${var.name}-${random_integer.default.result}"
  description  = var.name
}

resource "alicloud_log_store" "default" {
  project_name          = alicloud_log_project.default.name
  logstore_name         = "${var.name}-${random_integer.default.result}"
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "time_sleep" "wait_10_minutes" {
  depends_on = [alicloud_log_store.default]

  create_duration = "10m"
}

resource "alicloud_fcv3_function" "function" {
  memory_size   = "512"
  cpu           = 0.5
  handler       = "index.handler"
  function_name = "${var.name}-${random_integer.default.result}"
  runtime       = "python3.10"
  disk_size     = "512"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
  log_config {
    enable_instance_metrics = true
    enable_request_metrics  = true
    project                 = alicloud_log_project.default.project_name
    logstore                = alicloud_log_store.default.logstore_name
    log_begin_rule          = "None"
  }

  depends_on = [time_sleep.wait_10_minutes]
}

resource "alicloud_fcv3_provision_config" "default" {
  target = "1"
  target_tracking_policies {
    name          = "t1"
    start_time    = "2030-10-10T10:10:10Z"
    end_time      = "2035-10-10T10:10:10Z"
    min_capacity  = "0"
    max_capacity  = "1"
    metric_target = "1"
    metric_type   = "ProvisionedConcurrencyUtilization"
  }
  target_tracking_policies {
    metric_target = "1"
    metric_type   = "ProvisionedConcurrencyUtilization"
    name          = "t2"
    start_time    = "2030-10-10T10:10:10Z"
    end_time      = "2035-10-10T10:10:10Z"
    min_capacity  = "0"
    max_capacity  = "1"
  }
  target_tracking_policies {
    metric_type   = "ProvisionedConcurrencyUtilization"
    time_zone     = "Asia/Shanghai"
    name          = "t3"
    start_time    = "2030-10-10T10:10:10"
    end_time      = "2035-10-10T10:10:10"
    min_capacity  = "0"
    max_capacity  = "1"
    metric_target = "1"
  }

  scheduled_actions {
    target              = "0"
    name                = "s1"
    start_time          = "2030-10-10T10:10:10Z"
    end_time            = "2035-10-10T10:10:10Z"
    schedule_expression = "cron(0 0 4 * * *)"
  }
  scheduled_actions {
    name                = "s2"
    start_time          = "2030-10-10T10:10:10Z"
    end_time            = "2035-10-10T10:10:10Z"
    schedule_expression = "cron(0 0 6 * * *)"
    target              = "1"
  }
  scheduled_actions {
    start_time          = "2030-10-10T10:10:10"
    end_time            = "2035-10-10T10:10:10"
    schedule_expression = "cron(0 0 7 * * *)"
    target              = "0"
    time_zone           = "Asia/Shanghai"
    name                = "s3"
  }

  qualifier           = "LATEST"
  always_allocate_gpu = "true"
  function_name       = alicloud_fcv3_function.function.function_name
  always_allocate_cpu = "true"
}
```

## Argument Reference

The following arguments are supported:
* `always_allocate_cpu` - (Optional) Whether the CPU is always allocated. The default value is true.
* `always_allocate_gpu` - (Optional) Whether to always assign GPU to function instance
* `function_name` - (Required, ForceNew) The name of the function. If this parameter is not specified, the provisioned configurations of all functions are listed.
* `qualifier` - (Optional) The function alias or LATEST.
* `scheduled_actions` - (Optional, List) Timing policy configuration See [`scheduled_actions`](#scheduled_actions) below.
* `target` - (Optional, Int) Number of reserved target resources. The value range is [0,10000].
* `target_tracking_policies` - (Optional, List) Metric tracking scaling policy configuration See [`target_tracking_policies`](#target_tracking_policies) below.

### `scheduled_actions`

The scheduled_actions supports the following:
* `end_time` - (Optional) Policy expiration time
* `name` - (Optional) Policy Name
* `schedule_expression` - (Optional) Timing Configuration
* `start_time` - (Optional) Policy effective time
* `target` - (Optional, Int) Number of reserved target resources
* `time_zone` - (Optional) Time zone.

### `target_tracking_policies`

The target_tracking_policies supports the following:
* `end_time` - (Optional) Policy expiration time
* `max_capacity` - (Optional, Int) Maximum value of expansion
* `metric_target` - (Optional, Float) Tracking value of the indicator
* `metric_type` - (Optional) Provisionedconcurrency utilization: Concurrency utilization of reserved mode instances. CPU utilization: CPU utilization. GPUMemUtilization:GPU utilization
* `min_capacity` - (Optional, Int) Minimum Shrinkage
* `name` - (Optional) Policy Name
* `start_time` - (Optional) Policy Effective Time
* `time_zone` - (Optional) Time zone.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `current` - (Available since v1.234.0) Number of actual resources
* `current_error` - (Available since v1.234.0) Error message when a Reserved Instance creation fails
* `function_arn` - (Available since v1.234.0) Resource Description of the function

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Provision Config.
* `delete` - (Defaults to 5 mins) Used when delete the Provision Config.
* `update` - (Defaults to 5 mins) Used when update the Provision Config.

## Import

FCV3 Provision Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_fcv3_provision_config.example <id>
```