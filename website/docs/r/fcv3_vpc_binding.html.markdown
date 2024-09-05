---
subcategory: "Function Compute Service V3 (FCV3)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fcv3_vpc_binding"
description: |-
  Provides a Alicloud FCV3 Vpc Binding resource.
---

# alicloud_fcv3_vpc_binding

Provides a FCV3 Vpc Binding resource.



For information about FCV3 Vpc Binding and how to use it, see [What is Vpc Binding](https://www.alibabacloud.com/help/en/functioncompute/fc-3-0/developer-reference/api-fc-2023-03-30-createvpcbinding).

-> **NOTE:** Available since v1.230.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
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

resource "alicloud_fcv3_vpc_binding" "default" {
  function_name = alicloud_fcv3_function.function.function_name
  vpc_id        = alicloud_vpc.vpc.id
}
```

## Argument Reference

The following arguments are supported:
* `function_name` - (Required, ForceNew) Function Name
* `vpc_id` - (Optional, ForceNew, Computed) VPC instance ID

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<function_name>:<vpc_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vpc Binding.
* `delete` - (Defaults to 5 mins) Used when delete the Vpc Binding.

## Import

FCV3 Vpc Binding can be imported using the id, e.g.

```shell
$ terraform import alicloud_fcv3_vpc_binding.example <function_name>:<vpc_id>
```