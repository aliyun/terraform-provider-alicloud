---
subcategory: "File Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_lifecycle_policy"
sidebar_current: "docs-alicloud-resource-nas-lifecycle-policy"
description: |-
  Provides a Alicloud File Storage (NAS) Lifecycle Policy resource.
---

# alicloud\_nas\_lifecycle\_policy

Provides a File Storage (NAS) Lifecycle Policy resource.

For information about File Storage (NAS) Lifecycle Policy and how to use it, see [What is Lifecycle Policy](https://www.alibabacloud.com/help/en/doc-detail/169362.html).

-> **NOTE:** Available in v1.153.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nas_lifecycle_policy&exampleId=53f1763f-e56f-8458-40ec-937ff0fe4064c4c33cc1&activeTab=example&spm=docs.r.nas_lifecycle_policy.0.53f1763fe5&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_nas_file_system" "example" {
  protocol_type = "NFS"
  storage_type  = "Capacity"
}

resource "alicloud_nas_lifecycle_policy" "example" {
  file_system_id        = alicloud_nas_file_system.example.id
  lifecycle_policy_name = "terraform-example"
  lifecycle_rule_name   = "DEFAULT_ATIME_14"
  storage_type          = "InfrequentAccess"
  paths                 = ["/"]
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_nas_lifecycle_policy&spm=docs.r.nas_lifecycle_policy.example&intl_lang=EN_US)

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

File Storage (NAS) Lifecycle Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_nas_lifecycle_policy.example <file_system_id>:<lifecycle_policy_name>
```