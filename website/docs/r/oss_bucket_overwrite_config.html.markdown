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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_overwrite_config&exampleId=7dbdb3b4-4f39-6d7a-cbd5-5e656b7cb55b6a05d9b5&activeTab=example&spm=docs.r.oss_bucket_overwrite_config.0.7dbdb3b44f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_oss_bucket_overwrite_config&spm=docs.r.oss_bucket_overwrite_config.example&intl_lang=EN_US)

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