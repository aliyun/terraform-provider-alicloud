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

-> **NOTE:** Available since v1.2.0.

## Example Usage

Private Bucket

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "bucket-acl" {
  bucket = "example-value-${random_integer.default.result}"
  acl    = "private"
}
```

Static Website

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "bucket-website" {
  bucket = "example-value-${random_integer.default.result}"
  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}
```

Enable Logging

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "bucket-target" {
  bucket = "example-value-${random_integer.default.result}"
  acl    = "public-read"
}

resource "alicloud_oss_bucket" "bucket-logging" {
  bucket = "example-logging-${random_integer.default.result}"
  logging {
    target_bucket = alicloud_oss_bucket.bucket-target.id
    target_prefix = "log/"
  }
}
```

Referer configuration

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "bucket-referer" {
  bucket = "example-value-${random_integer.default.result}"
  acl    = "private"
  referer_config {
    allow_empty = false
    referers    = ["http://www.aliyun.com", "https://www.aliyun.com"]
  }
}
```

Set lifecycle rule

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "bucket-lifecycle1" {
  bucket = "example-lifecycle1-${random_integer.default.result}"
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

resource "alicloud_oss_bucket" "bucket-lifecycle2" {
  bucket = "example-lifecycle2-${random_integer.default.result}"
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

resource "alicloud_oss_bucket" "bucket-lifecycle3" {
  bucket = "example-lifecycle3-${random_integer.default.result}"
  acl    = "public-read"

  lifecycle_rule {
    id      = "rule-days-transition"
    prefix  = "path3/"
    enabled = true

    transitions {
      created_before_date = "2022-11-11"
      storage_class       = "IA"
    }
    transitions {
      created_before_date = "2021-11-11"
      storage_class       = "Archive"
    }
  }
}

resource "alicloud_oss_bucket" "bucket-lifecycle4" {
  bucket = "example-lifecycle4-${random_integer.default.result}"
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
  bucket = "example-lifecycle5-${random_integer.default.result}"
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
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "bucket-policy" {
  bucket = "example-policy-${random_integer.default.result}"
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
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "default" {
  bucket        = "example-${random_integer.default.result}"
  storage_class = "IA"
}
```

Set bucket server-side encryption rule 

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "bucket-sserule" {
  bucket = "terraform-example-${random_integer.default.result}"
  acl    = "private"

  server_side_encryption_rule {
    sse_algorithm = "AES256"
  }
}

resource "alicloud_kms_key" "kms" {
  description            = "terraform-example"
  pending_window_in_days = "7"
  status                 = "Enabled"
}

resource "alicloud_oss_bucket" "bucket-kms" {
  bucket = "terraform-example-kms-${random_integer.default.result}"
  acl    = "private"

  server_side_encryption_rule {
    sse_algorithm     = "KMS"
    kms_master_key_id = alicloud_kms_key.kms.id
  }
}
```

Set bucket tags 

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "bucket-tags" {
  bucket = "terraform-example-${random_integer.default.result}"
  acl    = "private"

  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}
```

Enable bucket versioning 

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "bucket-versioning" {
  bucket = "terraform-example-${random_integer.default.result}"
  acl    = "private"

  versioning {
    status = "Enabled"
  }
}
```

Set bucket redundancy type

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "bucket-redundancytype" {
  bucket          = "terraform-example-${random_integer.default.result}"
  redundancy_type = "ZRS"

  # ... other configuration ...
}
```

Set bucket accelerate configuration

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "bucket-accelerate" {
  bucket = "terraform-example-${random_integer.default.result}"

  transfer_acceleration {
    enabled = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Optional, ForceNew) The name of the bucket. If omitted, Terraform will assign a random and unique name.
* `acl` - (Optional) The [canned ACL](https://www.alibabacloud.com/help/doc-detail/31898.htm) to apply. Can be "private", "public-read" and "public-read-write". Defaults to "private".
* `cors_rule` - (Optional) A rule of  [Cross-Origin Resource Sharing](https://www.alibabacloud.com/help/doc-detail/31903.htm). The items of core rule are no more than 10 for every OSS bucket. See [`cors_rule`](#cors_rule) below.
* `website` - (Optional) A website configuration. See [`website`](#website) below.
* `logging` - (Optional) A Settings of [bucket logging](https://www.alibabacloud.com/help/doc-detail/31900.htm). See [`logging`](#logging) below.
* `logging_isenable` - (Optional, Deprecated from 1.37.0.) The flag of using logging enable container. Defaults true.
* `referer_config` - (Optional) The configuration of [referer](https://www.alibabacloud.com/help/doc-detail/31901.htm). See [`referer_config`](#referer_config) below.
* `lifecycle_rule` - (Optional) A configuration of [object lifecycle management](https://www.alibabacloud.com/help/doc-detail/31904.htm). See [`lifecycle_rule`](#lifecycle_rule) below.
* `policy` - (Optional, Available since 1.41.0) Json format text of bucket policy [bucket policy management](https://www.alibabacloud.com/help/doc-detail/100680.htm).
* `storage_class` - (Optional, ForceNew) The [storage class](https://www.alibabacloud.com/help/doc-detail/51374.htm) to apply. Can be "Standard", "IA", "Archive" and "ColdArchive". Defaults to "Standard". "ColdArchive" is available since 1.203.0.
* `redundancy_type` - (Optional, ForceNew, Available since 1.91.0) The [redundancy type](https://www.alibabacloud.com/help/doc-detail/90589.htm) to enable. Can be "LRS", and "ZRS". Defaults to "LRS".
* `server_side_encryption_rule` - (Optional, Available since 1.45.0) A configuration of server-side encryption. See [`server_side_encryption_rule`](#server_side_encryption_rule) below.
* `tags` - (Optional, Available since 1.45.0) A mapping of tags to assign to the bucket. The items are no more than 10 for a bucket.
* `versioning` - (Optional, Available since 1.45.0) A state of versioning. See [`versioning`](#versioning) below.
* `force_destroy` - (Optional, Available since 1.45.0) A boolean that indicates all objects should be deleted from the bucket so that the bucket can be destroyed without error. These objects are not recoverable. Defaults to "false".
* `transfer_acceleration` - (Optional, Available since 1.123.1) A transfer acceleration status of a bucket. See [`transfer_acceleration`](#transfer_acceleration) below.
* `lifecycle_rule_allow_same_action_overlap` - (Optional, Available since 1.208.1) A boolean that indicates lifecycle rules allow prefix overlap.

### `cors_rule`

The cors_rule configuration block supports the following:

* `allowed_headers` - (Optional) Specifies which headers are allowed.
* `allowed_methods` - (Required) Specifies which methods are allowed. Can be GET, PUT, POST, DELETE or HEAD.
* `allowed_origins` - (Required) Specifies which origins are allowed.
* `expose_headers` - (Optional) Specifies expose header in the response.
* `max_age_seconds` - (Optional) Specifies time in seconds that browser can cache the response for a preflight request.

### `website`

The website configuration block supports the following:

* `index_document` - (Required) Alicloud OSS returns this index document when requests are made to the root domain or any of the subfolders.
* `error_document` - (Optional) An absolute path to the document to return in case of a 4XX error.

### `logging`

The logging configuration block supports the following:

* `target_bucket` - (Required) The name of the bucket that will receive the log objects.
* `target_prefix` - (Optional) To specify a key prefix for log objects.

### `referer_config`

The referer_config configuration block supports the following:

* `allow_empty` - (Optional, Type: bool) Allows referer to be empty. Defaults false.
* `referers` - (Required, Type: list) The list of referer.

### `lifecycle_rule`

The lifecycle_rule configuration block supports the following:

* `id` - (Optional) Unique identifier for the rule. If omitted, OSS bucket will assign a unique name.
* `prefix` - (Optional, Available since v1.90.0) Object key prefix identifying one or more objects to which the rule applies. Default value is null, the rule applies to all objects in a bucket.
* `enabled` - (Required, Type: bool) Specifies lifecycle rule status.
* `expiration` - (Optional, Type: set) Specifies a period in the object's expire. See [`expiration`](#lifecycle_rule-expiration) below.
* `transitions` - (Optional, Type: set, Available since 1.62.1) Specifies the time when an object is converted to the IA or archive storage class during a valid life cycle. See [`transitions`](#lifecycle_rule-transitions) below.
* `abort_multipart_upload` - (Optional, Type: set, Available since 1.121.2) Specifies the number of days after initiating a multipart upload when the multipart upload must be completed. See [`abort_multipart_upload`](#lifecycle_rule-abort_multipart_upload) below.
* `noncurrent_version_expiration` - (Optional, Type: set, Available since 1.121.2) Specifies when noncurrent object versions expire. See [`noncurrent_version_expiration`](#lifecycle_rule-noncurrent_version_expiration) below.
* `noncurrent_version_transition` - (Optional, Type: set, Available since 1.121.2) Specifies when noncurrent object versions transitions. See [`noncurrent_version_transition`](#lifecycle_rule-noncurrent_version_transition) below.

`NOTE`: At least one of expiration, transitions, abort_multipart_upload, noncurrent_version_expiration and noncurrent_version_transition should be configured.

### `lifecycle_rule-expiration`

The expiration configuration block supports the following:

* `date` - (Optional) Specifies the date after which you want the corresponding action to take effect. The value obeys ISO8601 format like `2017-03-09`.
* `days` - (Optional, Type: int) Specifies the number of days after object creation when the specific rule action takes effect.
* `created_before_date` - (Optional, Available since 1.121.2) Specifies the time before which the rules take effect. The date must conform to the ISO8601 format and always be UTC 00:00. For example: 2002-10-11T00:00:00.000Z indicates that objects updated before 2002-10-11T00:00:00.000Z are deleted or converted to another storage class, and objects updated after this time (including this time) are not deleted or converted.
* `expired_object_delete_marker` - (Optional, Type: bool, Available since 1.121.2) On a versioned bucket (versioning-enabled or versioning-suspended bucket), you can add this element in the lifecycle configuration to direct OSS to delete expired object delete markers. This cannot be specified with Days, Date or CreatedBeforeDate in a Lifecycle Expiration Policy.

`NOTE`: One and only one of "date", "days", "created_before_date" and "expired_object_delete_marker" can be specified in one expiration configuration.

### `lifecycle_rule-transitions`

The transitions configuration block supports the following:

* `created_before_date` - (Optional) Specifies the time before which the rules take effect. The date must conform to the ISO8601 format and always be UTC 00:00. For example: 2002-10-11T00:00:00.000Z indicates that objects updated before 2002-10-11T00:00:00.000Z are deleted or converted to another storage class, and objects updated after this time (including this time) are not deleted or converted.
* `days` - (Optional, Type: int) Specifies the number of days after object creation when the specific rule action takes effect.
* `storage_class` - (Required) Specifies the storage class that objects that conform to the rule are converted into. The storage class of the objects in a bucket of the IA storage class can be converted into Archive but cannot be converted into Standard. Values: `IA`, `Archive`, `ColdArchive`. ColdArchive is available since 1.203.0.

`NOTE`: One and only one of "created_before_date" and "days" can be specified in one transition configuration.

### `lifecycle_rule-abort_multipart_upload`

The abort_multipart_upload configuration block supports the following:

* `created_before_date` - (Optional) Specifies the time before which the rules take effect. The date must conform to the ISO8601 format and always be UTC 00:00. For example: 2002-10-11T00:00:00.000Z indicates that parts created before 2002-10-11T00:00:00.000Z are deleted, and parts created after this time (including this time) are not deleted.
* `days` - (Optional, Type: int) Specifies the number of days after object creation when the specific rule action takes effect.

`NOTE`: One and only one of "created_before_date" and "days" can be specified in one abort_multipart_upload configuration.

### `lifecycle_rule-noncurrent_version_expiration`

The noncurrent_version_expiration configuration block supports the following:

* `days` - (Required, Type: int) Specifies the number of days noncurrent object versions expire.

### `lifecycle_rule-noncurrent_version_transition`

The noncurrent_version_transition configuration block supports the following:

* `days` - (Required, Type: int) Specifies the number of days noncurrent object versions transition.
* `storage_class` - (Required) Specifies the storage class that objects that conform to the rule are converted into. The storage class of the objects in a bucket of the IA storage class can be converted into Archive but cannot be converted into Standard. Values: `IA`, `Archive`, `CodeArchive`. ColdArchive is available since 1.203.0.


### `server_side_encryption_rule`

The server_side_encryption_rule configuration block supports the following:

* `sse_algorithm` - (Required) The server-side encryption algorithm to use. Possible values: `AES256` and `KMS`.
* `kms_master_key_id` - (Optional, Available since 1.92.0) The alibaba cloud KMS master key ID used for the SSE-KMS encryption.

### `versioning`

The versioning configuration block supports the following:

* `status` - (Required) Specifies the versioning state of a bucket. Valid values: `Enabled` and `Suspended`.


### `transfer_acceleration`

The transfer_acceleration configuration block supports the following:

* `enabled` - (Required, Type: bool) Specifies the accelerate status of a bucket.

## Attributes Reference

The following attributes are exported:

* `id` - The name of the bucket.
* `creation_date` - The creation date of the bucket.
* `extranet_endpoint` - The extranet access endpoint of the bucket.
* `intranet_endpoint` - The intranet access endpoint of the bucket.
* `location` - The location of the bucket.
* `owner` - The bucket owner.

## Import

OSS bucket can be imported using the bucket name, e.g.

```shell
$ terraform import alicloud_oss_bucket.bucket bucket-12345678
```
