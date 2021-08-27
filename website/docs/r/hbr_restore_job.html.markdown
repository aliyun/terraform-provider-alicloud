---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_restore_job"
sidebar_current: "docs-alicloud-resource-hbr-restore-job"
description: |-
  Provides a Alicloud Hybrid Backup Recovery (HBR) Restore Job resource.
---

# alicloud\_hbr\_restore\_job

Provides a Hybrid Backup Recovery (HBR) Restore Job resource.

For information about Hybrid Backup Recovery (HBR) Restore Job and how to use it, see [What is Restore Job](https://help.aliyun.com/document_detail/62361.html).

-> **NOTE:** Available in v1.133.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_hbr_nas_backup_plans" "default" {
	name_regex = "plan-tf-used-dont-delete"
}

data "alicloud_hbr_snapshots" "nas_snapshots" {
    source_type     = "NAS"
    vault_id        =  data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id
    file_system_id  =  data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id
    create_time     =  data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time
}

resource "alicloud_hbr_restore_job" "default" {
    restore_job_id =        "tftestacc112358"
    snapshot_hash =         data.alicloud_hbr_snapshots.nas_snapshots.snapshots.0.snapshot_hash
    vault_id =              data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id
    source_type =          "NAS"
    restore_type =         "NAS"
    snapshot_id =           data.alicloud_hbr_snapshots.nas_snapshots.snapshots.0.snapshot_id
    target_file_system_id = data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id
    target_create_time =    data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time
    target_path =           "/"
    options = <<EOF
    {"includes":[], "excludes":[]}
    EOF
}
```

-> **NOTE:** This resource can only be created, cannot be modified or deleted. Therefore, any modification of the resource attribute will not affect exist resource.

## Argument Reference

The following arguments are supported:

* `restore_job_id` - (Required, ForceNew) Restore Job ID. It's the unique key of this resource, you must specify a unique keyword.
* `vault_id` - (Required, ForceNew) The ID of Vault.
* `source_type` - (Required, ForceNew) The Type of Data Source. Valid values: `ECS_FILE`, `NAS`, `OSS`.
* `restore_type` - (Required, ForceNew) The Recovery Destination Types. Valid values: `ECS_FILE`, `NAS`, `OSS`. **Note**: Currently, there is a one-to-one correspondence between the data source type with the recovery destination type.
* `snapshot_hash` - (Required, ForceNew) Restore Snapshot of HashCode.
* `snapshot_id` - (Required, ForceNew) The ID of Snapshot.
* `options` - (Optional, ForceNew) Recovery Options. It's a json string with format:`"{"includes":[],"excludes":[]}",`.
* `exclude` - (Optional while source_type equals `ECS_FILE`, ForceNew) The exclude path. It's a json string with format:`["/home", "/exclude"]`. 
* `include` - (Optional while source_type equals `ECS_FILE`, ForceNew) The include path. It's a json string with format:`["/home", "/include"]`.
* `target_bucket` - (Required while source_type equals `OSS`, ForceNew) The Target ofo OSS Bucket Name.
* `target_prefix` - (Required while source_type equals `OSS`, ForceNew) The Target of the OSS Object Prefix.
* `target_file_system_id` - (Required while source_type equals `NAS`, ForceNew) The Destination File System ID.
* `target_create_time` - (Required while source_type equals `NAS`, ForceNew) The Destination File System Creation Time.
* `target_path` - (Required while source_type equals `NAS` or `ECS_FILE`, ForceNew) The Target of (ECS) Instance Changes the ECS File Path.
* `target_instance_id` - (Required while source_type equals `ECS_FILE`, ForceNew)  Objective to ECS Instance Id.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Restore Job. The value formats as `<restore_job_id>:<restore_type>`.
* `status` - The Restore Job Status.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Restore Job.

## Import

Hybrid Backup Recovery (HBR) Restore Job can be imported using the id, e.g.

```
$ terraform import alicloud_hbr_restore_job.example <restore_job_id>:<restore_type>
```