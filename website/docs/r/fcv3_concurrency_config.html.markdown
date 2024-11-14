---
subcategory: "Function Compute Service V3 (FCV3)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fcv3_concurrency_config"
description: |-
  Provides a Alicloud FCV3 Concurrency Config resource.
---

# alicloud_fcv3_concurrency_config

Provides a FCV3 Concurrency Config resource.

Function concurrency configuration.

For information about FCV3 Concurrency Config and how to use it, see [What is Concurrency Config](https://www.alibabacloud.com/help/en/functioncompute/developer-reference/api-fc-2023-03-30-putconcurrencyconfig).

-> **NOTE:** Available since v1.228.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
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

resource "alicloud_fcv3_concurrency_config" "default" {
  function_name        = alicloud_fcv3_function.function.function_name
  reserved_concurrency = "100"
}
```

## Argument Reference

The following arguments are supported:
* `function_name` - (Required, ForceNew) Function Name
* `reserved_concurrency` - (Optional, Int) Reserved Concurrency. Functions reserve a part of account concurrency. Other functions cannot use this part of concurrency. Reserved concurrency includes the total concurrency of Reserved Instances and As-You-go instances.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `function_arn` - (Available since v1.234.0) Resource identity of the function

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Concurrency Config.
* `delete` - (Defaults to 5 mins) Used when delete the Concurrency Config.
* `update` - (Defaults to 5 mins) Used when update the Concurrency Config.

## Import

FCV3 Concurrency Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_fcv3_concurrency_config.example <id>
```