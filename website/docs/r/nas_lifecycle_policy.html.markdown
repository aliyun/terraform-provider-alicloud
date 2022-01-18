---
subcategory: "Network Attached Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_lifecycle_policy"
sidebar_current: "docs-alicloud-resource-nas-lifecycle-policy"
description: |-
  Provides a Alicloud Network Attached Storage (NAS) Lifecycle Policy resource.
---

# alicloud\_nas\_lifecycle\_policy

Provides a Network Attached Storage (NAS) Lifecycle Policy resource.

For information about Network Attached Storage (NAS) Lifecycle Policy and how to use it, see [What is Lifecycle Policy](https://www.alibabacloud.com/help/en/doc-detail/169362.html).

-> **NOTE:** Available in v1.153.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_nas_file_system" "example" {
  protocol_type = "NFS"
  storage_type  = "Capacity"
}

resource "alicloud_nas_lifecycle_policy" "example" {
  file_system_id        = alicloud_nas_file_system.example.id
  lifecycle_policy_name = "my-LifecyclePolicy"
  lifecycle_rule_name   = "DEFAULT_ATIME_14"
  storage_type          = "InfrequentAccess"
}
```

## Argument Reference

The following arguments are supported:

* `file_system_id` - (Required, ForceNew) The ID of the file system.
* `lifecycle_policy_name` - (Required, ForceNew) The name of the lifecycle management policy.
* `lifecycle_rule_name` - (Required) The rules in the lifecycle management policy. Valid values: `DEFAULT_ATIME_14`, `DEFAULT_ATIME_30`, `DEFAULT_ATIME_60`, `DEFAULT_ATIME_90`.
* `paths` - (Required, ForceNew) The absolute path of the directory for which the lifecycle management policy is configured. Set a maximum of `10` path. The path value must be prefixed by a forward slash (/) and must be an existing path in the mount target.
* `storage_type` - (Required, ForceNew) The storage type of the data that is dumped to the IA storage medium. Valid values: `InfrequentAccess`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Lifecycle Policy. The value formats as `<file_system_id>:<lifecycle_policy_name>`.

## Import

Network Attached Storage (NAS) Lifecycle Policy can be imported using the id, e.g.

```
$ terraform import alicloud_nas_lifecycle_policy.example <file_system_id>:<lifecycle_policy_name>
```