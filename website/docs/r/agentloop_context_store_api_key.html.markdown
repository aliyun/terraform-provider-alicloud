---
subcategory: "AgentLoop"
layout: "alicloud"
page_title: "Alicloud: alicloud_agentloop_context_store_api_key"
description: |-
  Provides a Alicloud AgentLoop Context Store Api Key resource.
---

# alicloud_agentloop_context_store_api_key

Provides a AgentLoop Context Store Api Key resource.

Api Key is the access credential of the Context Store, which is used to access the context data of the Context Store.

For information about AgentLoop Context Store Api Key and how to use it, see [What is Context Store Api Key](https://next.api.alibabacloud.com/document/AgentLoop/2026-05-20/CreateContextStoreAPIKey).

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

resource "alicloud_agentloop_context_store" "default" {
  agent_space        = alicloud_agentloop_agent_space.default.agent_space
  context_store_name = "terraform_example_cs"
  context_type       = "experience"
  config {
    service_names = ["svc-tf-acc"]
    metadata_field = {
      userId    = "user_id"
      sessionId = "session_id"
    }
  }
}

resource "random_uuid" "default" {
}

resource "alicloud_agentloop_context_store_api_key" "default" {
  agent_space        = alicloud_agentloop_context_store.default.agent_space
  context_store_name = alicloud_agentloop_context_store.default.context_store_name
  name               = substr("${var.name}-${replace(random_uuid.default.result, "-", "")}", 0, 32)
}
```

## Argument Reference

The following arguments are supported:
* `agent_space` - (Required, ForceNew) The name of the Agent Space to which the Context Store belongs.
* `context_store_name` - (Required, ForceNew) The name of the Context Store.
* `name` - (Required, ForceNew) The name of the Api Key. **NOTE:** The backend reserves the name after an Api Key is deleted by design: a deleted name cannot be reused (recreating it returns `APIKeyNameConflict`), it is invisible to `List`, and the reservation cannot be cleared through the API. Always use a unique name, e.g. suffixed with `random_uuid` as shown in the example.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Context Store Api Key. It formats as `<agent_space>:<context_store_name>:<name>`.
* `api_key` - (Sensitive) The Api Key of the Context Store. The complete Api Key is only returned at creation time and is never refreshed from the Get API (which only returns a masked value), so it cannot be recovered by `terraform import`.
* `create_time` - The creation time of the Api Key, in ISO 8601 UTC format (`yyyy-MM-ddTHH:mm:ssZ`).
* `region_id` - The region ID of the Api Key.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Context Store Api Key.
* `delete` - (Defaults to 5 mins) Used when delete the Context Store Api Key.

## Import

AgentLoop Context Store Api Key can be imported using the id, e.g.

```shell
$ terraform import alicloud_agentloop_context_store_api_key.example <agent_space>:<context_store_name>:<name>
```

-> **NOTE:** `api_key` cannot be imported and is empty in the state of an imported Api Key. Recreate the Api Key if the complete value is required.
