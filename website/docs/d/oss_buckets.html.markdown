---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_buckets"
sidebar_current: "docs-alicloud-datasource-oss-buckets"
description: |-
    Provides a list of OSS buckets to the user.
---

# alicloud_oss_buckets

This data source provides the OSS buckets of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.17.0.

## Example Usage

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "bucket" {
  bucket = "oss-tf-example-${random_integer.default.result}"
}

data "alicloud_oss_buckets" "oss_buckets_ds" {
  name_regex = alicloud_oss_bucket.bucket.bucket
}

output "first_oss_bucket_name" {
  value = data.alicloud_oss_buckets.oss_buckets_ds.buckets.0.name
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, ForceNew) A regex string to filter results by bucket name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of bucket names. 
* `buckets` - A list of buckets. Each element contains the following attributes:
  * `name` - Bucket name.
  * `acl` - Bucket access control list. Possible values: `private`, `public-read` and `public-read-write`.
  * `extranet_endpoint` - Internet domain name for accessing the bucket from outside.
  * `intranet_endpoint` - Intranet domain name for accessing the bucket from an ECS instance in the same region.
  * `location` - Region of the data center where the bucket is located.
  * `owner` - Bucket owner.
  * `storage_class` - Object storage type. Possible values: `Standard`, `IA`, `Archive` and `ColdArchive`.
  * `redundancy_type` - Redundancy type. Possible values: `LRS`, and `ZRS`.
  * `creation_date` - Bucket creation date.
  * `policy` - The policies configured for a specified bucket.
  * `cors_rules` - A list of CORS rule configurations. Each element contains the following attributes:
    * `allowed_origins` - The origins allowed for cross-domain requests. Multiple elements can be used to specify multiple allowed origins. Each rule allows up to one wildcard "\*". If "\*" is specified, cross-domain requests of all origins are allowed.
    * `allowed_methods` - Specify the allowed methods for cross-domain requests. Possible values: `GET`, `PUT`, `DELETE`, `POST` and `HEAD`.
    * `allowed_headers` - Control whether the headers specified by Access-Control-Request-Headers in the OPTIONS prefetch command are allowed. Each header specified by Access-Control-Request-Headers must match a value in AllowedHeader. Each rule allows up to one wildcard “*” .
    * `expose_headers` - Specify the response headers allowing users to access from an application (for example, a Javascript XMLHttpRequest object). The wildcard "\*" is not allowed.
    * `max_age_seconds` - Specify the cache time for the returned result of a browser prefetch (OPTIONS) request to a specific resource.
  * `website` - A list of one element containing configuration parameters used when the bucket is used as a website. It contains the following attributes:
    * `index_document` - Key of the HTML document containing the home page.
    * `error_document` - Key of the HTML document containing the error page.
  * `logging` - A list of one element containing configuration parameters used for storing access log information. It contains the following attributes:
    * `target_bucket` - Bucket for storing access logs.
    * `target_prefix` - Prefix of the saved access log file paths.
  * `referer_config` - A list of one element containing referer configuration. It contains the following attributes:
    * `allow_empty` - Indicate whether the access request referer field can be empty.
    * `referers` - Referer access whitelist.
  * `lifecycle_rule` - A list CORS of lifecycle configurations. When Lifecycle is enabled, OSS automatically deletes the objects or transitions the objects (to another storage class) corresponding the lifecycle rules on a regular basis. Each element contains the following attributes:
    * `id` - Unique ID of the rule.
    * `prefix` - Prefix applicable to a rule. Only those objects with a matching prefix can be affected by the rule.
    * `enabled` - Indicate whether the rule is enabled or not.
    * `expiration` - A list of one element containing expiration attributes of an object. It contains the following attributes:
      * `date` - Date after which the rule to take effect. The format is like 2017-03-09.
      * `days` - Indicate the number of days after the last object update until the rules take effect.
  * `server_side_encryption_rule` - A configuration of default encryption for a bucket. It contains the following attributes:
    * `sse_algorithm` - The server-side encryption algorithm to use.
    * `kms_master_key_id` -  The alibaba cloud KMS master key ID used for the SSE-KMS encryption. 
  * `tags` - A mapping of tags.
  * `versioning` - If present , the versioning state has been set on the bucket. It contains the following attribute.
      * `status` - A bucket versioning state. Possible values:`Enabled` and `Suspended`.
