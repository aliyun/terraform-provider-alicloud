---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_replication"
sidebar_current: "docs-alicloud-resource-oss-bucket-replication"
description: |-
  Provides a OSS bucket replication configuration resource.
---

# alicloud\_oss\_bucket\_replication

Provides an independent replication configuration resource for OSS bucket.

For information about OSS replication and how to use it, see [What is cross-region replication](https://www.alibabacloud.com/help/doc-detail/31864.html) and [What is same-region replication](https://www.alibabacloud.com/help/doc-detail/254865.html).

-> **NOTE:** Available in v1.161.0+.

## Example Usage

Set bucket replication configuration

```
resource "alicloud_oss_bucket_replication" "cross-region-replication" {
    bucket = "bucket-in-hangzhou"

    destination {
        bucket = "bucket-in-beijing"
        location = "oss-cn-beijing"
    }
    action = "ALL"
}

resource "alicloud_oss_bucket_replication" "same-region-replication" {
    bucket = "bucket-in-hangzhou"

    destination {
        bucket = "bucket-in-hangzhou-1"
        location = "oss-cn-hangzhou"
    }
    action = "ALL"
}

resource "alicloud_oss_bucket_replication" "replication-with-prefix" {
    bucket = "bucket-1"

    prefix_set {
        prefixes = ["prefix1/", "prefix2/"]
    }

    destination {
        bucket = "bucket-2"
        location = "oss-cn-hangzhou"
    }
    action = "ALL"
    historical_object_replication = "disabled"
}

resource "alicloud_oss_bucket_replication" "replication-with-specific-action" {
    bucket = "bucket-1"

    destination {
        bucket = "bucket-2"
        location = "oss-cn-hangzhou"
    }
    action = "PUT"
    historical_object_replication = "disabled"
}

resource "alicloud_oss_bucket_replication" "replication-with-kms-encryption" {
    bucket = "bucket-1"

    destination {
        bucket = "bucket-2"
        location = "oss-cn-hangzhou"
    }

    action = "ALL"
    
    historical_object_replication = "disabled"
    
    sync_role = "<your ram role>"
    
    source_selection_criteria {
        sse_kms_encrypted_objects {
            status = "Enabled"
        }
    }
    
    encryption_configuration {
        replica_kms_key_id = "<your kms key id>"
    }
}

```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, ForceNew) The name of the bucket.
* `prefix_set` - (Optional, ForceNew) The prefixes used to specify the object to replicate. Only objects that match the prefix are replicated to the destination bucket(See the following block `prefix_set`).
* `destination` - (Required, ForceNew) Specifies the destination for the rule(See the following block `destination`).
* `action` - (Optional, ForceNew) The operations that can be synchronized to the destination bucket. You can set action to one or more of the following operation types. Valid values: `ALL`(contains PUT, DELETE, and ABORT), `PUT`, `DELETE` and `ABORT`. Defaults to `ALL`.    
* `historical_object_replication` - (Optional, ForceNew) Specifies whether to replicate historical data from the source bucket to the destination bucket before data replication is enabled. Can be `enabled` or `disabled`. Defaults to `enabled`.
* `sync_role` - (Optional, ForceNew) Specifies the role that you authorize OSS to use to replicate data. If SSE-KMS is specified to encrypt the objects replicated to the destination bucket, it must be specified.
* `source_selection_criteria` - (Optional, ForceNew) Specifies other conditions used to filter the source objects to replicate(See the following block `source_selection_criteria`).
* `encryption_configuration` - (Optional, ForceNew) Specifies the encryption configuration for the objects replicated to the destination bucket(See the following block `encryption_configuration`).


#### Block prefix_set

The `prefix_set` configuration block supports the following:

* `prefixes` - (Required, ForceNew) The list of object key name prefix identifying one or more objects to which the rule applies.

`NOTE`: The prefix must be less than or equal to 1024 characters in length.

#### Block destination

The `destination` configuration block supports the following:

* `bucket` - (Required, ForceNew) The destination bucket to which the data is replicated.
* `loaction` - (Required, ForceNew) The region in which the destination bucket is located.
* `transfer_type` - (Optional, ForceNew) The link used to transfer data in data replication.. Can be `internal` or `oss_acc`. Defaults to `internal`.

`NOTE`: You can set transfer_type to oss_acc only when you create cross-region replication (CRR) rules.

#### Block source_selection_criteria

The `source_selection_criteria` configuration block supports the following:

* `sse_kms_encrypted_objects` - (Required, ForceNew) Filter source objects encrypted by using SSE-KMS(See the following block `sse_kms_encrypted_objects`).

#### Block sse_kms_encrypted_objects

The `sse_kms_encrypted_objects` configuration block supports the following:

* `status` - (Required, ForceNew) Specifies whether to replicate objects encrypted by using SSE-KMS. Can be `Enabled` or `Disabled`.

#### Block encryption_configuration

The `encryption_configuration` configuration block supports the following:

* `replica_kms_key_id` - (Required, ForceNew) The CMK ID used in SSE-KMS.

`NOTE`: If the status of sse_kms_encrypted_objects is set to Enabled, you must specify the replica_kms_key_id.

## Attributes Reference

The following attributes are exported:

* `id` - The current replication configuration resource ID. Composed of bucket name and rule_id with format `<bucket>:<rule_id>`.
* `rule_id` - The ID of the data replication rule.
* `status` - The status of the data replication task. Can be starting, doing and closing.
* `progress` - Retrieves the progress of the data replication task. This status is returned only when the data replication task is in the doing state.
    * `historical_object` - The percentage of the replicated historical data. This element is valid only when historical_object_replication is set to enabled.
    * `new_object` - The time used to distinguish new data from historical data. Data that is written to the source bucket before the time is replicated to the destination bucket as new data. The value of this element is in GMT.

## Import

Oss Bucket Replication can be imported using the id, e.g.

```
$ terraform import alicloud_oss_bucket_replication.example
```

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `delete` - (Defaults to 30 mins) Used when delete a data replication rule (until the data replication task is cleared). 

