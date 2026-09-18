---
subcategory: "AgentLoop"
layout: "alicloud"
page_title: "Alicloud: alicloud_agentloop_optimize_task"
description: |-
  Provides a Alicloud AgentLoop Optimize Task resource.
---

# alicloud_agentloop_optimize_task

Provides a AgentLoop Optimize Task resource.

Optimize Task is used to create an optimization task in the Agent Space to analyze and optimize the agents.

For information about AgentLoop Optimize Task and how to use it, see [What is Optimize Task](https://next.api.alibabacloud.com/document/AgentLoop/2026-05-20/CreateOptimizeTask).

-> **NOTE:** Available since v1.294.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_agentloop_agent_space" "default" {
  agent_space = "${var.name}-as"
}

resource "alicloud_agentloop_optimize_task" "default" {
  agent_space        = alicloud_agentloop_agent_space.default.agent_space
  optimize_task_name = var.name
  type               = "EfficiencyAnalysis"
  status             = "Enabled"
  config = {
    config1 = "value1"
  }
}
```

## Argument Reference

The following arguments are supported:
* `agent_space` - (Required, ForceNew) The name of the Agent Space to which the Optimize Task belongs.
* `config` - (Optional, ForceNew, Map) The configuration of the Optimize Task.
* `optimize_task_name` - (Required, ForceNew) The name of the Optimize Task.
* `status` - (Required, ForceNew) The status of the Optimize Task. Valid values: `Enabled`, `Disabled`.
* `type` - (Required, ForceNew) The type of the Optimize Task.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Optimize Task. It formats as `<agent_space>:<optimize_task_name>`.
* `create_time` - The creation time of the Optimize Task.
* `region_id` - The region ID of the Optimize Task.
* `update_time` - The last update time of the Optimize Task.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Optimize Task.
* `delete` - (Defaults to 5 mins) Used when delete the Optimize Task.

## Import

AgentLoop Optimize Task can be imported using the id, e.g.

```shell
$ terraform import alicloud_agentloop_optimize_task.example <agent_space>:<optimize_task_name>
```
