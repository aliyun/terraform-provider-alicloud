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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket&exampleId=b1523857-e8b4-b770-488f-623b9902e4060adc3742&activeTab=example&spm=docs.r.oss_bucket.0.b1523857e8&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "bucket-acl" {
  bucket = "example-value-${random_integer.default.result}"
}

resource "alicloud_oss_bucket_acl" "bucket-acl" {
  bucket = alicloud_oss_bucket.bucket-acl.bucket
  acl    = "private"
}
```

Static Website

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket&exampleId=5092ff58-0422-ce90-478a-1b7d6a4cf2e8545e3e5c&activeTab=example&spm=docs.r.oss_bucket.1.5092ff5804&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket&exampleId=c9a315e1-ba24-dc11-fd7e-f3de5840bceb09920aa6&activeTab=example&spm=docs.r.oss_bucket.2.c9a315e1ba&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "bucket-target" {
  bucket = "example-value-${random_integer.default.result}"
}

resource "alicloud_oss_bucket_acl" "bucket-target" {
  bucket = alicloud_oss_bucket.bucket-target.bucket
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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket&exampleId=159941e6-def7-9c73-64b7-6c9218903396c5a86caf&activeTab=example&spm=docs.r.oss_bucket.3.159941e6de&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "bucket-referer" {
  bucket = "example-value-${random_integer.default.result}"
  referer_config {
    allow_empty = false
    referers    = ["http://www.aliyun.com", "https://www.aliyun.com"]
  }
}

resource "alicloud_oss_bucket_acl" "default" {
  bucket = alicloud_oss_bucket.bucket-referer.bucket
  acl    = "private"
}
```

