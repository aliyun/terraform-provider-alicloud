---
subcategory: "OSS"
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

```terraform
resource "alicloud_oss_bucket" "bucket-acl" {
  bucket = "bucket-170309-acl"
  acl    = "private"
}
```

Static Website

```terraform
resource "alicloud_oss_bucket" "bucket-website" {
  bucket = "bucket-170309-website"

  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}
```

Enable Logging

```terraform
resource "alicloud_oss_bucket" "bucket-target" {
  bucket = "bucket-170309-acl"
  acl    = "public-read"
}

resource "alicloud_oss_bucket" "bucket-logging" {
  bucket = "bucket-170309-logging"

  logging {
    target_bucket = alicloud_oss_bucket.bucket-target.id
    target_prefix = "log/"
  }
}
```

Referer configuration

```terraform
resource "alicloud_oss_bucket" "bucket-referer" {
  bucket = "bucket-170309-referer"
  acl    = "private"

  referer_config {
    allow_empty = false
    referers    = ["http://www.aliyun.com", "https://www.aliyun.com"]
  }
}
```

Set lifecycle rule

```terraform
resource "alicloud_oss_bucket" "bucket-lifecycle" {
  bucket = "bucket-170309-lifecycle"
  acl    = "public-read"

  lifecycle_rule {
    id      = "rule-days"
    prefix  = "path1/"
    enabled = true

    expiration {
      days = 365
    }
  }
  lifecycle_rule {
    id      = "rule-date"
    prefix  = "path2/"
    enabled = true

    expiration {
      date = "2018-01-12"
    }
  }
}

resource "alicloud_oss_bucket" "bucket-lifecycle" {
  bucket = "bucket-170309-lifecycle"
  acl    = "public-read"

  lifecycle_rule {
    id      = "rule-days-transition"
    prefix  = "path3/"
    enabled = true

    transitions {
      days          = "3"
      storage_class = "IA"
    }
    transitions {
      days          = "30"
      storage_class = "Archive"
    }
  }
}

resource "alicloud_oss_bucket" "bucket-lifecycle" {
  bucket = "bucket-170309-lifecycle"
  acl    = "public-read"

  lifecycle_rule {
    id      = "rule-days-transition"
    prefix  = "path3/"
    enabled = true

    transitions {
      created_before_date = "2020-11-11"
      storage_class       = "IA"
    }
    transitions {
      created_before_date = "2021-11-11"
      storage_class       = "Archive"
    }
  }
}

resource "alicloud_oss_bucket" "bucket-lifecycle" {
  bucket = "bucket-170309-lifecycle"
  acl    = "public-read"

  lifecycle_rule {
    id      = "rule-abort-multipart-upload"
    prefix  = "path3/"
    enabled = true

    abort_multipart_upload {
      days = 128
    }
  }
}

resource "alicloud_oss_bucket" "bucket-versioning-lifecycle" {
  bucket = "bucket-170309-lifecycle"
  acl    = "private"

  versioning {
    status = "Enabled"
  }

  lifecycle_rule {
    id      = "rule-versioning"
    prefix  = "path1/"
    enabled = true

    expiration {
      expired_object_delete_marker = true
    }

    noncurrent_version_expiration {
      days = 240
    }

    noncurrent_version_transition {
      days          = 180
      storage_class = "Archive"
    }

    noncurrent_version_transition {
      days          = 60
      storage_class = "IA"
    }
  }
}

```

Set bucket policy 

```terraform
resource "alicloud_oss_bucket" "bucket-policy" {
  bucket = "bucket-170309-policy"
  acl    = "private"

  policy = <<POLICY
  {"Statement":
      [{"Action":
          ["oss:PutObject", "oss:GetObject", "oss:DeleteBucket"],
        "Effect":"Allow",
        "Resource":
            ["acs:oss:*:*:*"]}],
   "Version":"1"}
  POLICY
}
```

IA Bucket

```terraform
resource "alicloud_oss_bucket" "bucket-storageclass" {
  bucket        = "bucket-170309-storageclass"
  storage_class = "IA"
}
```

Set bucket server-side encryption rule 

```terraform
resource "alicloud_oss_bucket" "bucket-sserule" {
  bucket = "bucket-170309-sserule"
  acl    = "private"

  server_side_encryption_rule {
    sse_algorithm = "AES256"
  }
}

resource "alicloud_oss_bucket" "bucket-sserule" {
  bucket = "bucket-170309-sserule"
  acl    = "private"

  server_side_encryption_rule {
    sse_algorithm     = "KMS"
    kms_master_key_id = "your kms key id"
  }
}
```

Set bucket tags 

```terraform
resource "alicloud_oss_bucket" "bucket-tags" {
  bucket = "bucket-170309-tags"
  acl    = "private"

  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}
```

Enable bucket versioning 

