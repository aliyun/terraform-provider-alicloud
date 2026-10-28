---
subcategory: "File Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_lifecycle_policy"
sidebar_current: "docs-alicloud-resource-nas-lifecycle-policy"
description: |-
  Provides a Alicloud File Storage (NAS) Lifecycle Policy resource.
---

# alicloud_nas_lifecycle_policy

Provides a File Storage (NAS) Lifecycle Policy resource.

For information about File Storage (NAS) Lifecycle Policy and how to use it, see [What is Lifecycle Policy](https://www.alibabacloud.com/help/en/doc-detail/169362.html).

-> **NOTE:** Available since v1.153.0.

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
* `paths` - (Optional, ForceNew) The absolute path of the directory for which the lifecycle management policy is configured. Set a maximum of `10` path. The path value must be prefixed by a forward slash (/) and must be an existing path in the mount target. `paths` and `path` are mutually exclusive, and at least one of them must be specified.
* `path` - (Optional, ForceNew) The absolute path of the directory for which the lifecycle management policy is configured. The path value must be prefixed by a forward slash (/) and must be an existing path in the mount target. `path` and `paths` are mutually exclusive, and at least one of them must be specified.
* `storage_type` - (Required, ForceNew) The storage type of the data that is dumped to the IA storage medium. Valid values: `InfrequentAccess`.
* `description` - (Optional) The description of the lifecycle management policy. **NOTE:** `DescribeLifecyclePolicies` does not echo `description` on regular (non-CPFS) NAS file systems, so setting it may show a plan diff on non-CPFS file systems. Its read-back is not covered by acceptance tests; see [Unverified boundary](#unverified-boundary).
* `lifecycle_policy_type` - (Optional, ForceNew) The type of the lifecycle management policy. Valid values: `Auto`, `OnDemand`. Default to `Auto`. **NOTE:** `OnDemand` is available on CPFS智算版 file systems only.
* `retrieve_rules` - (Optional, ForceNew) The retrieve rules of the lifecycle management policy. **NOTE:** This parameter is available on CPFS智算版 file systems only and is not covered by acceptance tests. See [Unverified boundary](#unverified-boundary). See [`retrieve_rules`](#retrieve_rules) below.
* `transit_rules` - (Optional, ForceNew) The transit rules of the lifecycle management policy. **NOTE:** This parameter is available on CPFS智算版 file systems only and is not covered by acceptance tests. See [Unverified boundary](#unverified-boundary). See [`transit_rules`](#transit_rules) below.

### `retrieve_rules`

The `retrieve_rules` block supports the following:

* `attribute` - (Optional, ForceNew) The attribute of the retrieve rule. Valid values: `RetrieveType`.
* `threshold` - (Optional, ForceNew) The threshold of the retrieve rule. Valid values: `AfterVisit`, `All`.

### `transit_rules`

The `transit_rules` block supports the following:

* `attribute` - (Optional, ForceNew) The attribute of the transit rule. Valid values: `Atime`.
* `threshold` - (Optional, ForceNew) The threshold of the transit rule.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Lifecycle Policy. The value formats as `<file_system_id>:<lifecycle_policy_name>`.
* `lifecycle_policy_id` - The ID of the lifecycle management policy.

## Unverified boundary

The `retrieve_rules`, `transit_rules` and `lifecycle_policy_type = OnDemand` parameters are specific to CPFS智算版 file systems. They are validated only by unit tests (`TestUnitAlicloudNASLifecyclePolicy`) and are **not** covered by real-cloud acceptance tests, because the test account does not provide CPFS智算版 resources. Use them only on CPFS智算版 file systems and verify against the actual NAS OpenAPI behavior.

`description` is accepted by `CreateLifecyclePolicy`/`ModifyLifecyclePolicy`, but `DescribeLifecyclePolicies` does not echo it on regular (non-CPFS) NAS file systems in the test account. Its read-back is therefore **not** covered by real-cloud acceptance tests and is exercised only by the unit test; on non-CPFS file systems setting `description` may produce a non-empty plan diff until the API echoes it.

## Import

File Storage (NAS) Lifecycle Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_nas_lifecycle_policy.example <file_system_id>:<lifecycle_policy_name>
```
