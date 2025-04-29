---
subcategory: "Function Compute Service (FC)"
layout: "alicloud"
page_title: "Alicloud: alicloud_fc_layer_version"
sidebar_current: "docs-alicloud-resource-fc-layer-version"
description: |-
  Provides a Alicloud FC Layer Version resource.
---

# alicloud_fc_layer_version

Provides a Function Compute Layer Version resource.

For information about FC Layer Version and how to use it, see [What is Layer Version](https://www.alibabacloud.com/help/en/fc/developer-reference/api-fc-open-2021-04-06-createlayerversion).

-> **NOTE:** Available since v1.180.0.

-> **NOTE: Setting `skip_destroy` to `true` means that the Alicloud Provider will not destroy any layer version, even when running `terraform destroy`. Layer versions are thus intentional dangling resources that are not managed by Terraform and may incur extra expense in your Alicloud account.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_fc_layer_version&exampleId=91eb8b93-74c4-09dd-c73e-1c8b2aabeca4ffe5f7ed&activeTab=example&spm=docs.r.fc_layer_version.0.91eb8b9374&intl_lang=EN_US" target="_blank">
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
resource "alicloud_oss_bucket" "default" {
  bucket = "terraform-example-${random_integer.default.result}"
}
# If you upload the function by OSS Bucket, you need to specify path can't upload by content.
resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.id
  key     = "index.py"
  content = "import logging \ndef handler(event, context): \nlogger = logging.getLogger() \nlogger.info('hello world') \nreturn 'hello world'"
}

resource "alicloud_fc_layer_version" "example" {
  layer_name         = "terraform-example-${random_integer.default.result}"
  compatible_runtime = ["python2.7"]
  oss_bucket_name    = alicloud_oss_bucket.default.bucket
  oss_object_name    = alicloud_oss_bucket_object.default.key
}
```

## Argument Reference

The following arguments are supported:

* `layer_name` - (Required, ForceNew) The name of the layer.
* `description` - (Optional, ForceNew) The description of the layer version.
* `skip_destroy` - (Optional) Whether to retain the old version of a previously deployed Lambda Layer. Default is `false`. When this is not set to `true`, changing any of `compatible_runtimes`, `description`, `layer_name`, `oss_bucket_name`,  `oss_object_name`, or `zip_file` forces deletion of the existing layer version and creation of a new layer version.
* `compatible_runtime` - (Required, ForceNew) The list of runtime environments that are supported by the layer. Valid values: `nodejs14`, `nodejs12`, `nodejs10`, `nodejs8`, `nodejs6`, `python3.9`, `python3`, `python2.7`, `java11`, `java8`, `php7.2`, `go1`,`dotnetcore2.1`, `custom`.
* `oss_bucket_name` - (Optional, ForceNew) The name of the OSS bucket that stores the ZIP package of the function code.
* `oss_object_name` - (Optional, ForceNew) The name of the OSS object (ZIP package) that contains the function code.
* `zip_file` - (Optional, ForceNew) The ZIP package of the function code that is encoded in the Base64 format.

-> **NOTE:** `zip_file` and `oss_bucket_name`, `oss_object_name` cannot be used together.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Layer Version. The value formats as `<layer_name>:<version>`.
* `version` - The version of Layer Version.
* `acl` - The access mode of Layer Version.
* `arn` - The arn of Layer Version.
* `code_check_sum` - The checksum of the layer code package.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the FC Layer Version.
* `delete` - (Defaults to 1 mins) Used when delete the FC Layer Version.

## Import

Function Compute Layer Version can be imported using the id, e.g.

```shell
$ terraform import alicloud_fc_layer_version.example my_function
```