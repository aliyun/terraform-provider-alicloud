---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_overwrite_config"
description: |-
  Provides a Alicloud OSS Bucket Overwrite Config resource.
---

# alicloud_oss_bucket_overwrite_config

Provides a OSS Bucket Overwrite Config resource.

Bucket Overwrite Configuration.

For information about OSS Bucket Overwrite Config and how to use it, see [What is Bucket Overwrite Config](https://next.api.alibabacloud.com/document/Oss/2019-05-17/PutBucketOverwriteConfig).

-> **NOTE:** Available since v1.271.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = ""
}

resource "alicloud_oss_bucket" "defaultrdrM3m" {
  storage_class = "Standard"
}


resource "alicloud_oss_bucket_overwrite_config" "default" {
  bucket = alicloud_oss_bucket.defaultrdrM3m.id
  rule {
    id     = "rule1"
    action = "forbid"
    prefix = "rule1-prefix/"
    suffix = "rule1-suffix/"
    principals {
      principal = ["a", "b", "c"]
    }
  }
  rule {
    id     = "rule2"
    action = "forbid"
    prefix = "rule2-prefix/"
    suffix = "rule2-suffix/"
    principals {
      principal = ["d", "e", "f"]
    }
  }
  rule {
    id     = "rule3"
    action = "forbid"
    prefix = "rule3-prefix/"
    suffix = "rule3-suffix/"
    principals {
      principal = ["1", "2", "3"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) The name of the bucket
* `rule` - (Optional, List) Forbid overwrite rule See [`rule`](#rule) below.

### `rule`

The rule supports the following:
* `action` - (Optional) The operation type. Currently, only "forbid" is supported.
* `id` - (Optional) Rule ID
* `prefix` - (Optional) The prefix of the Object name, which is used to filter objects to be processed.
* `principals` - (Optional, Set) A collection of authorized principals. The usage is similar to that of the Principal of the Bucket Policy. You can enter the primary account, sub-account, or role. If this parameter is empty or not configured, overwriting is not allowed for objects that meet the preceding and suffix conditions. See [`principals`](#rule-principals) below.
* `suffix` - (Optional) The suffix of the Object name, which is used to filter objects to be processed.

### `rule-principals`

The rule-principals supports the following:
* `principal` - (Optional, List) Authorized subject. Supports the input of primary accounts, sub-accounts, or roles. Invalid setting if the value is empty.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Overwrite Config.
* `delete` - (Defaults to 5 mins) Used when delete the Bucket Overwrite Config.
* `update` - (Defaults to 5 mins) Used when update the Bucket Overwrite Config.

## Import

OSS Bucket Overwrite Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_overwrite_config.example <bucket>
```