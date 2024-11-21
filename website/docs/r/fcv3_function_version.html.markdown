---
subcategory: "Function Compute Service V3 (FCV3)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fcv3_function_version"
description: |-
  Provides a Alicloud FCV3 Function Version resource.
---

# alicloud_fcv3_function_version

Provides a FCV3 Function Version resource.

Version of the function.

For information about FCV3 Function Version and how to use it, see [What is Function Version](https://www.alibabacloud.com/help/en/functioncompute/api-fc-2023-03-30-listfunctionversions).

-> **NOTE:** Available since v1.228.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_fcv3_function_version&exampleId=0f46471e-ba97-ab68-f33b-50c7e9474f2fb6259661&activeTab=example&spm=docs.r.fcv3_function_version.0.0f46471eba&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_uuid" "default" {
}

resource "alicloud_fcv3_function" "function" {
  memory_size = "512"
  cpu         = 0.5
  handler     = "index.Handler"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
  function_name = "${var.name}-${random_uuid.default.result}"
  runtime       = "python3.9"
  disk_size     = "512"
  log_config {
    log_begin_rule = "None"
  }
}

resource "alicloud_fcv3_function_version" "default" {
  function_name = alicloud_fcv3_function.function.function_name
  description   = "version1"
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional, ForceNew) Description of the function version
* `function_name` - (Required, ForceNew) Function Name

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<function_name>:<version_id>`.
* `create_time` - The creation time of the resource
* `last_modified_time` - (Available since v1.234.0) Update time

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Function Version.
* `delete` - (Defaults to 5 mins) Used when delete the Function Version.

## Import

FCV3 Function Version can be imported using the id, e.g.

```shell
$ terraform import alicloud_fcv3_function_version.example <function_name>:<version_id>
```