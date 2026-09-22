---
subcategory: "AgentLoop"
layout: "alicloud"
page_title: "Alicloud: alicloud_agentloop_agent_spaces"
sidebar_current: "docs-alicloud-datasource-agentloop-agent-spaces"
description: |-
  Provides a list of Agent Loop Agent Space owned by an Alibaba Cloud account.
---

# alicloud_agentloop_agent_spaces

This data source provides Agent Loop Agent Space available to the user.

[What is Agent Space](https://next.api.alibabacloud.com/document/AgentLoop/2026-05-20/CreateAgentSpace)

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
data "alicloud_agentloop_agent_spaces" "default" {
  agent_space = "terraform-example"
}

output "first_agent_space" {
  value = data.alicloud_agentloop_agent_spaces.default.spaces.0
}
```

## Argument Reference

The following arguments are supported:
* `agent_space` - (Optional) AgentSpace name.
* `ids` - (Optional, Computed) A list of Agent Space IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Agent Space IDs.
* `spaces` - A list of Agent Space Entries. Each element contains the following attributes:
  * `agent_space` - AgentSpace name.
  * `cms_workspace` - Associated CMS Workspace.
  * `create_time` - A resource property field representing the creation time.
  * `description` - Description.
  * `mse_workspace` - Associated MSE Workspace.
  * `region_id` - A resource property field representing the region ID.
  * `sls_project` - Automatically created SLS Project.
  * `update_time` - Last modification time.
  * `id` - The ID of the resource supplied above.
