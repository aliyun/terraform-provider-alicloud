---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_objects"
sidebar_current: "docs-alicloud-datasource-oss-bucket-objects"
description: |-
    Provides a list of bucket objects to the user.
---

# alicloud\_oss_bucket_objects

This data source provides the objects of an OSS bucket.

## Example Usage

```
data "alicloud_oss_bucket_objects" "bucket_objects_ds" {
  bucket_name = "sample_bucket"
  key_regex   = "sample/sample_object.txt"
}

output "first_object_key" {
  value = "${data.alicloud_oss_bucket_objects.bucket_objects_ds.objects.0.key}"
}
```

## Argument Reference

The following arguments are supported:

* `bucket_name` - Name of the bucket that contains the objects to find.
* `key_regex` - (Optional) A regex string to filter results by key.
* `key_prefix` - (Optional) Filter results by the given key prefix (such as "path/to/folder/logs-").
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `objects` - A list of bucket objects. Each element contains the following attributes:
  * `key` - Object key.
  * `acl` - Object access control list. Possible values: `default`, `private`, `public-read` and `public-read-write`.
  * `content_type` - Standard MIME type describing the format of the object data, e.g. "application/octet-stream".
  * `content_length` - Size of the object in bytes.
  * `cache_control` - Caching behavior along the request/reply chain. Read [RFC2616 Cache-Control](https://www.ietf.org/rfc/rfc2616.txt) for further details.
  * `content_disposition` - Presentational information for the object. Read [RFC2616 Content-Disposition](https://www.ietf.org/rfc/rfc2616.txt) for further details.
  * `content_encoding` - Content encodings that have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field. Read [RFC2616 Content-Encoding](https://www.ietf.org/rfc/rfc2616.txt) for further details.
  * `content_md5` - MD5 value of the content. Read [MD5](https://www.alibabacloud.com/help/doc-detail/31978.htm) for computing method.
  * `expires` - Expiration date for the the request/response. Read [RFC2616 Expires](https://www.ietf.org/rfc/rfc2616.txt) for further details.
  * `server_side_encryption` - Server-side encryption of the object in OSS. It can be empty or `AES256`.
  * `sse_kms_key_id` - If present, specifies the ID of the Key Management Service(KMS) master encryption key that was used for the object.
  * `etag` - ETag generated for the object (MD5 sum of the object content).
  * `storage_class` - Object storage type. Possible values: `Standard`, `IA`, `Archive` and `ColdArchive`.
  * `last_modification_time` - Last modification time of the object.
