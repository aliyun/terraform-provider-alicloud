---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_session_manager_status"
sidebar_current: "docs-alicloud-resource-ecs-session-manager-status"
description: |-
  Provides a Alicloud ECS Session Manager Status resource.
---

# alicloud\_ecs\_session\_manager\_status

Provides a ECS Session Manager Status resource.

For information about ECS Session Manager Status and how to use it, see [What is Session Manager Status](https://www.alibabacloud.com/help/zh/doc-detail/337915.html).

-> **NOTE:** Available in v1.148.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecs_session_manager_status" "default" {
  session_manager_status_name = "sessionManagerStatus"
  status                      = "Disabled"
}
```

## Argument Reference

The following arguments are supported:

* `session_manager_status_name` - (Required, ForceNew) The name of the resource. Valid values: `sessionManagerStatus`.
* `status` - (Required) The status of the resource. Valid values: `Disabled`, `Enabled`.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Session Manager Status. Its value is same as `session_manager_status_name`.

## Import

ECS Session Manager Status can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_session_manager_status.example <session_manager_status_name>
```