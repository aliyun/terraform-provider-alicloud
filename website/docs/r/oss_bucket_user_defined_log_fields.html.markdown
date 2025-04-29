---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_user_defined_log_fields"
description: |-
  Provides a Alicloud OSS Bucket User Defined Log Fields resource.
---

# alicloud_oss_bucket_user_defined_log_fields

Provides a OSS Bucket User Defined Log Fields resource. Used to personalize the user_defined_log_fields field in the Bucket real-time log.

For information about OSS Bucket User Defined Log Fields and how to use it, see [What is Bucket User Defined Log Fields](https://www.alibabacloud.com/help/en/oss/developer-reference/putuserdefinedlogfieldsconfig).

-> **NOTE:** Available since v1.224.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_user_defined_log_fields&exampleId=a22fcf03-4ef5-2916-0e68-e9e77a61b5ce8f914c3a&activeTab=example&spm=docs.r.oss_bucket_user_defined_log_fields.0.a22fcf034e&intl_lang=EN_US" target="_blank">
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

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket User Defined Log Fields.
* `delete` - (Defaults to 5 mins) Used when delete the Bucket User Defined Log Fields.
* `update` - (Defaults to 5 mins) Used when update the Bucket User Defined Log Fields.

## Import

OSS Bucket User Defined Log Fields can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_user_defined_log_fields.example <id>
```