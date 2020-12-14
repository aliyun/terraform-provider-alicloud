---
subcategory: "Actiontrail"
layout: "alicloud"
page_title: "Alicloud: alicloud_actiontrail_trail"
sidebar_current: "docs-alicloud-resource-actiontrail-trail"
description: |-
  Provides Alibaba Cloud ActionTrail Trail Resource
---

# alicloud\_actiontrail\_trail

Provides a ActionTrail Trail resource. For information about alicloud actiontrail trail and how to use it, see [What is Resource Alicloud ActionTrail Trail](https://www.alibabacloud.com/help/doc-detail/28804.htm).

-> **NOTE:** Available in 1.95.0+

## Example Usage

```terraform
# Create a new actiontrail trail.
resource "alicloud_actiontrail_trail" "default" {
  trail_name      = "action-trail"
  role_name       = "aliyunserviceroleforactiontrail"
  oss_bucket_name = "bucket_name"
  event_rw        = "All"
  trail_region    = "cn-hangzhou"
}
```

## Argument Reference

The following arguments are supported:

* `trail_name` - (Optional, ForceNew) The name of the trail to be created, which must be unique for an account.
* `name` - (Optional, ForceNew) Field `name` has been deprecated from version 1.95.0. Use `trail_name` instead. 
* `event_rw` - (Optional) Indicates whether the event is a read or a write event. Valid values: `Read`, `Write`, and `All`. Default to `Write`.
* `oss_bucket_name` - (Optional) The OSS bucket to which the trail delivers logs. Ensure that this is an existing OSS bucket.
* `role_name` - (Optional) The RAM role in ActionTrail permitted by the user.
* `oss_key_prefix` - (Optional) The prefix of the specified OSS bucket name. This parameter can be left empty.
* `sls_project_arn` - (Optional) The unique ARN of the Log Service project.
* `sls_write_role_arn` - (Optional) The unique ARN of the Log Service role.
* `trail_region` - (Optional) The regions to which the trail is applied. Valid values: `cn-beijing`, `cn-hangzhou`, and `All`. Default to `All`.
* `mns_topic_arn` - (Optional) The ARN of the Message Service (MNS) topic to which ActionTrail sends messages. If the ARN is specified, a message is generated and delivered to the MNS topic whenever an event is delivered to OSS.
* `status` - (Optional) The status of ActionTrail Trail. After creation, tracking is turned on by default, and you can set the status value to `Disable` to turn off tracking. Valid values: `Enable`, `Disable`. Default to `Enable`.

-> **NOTE:** `sls_project_arn` and `sls_write_role_arn` should be set or not set at the same time when actiontrail delivers logs.

## Attributes Reference

The following attributes are exported:

* `id` - The id of ActionTrail Trail. The value is the same as trail_name.

## Import

Action trail can be imported using the id or trail_name, e.g.

```
$ terraform import alicloud_actiontrail_trail.default abc12345678
```
