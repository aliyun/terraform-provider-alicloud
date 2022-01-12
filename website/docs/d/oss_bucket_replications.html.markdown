---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_replications"
sidebar_current: "docs-alicloud-datasource-oss-bucket-replications"
description: |-
    Provides a list of OSS bucket replication rule to the user.
---

# alicloud\_oss_buckets

This data source provides the OSS buckets of the current Alibaba Cloud user.

## Example Usage

```
data "alicloud_oss_bucket_replications" "oss_buckets_ds" {
  name_regex = "sample_oss_bucket"
}

output "first_oss_bucket_replication" {
  value = "${data.alicloud_oss_bucket_replications.oss_buckets_ds.buckets.0.replication_rule}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by bucket name.
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
  * `storage_class` - Object storage type. Possible values: `Standard`, `IA` and `Archive`.
  * `redundancy_type` - Redundancy type. Possible values: `LRS`, and `ZRS`.
  * `creation_date` - Bucket creation date.
  * `cross_region_replication` - Bucket replication type. Possible values: `Enabled`, and `Disabled`.
  * `transfer_acceleration` - Bucket replication accelerator. Possible values: `Enabled`, and `Disabled`.
  * `replictation_rule` - A configuration of replication for a bucket. It contains the following attributes:
    * `id` - (Computed, Type: string) Unique identifier for the rule. OSS bucket will assign a unique name.
    * `status` - (Computed, Type: string) Specifies replication rule status. Possible values:`starting`, `doing` and `closing`.
    * `action` - Defined which operation of objects in bucket will be sync. Possible values:`ALL`, `DELETE`, `PUT` or any combination of them.
    * `prefix_set` - A container for holding prefixes.
      * `prefixes` - The prefixes for the object to be copied.
    * `destination` - Defined the target bucket of replication
      * `bucket` - Target bucket
      * `location` - The location of target bucket
      * `transfer_type` - The transfer type of replication. Possible values:`internal` and `oss_acc`.
    * `historical_object_replication` - Defined whether the objects already exits will be sync. Possible values:`enabled` and `disabled`.
    * `sync_role` - Authorizes OSS to use a role for data replication.
