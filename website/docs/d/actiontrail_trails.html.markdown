---
subcategory: "Actiontrail"
layout: "alicloud"
page_title: "Alicloud: alicloud_actiontrail_trails"
sidebar_current: "docs-alicloud-datasource-actiontrail-trails"
description: |-
  Provides a list of ActionTrail Trails to the specified filters.
---

# alicloud\_actiontrail\_trails

This data source provides a list of ActionTrail Trails in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available in 1.95.0+

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
* `actiontrails` - Field `actiontrails` has been deprecated from version 1.95.0. Use `trails` instead."
* `trails` - A list of ActionTrail Trails. Each element contains the following attributes:
  * `trail_name` - The name of the ActionTrail Trail.
  * `event_rw` - Indicates whether the event is a read or a write event.
  * `oss_bucket_name` - The name of the specified OSS bucket.
  * `oss_key_prefix` - The prefix of the specified OSS bucket name.
  * `sls_project_arn` - The unique ARN of the Log Service project.
  * `sls_write_role_arn` - The unique ARN of the Log Service role.
  * `status` - The status of the ActionTrail Trail.
  * `id` - The id of the ActionTrail Trail. It is the same as trail name.
  * `trail_region` - The regions to which the trail is applied.
