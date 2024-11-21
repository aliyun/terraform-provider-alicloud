---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_deployment_set"
sidebar_current: "docs-alicloud-resource-ecs-deployment-set"
description: |-
  Provides a Alicloud ECS Deployment Set resource.
---

# alicloud_ecs_deployment_set

Provides a ECS Deployment Set resource.

For information about ECS Deployment Set and how to use it, see [What is Deployment Set](https://www.alibabacloud.com/help/en/doc-detail/91269.htm).

-> **NOTE:** Available since v1.140.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_deployment_set&exampleId=53b51b98-dccf-878a-fb01-554a57d62e90fb9bca66&activeTab=example&spm=docs.r.ecs_deployment_set.0.53b51b98dc&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
* `strategy` - (Optional, ForceNew) The deployment strategy. Valid values: `Availability`(Default), `AvailabilityGroup`, `LowLatency`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Deployment Set.

## Import

ECS Deployment Set can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_deployment_set.example <id>
```
