---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_user_defined_log_fields"
description: |-
  Provides a Alicloud OSS Bucket User Defined Log Fields resource.
---

# alicloud_oss_bucket_user_defined_log_fields

Provides a OSS Bucket User Defined Log Fields resource. Used to personalize the user_defined_log_fields field in the Bucket real-time log.

For information about OSS Bucket User Defined Log Fields and how to use it, see [What is Bucket User Defined Log Fields](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.224.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = "${var.name}-${random_integer.default.result}"
}


resource "alicloud_oss_bucket_user_defined_log_fields" "default" {
  bucket     = alicloud_oss_bucket.CreateBucket.bucket
  param_set  = ["oss-example", "example-para", "abc"]
  header_set = ["def", "example-header"]
}
```

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) The name of the bucket.
* `header_set` - (Optional) Container for custom request header configuration information.
* `param_set` - (Optional) Container for custom request parameters configuration information.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket User Defined Log Fields.
* `delete` - (Defaults to 5 mins) Used when delete the Bucket User Defined Log Fields.
* `update` - (Defaults to 5 mins) Used when update the Bucket User Defined Log Fields.

## Import

OSS Bucket User Defined Log Fields can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_user_defined_log_fields.example <id>
```