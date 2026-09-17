---
subcategory: "AgentLoop"
layout: "alicloud"
page_title: "Alicloud: alicloud_agentloop_pipeline"
description: |-
  Provides a Alicloud AgentLoop Pipeline resource.
---

# alicloud_agentloop_pipeline

Provides a AgentLoop Pipeline resource.

Pipeline defines a data processing pipeline in the Agent Space, including the processing nodes, the data source (logstore or dataset), the data sink (dataset or condition routes) and the execute policy.

For information about AgentLoop Pipeline and how to use it, see [What is Pipeline](https://next.api.alibabacloud.com/document/AgentLoop/2026-05-20/CreatePipeline).

-> **NOTE:** Available since v1.295.0.

-> **NOTE:** A Pipeline whose schedule status is `Active` cannot be deleted directly. When `terraform destroy` meets this state, the provider automatically pauses the Pipeline through the `PausePipeline` API and then deletes it.

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

resource "alicloud_agentloop_dataset" "default" {
  agent_space  = alicloud_agentloop_agent_space.default.agent_space
  dataset_name = "terraform_example_ds1"
  schema = {
    input = jsonencode({ type = "text", chn = true })
  }
}

resource "alicloud_agentloop_dataset" "default2" {
  agent_space  = alicloud_agentloop_agent_space.default.agent_space
  dataset_name = "terraform_example_ds2"
  schema = {
    input = jsonencode({ type = "text", chn = true })
  }
}

resource "alicloud_agentloop_pipeline" "default" {
  agent_space   = alicloud_agentloop_agent_space.default.agent_space
  pipeline_name = var.name
  description   = "terraform-example"
  pipeline {
    nodes {
      id   = "node-1"
      type = "where"
      parameters = {
        filter = "input IS NOT NULL"
      }
    }
  }
  source {
    type = "dataset"
    input_fields {
      name = "input"
      type = "text"
    }
    dataset {
      dataset = alicloud_agentloop_dataset.default.dataset_name
      filter  = "status = 'pending'"
    }
  }
  sink {
    type = "dataset"
    dataset {
      agent_space = alicloud_agentloop_agent_space.default.agent_space
      dataset     = alicloud_agentloop_dataset.default2.dataset_name
    }
  }
  execute_policy {
    mode = "RunOnce"
    run_once {
      from_time = 1735660800
      to_time   = 1735747200
    }
  }
}
```

## Argument Reference

The following arguments are supported:
* `agent_space` - (Required, ForceNew) The name of the Agent Space to which the Pipeline belongs.
* `description` - (Optional) The description of the Pipeline. It can be up to 256 characters in length and can be updated in place.
* `execute_policy` - (Optional, ForceNew, Set) The execute policy of the Pipeline. See [`execute_policy`](#execute_policy) below.
* `pipeline` - (Required, Set) The definition of the Pipeline. It can be updated in place. See [`pipeline`](#pipeline) below.
* `pipeline_name` - (Required, ForceNew) The name of the Pipeline. It must be 3 to 63 characters in length and can contain only lowercase letters, digits, and hyphens (`-`).
* `sink` - (Optional, ForceNew, Set) The data sink of the Pipeline. See [`sink`](#sink) below.
* `source` - (Optional, Set) The data source of the Pipeline. The dataset source filter can be updated in place; changes to the source type, input fields, logstore source, or source dataset recreate the resource. See [`source`](#source) below.

### `execute_policy`

The execute_policy supports the following:
* `mode` - (Optional, ForceNew) The execute mode of the Pipeline. Valid values: `RunOnce`, `Scheduled`.
* `run_once` - (Optional, ForceNew, Set) The run-once policy. It is required when `mode` is `RunOnce`. See [`run_once`](#execute_policy-run_once) below.
* `scheduled` - (Optional, ForceNew, Set) The scheduled policy. It is required when `mode` is `Scheduled`. See [`scheduled`](#execute_policy-scheduled) below.

### `execute_policy-run_once`

The run_once supports the following:
* `from_time` - (Optional, ForceNew, Int) The start time of the run, in Unix seconds. It must be less than `to_time`.
* `to_time` - (Optional, ForceNew, Int) The end time of the run, in Unix seconds. It must be greater than `from_time`.

### `execute_policy-scheduled`

The scheduled supports the following:
* `from_time` - (Optional, ForceNew, Int) The start time of the schedule, in Unix milliseconds.
* `interval` - (Optional, ForceNew) The run interval. Valid values: `1h`, `6h`, `12h`, `1d`.

### `pipeline`

The pipeline supports the following:
* `nodes` - (Required, Set) The processing nodes of the Pipeline. See [`nodes`](#pipeline-nodes) below.

### `pipeline-nodes`

The nodes supports the following:
* `id` - (Optional) The ID of the node.
* `parameters` - (Optional, Map) The parameters of the node. The parameters vary depending on the node type. A value that is a valid JSON document (e.g. `jsonencode(["a", "b"])` for an array, `"10"` for a number, `"true"` for a boolean) is sent in its parsed form; any other value is sent verbatim as a string.
* `type` - (Optional, ForceNew) The type of the node.

### `sink`

The sink supports the following:
* `condition` - (Optional, ForceNew, Set) The condition sink. It is valid when `type` is `condition`. See [`condition`](#sink-condition) below.
* `dataset` - (Optional, ForceNew, Set) The dataset sink. It is valid when `type` is `dataset`. See [`dataset`](#sink-dataset) below.
* `type` - (Optional, ForceNew) The type of the sink. Valid values: `dataset`, `condition`.

### `sink-condition`

The condition supports the following:
* `default_sink` - (Optional, ForceNew, Set) The default sink target when no route matches. See [`default_sink`](#sink-condition-default_sink) below.
* `match_mode` - (Optional, ForceNew) The match mode of the routes. Valid values: `all`.
* `routes` - (Optional, ForceNew, Set) The condition routes. See [`routes`](#sink-condition-routes) below.

### `sink-condition-default_sink`

The default_sink supports the following:
* `dataset` - (Optional, ForceNew, Set) The dataset sink target. See [`dataset`](#sink-condition-default_sink-dataset) below.
* `type` - (Optional, ForceNew) The type of the sink target. Valid values: `dataset`.

### `sink-condition-default_sink-dataset`

The dataset supports the following:
* `agent_space` - (Optional, ForceNew) The name of the Agent Space to which the target dataset belongs.
* `dataset` - (Optional, ForceNew) The name of the target dataset.

### `sink-condition-routes`

The routes supports the following:
* `expression` - (Optional, ForceNew) The filter expression of the route, written in SPL. Only the `where`, `project` and `extend` instructions are supported.
* `id` - (Optional, ForceNew) The ID of the route.
* `sink` - (Optional, ForceNew, Set) The sink target of the route. See [`sink`](#sink-condition-routes-sink) below.

### `sink-condition-routes-sink`

The sink supports the following:
* `dataset` - (Optional, ForceNew, Set) The dataset sink target. See [`dataset`](#sink-condition-routes-sink-dataset) below.
* `type` - (Optional, ForceNew) The type of the sink target. Valid values: `dataset`.

### `sink-condition-routes-sink-dataset`

The dataset supports the following:
* `agent_space` - (Optional, ForceNew) The name of the Agent Space to which the target dataset belongs.
* `dataset` - (Optional, ForceNew) The name of the target dataset.

### `sink-dataset`

The dataset supports the following:
* `agent_space` - (Optional, ForceNew) The name of the Agent Space to which the target dataset belongs.
* `dataset` - (Optional, ForceNew) The name of the target dataset.

### `source`

The source supports the following:
* `dataset` - (Optional, Set) The dataset source. It is valid when `type` is `dataset`. See [`dataset`](#source-dataset) below.
* `input_fields` - (Optional, ForceNew, Set) The input fields of the source. See [`input_fields`](#source-input_fields) below.
* `logstore` - (Optional, ForceNew, Set) The logstore source. It is valid when `type` is `logstore`. See [`logstore`](#source-logstore) below.
* `type` - (Optional, ForceNew) The type of the source. Valid values: `dataset`, `logstore`.

### `source-dataset`

The dataset supports the following:
* `dataset` - (Optional, ForceNew) The name of the source dataset. The dataset must belong to the Agent Space of the Pipeline.
* `filter` - (Optional) The filter expression of the source dataset. It can be updated in place.

### `source-input_fields`

The input_fields supports the following:
* `name` - (Optional, ForceNew) The name of the input field.
* `type` - (Optional, ForceNew) The type of the input field. Valid values: `text`, `long`, `double`, `json`.

### `source-logstore`

The logstore supports the following:
* `logstore` - (Optional) The name of the SLS logstore.
* `project` - (Optional) The name of the SLS project.
* `query` - (Optional) The query statement used to filter the logstore data, written in the SLS query and analysis syntax.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Pipeline. It formats as `<agent_space>:<pipeline_name>`.
* `create_time` - The creation time of the Pipeline, in ISO 8601 UTC format (`yyyy-MM-ddTHH:mm:ssZ`).
* `region_id` - The region ID of the Pipeline.
* `update_time` - The last update time of the Pipeline, in ISO 8601 UTC format (`yyyy-MM-ddTHH:mm:ssZ`).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Pipeline.
* `delete` - (Defaults to 5 mins) Used when delete the Pipeline.
* `update` - (Defaults to 5 mins) Used when update the Pipeline.

## Import

AgentLoop Pipeline can be imported using the id, e.g.

```shell
$ terraform import alicloud_agentloop_pipeline.example <agent_space>:<pipeline_name>
```
