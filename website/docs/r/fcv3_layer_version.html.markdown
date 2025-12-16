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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_fcv3_layer_version&exampleId=c3a0de1e-1e03-c7b9-d520-cf23650d86d3863042cf&activeTab=example&spm=docs.r.fcv3_layer_version.0.c3a0de1e1e&intl_lang=EN_US" target="_blank">
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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_fcv3_layer_version&spm=docs.r.fcv3_layer_version.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `acl` - (Optional) The access permission of the layer, 1: public, 0: private, default is private
* `code` - (Optional, ForceNew, List) Layer code configuration See [`code`](#code) below.
* `compatible_runtime` - (Optional, ForceNew, List) List of runtime environments supported by the layer
* `description` - (Optional, ForceNew) Description of the version
* `layer_name` - (Required, ForceNew) Name of the layer
* `license` - (Optional, ForceNew) Layer License Agreement
* `public` - (Optional, ForceNew, Available since v1.234.0) Whether to expose the layer. Enumeration values: true, false. (Deprecated, please use acl instead)

### `code`

The code supports the following:
* `checksum` - (Optional, ForceNew) The CRC-64 value of the code package. If checksum is provided, Function Compute checks whether the checksum of the code package is consistent with the provided checksum.
* `oss_bucket_name` - (Optional, ForceNew) Name of the OSS Bucket where the user stores the Layer Code ZIP package
* `oss_object_name` - (Optional, ForceNew) Name of the OSS Object where the user stores the Layer Code ZIP package
* `zip_file` - (Optional, ForceNew) Base 64 encoding of Layer Code ZIP package

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<layer_name>:<version>`.
* `code_size` - (Available since v1.234.0) The code package size of the layer, in bytes.
* `create_time` - The creation time of the resource
* `layer_version_arn` - (Available since v1.234.0) Layer version ARN. The format is acs:fc:{region }:{ accountID}:layers/{layerName}/versions/{layerVersion}.
* `version` - The version of the layer

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Layer Version.
* `delete` - (Defaults to 5 mins) Used when delete the Layer Version.
* `update` - (Defaults to 5 mins) Used when update the Layer Version.

## Import

FCV3 Layer Version can be imported using the id, e.g.

```shell
$ terraform import alicloud_fcv3_layer_version.example <layer_name>:<version>
```