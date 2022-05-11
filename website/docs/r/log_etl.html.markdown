---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_etl"
sidebar_current: "docs-alicloud-resource-log-etl"
description: |-
  Provides a Alicloud log etl resource.
---

# alicloud\_log\_etl

The data transformation of the log service is a hosted, highly available, and scalable data processing service, 
which is widely applicable to scenarios such as data regularization, enrichment, distribution, aggregation, and index reconstruction.
[Refer to details](https://www.alibabacloud.com/help/zh/doc-detail/125384.htm).

-> **NOTE:** Available in 1.120.0

## Example Usage

Basic Usage

```
resource "alicloud_log_project" "example" {
  name        = "tf-log"
  description = "created by terraform"

}
resource "alicloud_log_store" "example" {
  project               = alicloud_log_project.example.name
  name                  = "tf-test-logstore"
  retention_period      = 3650
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}
resource "alicloud_log_store" "example2" {
  project               = alicloud_log_project.example.name
  name                  = "tf-test-logstore-2"
  retention_period      = 3650
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}
resource "alicloud_log_store" "example3" {
  project               = alicloud_log_project.example.name
  name                  = "tf-test-logstore-3"
  retention_period      = 3650
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}
resource "alicloud_log_etl" "example" {
  etl_name          = "etl_name"
  project           = alicloud_log_project.example.name
  display_name      = "display_name"
  description       = "etl_description"
  access_key_id     = "access_key_id"
  access_key_secret = "access_key_secret"
  script            = "e_set('new','key')"
  logstore          = alicloud_log_store.example.name
  etl_sinks {
    name              = "target_name"
    access_key_id     = "example2_access_key_id"
    access_key_secret = "example2_access_key_secret"
    endpoint          = "cn-hangzhou.log.aliyuncs.com"
    project           = alicloud_log_project.example.name
    logstore          = alicloud_log_store.example2.name
  }
  etl_sinks {
    name              = "target_name2"
    access_key_id     = "example3_access_key_id"
    access_key_secret = "example3_access_key_secret"
    endpoint          = "cn-hangzhou.log.aliyuncs.com"
    project           = alicloud_log_project.example.name
    logstore          = alicloud_log_store.example3.name
  }
}
```
Stop the task in progress
```
resource "alicloud_log_etl" "example" {
  status            = STOPPED
  etl_name          = "etl_name"
  project           = alicloud_log_project.example.name
  display_name      = "display_name"
  description       = "etl_description"
  access_key_id     = "access_key_id"
  access_key_secret = "access_key_secret"
  script            = "e_set('new','key')"
  logstore          = alicloud_log_store.example.name
  etl_sinks {
    name              = "target_name"
    access_key_id     = "example2_access_key_id"
    access_key_secret = "example2_access_key_secret"
    endpoint          = "cn-hangzhou.log.aliyuncs.com"
    project           = alicloud_log_project.example.name
    logstore          = alicloud_log_store.example2.name
  }
  etl_sinks {
    name              = "target_name2"
    access_key_id     = "example3_access_key_id"
    access_key_secret = "example3_access_key_secret"
    endpoint          = "cn-hangzhou.log.aliyuncs.com"
    project           = alicloud_log_project.example.name
    logstore          = alicloud_log_store.example3.name
  }
}
```
ReStart the stopped task
```
resource "alicloud_log_etl" "example" {
  status            = RUNNING
  etl_name          = "etl_name"
  project           = alicloud_log_project.example.name
  display_name      = "display_name"
  description       = "etl_description"
  access_key_id     = "access_key_id"
  access_key_secret = "access_key_secret"
  script            = "e_set('new','key')"
  logstore          = alicloud_log_store.example.name
  etl_sinks {
    name              = "target_name"
    access_key_id     = "example2_access_key_id"
    access_key_secret = "example2_access_key_secret"
    endpoint          = "cn-hangzhou.log.aliyuncs.com"
    project           = alicloud_log_project.example.name
    logstore          = alicloud_log_store.example2.name
  }
  etl_sinks {
    name              = "target_name2"
    access_key_id     = "example3_access_key_id"
    access_key_secret = "example3_access_key_secret"
    endpoint          = "cn-hangzhou.log.aliyuncs.com"
    project           = alicloud_log_project.example.name
    logstore          = alicloud_log_store.example3.name
  }
}
```

## Argument Reference

The following arguments are supported:

* `etl_name` - (Required, ForceNew) The name of the log etl job.
* `description` - (Optional) Description of the log etl job.
* `project` - (Required, ForceNew) The name of the project where the etl job is located.
* `display_name` - (Required) Log service etl job alias.
* `schedule` - (Optional) Job scheduling type, the default value is Resident.
* `etl_type` - (Optional) Log service etl type, the default value is `ETL`.
* `status` - (Optional) Log project tags. the default value is RUNNING, Only 4 values are supported: `STARTING`，`RUNNING`，`STOPPING`，`STOPPED`.
* `create_time` - (Optional) The etl job create time.
* `last_modified_time` - (Optional) ETL job last modified time.
* `access_key_id` - (Optional,Sensitive) Source logstore access key id.
* `kms_encrypted_access_key_id` - (Optional) An KMS encrypts access key id used to a log etl job. If the `access_key_id` is filled in, this field will be ignored.
* `kms_encryption_access_key_id_context` - (Optional) An KMS encryption context used to decrypt `kms_encrypted_access_key_id` before creating or updating an instance with `kms_encrypted_access_key_id`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set. When it is changed, the instance will reboot to make the change take effect.
* `access_key_secret` - (Optional,Sensitive) Source logstore access key secret.
* `kms_encrypted_access_key_secret` - (Optional) An KMS encrypts access key secret used to a log etl job. If the `access_key_secret` is filled in, this field will be ignored.
* `kms_encryption_access_key_secret_context` - (Optional) An KMS encryption context used to decrypt `kms_encrypted_access_key_secret` before creating or updating an instance with `kms_encrypted_access_key_secret`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set. When it is changed, the instance will reboot to make the change take effect.

* `from_time` - (Optional) The start time of the processing job, if not set the value is 0, indicates to start processing from the oldest data.
* `to_time` - (Optional) Deadline of processing job, if not set the value is 0, indicates that new data will be processed continuously.
* `script` - (Required) Processing operation grammar.
* `version` - (Optional) Log etl job version. the default value is `2`.
* `logstore` - (Required) The source logstore of the processing job.
* `parameters` - (Optional) Advanced parameter configuration of processing operations.
* `role_arn` - (Optional) Sts role info under source logstore. `role_arn` and `(access_key_id, access_key_secret)` fill in at most one. If you do not fill in both, then you must fill in `(kms_encrypted_access_key_id, kms_encrypted_access_key_secret, kms_encryption_access_key_id_context, kms_encryption_access_key_secret_context)` to use KMS to get the key pair.
* `etl_sinks` - (Required) Target logstore configuration for delivery after data processing.
    * `access_key_id` - (Optional,Sensitive) Delivery target logstore access key id.
    * `kms_encrypted_access_key_id` - (Optional) An KMS encrypts access key id used to a log etl job. If the `access_key_id` is filled in, this field will be ignored.
    * `access_key_secret`- (Optional,Sensitive) Delivery target logstore access key secret.
    * `kms_encrypted_access_key_secret` - (Optional) An KMS encrypts access key secret used to a log etl job. If the `access_key_secret` is filled in, this field will be ignored.
    * `endpoint` - (Required) Delivery target logstore region.
    * `name` - (Required) Delivery target name.
    * `project` - (Required) The project where the target logstore is delivered.
    * `logstore` - (Required) Delivery target logstore.
    * `role_arn` - (Optional) Sts role info under delivery target logstore. `role_arn` and `(access_key_id, access_key_secret)` fill in at most one. If you do not fill in both, then you must fill in `(kms_encrypted_access_key_id, kms_encrypted_access_key_secret, kms_encryption_access_key_id_context, kms_encryption_access_key_secret_context)` to use KMS to get the key pair.
    * `type` - (Optional)  ETL sinks type, the default value is AliyunLOG.
    
-> **Note:** `from_time` and `to_time` no modification allowed after successful creation.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the log etl. It formats of `<project>:<etl_name>`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when Creating LogEtl instance. 
* `delete` - (Defaults to 3 mins) Used when terminating the LogEtl instance. 
* `update` - (Defaults to 5 mins) Used when Updating LogEtl instance. 


## Import

Log etl can be imported using the id, e.g.

```
$ terraform import alicloud_log_etl.example tf-log-project:tf-log-etl-name
```
