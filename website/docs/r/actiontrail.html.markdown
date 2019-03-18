---
layout: "alicloud"
page_title: "Alicloud: alicloud_actiontrail"
sidebar_current: "docs-alicloud-resource-actiontrail"
description: |-
  It provides descriptions of operations, common parameters, and common errors of ActionTrail.
---

# alicloud\_actiontrail

Provides descriptions of operations, common parameters, and common errors of ActionTrail.

~> **NOTE:** Make sure that you are familiar with the working nature of ActionTrail and fully aware of the SLA before using the APIs.

-> NOTE: Available in 1.35.0+

## Example Usage

```
# Create a new action trail.
resource "alicloud_actiontrail" "foo" {
	name = "action-trail"
	event_rw = "Write-test"
	oss_bucket_name = "${alicloud_oss_bucket.bucket.id}"
	role_name = "${alicloud_ram_role_policy_attachment.attach.role_name}"
	oss_key_prefix = "at-product-account-audit-B"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required,ForceNew) The name of the trail to be created, which must be unique for an account.
* `event_rw` - (Optional) Indicates whether the event is a read or a write event. Valid values: Read, Write, and All. Default value: Write.
* `oss_bucket_name` - (Required) The OSS bucket to which the trail delivers logs. Ensure that this is an existing OSS bucket.
* `role_name` - (Required) The RAM role in ActionTrail permitted by the user.
* `oss_key_prefix` - (Optional) The prefix of the specified OSS bucket name. This parameter can be left empty.
* `sls_project_arn` - (Optional) The unique ARN of the Log Service project.
* `sls_write_role_arn` - (Optional) The unique ARN of the Log Service role.

## Attributes Reference

The following attributes are exported:

* `id` - The action trail id. The value is same as its name.

## Import

Action trail can be imported using the id, e.g.

```
$ terraform import alicloud_actiontrail.foo abc12345678
```
