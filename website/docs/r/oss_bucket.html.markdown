---
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket"
sidebar_current: "docs-alicloud-resource-oss-bucket"
description: |-
  Provides a resource to create a oss bucket.
---

# alicloud\_oss\_bucket

Provides a resource to create a oss bucket and set its attribution.

-> **NOTE:** The bucket namespace is shared by all users of the OSS system. Please set bucket name as unique as possible.


## Example Usage

Private Bucket

```
resource "alicloud_oss_bucket" "bucket-acl"{
  bucket = "bucket-170309-acl"
  acl = "private"
}
```

Static Website

```
resource "alicloud_oss_bucket" "bucket-website" {
  bucket = "bucket-170309-website"

  website = {
    index_document = "index.html"
    error_document = "error.html"
  }
}
```

Enable Logging

```
resource "alicloud_oss_bucket" "bucket-target"{
  bucket = "bucket-170309-acl"
  acl = "public-read"
}

resource "alicloud_oss_bucket" "bucket-logging" {
  bucket = "bucket-170309-logging"

  logging {
    target_bucket = "${alicloud_oss_bucket.bucket-target.id}"
    target_prefix = "log/"
  }
}
```

Referer configuration

```
resource "alicloud_oss_bucket" "bucket-referer" {
  bucket = "bucket-170309-referer"
  acl = "private"

  referer_config {
      allow_empty = false
      referers = ["http://www.aliyun.com", "https://www.aliyun.com"]
  }
}
```

Set lifecycle rule

```
resource "alicloud_oss_bucket" "bucket-lifecycle" {
  bucket = "bucket-170309-lifecycle"
  acl = "public-read"

  lifecycle_rule {
    id = "rule-days"
    prefix = "path1/"
    enabled = true

    expiration {
      days = 365
    }
  }
  lifecycle_rule {
    id = "rule-date"
    prefix = "path2/"
    enabled = true

    expiration {
      date = "2018-01-12"
    }
  }
}
```
## Argument Reference

The following arguments are supported:

* `bucket` - (Optional, ForceNew) The name of the bucket. If omitted, Terraform will assign a random and unique name.
* `acl` - (Optional) The [canned ACL](https://www.alibabacloud.com/help/doc-detail/31898.htm) to apply. Defaults to "private".
* `cors_rule` - (Optional) A list rules of [Cross-Origin Resource Sharing](https://www.alibabacloud.com/help/doc-detail/31903.htm) (documented below). The items of cors rule are no more than 10 for every OSS bucket.
* `website` - (Optional) A list website objects(documented below). The items of website are no more than 1 for every OSS bucket.
* `logging` - (Optional) A list settings of [bucket logging](https://www.alibabacloud.com/help/doc-detail/31900.htm) (documented below). The items of logging are no more than 1 for every OSS bucket.
* `logging_isenable` - (Deprecated) It has been deprecated from 1.37.0. When `logging` is set, the bucket logging will be able.
* `referer_config` - (Optional) A list configurations of [referer](https://www.alibabacloud.com/help/doc-detail/31901.htm) (documented below). The items of referer_config are no more than 1 for every OSS bucket.
* `lifecycle_rule` - (Optional) A list configurations of [object lifecycle management](https://www.alibabacloud.com/help/doc-detail/31904.htm) (documented below). The items of rules are no more than 1000 for every OSS bucket.

### Block cors_rule

The cors_rule mapping supports the following:

* `allowed_headers` - (Optional) Specifies which headers are allowed.
* `allowed_methods` - (Required) Specifies which methods are allowed. Can be GET, PUT, POST, DELETE or HEAD.
* `allowed_origins` - (Required) Specifies which origins are allowed.
* `expose_headers` - (Optional) Specifies expose header in the response.
* `max_age_seconds` - (Optional) Specifies time in seconds that browser can cache the response for a preflight request.

### Block website

The website mapping supports the following:

* `index_document` - (Required) Alicloud OSS returns this index document when requests are made to the root domain or any of the subfolders.
* `error_document` - (Optional) An absolute path to the document to return in case of a 4XX error.

### Block logging

The logging object supports the following:

* `target_bucket` - (Required) The name of the bucket that will receive the log objects.
* `target_prefix` - (Optional) To specify a key prefix for log objects.

### Block referer configuration

The referer configuration supports the following:

* `allow_empty` - (Optional, Type: bool) Allows referer to be empty. Defaults true.
* `referers` - (Required, Type: list) The list of referer.

### Block lifecycle_rule

The lifecycle_rule object supports the following:

* `id` - (Optional) Unique identifier for the rule. If omitted, OSS bucket will assign a unique name.
* `prefix` - (Required) Object key prefix identifying one or more objects to which the rule applies.
* `enabled` - (Required, Type: bool) Specifies lifecycle rule status.
* `expiration` - (Optional, Required, Type: set) Specifies a period in the object's expire (documented below).

#### Block expiration

The lifecycle_rule expiration object supports the following:

* `date` - (Optional) Specifies the date after which you want the corresponding action to take effect. The value obeys ISO8601 format like `2017-03-09`.
* `days` - (Optional, Type: int) Specifies the number of days after object creation when the specific rule action takes effect.

`NOTE`: One and only one of "date" and "days" can be specified in one expiration configuration.

## Attributes Reference

The following attributes are exported:

* `id` - The name of the bucket.
* `acl` - The acl of the bucket.
* `creation_date` - The creation date of the bucket.
* `extranet_endpoint` - The extranet access endpoint of the bucket.
* `intranet_endpoint` - The intranet access endpoint of the bucket.
* `location` - The location of the bucket.
* `owner` - The bucket owner.
* `storage_class` - The bucket storage type.

## Import

OSS bucket can be imported using the bucket name, e.g.

```
$ terraform import alicloud_oss_bucket.bucket bucket-12345678
```