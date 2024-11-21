---
subcategory: "Function Compute Service V3 (FCV3)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fcv3_async_invoke_config"
description: |-
  Provides a Alicloud FCV3 Async Invoke Config resource.
---

# alicloud_fcv3_async_invoke_config

Provides a FCV3 Async Invoke Config resource.

Function Asynchronous Configuration.

For information about FCV3 Async Invoke Config and how to use it, see [What is Async Invoke Config](https://www.alibabacloud.com/help/en/functioncompute/developer-reference/api-fc-2023-03-30-getasyncinvokeconfig).

-> **NOTE:** Available since v1.228.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_fcv3_async_invoke_config&exampleId=be7c16e0-04af-42e5-9682-c2abbd16dfcea05a8929&activeTab=example&spm=docs.r.fcv3_async_invoke_config.0.be7c16e004&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "eu-central-1"
}

data "alicloud_account" "current" {
}

resource "alicloud_fcv3_function" "function" {
  memory_size = "512"
  cpu         = 0.5
  handler     = "index.Handler"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
  function_name = var.name
  runtime       = "python3.9"
  disk_size     = "512"
  log_config {
    log_begin_rule = "None"
  }
}

resource "alicloud_fcv3_function" "function1" {
  memory_size = "512"
  cpu         = 0.5
  handler     = "index.Handler"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
  function_name = format("%s_%s", var.name, "update1")
  runtime       = "python3.9"
  disk_size     = "512"
  log_config {
    log_begin_rule = "None"
  }
}

resource "alicloud_fcv3_function" "function2" {
  memory_size = "512"
  cpu         = 0.5
  handler     = "index.Handler"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
  function_name = format("%s_%s", var.name, "update2")
  runtime       = "python3.9"
  disk_size     = "512"
  log_config {
    log_begin_rule = "None"
  }
}

resource "alicloud_fcv3_async_invoke_config" "default" {
  max_async_retry_attempts       = "1"
  max_async_event_age_in_seconds = "1"
  async_task                     = "true"
  function_name                  = alicloud_fcv3_function.function.function_name
  destination_config {
    on_failure {
      destination = "acs:fc:eu-central-1:${data.alicloud_account.current.id}:functions/${alicloud_fcv3_function.function1.function_name}"
    }

    on_success {
      destination = "acs:fc:eu-central-1:${data.alicloud_account.current.id}:functions/${alicloud_fcv3_function.function1.function_name}"
    }

  }

  qualifier = "LATEST"
}
```

## Argument Reference

The following arguments are supported:
* `async_task` - (Optional) Whether to enable an asynchronous task
* `destination_config` - (Optional, List) Target Configuration See [`destination_config`](#destination_config) below.
* `function_name` - (Required, ForceNew) Function Name
* `max_async_event_age_in_seconds` - (Optional, Int) Event maximum survival time
* `max_async_retry_attempts` - (Optional, Int) Number of Asynchronous call retries
* `qualifier` - (Optional) Function version or alias

### `destination_config`

The destination_config supports the following:
* `on_failure` - (Optional, List) Failed callback target structure See [`on_failure`](#destination_config-on_failure) below.
* `on_success` - (Optional, List) Successful callback target structure See [`on_success`](#destination_config-on_success) below.

### `destination_config-on_failure`

The destination_config-on_failure supports the following:
* `destination` - (Optional) Asynchronous call target Resource Descriptor

### `destination_config-on_success`

The destination_config-on_success supports the following:
* `destination` - (Optional) Asynchronous call target Resource Descriptor

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource
* `function_arn` - (Available since v1.234.0) Function resource identification
* `last_modified_time` - (Available since v1.234.0) Last modification time

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Async Invoke Config.
* `delete` - (Defaults to 5 mins) Used when delete the Async Invoke Config.
* `update` - (Defaults to 5 mins) Used when update the Async Invoke Config.

## Import

FCV3 Async Invoke Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_fcv3_async_invoke_config.example <id>
```