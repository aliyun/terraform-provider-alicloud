---
subcategory: "Actiontrail"
layout: "alicloud"
page_title: "Alicloud: alicloud_actiontrails"
sidebar_current: "docs-alicloud-datasource-actiontrails"
description: |-
  Provides a list of action trail to the user.
---

# alicloud\_actiontrails

-> **DEPRECATED:**  This datasource has been renamed to [alicloud_actiontrail_trails](https://www.terraform.io/docs/providers/alicloud/d/actiontrail_trails) from version 1.95.0.

This data source provides a list of action trail of the current Alibaba Cloud user.

## Example Usage

```
data "alicloud_actiontrails" "trails" {
  name_regex = "tf-testacc-actiontrail"
}

output "first_trail_name" {
  value = "${data.alicloud_actiontrails.trails.actiontrails.0.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results action trail name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of trail names.
* `actiontrails` - A list of actiontrails. Each element contains the following attributes:
  * `name` - The name of the trail.
  * `event_rw` - Indicates whether the event is a read or a write event.
  * `oss_bucket_name` - The name of the specified OSS bucket.
  * `oss_key_prefix` - The prefix of the specified OSS bucket name.
  * `role_name` - The role in ActionTrail.
  * `sls_project_arn` - The unique ARN of the Log Service project.
  * `sls_write_role_arn` - The unique ARN of the Log Service role.
