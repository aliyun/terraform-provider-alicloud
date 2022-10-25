---
subcategory: "Function Compute Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_fc_layer_version"
sidebar_current: "docs-alicloud-resource-fc-layer-version"
description: |-
  Provides a Alicloud FC Layer Version resource.
---

# alicloud\_fc\_layer\_version

Provides a Function Compute Layer Version resource.

For information about FC Layer Version and how to use it, see [What is Layer Version](https://www.alibabacloud.com/help/en/icms-test/latest/api-doc-pre-fc-open-2021-04-06-api-doc-createlayerversion).

-> **NOTE:** Available in v1.180.0+.

-> **NOTE: Setting `skip_destroy` to `true` means that the Alicloud Provider will not destroy any layer version, even when running `terraform destroy`. Layer versions are thus intentional dangling resources that are not managed by Terraform and may incur extra expense in your Alicloud account.

## Example Usage

Basic Usage

```terraform
resource "alicloud_fc_layer_version" "example" {
  layer_name         = "your_layer_name"
  compatible_runtime = ["nodejs12"]
  oss_bucket_name    = "your_code_oss_bucket_name"
  oss_object_name    = "your_code_oss_object_name"
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

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the FC Layer Version.
* `delete` - (Defaults to 1 mins) Used when delete the FC Layer Version.

## Import

Function Compute Layer Version can be imported using the id, e.g.

```
$ terraform import alicloud_fc_layer_version.example my_function
```