---
subcategory: "Actiontrail"
layout: "alicloud"
page_title: "Alicloud: alicloud_actiontrails"
sidebar_current: "docs-alicloud-datasource-actiontrails"
description: |-
  Provides a list of ActionTrail Trails to the specified filters.
---

# alicloud\_actiontrails

-> **DEPRECATED:**  This datasource has been renamed to [alicloud_actiontrail_trails](https://www.terraform.io/docs/providers/alicloud/d/actiontrail_trails) from version 1.95.0.

This data source provides a list of ActionTrail Trails in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available since v1.95.0.

## Example Usage

```terraform
data "alicloud_actiontrail_trails" "default" {
  name_regex = "tf-testacc-actiontrail"
}

output "trail_name" {
  value = data.alicloud_actiontrail_trails.default.trails.0.id
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by trail name.
* `include_shadow_trails` - (Optional) Whether to show shadow tracking. Default to `false`.
* `include_organization_trail` - (Optional, Available in 1.112+) Whether to show organization tracking. Default to `false`.
* `status` - (Optional) Filter the results by status of the ActionTrail Trail. Valid values: `Disable`, `Enable`, `Fresh`.
* `ids` - (Optional) A list of ActionTrail Trail IDs. It is the same as trail name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of ActionTrail Trail ids. It is the same as trail name.
* `names` - A list of trail names.
* `actiontrails` - (Deprecated since v1.95.0) Field `actiontrails` has been deprecated from version 1.95.0. Use `trails` instead. Each element contains the following attributes:
* `trails` - A list of ActionTrail Trails. Each element contains the following attributes:

### `actiontrails`

The actiontrails supports the following attributes:
* `trail_name` - The name of the ActionTrail Trail.
* `event_rw` - Indicates whether the event is a read or a write event.
* `is_organization_trail` - Indicates whether the trail is a multi-account trail.
* `oss_bucket_name` - The name of the specified OSS bucket.
* `oss_key_prefix` - The prefix of the specified OSS bucket name.
* `oss_write_role_arn` - The RAM role ARN that ActionTrail assumes to access the OSS bucket.
* `sls_project_arn` - The unique ARN of the Log Service project.
* `sls_write_role_arn` - The unique ARN of the Log Service role.
* `status` - The status of the ActionTrail Trail.
* `id` - The id of the ActionTrail Trail. It is the same as trail name.
* `trail_region` - The regions to which the trail is applied.

### `trails`

The trails supports the following attributes:
* `trail_name` - The name of the ActionTrail Trail.
* `event_rw` - Indicates whether the event is a read or a write event.
* `is_organization_trail` - Indicates whether the trail is a multi-account trail.
* `oss_bucket_name` - The name of the specified OSS bucket.
* `oss_key_prefix` - The prefix of the specified OSS bucket name.
* `oss_write_role_arn` - The RAM role ARN that ActionTrail assumes to access the OSS bucket.
* `sls_project_arn` - The unique ARN of the Log Service project.
* `sls_write_role_arn` - The unique ARN of the Log Service role.
* `status` - The status of the ActionTrail Trail.
* `id` - The id of the ActionTrail Trail. It is the same as trail name.
* `trail_region` - The regions to which the trail is applied.
