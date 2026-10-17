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

-> **NOTE:** Available since v1.295.0.

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
  name        = "terraform_example_evaluator"
  type        = "AGENT"
  metric_name = "example_metric"
  version     = "v1"
  # The variable_mapping keys of the Evaluation Task evaluators must be
  # defined in the referenced evaluator's variables.
  config = {
    variables = jsonencode([{ name = "input", type = "text" }])
  }
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
      interval_unit      = "HOUR"
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
    # Reference mode: the registered evaluator is loaded by evaluator_ref.
    evaluator_ref = alicloud_agentloop_evaluator.default.name
    result_name   = "result_a"
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
* `channel` - (Optional, Computed, ForceNew) The channel of the Evaluation Task. It cannot be changed after creation (the update API treats it as an internal-only parameter and silently ignores changes). If not specified, the backend defaults it to `default`.
* `config` - (Optional, Map) The configuration of the Evaluation Task. The update API replaces the whole map. The backend injects derived data source keys (`dataScope`, `project`, `storeName`, `traceFormat`) resolved from the Agent Space; they are not tracked unless configured, and they are dropped by the backend when `config` is changed. All other keys are read back, so out-of-band changes are detected as drift.
* `data_filter` - (Optional) The filter expression of the data to be evaluated, as a JSON object string (e.g. `jsonencode({ query = "checkout-service", maxRecords = 10 })`). Semantically equivalent JSON (different key order or whitespace) does not produce a diff.
* `data_type` - (Optional, ForceNew) The type of the data to be evaluated. Valid values: `trace`, `log`, `atif`, `dataset`.
* `description` - (Optional) The description of the Evaluation Task.
* `evaluators` - (Required, Set) The evaluators of the Evaluation Task. It cannot be empty, and each `evaluator_ref` (or `name` for inline evaluators) must be unique within the task. The update API replaces the whole list. The evaluators are read back from the API, so out-of-band changes are detected as drift. See [`evaluators`](#evaluators) below.
* `run_strategies` - (Optional, Set) The run strategies of the Evaluation Task. It can be updated in place; the update API replaces the whole strategy, and removing the block clears the run strategies. When a terminated task enables a strategy again, the backend moves it back to `Pending` and schedules it. See [`run_strategies`](#run_strategies) below.
* `status` - (Optional, Computed) The status of the Evaluation Task. Valid values: `Pending`, `Running`, `Completed`, `Scheduling`, `Failed`, `Terminated`. The status is only sent to the update API when it is changed in the configuration, so updating other attributes never echoes back the backend-managed status. `Deleted` is intentionally not configurable: the backend accepts but silently ignores status transitions to `Deleted` (the task keeps its previous status), so waiting for it would never succeed. To delete the task, remove the resource with `terraform destroy` instead.
* `tags` - (Optional, Map) The tags of the Evaluation Task. Removing all tags clears them.
* `task_mode` - (Optional, ForceNew) The mode of the Evaluation Task, for example `batch` (a persistent evaluation task).
* `task_name` - (Required, ForceNew) The name of the Evaluation Task. It can be up to 256 characters in length and must be unique among the undeleted tasks in the Agent Space.

### `evaluators`

An evaluator works in one of the following modes:
- Reference mode (`evaluator_ref` is set): the registered built-in or custom evaluator is loaded. `evaluator_ref` takes precedence over the inline fields. The API returns `evaluator_ref` as the evaluator name (a configured `name` is kept in the state), `type` defaults to `LLM`, and `result_name` falls back to the metric name of the referenced evaluator.
- Inline mode (`evaluator_ref` is not set): `name` and `result_type` are required, and `variable_mapping` is required for `LLM` and `AGENT` evaluators. `CODE` evaluators cannot be defined inline and must be referenced via `evaluator_ref`.

The evaluators supports the following:
* `config` - (Optional, Map) The configuration of the evaluator. For inline LLM evaluators it carries the evaluation settings such as `prompt`; for referenced evaluators it is usually omitted and only carries run parameters such as `version`. The backend injects the resolved `version`, which is only tracked when configured.
* `evaluator_ref` - (Optional) The name of the registered (built-in or custom) evaluator to reference. When set, it takes precedence over the inline definition.
* `filters` - (Optional, Map) The filters of the evaluator. Its `query` takes effect together with the `query` of the task-level `data_filter`.
* `name` - (Optional, Computed) The name of the evaluator. It is required when `evaluator_ref` is not set.
* `result_name` - (Optional, Computed) The name of the evaluation result. In reference mode it falls back to the metric name of the referenced evaluator.
* `result_type` - (Optional, Computed) The type of the evaluation result. Valid values: `score`, `binary`, `text`. It is required when `evaluator_ref` is not set; in reference mode it defaults to `score`.
* `type` - (Optional, Computed) The type of the evaluator. Valid values: `LLM`, `AGENT`, `CODE`. Defaults to `LLM`. `CODE` is only supported in reference mode.
* `variable_mapping` - (Optional, Map) The variable mapping of the evaluator. It is required for inline `LLM` and `AGENT` evaluators. In reference mode, each key must be a variable defined by the referenced evaluator.

### `run_strategies`

A strategy block that is present is enabled unless its `enabled` is explicitly set to `false`.

The run_strategies supports the following:
* `backfill` - (Optional, Set) The backfill strategy. See [`backfill`](#run_strategies-backfill) below.
* `continuous` - (Optional, Set) The continuous strategy. See [`continuous`](#run_strategies-continuous) below.

### `run_strategies-backfill`

The backfill supports the following:
* `enabled` - (Optional) Whether to enable the backfill strategy. Default value: `true` (matching the API, which treats an omitted value as enabled). `false` disables the strategy but keeps its configuration.
* `end_time` - (Optional, Int) The end time of the backfill range, as a Unix timestamp in milliseconds. To start a backfill manually, provide both `start_time` and `end_time`.
* `immediate` - (Optional) Whether to run the backfill immediately. Default value: `false`.
* `start_time` - (Optional, Int) The start time of the backfill range, as a Unix timestamp in milliseconds. To start a backfill manually, provide both `start_time` and `end_time`.

### `run_strategies-continuous`

The continuous supports the following:
* `data_delay_minutes` - (Optional, Int) The data arrival delay, in minutes. Each run is created this many minutes after its window ends, so that the data of the window has fully arrived. Default value: `0`.
* `enabled` - (Optional) Whether to enable the continuous strategy. Default value: `true` (matching the API, which treats an omitted value as enabled). `false` disables the strategy but keeps its configuration.
* `interval_unit` - (Required) The unit of the run interval. Valid values: `HOUR`, `MINUTE`, `DAY`.
* `interval_value` - (Required, Int) The value of the run interval, used together with `interval_unit`. It must be greater than or equal to `1`.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Evaluation Task. It formats as `<agent_space>:<task_id>`.
* `created_at` - The creation time of the Evaluation Task, as a Unix timestamp in seconds.
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
