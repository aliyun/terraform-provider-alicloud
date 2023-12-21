---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_referer"
description: |-
  Provides a Alicloud OSS Bucket Referer resource.
---

# alicloud_oss_bucket_referer

Provides a OSS Bucket Referer resource. Bucket Referer configuration.

For information about OSS Bucket Referer and how to use it, see [What is Bucket Referer](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.214.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_oss_bucket_referer" "default" {
  bucket_name         = "qd-api-test"
  allow_empty_referer = true
  referer_list {
    referer = ["www.abc.com", "www.aliyun.com"]
  }
}
```

### Deleting `alicloud_oss_bucket_referer` or removing it from your configuration

Terraform cannot destroy resource `alicloud_oss_bucket_referer`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `allow_empty_referer` - (Required) Whether to allow empty Referer request headers.
* `allow_truncate_query_string` - (Required) Allow phase request parameters.
* `bucket_name` - (Required, ForceNew) Name of the Bucket.
* `referer_list` - (Optional) The container that holds the Referer whitelist. See [`referer_list`](#referer_list) below.
* `truncate_path` - (Optional) Name of the bucket.

### `referer_list`

The referer_list supports the following:
* `referer` - (Optional) Referer列表.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Referer.
* `update` - (Defaults to 5 mins) Used when update the Bucket Referer.

## Import

OSS Bucket Referer can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_referer.example <id>
```