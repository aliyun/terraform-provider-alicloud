---
subcategory: "AgentLoop"
layout: "alicloud"
page_title: "Alicloud: alicloud_agentloop_agent_space"
description: |-
  Provides a Alicloud Agent Loop Agent Space resource.
---

# alicloud_agentloop_agent_space

Provides a Agent Loop Agent Space resource.

AgentSpace is the basic unit for organizing and managing resources in AgentLoop.

For information about Agent Loop Agent Space and how to use it, see [What is Agent Space](https://next.api.alibabacloud.com/document/AgentLoop/2026-05-20/CreateAgentSpace).

-> **NOTE:** Available since v1.294.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_agentloop_agent_space" "default" {
  agent_space = var.name
  description = "AgentSpace created by terraform example"
}
```

## Argument Reference

The following arguments are supported:
* `agent_space` - (Required, ForceNew) AgentSpace name.
* `cms_workspace` - (Optional, ForceNew) Associated CMS Workspace.
* `description` - (Optional) Description.
* `mse_workspace` - (Optional, ForceNew) Associated MSE Workspace.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - A resource property field representing the creation time.
* `region_id` - A resource property field representing the region ID.
* `sls_project` - Automatically created SLS Project.
* `update_time` - Last modification time.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Agent Space.
* `delete` - (Defaults to 5 mins) Used when delete the Agent Space.
* `update` - (Defaults to 5 mins) Used when update the Agent Space.

## Import

Agent Loop Agent Space can be imported using the id, e.g.

```shell
$ terraform import alicloud_agentloop_agent_space.example <agent_space>
```