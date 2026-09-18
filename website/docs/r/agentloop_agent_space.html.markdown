---
subcategory: "AgentLoop"
layout: "alicloud"
page_title: "Alicloud: alicloud_agentloop_agent_space"
description: |-
  Provides a Alicloud AgentLoop Agent Space resource.
---

# alicloud_agentloop_agent_space

Provides a AgentLoop Agent Space resource.

Agent Space is the basic management unit of AgentLoop, which is used to isolate resources such as agents, sessions and trajectory data between different businesses or environments.

For information about AgentLoop Agent Space and how to use it, see [What is Agent Space](https://next.api.alibabacloud.com/document/AgentLoop/2026-05-20/CreateAgentSpace).

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
  agent_space = var.name
  description = "terraform-example"
}
```

## Argument Reference

The following arguments are supported:
* `agent_space` - (Required, ForceNew) The name of the Agent Space. It must be 2 to 64 characters in length.
* `cms_workspace` - (Optional) The name of an existing CMS (CloudMonitor) workspace to bind to the Agent Space, used for observability data. If omitted, the backend automatically creates a dedicated CMS workspace (recommended).
* `description` - (Optional) The description of the Agent Space.
* `mse_namespace_id` - (Optional, ForceNew) The ID of an existing MSE Nacos namespace to bind to the Agent Space. If omitted, the backend automatically creates a dedicated MSE namespace (recommended).
* `trajectory_store_enabled` - (Optional, ForceNew) Whether to enable trajectory storage for the Agent Space. Valid values: `true`, `false`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is the same as `agent_space`.
* `cms_workspace_bind_type` - The bind type of the CMS workspace. Valid values: `AutoCreated`, `UserSelected`.
* `create_time` - The creation time of the Agent Space.
* `region_id` - The region ID of the Agent Space.
* `sls_project` - The name of the SLS project used to store the trajectory data of the Agent Space.
* `update_time` - The last update time of the Agent Space.

-> **NOTE** Deleting the Agent Space also deletes the auto-created companion resources (SLS project, CMS workspace and MSE namespace). Resources you selected and bound yourself via `cms_workspace` / `mse_namespace_id` are kept.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Agent Space.
* `delete` - (Defaults to 5 mins) Used when delete the Agent Space.
* `update` - (Defaults to 5 mins) Used when update the Agent Space.

## Import

AgentLoop Agent Space can be imported using the id, e.g.

```shell
$ terraform import alicloud_agentloop_agent_space.example <agent_space>
```
