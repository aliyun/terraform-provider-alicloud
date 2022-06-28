---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_deployment_set"
sidebar_current: "docs-alicloud-resource-ecs-deployment-set"
description: |-
  Provides a Alicloud ECS Deployment Set resource.
---

# alicloud\_ecs\_deployment\_set

Provides a ECS Deployment Set resource.

For information about ECS Deployment Set and how to use it, see [What is Deployment Set](https://www.alibabacloud.com/help/en/doc-detail/91269.htm).

-> **NOTE:** Available in v1.140.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecs_deployment_set" "default" {
  strategy            = "Availability"
  domain              = "Default"
  granularity         = "Host"
  deployment_set_name = "example_value"
  description         = "example_value"
}
```

## Argument Reference

The following arguments are supported:

* `deployment_set_name` - (Optional) The name of the deployment set. The name must be 2 to 128 characters in length and can contain letters, digits, colons (:), underscores (_), and hyphens (-). It must start with a letter and cannot start with `http://` or `https://`.
* `description` - (Optional) The description of the deployment set. The description must be 2 to 256 characters in length and cannot start with `http://` or `https://`.
* `domain` - (Optional, ForceNew) The deployment domain. Valid values: `Default`.
* `granularity` - (Optional, ForceNew) The deployment granularity. Valid values: `Host`.
* `on_unable_to_redeploy_failed_instance` - (Optional) The on unable to redeploy failed instance. Valid values: `CancelMembershipAndStart`, `KeepStopped`.
  * `CancelMembershipAndStart` - Removes the instances from the deployment set and restarts the instances immediately after the failover is complete.
  * `KeepStopped`- Keeps the instances in the abnormal state and restarts them after ECS resources are replenished. 
* `strategy` - (Optional, ForceNew) The deployment strategy. Valid values: `Availability`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Deployment Set.

## Import

ECS Deployment Set can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_deployment_set.example <id>
```
