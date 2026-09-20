---
subcategory: "AgentLoop"
layout: "alicloud"
page_title: "Alicloud: alicloud_agentloop_evaluator"
description: |-
  Provides a Alicloud AgentLoop Evaluator resource.
---

# alicloud_agentloop_evaluator

Provides a AgentLoop Evaluator resource.

Evaluator defines how to evaluate the quality of agents, including the evaluation type, metric and configurations.

For information about AgentLoop Evaluator and how to use it, see [What is Evaluator](https://next.api.alibabacloud.com/document/AgentLoop/2026-05-20/CreateEvaluator).

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
  agent_space         = alicloud_agentloop_agent_space.default.agent_space
  name                = var.name
  type                = "custom"
  metric_name         = "example_metric"
  version             = "v1"
  version_description = "initial version"
  display_name        = var.name
  description         = "terraform-example"
  annotations         = ["anno1", "anno2"]
  config = {
    config1 = "value1"
  }
  properties = {
    property1 = "value1"
  }
}
```

## Argument Reference

The following arguments are supported:
* `agent_space` - (Required, ForceNew) The name of the Agent Space to which the Evaluator belongs.
* `annotations` - (Optional, List) The annotations of the Evaluator.
* `config` - (Optional, Map) The configuration of the Evaluator.
* `description` - (Optional) The description of the Evaluator.
* `display_name` - (Optional) The display name of the Evaluator.
* `metric_name` - (Required, ForceNew) The metric name of the Evaluator.
* `name` - (Optional, Computed, ForceNew) The name of the Evaluator.
* `properties` - (Optional, Map) The extended properties of the Evaluator.
* `type` - (Required, ForceNew) The type of the Evaluator.
* `version` - (Optional) The version of the Evaluator. **NOTE:** The Get API does not return `version` (it only returns `current_version`), so the configured value is kept in the state. Changing `version` creates a new evaluator version. `terraform destroy` deletes the whole evaluator including all of its versions.
* `version_description` - (Optional, ForceNew) The description of the version.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Evaluator. It formats as `<agent_space>:<name>`.
* `created_at` - The creation time of the Evaluator.
* `current_version` - The current version of the Evaluator.
* `updated_at` - The last update time of the Evaluator.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Evaluator.
* `delete` - (Defaults to 5 mins) Used when delete the Evaluator.
* `update` - (Defaults to 5 mins) Used when update the Evaluator.

## Import

AgentLoop Evaluator can be imported using the id, e.g.

```shell
$ terraform import alicloud_agentloop_evaluator.example <agent_space>:<name>
```
