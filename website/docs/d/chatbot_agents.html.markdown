---
subcategory: "Chatbot"
layout: "alicloud"
page_title: "Alicloud: alicloud_chatbot_agents"
sidebar_current: "docs-alicloud-resource-chatbot-agents"
description: |-
  Provides a list of Chatbot Agents to the user.
---

# alicloud_chatbot_agents

This data source provides the Chatbot Agents of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.203.0+.

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
* `agent_name` - (Optional) The name of the agent.
* `name_regex` - (Optional) A regex string to filter resulting chatbot agents by name.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of chatbot agents IDs.
* `names` - A list of chatbot agents names.
* `agents` - A list of availability zones.

---

### `agents`

The agents supports the following attributes:

* `id` - ID of the agent.
* `agent_id` - The agent id.
* `agent_key` - Service space signature, which is used when PAAS interface specifies the service space.
* `agent_name` - The agent Name.