```terraform
resource "alicloud_oss_bucket" "bucket-versioning" {
  bucket = "bucket-170309-versioning"
  acl    = "private"

  versioning {
    status = "Enabled"
  }
}
```

Set bucket redundancy type

```terraform
resource "alicloud_oss_bucket" "bucket-redundancytype" {
  bucket          = "bucket_name"
  redundancy_type = "ZRS"

  # ... other configuration ...
}
```

Set bucket accelerate configuration

```terraform
resource "alicloud_oss_bucket" "bucket-accelerate" {
  bucket = "bucket_name"

  transfer_acceleration {
    enabled = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Optional, ForceNew) The name of the bucket. If omitted, Terraform will assign a random and unique name.
* `acl` - (Optional) The [canned ACL](https://www.alibabacloud.com/help/doc-detail/31898.htm) to apply. Can be "private", "public-read" and "public-read-write". Defaults to "private".
* `cors_rule` - (Optional) A rule of [Cross-Origin Resource Sharing](https://www.alibabacloud.com/help/doc-detail/31903.htm) (documented below). The items of core rule are no more than 10 for every OSS bucket.
* `website` - (Optional) A website object(documented below).
* `logging` - (Optional) A Settings of [bucket logging](https://www.alibabacloud.com/help/doc-detail/31900.htm) (documented below).
* `logging_isenable` - (Optional) The flag of using logging enable container. Defaults true.
* `referer_config` - (Optional) The configuration of [referer](https://www.alibabacloud.com/help/doc-detail/31901.htm) (documented below).
* `lifecycle_rule` - (Optional) A configuration of [object lifecycle management](https://www.alibabacloud.com/help/doc-detail/31904.htm) (documented below).
* `policy` - (Optional, Available in 1.41.0) Json format text of bucket policy [bucket policy management](https://www.alibabacloud.com/help/doc-detail/100680.htm).
* `storage_class` - (Optional, ForceNew) The [storage class](https://www.alibabacloud.com/help/doc-detail/51374.htm) to apply. Can be "Standard", "IA" and "Archive". Defaults to "Standard".
* `redundancy_type` - (Optional, ForceNew, Available in 1.91.0) The [redundancy type](https://www.alibabacloud.com/help/doc-detail/90589.htm) to enable. Can be "LRS", and "ZRS". Defaults to "LRS".
* `server_side_encryption_rule` - (Optional, Available in 1.45.0+) A configuration of server-side encryption (documented below).
* `tags` - (Optional, Available in 1.45.0+) A mapping of tags to assign to the bucket. The items are no more than 10 for a bucket.
* `versioning` - (Optional, Available in 1.45.0+) A state of versioning (documented below).
* `force_destroy` - (Optional, Available in 1.45.0+) A boolean that indicates all objects should be deleted from the bucket so that the bucket can be destroyed without error. These objects are not recoverable. Defaults to "false".
* `transfer_acceleration` - (Optional, Available in 1.123.1+) A transfer acceleration status of a bucket (documented below).

#### Block cors_rule

The cors_rule mapping supports the following:

* `allowed_headers` - (Optional) Specifies which headers are allowed.
* `allowed_methods` - (Required) Specifies which methods are allowed. Can be GET, PUT, POST, DELETE or HEAD.
* `allowed_origins` - (Required) Specifies which origins are allowed.
* `expose_headers` - (Optional) Specifies expose header in the response.
* `max_age_seconds` - (Optional) Specifies time in seconds that browser can cache the response for a preflight request.

#### Block website

The website mapping supports the following:

* `index_document` - (Required) Alicloud OSS returns this index document when requests are made to the root domain or any of the subfolders.
* `error_document` - (Optional) An absolute path to the document to return in case of a 4XX error.

#### Block logging

The logging object supports the following:

* `target_bucket` - (Required) The name of the bucket that will receive the log objects.
* `target_prefix` - (Optional) To specify a key prefix for log objects.

#### Block referer configuration

The referer configuration supports the following:

* `allow_empty` - (Optional, Type: bool) Allows referer to be empty. Defaults false.
* `referers` - (Required, Type: list) The list of referer.

#### Block lifecycle_rule

The lifecycle_rule object supports the following:

* `id` - (Optional) Unique identifier for the rule. If omitted, OSS bucket will assign a unique name.
* `prefix` - (Optional, Available in v1.90.0+) Object key prefix identifying one or more objects to which the rule applies. Default value is null, the rule applies to all objects in a bucket.
* `enabled` - (Required, Type: bool) Specifies lifecycle rule status.
* `expiration` - (Optional, Type: set) Specifies a period in the object's expire (documented below).
* `transitions` - (Optional, Type: set, Available in 1.62.1+) Specifies the time when an object is converted to the IA or archive storage class during a valid life cycle. (documented below).
* `abort_multipart_upload` - (Optional, Type: set, Available in 1.121.2+) Specifies the number of days after initiating a multipart upload when the multipart upload must be completed (documented below).
* `noncurrent_version_expiration` - (Optional, Type: set, Available in 1.121.2+) Specifies when noncurrent object versions expire (documented below).
* `noncurrent_version_transition` - (Optional, Type: set, Available in 1.121.2+) Specifies when noncurrent object versions transitions (documented below).

`NOTE`: At least one of expiration, transitions, abort_multipart_upload, noncurrent_version_expiration and noncurrent_version_transition should be configured.

#### Block expiration

The lifecycle_rule expiration object supports the following:

* `date` - (Optional) Specifies the date after which you want the corresponding action to take effect. The value obeys ISO8601 format like `2017-03-09`.
* `days` - (Optional, Type: int) Specifies the number of days after object creation when the specific rule action takes effect.
* `created_before_date` - (Optional, Available in 1.121.2+) Specifies the time before which the rules take effect. The date must conform to the ISO8601 format and always be UTC 00:00. For example: 2002-10-11T00:00:00.000Z indicates that objects updated before 2002-10-11T00:00:00.000Z are deleted or converted to another storage class, and objects updated after this time (including this time) are not deleted or converted.
* `expired_object_delete_marker` - (Optional, Type: bool, Available in 1.121.2+) On a versioned bucket (versioning-enabled or versioning-suspended bucket), you can add this element in the lifecycle configuration to direct OSS to delete expired object delete markers. This cannot be specified with Days, Date or CreatedBeforeDate in a Lifecycle Expiration Policy.

`NOTE`: One and only one of "date", "days", "created_before_date" and "expired_object_delete_marker" can be specified in one expiration configuration.

#### Block transitions

The lifecycle_rule transitions object supports the following:

* `created_before_date` - (Optional) Specifies the time before which the rules take effect. The date must conform to the ISO8601 format and always be UTC 00:00. For example: 2002-10-11T00:00:00.000Z indicates that objects updated before 2002-10-11T00:00:00.000Z are deleted or converted to another storage class, and objects updated after this time (including this time) are not deleted or converted.
* `days` - (Optional, Type: int) Specifies the number of days after object creation when the specific rule action takes effect.
* `storage_class` - (Required) Specifies the storage class that objects that conform to the rule are converted into. The storage class of the objects in a bucket of the IA storage class can be converted into Archive but cannot be converted into Standard. Values: `IA`, `Archive`. 

`NOTE`: One and only one of "created_before_date" and "days" can be specified in one transition configuration.

#### Block abort_multipart_upload

The lifecycle_rule abort_multipart_upload object supports the following:

* `created_before_date` - (Optional) Specifies the time before which the rules take effect. The date must conform to the ISO8601 format and always be UTC 00:00. For example: 2002-10-11T00:00:00.000Z indicates that parts created before 2002-10-11T00:00:00.000Z are deleted, and parts created after this time (including this time) are not deleted.
* `days` - (Optional, Type: int) Specifies the number of days after object creation when the specific rule action takes effect.

`NOTE`: One and only one of "created_before_date" and "days" can be specified in one abort_multipart_upload configuration.

#### Block noncurrent_version_expiration

The lifecycle_rule noncurrent_version_expiration object supports the following:

* `days` - (Required, Type: int) Specifies the number of days noncurrent object versions expire.

#### Block noncurrent_version_transition

The lifecycle_rule noncurrent_version_transition object supports the following:

* `days` - (Required, Type: int) Specifies the number of days noncurrent object versions transition.
* `storage_class` - (Required) Specifies the storage class that objects that conform to the rule are converted into. The storage class of the objects in a bucket of the IA storage class can be converted into Archive but cannot be converted into Standard. Values: `IA`, `Archive`. 


#### Block server-side encryption rule

The server-side encryption rule supports the following:

* `sse_algorithm` - (Required) The server-side encryption algorithm to use. Possible values: `AES256` and `KMS`.
* `kms_master_key_id` -  (optional, Available in 1.92.0+) The alibaba cloud KMS master key ID used for the SSE-KMS encryption. 

#### Block versioning

The versioning supports the following:

* `status` - (Required) Specifies the versioning state of a bucket. Valid values: `Enabled` and `Suspended`.

`NOTE`: Currently, the `versioning` feature is only available in ap-south-1 and with white list. If you want to use it, please contact us.

#### Block transfer_acceleration

The transfer_acceleration supports the following:

* `enabled` - (Required, Type: bool) Specifies the accelerate status of a bucket.

## Attributes Reference

The following attributes are exported:

* `id` - The name of the bucket.
* `acl` - The acl of the bucket.
* `creation_date` - The creation date of the bucket.
* `extranet_endpoint` - The extranet access endpoint of the bucket.
* `intranet_endpoint` - The intranet access endpoint of the bucket.
* `location` - The location of the bucket.
* `owner` - The bucket owner.

## Import

OSS bucket can be imported using the bucket name, e.g.

```
$ terraform import alicloud_oss_bucket.bucket bucket-12345678
```
