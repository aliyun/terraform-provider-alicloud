---
subcategory: "unreleased"
layout: "alicloud"
page_title: "Alicloud: alicloud_agentloop_agent_space"
description: |-
  Provides a Alicloud Agent Loop Agent Space resource.
---

# alicloud_agentloop_agent_space

Provides a Agent Loop Agent Space resource.

AgentSpace is the top-level resource container for AgentLoop. It manages the association with CMS Workspace, MSE Namespace and SLS Project, and controls the lifecycle of those linked resources on deletion.

For information about Agent Loop Agent Space and how to use it, see [What is Agent Space](https://next.api.alibabacloud.com/document/AgentLoop/2026-05-20/CreateAgentSpace).

-> **NOTE:** Available since v1.287.0.

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
  agent_space   = var.name
  description   = "created by terraform"
  cms_workspace = "cms-workspace-example"
  mse_workspace = "mse-namespace-example"
}
```

## Argument Reference

The following arguments are supported:

* `agent_space` - (Required, ForceNew) AgentSpace name.
* `cms_workspace` - (Optional, ForceNew) Associated CMS Workspace.
* `description` - (Optional) Description of the AgentSpace.
* `mse_workspace` - (Optional, ForceNew) Associated MSE Workspace.
* `delete_cms_workspace` - (Optional) Whether to delete the associated CMS Workspace when deleting the AgentSpace. Default to false.
* `delete_mse_namespace` - (Optional) Whether to delete the associated MSE Namespace when deleting the AgentSpace. Default to false.
* `delete_sls_project` - (Optional) Whether to delete the associated SLS Project when deleting the AgentSpace. Default to false.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource supplied above. It is the same as `agent_space`.
* `create_time` - The creation time of the resource.
* `region_id` - The region ID of the resource.
* `sls_project` - Automatically created SLS Project.
* `update_time` - The last modification time of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Agent Space.
* `update` - (Defaults to 5 mins) Used when update the Agent Space.
* `delete` - (Defaults to 5 mins) Used when delete the Agent Space.

## Import

Agent Loop Agent Space can be imported using the id, e.g.

```shell
$ terraform import alicloud_agentloop_agent_space.example <agent_space>
```
