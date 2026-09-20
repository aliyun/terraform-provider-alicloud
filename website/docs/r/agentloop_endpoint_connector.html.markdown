---
subcategory: "AgentLoop"
layout: "alicloud"
page_title: "Alicloud: alicloud_agentloop_endpoint_connector"
description: |-
  Provides a Alicloud AgentLoop Endpoint Connector resource.
---

# alicloud_agentloop_endpoint_connector

Provides a AgentLoop Endpoint Connector resource.

Endpoint Connector is used to connect the Agent Space to external endpoints, such as model services, with credentials and custom headers.

For information about AgentLoop Endpoint Connector and how to use it, see [What is Endpoint Connector](https://next.api.alibabacloud.com/document/AgentLoop/2026-05-20/CreateEndpointConnector).

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

resource "alicloud_agentloop_endpoint_connector" "default" {
  agent_space = alicloud_agentloop_agent_space.default.agent_space
  name        = var.name
  type        = "model_service"
  endpoint    = "https://dashscope.aliyuncs.com/compatible-mode/v1"
  credential = {
    api_key = "sk-your-api-key"
  }
  alias       = "${var.name}-alias"
  description = "terraform-example"
  headers {
    key   = "X-Header-1"
    value = "value1"
  }
  properties = {
    property1 = "value1"
  }
  tags = ["tag1"]
}
```

## Argument Reference

The following arguments are supported:
* `agent_space` - (Required, ForceNew) The name of the Agent Space to which the Endpoint Connector belongs.
* `alias` - (Optional) The alias of the Endpoint Connector.
* `credential` - (Required, Sensitive, Map) The credential used to access the endpoint, e.g. `{"api_key" = "sk-..."}`. The create and update APIs accept different credential shapes: create takes the bare map as-is, while update requires a `providerType` (`plain` or `kms`) plus the key material wrapped in a JSON-encoded `sensitiveInfo` string. When the map does not contain `providerType`, the provider automatically converts the `api_key` shorthand into `{"providerType" = "plain", "sensitiveInfo" = "{\"apiKey\": \"...\"}"}` on updates. If you manage the credential through KMS, include `"providerType" = "kms"` in the map explicitly so it is passed through unchanged.
* `description` - (Optional) The description of the Endpoint Connector.
* `endpoint` - (Required) The endpoint URL to connect to.
* `headers` - (Optional) The custom request headers. See [`headers`](#headers) below.
* `name` - (Required) The name of the Endpoint Connector.
* `properties` - (Optional, Map) The extended properties of the Endpoint Connector.
* `tags` - (Optional, ForceNew, List) The tags of the Endpoint Connector.
* `type` - (Required, ForceNew) The type of the Endpoint Connector.

### `headers`

The headers supports the following:
* `key` - (Optional) The key of the header.
* `value` - (Optional) The value of the header.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Endpoint Connector. It formats as `<agent_space>:<connector_id>`.
* `connector_id` - The ID of the Endpoint Connector.
* `created_at` - The creation time of the Endpoint Connector.
* `region_id` - The region ID of the Endpoint Connector.
* `updated_at` - The last update time of the Endpoint Connector.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Endpoint Connector.
* `delete` - (Defaults to 5 mins) Used when delete the Endpoint Connector.
* `update` - (Defaults to 5 mins) Used when update the Endpoint Connector.

## Import

AgentLoop Endpoint Connector can be imported using the id, e.g.

```shell
$ terraform import alicloud_agentloop_endpoint_connector.example <agent_space>:<connector_id>
```
