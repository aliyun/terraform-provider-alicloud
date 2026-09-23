---
subcategory: "Chatbot"
layout: "alicloud"
page_title: "Alicloud: alicloud_chatbot_agents"
sidebar_current: "docs-alicloud-resource-chatbot-agents"
description: |-
  Provides a list of Chatbot Agents to the user.
---

# alicloud\_chatbot\_agents

This data source provides the Chatbot Agents of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.203.0+.

## Example Usage

```terraform
data "alicloud_chatbot_agents" "nameRegex" {
  name_regex = "^my-Agent"
}
output "alicloud_chatbot_agents_id_1" {
  value = data.alicloud_chatbot_agents.nameRegex.agents.0.id
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `agent_name` - (Optional,ForceNew) The name of the agent.
* `name_regex` - (Optional,ForceNew) A regex string to filter resulting chatbot agents by name.
* `ids` - (Optional,ForceNew,Computed) A list of chatbot agents IDs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of chatbot agents names.
* `agents` - A list of availability zones. Each element contains the following attributes:
  * `id` - ID of the agent.
  * `agent_id` - The agent id.
  * `agent_key` - Service space signature, which is used when PAAS interface specifies the service space.
  * `agent_name` - The agent Name.

