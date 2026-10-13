---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_gateway_model_api"
sidebar_current: "docs-alicloud-resource-polardb-gateway-model-api"
description: |-
  Provides a PolarDB AI Gateway Model API resource.
---

# alicloud_polardb_gateway_model_api

Provides a PolarDB AI Gateway Model API resource.

-> **NOTE:** Available since v2.1.0.

## Example Usage

```terraform
resource "alicloud_polardb_gateway_model_api" "example" {
  gateway_id     = alicloud_polardb_gateway.example.id
  name           = "chat"
  model_category = "text"
  path_prefix    = "/chat"
  protocol       = "openai"
  route_rules = jsonencode([{
    RuleName  = "primary"
    Providers = [{ ModelServiceName = alicloud_polardb_gateway_model_service.example.name, Weight = "100" }]
  }])
}
```

## Argument Reference

* `gateway_id` - (Required, ForceNew, Available since v2.1.0) The ID of the PolarDB AI gateway.
* `name` - (Required, ForceNew, Available since v2.1.0) The name of the model API.
* `model_category` - (Required, Available since v2.1.0) The model category. Valid values: `text`, `embedding`, `rerank`.
* `path_prefix` - (Required, Available since v2.1.0) The API path prefix.
* `protocol` - (Required, Available since v2.1.0) The protocol. Valid values: `openai`, `anthropic`, `bailian`, `vllm`.
* `record_input` - (Optional, Available since v2.1.0) Whether or how much request input is recorded for billing. The service does not support clearing this value after it is set; removing it from the configuration retains the remote value.
* `record_output` - (Optional, Available since v2.1.0) Whether or how much response output is recorded for billing. The service does not support clearing this value after it is set; removing it from the configuration retains the remote value.
* `route_rules` - (Required, Available since v2.1.0) A non-empty JSON array containing routing rules. JSON formatting and object key order are normalized in state.
* `force_model` - (Optional, ForceNew, Available since v2.1.0) The model to which requests are forcibly routed.

## Attributes Reference

* `id` - The resource ID in the format `<gateway_id>:<model_api_id>`.
* `model_api_id` - (Available since v2.1.0) The model API ID returned by the service API.
* `invoke_endpoint` - (Available since v2.1.0) The invocation endpoint.
* `status` - (Available since v2.1.0) The model API status.
* `create_time` - (Available since v2.1.0) The creation time.

## Import

The service API does not return `route_rules`, `force_model`, `record_input`, or `record_output`. Configure these values in Terraform after import.

```shell
$ terraform import alicloud_polardb_gateway_model_api.example pg-abc:mi-abc
```
