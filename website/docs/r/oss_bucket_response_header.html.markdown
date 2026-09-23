---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_response_header"
description: |-
  Provides a Alicloud OSS Bucket Response Header resource.
---

# alicloud_oss_bucket_response_header

Provides a OSS Bucket Response Header resource.

Response header configuration of a bucket.

For information about OSS Bucket Response Header and how to use it, see [What is Bucket Response Header](https://next.api.alibabacloud.com/document/Oss/2019-05-17/PutBucketResponseHeader).

-> **NOTE:** Available since v1.271.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_response_header&exampleId=7dbdb3b4-4f39-6d7a-cbd5-5e656b7cb55b6a05d9b5&activeTab=example&spm=docs.r.oss_bucket_response_header.0.7dbdb3b44f&intl_lang=EN_US" target="_blank">
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


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_oss_bucket_response_header&spm=docs.r.oss_bucket_response_header.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) The name of the bucket
* `rule` - (Optional, List) The container that holds the response header rules. See [`rule`](#rule) below.

### `rule`

The rule supports the following:
* `filters` - (Optional, Set) The container that holds the operations that need to be apply rules. See [`filters`](#rule-filters) below.
* `hide_headers` - (Optional, Set) The container that holds the response headers that need to be hidden. See [`hide_headers`](#rule-hide_headers) below.
* `name` - (Optional) The response header rule name.

### `rule-filters`

The rule-filters supports the following:
* `operation` - (Optional, List) The operation to which the rule applies.

### `rule-hide_headers`

The rule-hide_headers supports the following:
* `header` - (Optional, List) The response header needs to be hidden.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Response Header.
* `delete` - (Defaults to 5 mins) Used when delete the Bucket Response Header.
* `update` - (Defaults to 5 mins) Used when update the Bucket Response Header.

## Import

OSS Bucket Response Header can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_response_header.example <bucket>
```