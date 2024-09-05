---
subcategory: "Function Compute Service V3 (FCV3)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fcv3_layer_version"
description: |-
  Provides a Alicloud FCV3 Layer Version resource.
---

# alicloud_fcv3_layer_version

Provides a FCV3 Layer Version resource.

Layer provides you with the ability to publish and deploy common dependency libraries, runtime environments, and function extensions.

For information about FCV3 Layer Version and how to use it, see [What is Layer Version](https://www.alibabacloud.com/help/en/functioncompute/api-fc-2023-03-30-createlayerversion).

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


resource "alicloud_fcv3_layer_version" "default" {
  description = var.name
  layer_name  = "FC3LayerResouceTest_ZIP_2024SepWed"
  license     = "Apache2.0"
  acl         = "0"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
}
```

## Argument Reference

The following arguments are supported:
* `acl` - (Optional) The access permission of the layer, 1: public, 0: private, default is private
* `code` - (Optional, ForceNew, List) Layer code configuration See [`code`](#code) below.
* `compatible_runtime` - (Optional, Computed, ForceNew, List) List of runtime environments supported by the layer
* `description` - (Optional, ForceNew) Description of the version
* `layer_name` - (Required, ForceNew) Name of the layer
* `license` - (Optional, ForceNew) Layer License Agreement

### `code`

The code supports the following:
* `checksum` - (Optional, ForceNew) The CRC-64 value of the code package. If checksum is provided, Function Compute checks whether the checksum of the code package is consistent with the provided checksum.
* `oss_bucket_name` - (Optional, ForceNew) Name of the OSS Bucket where the user stores the Layer Code ZIP package.
* `oss_object_name` - (Optional, ForceNew) Name of the OSS Object where the user stores the Layer Code ZIP package.
* `zip_file` - (Optional, ForceNew) Base 64 encoding of Layer Code ZIP package.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<layer_name>:<version>`.
* `create_time` - The creation time of the resource
* `version` - The version of the layer

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Layer Version.
* `delete` - (Defaults to 5 mins) Used when delete the Layer Version.
* `update` - (Defaults to 5 mins) Used when update the Layer Version.

## Import

FCV3 Layer Version can be imported using the id, e.g.

```shell
$ terraform import alicloud_fcv3_layer_version.example <layer_name>:<version>
```