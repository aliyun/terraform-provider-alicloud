---
subcategory: "AgentLoop"
layout: "alicloud"
page_title: "Alicloud: alicloud_agentloop_evaluation_task"
description: |-
  Provides a Alicloud AgentLoop Evaluation Task resource.
---

# alicloud_agentloop_evaluation_task

Provides a AgentLoop Evaluation Task resource.

Evaluation Task is used to evaluate the quality of agents in the Agent Space with a group of evaluators, and supports continuous and backfill run strategies.

For information about AgentLoop Evaluation Task and how to use it, see [What is Evaluation Task](https://next.api.alibabacloud.com/document/AgentLoop/2026-05-20/CreateEvaluationTask).

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

resource "alicloud_agentloop_evaluator" "default" {
  agent_space = alicloud_agentloop_agent_space.default.agent_space
  name        = "${var.name}-evaluator"
  type        = "custom"
  metric_name = "example_metric"
}

resource "alicloud_agentloop_evaluation_task" "default" {
  agent_space = alicloud_agentloop_agent_space.default.agent_space
  task_name   = var.name
  task_mode   = "batch"
  data_type   = "trace"
  channel     = "default"
  data_filter = "{\"query\":\"checkout-service\",\"maxRecords\":10}"
  description = "terraform-example"
  config = {
    config1 = "value1"
  }
  tags = {
    tag1 = "value1"
  }
  run_strategies {
    continuous {
      enabled            = true
      interval_unit      = "hours"
      interval_value     = 6
      data_delay_minutes = 30
    }
    backfill {
      enabled    = true
      immediate  = true
      start_time = 1735689600000
      end_time   = 1735776000000
    }
  }
  evaluators {
    name          = "evaluator-a"
    result_name   = "result-a"
    type          = "custom"
    evaluator_ref = alicloud_agentloop_evaluator.default.name
    filters = {
      filter1 = "value1"
    }
    config = {
      config1 = "value1"
    }
    variable_mapping = {
      input = "$.input"
    }
  }
}
```

## Argument Reference

The following arguments are supported:
* `agent_space` - (Required, ForceNew) The name of the Agent Space to which the Evaluation Task belongs.
* `channel` - (Optional, ForceNew) The channel of the Evaluation Task. It cannot be changed after creation (the update API treats it as an internal-only parameter and silently ignores changes).
* `config` - (Optional, Map) The configuration of the Evaluation Task.
* `data_filter` - (Optional) The filter expression of the data to be evaluated.
* `data_type` - (Optional, ForceNew) The type of the data to be evaluated.
* `description` - (Optional) The description of the Evaluation Task.
* `evaluators` - (Required) The evaluators of the Evaluation Task. See [`evaluators`](#evaluators) below.
* `run_strategies` - (Optional, ForceNew) The run strategies of the Evaluation Task. See [`run_strategies`](#run_strategies) below.
* `status` - (Optional, Computed) The status of the Evaluation Task. Valid values: `Pending`, `Running`, `Completed`, `Scheduling`, `Failed`, `Terminated`. `Deleted` is intentionally not configurable: the backend accepts but silently ignores status transitions to `Deleted` (the task keeps its previous status), so waiting for it would never succeed. To delete the task, remove the resource with `terraform destroy` instead.
* `tags` - (Optional, Map) The tags of the Evaluation Task.
* `task_mode` - (Optional, ForceNew) The mode of the Evaluation Task.
* `task_name` - (Required, ForceNew) The name of the Evaluation Task.

### `evaluators`

The evaluators supports the following:
* `config` - (Optional, Map) The configuration of the evaluator.
* `evaluator_ref` - (Optional) The reference name of the evaluator.
* `filters` - (Optional, Map) The filters of the evaluator.
* `name` - (Required) The name of the evaluator.
* `result_name` - (Required) The name of the evaluation result.
* `result_type` - (Optional) The type of the evaluation result.
* `type` - (Optional) The type of the evaluator.
* `variable_mapping` - (Required, Map) The variable mapping of the evaluator.

### `run_strategies`

The run_strategies supports the following:
* `backfill` - (Optional) The backfill strategy. See [`backfill`](#run_strategies-backfill) below.
* `continuous` - (Optional) The continuous strategy. See [`continuous`](#run_strategies-continuous) below.

### `run_strategies-backfill`

The backfill supports the following:
* `enabled` - (Optional) Whether to enable the backfill strategy.
* `end_time` - (Optional, Int) The end time of the backfill.
* `immediate` - (Optional) Whether to run the backfill immediately.
* `start_time` - (Optional, Int) The start time of the backfill.

### `run_strategies-continuous`

The continuous supports the following:
* `data_delay_minutes` - (Optional, Int) The data delay in minutes.
* `enabled` - (Optional) Whether to enable the continuous strategy.
* `interval_unit` - (Optional) The unit of the run interval.
* `interval_value` - (Optional, Int) The value of the run interval.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Evaluation Task. It formats as `<agent_space>:<task_id>`.
* `created_at` - The creation time of the Evaluation Task.
* `region_id` - The region ID of the Evaluation Task.
* `task_id` - The ID of the Evaluation Task.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Evaluation Task.
* `delete` - (Defaults to 5 mins) Used when delete the Evaluation Task.
* `update` - (Defaults to 5 mins) Used when update the Evaluation Task.

## Import

AgentLoop Evaluation Task can be imported using the id, e.g.

```shell
$ terraform import alicloud_agentloop_evaluation_task.example <agent_space>:<task_id>
```
