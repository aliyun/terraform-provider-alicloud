---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_object"
sidebar_current: "docs-alicloud-resource-oss-bucket-object"
description: |-
  Provides a resource to create a oss bucket object.
---

# alicloud\_oss\_bucket\_object

Provides a resource to put a object(content or file) to a oss bucket.

## Example Usage

### Uploading a file to a bucket

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_object&exampleId=61f00540-416c-d582-0484-f3dabd21d0a6b92739d4&activeTab=example&spm=docs.r.oss_bucket_object.0.61f0054041&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "default" {
  bucket = "terraform-example-${random_integer.default.result}"
}

resource "alicloud_oss_bucket_acl" "default" {
  bucket = alicloud_oss_bucket.default.bucket
  acl    = "private"
}

resource "alicloud_oss_bucket_object" "default" {
  bucket = alicloud_oss_bucket.default.bucket
  key    = "example_key"
  source = "./main.tf"
}
```

### Uploading a content to a bucket

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_object&exampleId=dade73b3-5a61-2529-6fe5-6dd66693916a63b3c7f5&activeTab=example&spm=docs.r.oss_bucket_object.1.dade73b35a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "default" {
  bucket = "terraform-example-${random_integer.default.result}"
}

resource "alicloud_oss_bucket_acl" "default" {
  bucket = alicloud_oss_bucket.default.bucket
  acl    = "private"
}


resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.bucket
  key     = "example_key"
  content = "the content that you want to upload."
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_oss_bucket_object&spm=docs.r.oss_bucket_object.example&intl_lang=EN_US)

## Argument Reference

-> **Note:** If you specify `content_encoding` you are responsible for encoding the body appropriately (i.e. `source` and `content` both expect already encoded/compressed bytes)

The following arguments are supported:

* `bucket` - (Required, ForceNew) The name of the bucket to put the file in.
* `key` - (Required, ForceNew) The name of the object once it is in the bucket.
* `source` - (Optional) The path to the source file being uploaded to the bucket.
* `content` - (Optional unless `source` given) The literal content being uploaded to the bucket.
* `acl` - (Optional) The [canned ACL](https://www.alibabacloud.com/help/doc-detail/52284.htm) to apply. Defaults to "private".
* `content_type` - (Optional) A standard MIME type describing the format of the object data, e.g. application/octet-stream. All Valid MIME Types are valid for this input.
* `cache_control` - (Optional) Specifies caching behavior along the request/reply chain. Read [RFC2616 Cache-Control](https://www.ietf.org/rfc/rfc2616.txt) for further details.
* `content_disposition` - (Optional) Specifies presentational information for the object. Read [RFC2616 Content-Disposition](https://www.ietf.org/rfc/rfc2616.txt) for further details.
* `content_encoding` - (Optional) Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field. Read [RFC2616 Content-Encoding](https://www.ietf.org/rfc/rfc2616.txt) for further details.
* `content_md5` - (Optional) The MD5 value of the content. Read [MD5](https://www.alibabacloud.com/help/doc-detail/31978.htm) for computing method.
* `expires` - (Optional) Specifies expire date for the the request/response. Read [RFC2616 Expires](https://www.ietf.org/rfc/rfc2616.txt) for further details.
* `server_side_encryption` - (Optional) Specifies server-side encryption of the object in OSS. Valid values are `AES256`, `KMS`. Default value is `AES256`.
* `kms_key_id` - (Optional, Available in 1.62.1+) Specifies the primary key managed by KMS. This parameter is valid when the value of `server_side_encryption` is set to KMS.

Either `source` or `content` must be provided to specify the bucket content.
These two arguments are mutually-exclusive.

## Attributes Reference

The following attributes are exported

* `id` - the `key` of the resource supplied above.
* `content_length` - the content length of request.
* `etag` - the ETag generated for the object (an MD5 sum of the object content).
* `version_id` - A unique version ID value for the object, if bucket versioning is enabled.
