---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_replication"
sidebar_current: "docs-alicloud-resource-oss-bucket-replication"
description: |-
  Provides a resource to create an oss bucket replication rule.
---

# alicloud\_oss\_bucket

Provides a resource to create a oss bucket and set its attribution.

-> **NOTE:** The bucket namespace is shared by all users of the OSS system. Please set bucket name as unique as possible.


## Example Usage

Create bucket

```
resource "alicloud_oss_bucket_replication" "bucket" {
    bucket = "bucket-173448-replication"
```

Set bucket replication rule

```
resource "alicloud_oss_bucket_replication" "bucket-replication" {
    bucket = "bucket-170309-replication"

    replication_rule {
        prefix_set {
            prefixes = ["xx/", "test/"]
        }
        destination {
            bucket = "guox-test-dst"
            location = "oss-cn-beijing"
            transfer_type = "internal"
        }
        action = "ALL"
        historical_object_replication = "disabled"
    }
}

resource "alicloud_oss_bucket_replication" "bucket-replication" {
    bucket = "bucket-170309-replication"

    replication_rule {
        destination {
            bucket = "guoxing-test-dst"
            location = "oss-cn-beijing"
            transfer_type = "oss_acc"
        }
        action = "PUT"
        historical_object_replication = "enabled"
        sync_role = "User"
        source_selection_criteria {
            sse_kms_encrypted_objects {
                status = "Enabled"
            }
        }
        encryption_configuration {
            replica_kms_key_id = "1a8d780d-0d34-49e5-ab45-7b9b*******4"
        }
    }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Optional, ForceNew) The name of the bucket. If omitted, Terraform will assign a random and unique name.
* `replication_rule` - (Optional) A configuration of [replication management](https://www.alibabacloud.com/help/zh/doc-detail/31905.html) (documented below).

#### Block replication_rule

The replication_rule object supports the following:

* `prefixSet` - (Optional, Type: set) The set of prefix. If value is nil, the rule applies to all objects in a bucket.
* `action` - (Optional, Type: string) Which operation will be replication. Value: `ALL`, `DELETE`, `PUT` or any combination of them.
* `destination` - (Optional, Type: set) Which operation will be replication.
* `historical_object_replication` - (Optional, Type: string) Whether the objects already exists would be replicated. Value: `enabled` or `disabled`.
* `sync_role` - (Optional, Type: string) The operator of replication.
* `source_selection_criteria` - (Optional, Type: set) Another filter of objects. Only support objects encrypted as "SSE-KMS".
* `encryption_configuration` - (Optional, Type: set) This params based on "source_selection_criteria", the objects which was encrypted would be replicated by a new kms-id.

`NOTE`: If "source_selection_criteria.sse_kms_encrypted_objects.status" is Enabled, then "sync_role", "encryption_configuration" and "encryption_configuration.replica_kms_key_id" should be defined.

#### Block prefix_set

The replication_rule prefix_set object supports the following:

* `prefix` - (Optional, Type: string) Object key prefix identifying one or more objects to which the rule applies. 

#### Block destination

The replication_rule destination object supports the following:

* `bucket` - (Optional, Type: string) The destination bucket of replication.
* `loaction` - (Optional, Type: string) The location of destination bucket.
* `transfer_type` - (Optional, Type: string) The transfer type of replication. Value: `internal`, `oss_acc`,

#### Block source_selection_criteria

The replication_rule source_selection_criteria object supports the following:

* `SseKmsEncryptedObjects` - (Optional, Type: set) The objects encrypted as "SSE-KMS".

#### Block sse_kms_encrypted_objects

The replication_rule sse_kms_encrypted_objects supports the following:

* `Status` - (Optional, Type: string) Whether the objects would be replicated. Value:`Enabled` \and `Disabled`.

#### Block encryption_configuration

The replication_rule encryption_configuration object supports the following:

* `replica_kms_key_id` - (Optional, Type: string) The new kms-id for the objects.

## Attributes Reference

The following attributes are exported:

* `id` - The name of the bucket.
* `acl` - The acl of the bucket.
* `creation_date` - The creation date of the bucket.
* `extranet_endpoint` - The extranet access endpoint of the bucket.
* `intranet_endpoint` - The intranet access endpoint of the bucket.
* `location` - The location of the bucket.
* `owner` - The bucket owner.
* `cross_region_replication` - Displays the cross-region replication status of the Bucket.
* `transfer_acceleration` - Display Bucket transfer acceleration status.
