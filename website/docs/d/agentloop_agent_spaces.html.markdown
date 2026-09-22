---
subcategory: "unreleased"
layout: "alicloud"
page_title: "Alicloud: alicloud_agentloop_agent_spaces"
sidebar_current: "docs-alicloud-datasource-agentloop-agent-spaces"
description: |-
  Provides a list of Agent Loop Agent Space owned by an Alibaba Cloud account.
---

# alicloud_agentloop_agent_spaces

This data source provides Agent Loop Agent Space available to the user.[What is Agent Space](https://next.api.alibabacloud.com/document/AgentLoop/2026-05-20/CreateAgentSpace)

-> **NOTE:** Available since v1.287.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_agentloop_agent_spaces" "default" {
  agent_space   = var.name
  biz_region_id = "cn-hangzhou"
}

output "first_space_id" {
  value = data.alicloud_agentloop_agent_spaces.default.spaces.0.id
}
```

## Argument Reference

The following arguments are supported:
* `agent_space` - (ForceNew, Optional) AgentSpace name.
* `biz_region_id` - (Optional) Filter results by business region ID.
* `ids` - (Optional, Computed) A list of Agent Space IDs.
* `max_results` - (Optional) The maximum number of results to return per page. If not set, the API default is used.
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