Set lifecycle rule

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket&exampleId=d537d136-1a72-6fbd-5c50-a2367e6ef480db3a6172&activeTab=example&spm=docs.r.oss_bucket.4.d537d1361a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "bucket-lifecycle1" {
  bucket = "example-lifecycle1-${random_integer.default.result}"

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

resource "alicloud_oss_bucket_acl" "bucket-lifecycle1" {
  bucket = alicloud_oss_bucket.bucket-lifecycle1.bucket
  acl    = "public-read"
}


resource "alicloud_oss_bucket" "bucket-lifecycle2" {
  bucket = "example-lifecycle2-${random_integer.default.result}"

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

resource "alicloud_oss_bucket_acl" "bucket-lifecycle2" {
  bucket = alicloud_oss_bucket.bucket-lifecycle2.bucket
  acl    = "public-read"
}


resource "alicloud_oss_bucket" "bucket-lifecycle3" {
  bucket = "example-lifecycle3-${random_integer.default.result}"

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

resource "alicloud_oss_bucket_acl" "bucket-lifecycle3" {
  bucket = alicloud_oss_bucket.bucket-lifecycle3.bucket
  acl    = "public-read"
}


resource "alicloud_oss_bucket" "bucket-lifecycle4" {
  bucket = "example-lifecycle4-${random_integer.default.result}"

  lifecycle_rule {
    id      = "rule-abort-multipart-upload"
    prefix  = "path3/"
    enabled = true

    abort_multipart_upload {
      days = 128
    }
  }
}

resource "alicloud_oss_bucket_acl" "bucket-lifecycle4" {
  bucket = alicloud_oss_bucket.bucket-lifecycle4.bucket
  acl    = "public-read"
}


resource "alicloud_oss_bucket" "bucket-versioning-lifecycle" {
  bucket = "example-lifecycle5-${random_integer.default.result}"

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

resource "alicloud_oss_bucket_acl" "bucket-versioning-lifecycle" {
  bucket = alicloud_oss_bucket.bucket-versioning-lifecycle.bucket
  acl    = "private"
}


resource "alicloud_oss_bucket" "bucket-access-monitor-lifecycle" {
  bucket = format("example-lifecycle6-%s", random_integer.default.result)

  access_monitor {
    status = "Enabled"
  }

  lifecycle_rule {
    id      = "rule-days-transition"
    prefix  = "path/"
    enabled = true

    transitions {
      days                     = 30
      storage_class            = "IA"
      is_access_time           = true
      return_to_std_when_visit = true
    }
  }
}

resource "alicloud_oss_bucket_acl" "bucket-access-monitor-lifecycle" {
  bucket = alicloud_oss_bucket.bucket-access-monitor-lifecycle.bucket
  acl    = "private"
}


resource "alicloud_oss_bucket" "bucket-tag-lifecycle" {
  bucket = format("example-lifecycle7-%s", random_integer.default.result)

  lifecycle_rule {
    id      = "rule-days-transition"
    prefix  = "path/"
    enabled = true
    transitions {
      created_before_date = "2022-11-11"
      storage_class       = "IA"
    }
  }

  tags = {
    Created = "TF",
    For     = "example",
  }
}

resource "alicloud_oss_bucket_acl" "bucket-tag-lifecycle" {
  bucket = alicloud_oss_bucket.bucket-tag-lifecycle.bucket
  acl    = "private"
}
```

Set bucket policy 

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket&exampleId=ac859beb-d9a3-1c87-9e63-a0b12dcc36541b9feae5&activeTab=example&spm=docs.r.oss_bucket.5.ac859bebd9&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "bucket-policy" {
  bucket = "example-policy-${random_integer.default.result}"

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

resource "alicloud_oss_bucket_acl" "default" {
  bucket = alicloud_oss_bucket.bucket-policy.bucket
  acl    = "private"
}
```

IA Bucket

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket&exampleId=d2bad287-91f7-8d9d-47ca-a2212bda0c9a38d522fc&activeTab=example&spm=docs.r.oss_bucket.6.d2bad28791&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket&exampleId=160d1480-f315-e1b1-53a4-f1f0ee5f48cbddb6294c&activeTab=example&spm=docs.r.oss_bucket.7.160d1480f3&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "bucket-sserule" {
  bucket = "terraform-example-${random_integer.default.result}"

  server_side_encryption_rule {
    sse_algorithm = "AES256"
  }
}

resource "alicloud_oss_bucket_acl" "bucket-sserule" {
  bucket = alicloud_oss_bucket.bucket-sserule.bucket
  acl    = "private"
}

resource "alicloud_kms_key" "kms" {
  description            = "terraform-example"
  pending_window_in_days = "7"
  status                 = "Enabled"
}

resource "alicloud_oss_bucket" "bucket-kms" {
  bucket = "terraform-example-kms-${random_integer.default.result}"

  server_side_encryption_rule {
    sse_algorithm     = "KMS"
    kms_master_key_id = alicloud_kms_key.kms.id
  }
}

resource "alicloud_oss_bucket_acl" "bucket-kms" {
  bucket = alicloud_oss_bucket.bucket-kms.bucket
  acl    = "private"
}
```

Set bucket tags 

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket&exampleId=2a2df177-e6f4-d509-595b-e4f9e6f2a96718df758b&activeTab=example&spm=docs.r.oss_bucket.8.2a2df177e6&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "bucket-tags" {
  bucket = "terraform-example-${random_integer.default.result}"

  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}

resource "alicloud_oss_bucket_acl" "bucket-tags" {
  bucket = alicloud_oss_bucket.bucket-tags.bucket
  acl    = "private"
}
```

Enable bucket versioning 

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket&exampleId=b8d407a9-6b97-1516-1050-161aaadf990c2ef94de1&activeTab=example&spm=docs.r.oss_bucket.9.b8d407a96b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "bucket-versioning" {
  bucket = "terraform-example-${random_integer.default.result}"
  versioning {
    status = "Enabled"
  }
}

resource "alicloud_oss_bucket_acl" "default" {
  bucket = alicloud_oss_bucket.bucket-versioning.bucket
  acl    = "private"
}
```

Set bucket redundancy type

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket&exampleId=5092ff58-0422-ce90-478a-1b7d6a4cf2e8545e3e5c&activeTab=example&spm=docs.r.oss_bucket.1.5092ff5804&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket&exampleId=5092ff58-0422-ce90-478a-1b7d6a4cf2e8545e3e5c&activeTab=example&spm=docs.r.oss_bucket.1.5092ff5804&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

Set bucket resource group id

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket&exampleId=5092ff58-0422-ce90-478a-1b7d6a4cf2e8545e3e5c&activeTab=example&spm=docs.r.oss_bucket.1.5092ff5804&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}

resource "alicloud_oss_bucket" "bucket-accelerate" {
  bucket            = "terraform-example-${random_integer.default.result}"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Optional, ForceNew) The name of the bucket. If omitted, Terraform will assign a random and unique name.
* `acl` - (Optional, Computed, Deprecated since 1.220.0) The [canned ACL](https://www.alibabacloud.com/help/doc-detail/31898.htm) to apply. Can be "private", "public-read" and "public-read-write". This property has been deprecated since 1.220.0, please use the resource `alicloud_oss_bucket_acl` instead.
* `cors_rule` - (Optional) A rule of  [Cross-Origin Resource Sharing](https://www.alibabacloud.com/help/doc-detail/31903.htm). The items of core rule are no more than 10 for every OSS bucket. See [`cors_rule`](#cors_rule) below.
* `website` - (Optional) A website configuration. See [`website`](#website) below.
* `logging` - (Optional) A Settings of [bucket logging](https://www.alibabacloud.com/help/doc-detail/31900.htm). See [`logging`](#logging) below.
* `logging_isenable` - (Optional, Deprecated from 1.37.0.) The flag of using logging enable container. Defaults true.
* `referer_config` - (Optional, Deprecated since 1.220.0) The configuration of [referer](https://www.alibabacloud.com/help/doc-detail/31901.htm). This property has been deprecated since 1.220.0, please use the resource `alicloud_oss_bucket_referer` instead. See [`referer_config`](#referer_config) below.
* `lifecycle_rule` - (Optional) A configuration of [object lifecycle management](https://www.alibabacloud.com/help/doc-detail/31904.htm). See [`lifecycle_rule`](#lifecycle_rule) below.
* `policy` - (Optional, Available since 1.41.0, Deprecated since 1.220.0) Json format text of bucket policy [bucket policy management](https://www.alibabacloud.com/help/doc-detail/100680.htm). This property has been deprecated since 1.220.0, please use the resource `alicloud_oss_bucket_policy` instead.
* `storage_class` - (Optional, ForceNew) The [storage class](https://www.alibabacloud.com/help/doc-detail/51374.htm) to apply. Can be "Standard", "IA", "Archive", "ColdArchive" and "DeepColdArchive". Defaults to "Standard". "ColdArchive" is available since 1.203.0. "DeepColdArchive" is available since 1.209.0.
* `redundancy_type` - (Optional, ForceNew, Available since 1.91.0) The [redundancy type](https://www.alibabacloud.com/help/doc-detail/90589.htm) to enable. Can be "LRS", and "ZRS". Defaults to "LRS".
* `server_side_encryption_rule` - (Optional, Available since 1.45.0) A configuration of server-side encryption. See [`server_side_encryption_rule`](#server_side_encryption_rule) below.
* `tags` - (Optional, Available since 1.45.0) A mapping of tags to assign to the bucket. The items are no more than 10 for a bucket.
* `versioning` - (Optional, Available since 1.45.0) A state of versioning. See [`versioning`](#versioning) below.
* `force_destroy` - (Optional, Available since 1.45.0) A boolean that indicates all objects should be deleted from the bucket so that the bucket can be destroyed without error. These objects are not recoverable. Defaults to "false".
* `transfer_acceleration` - (Optional, Available since 1.123.1) A transfer acceleration status of a bucket. See [`transfer_acceleration`](#transfer_acceleration) below.
* `lifecycle_rule_allow_same_action_overlap` - (Optional, Available since 1.208.1) A boolean that indicates lifecycle rules allow prefix overlap.
* `access_monitor` - (Optional, Available since 1.208.1) A access monitor status of a bucket. See [`access_monitor`](#access_monitor) below.
* `resource_group_id` - (Optional, Available since 1.219.0) The ID of the resource group to which the bucket belongs.


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
* `tags` - (Optional, Available since 1.209.0) Key-value map of resource tags. All of these tags must exist in the object's tag set in order for the rule to apply.
* `filter` - (Optional, Available since 1.209.1) Configuration block used to identify objects that a Lifecycle rule applies to. See [`filter`](#lifecycle_rule-filter) below.

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
* `storage_class` - (Required) Specifies the storage class that objects that conform to the rule are converted into. The storage class of the objects in a bucket of the IA storage class can be converted into Archive but cannot be converted into Standard. Values: `IA`, `Archive`, `ColdArchive`, `DeepColdArchive`. ColdArchive is available since 1.203.0. DeepColdArchive is available since 1.209.0.
* `is_access_time` - (Optional, Type: bool, Available since 1.208.1) Specifies whether the lifecycle rule applies to objects based on their last access time. If set to `true`, the rule applies to objects based on their last access time; if set to `false`, the rule applies to objects based on their last modified time. If configure the rule based on the last access time, please enable `access_monitor` first.
* `return_to_std_when_visit` - (Optional, Type: bool, Available since 1.208.1) Specifies whether to convert the storage class of non-Standard objects back to Standard after the objects are accessed. It takes effect only when the IsAccessTime parameter is set to true. If set to `true`, converts the storage class of the objects to Standard; if set to `false`, does not convert the storage class of the objects to Standard.
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
* `storage_class` - (Required) Specifies the storage class that objects that conform to the rule are converted into. The storage class of the objects in a bucket of the IA storage class can be converted into Archive but cannot be converted into Standard. Values: `IA`, `Archive`, `CodeArchive`, `DeepColdArchive`. ColdArchive is available since 1.203.0. DeepColdArchive is available since 1.209.0.
* `is_access_time` - (Optional, Type: bool, Available since 1.208.1) Specifies whether the lifecycle rule applies to objects based on their last access time. If set to `true`, the rule applies to objects based on their last access time; if set to `false`, the rule applies to objects based on their last modified time. If configure the rule based on the last access time, please enable `access_monitor` first.
* `return_to_std_when_visit` - (Optional, Type: bool, Available since 1.208.1) Specifies whether to convert the storage class of non-Standard objects back to Standard after the objects are accessed. It takes effect only when the IsAccessTime parameter is set to true. If set to `true`, converts the storage class of the objects to Standard; if set to `false`, does not convert the storage class of the objects to Standard.

### `lifecycle_rule-filter`

The filter configuration block supports the following:

* `not`- (Optional) The condition that is matched by objects to which the lifecycle rule does not apply. See [`not`](#lifecycle_rule-filter-not) below.
* `object_size_greater_than` - (Optional) Minimum object size (in bytes) to which the rule applies.
* `object_size_less_than` - (Optional) Maximum object size (in bytes) to which the rule applies.

### `lifecycle_rule-filter-not`

The not configuration block supports the following:

* `prefix` - (Optional) The prefix in the names of the objects to which the lifecycle rule does not apply.
* `tag` - (Optional) The tag of the objects to which the lifecycle rule does not apply. See [`tag`](#lifecycle_rule-filter-not-tag) below.

### `lifecycle_rule-filter-not-tag`

The tag configuration block supports the following:

* `key` - (Required) The key of the tag that is specified for the objects.
* `value` - (Required) The value of the tag that is specified for the objects.

### `server_side_encryption_rule`

The server_side_encryption_rule configuration block supports the following:

* `sse_algorithm` - (Required) The server-side encryption algorithm to use. Possible values: `AES256` and `KMS`.
* `kms_master_key_id` - (Optional, Available since 1.92.0) The alibaba cloud KMS master key ID used for the SSE-KMS encryption.
* `kms_data_encryption` - (Optional, Available since 1.246.0) The algorithm used to encrypt objects. If this element is not specified, objects are encrypted with AES256. This element is valid only when the value of SSEAlgorithm is set to KMS. Valid values: `SM4`.

### `versioning`

The versioning configuration block supports the following:

* `status` - (Required) Specifies the versioning state of a bucket. Valid values: `Enabled` and `Suspended`.


### `transfer_acceleration`

The transfer_acceleration configuration block supports the following:

* `enabled` - (Required, Type: bool) Specifies the accelerate status of a bucket.

### `access_monitor`

The access_monitor configuration block supports the following:

* `status` - (Optional) The access monitor state of a bucket. If you want to manage objects based on the last access time of the objects, specifies the status to `Enabled`. Valid values: `Enabled` and `Disabled`.

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
