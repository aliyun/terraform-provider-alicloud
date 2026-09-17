---
subcategory: "AgentLoop"
layout: "alicloud"
page_title: "Alicloud: alicloud_agentloop_context_store"
description: |-
  Provides a Alicloud AgentLoop Context Store resource.
---

# alicloud_agentloop_context_store

Provides a AgentLoop Context Store resource.

Context Store is used to store and manage the context data of agents in the Agent Space, such as experience and memory data imported from log data sources.

For information about AgentLoop Context Store and how to use it, see [What is Context Store](https://next.api.alibabacloud.com/document/AgentLoop/2026-05-20/CreateContextStore).

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
  context_store_name = "terraform_example"
  context_type       = "experience"
  description        = "terraform-example"
  config {
    service_names = ["svc-tf-acc"]
    metadata_field = {
      userId    = "user_id"
      sessionId = "session_id"
    }
    source {
      agent_space = alicloud_agentloop_agent_space.default.agent_space
      start_time  = "2026-01-01T00:00:00Z"
    }
  }
}
```

## Argument Reference

The following arguments are supported:
* `agent_space` - (Required, ForceNew) The name of the Agent Space to which the Context Store belongs. It must be 2 to 64 characters in length.
* `config` - (Optional, ForceNew, Set) The configuration of the Context Store. It is required when `context_type` is `experience`. `service_names` and `source` are read back from the API, so out-of-band changes are detected as drift. See [`config`](#config) below.
* `context_store_name` - (Required, ForceNew) The name of the Context Store. It must be 2 to 64 characters in length, can contain only lowercase letters, digits, and underscores (`_`), and must be unique within the Agent Space.
* `context_type` - (Required, ForceNew) The type of the Context Store, for example `experience` or `memory`.
* `description` - (Optional) The description of the Context Store.

-> **NOTE:** This resource only supports in-place update of `description`. All attributes under `config` and `context_type` are marked `ForceNew`, and changing them recreates the Context Store.

### `config`

The config supports the following:
* `metadata_field` - (Optional, ForceNew, Map) The mapping between business metadata fields and storage fields. The map key is the business field name and the map value (string) is the field name used in storage, for example `userId = "user_id"`. The Get API does not return the configured mapping, so this field is kept from the prior state and cannot be recovered by `terraform import`. After importing, add `lifecycle { ignore_changes = [config[0].metadata_field] }` to avoid replacing the Context Store.
* `service_names` - (Optional, ForceNew, List) The list of service names whose data is imported into the Context Store. It must be a non-empty list when `context_type` is `experience`.
* `source` - (Optional, ForceNew, Set) The log data source of the Context Store. See [`source`](#config-source) below.

### `config-source`

The source supports the following:
* `agent_space` - (Optional, ForceNew) The AgentSpace where the trace data source resides. If not specified, the AgentSpace of the Context Store is used. Cross-AgentSpace access is not supported: when specified, the value must match the AgentSpace of the Context Store.
* `start_time` - (Optional, ForceNew) The start time for data backfill, in ISO 8601 UTC format (`yyyy-MM-ddTHH:mm:ssZ`). If not specified, the current time is used.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Context Store. It formats as `<agent_space>:<context_store_name>`.
* `create_time` - The creation time of the Context Store, in ISO 8601 UTC format (`yyyy-MM-ddTHH:mm:ssZ`).
* `region_id` - The region ID of the Context Store.
* `status` - The status of the Context Store.
* `update_time` - The last update time of the Context Store, in ISO 8601 UTC format (`yyyy-MM-ddTHH:mm:ssZ`).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Context Store.
* `delete` - (Defaults to 5 mins) Used when delete the Context Store.
* `update` - (Defaults to 5 mins) Used when update the Context Store.

## Import

AgentLoop Context Store can be imported using the id, e.g.

```shell
$ terraform import alicloud_agentloop_context_store.example <agent_space>:<context_store_name>
```

-> **NOTE:** `config.metadata_field` cannot be imported. If it is set in the configuration, use `lifecycle { ignore_changes = [config[0].metadata_field] }` after importing.
