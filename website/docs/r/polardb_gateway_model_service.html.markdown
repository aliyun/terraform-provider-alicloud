---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_gateway_model_service"
sidebar_current: "docs-alicloud-resource-polardb-gateway-model-service"
description: |-
  Provides a PolarDB AI Gateway Model Service resource.
---

# alicloud_polardb_gateway_model_service

Provides a PolarDB AI Gateway Model Service resource.

-> **NOTE:** Available since v2.1.0.

-> **NOTE:** `api_key` is sensitive but is stored in Terraform state. Use an encrypted remote state backend and restrict access to the state.

## Example Usage

```terraform
resource "alicloud_polardb_gateway_model_service" "example" {
  gateway_id     = alicloud_polardb_gateway.example.id
  name           = "qwen"
  model_category = "text"
  protocol       = "openai"
  base_url       = "https://dashscope.aliyuncs.com/compatible-mode/v1"
  api_key        = var.model_service_api_key
  vendor         = "bailian"
}
```

## Argument Reference

* `gateway_id` - (Required, ForceNew, Available since v2.1.0) The ID of the PolarDB AI gateway.
* `name` - (Required, ForceNew, Available since v2.1.0) The name of the model service.
* `model_category` - (Required, Available since v2.1.0) The model category. Valid values: `text`, `embedding`, `rerank`.
* `protocol` - (Required, Available since v2.1.0) The upstream protocol. Valid values: `openai`, `anthropic`, `bailian`, `vllm`.
* `base_url` - (Required, Available since v2.1.0) The HTTP or HTTPS URL of the upstream model service.
* `api_key` - (Required, Sensitive, Available since v2.1.0) The API key of the upstream model service.
* `vendor` - (Required, ForceNew, Available since v2.1.0) The model service vendor.
* `input_cost_points_per_million` - (Optional, Available since v2.1.0) The input cost in points per million tokens. The service does not support clearing this value after it is set; removing it from the configuration retains the remote value.
* `output_cost_points_per_million` - (Optional, Available since v2.1.0) The output cost in points per million tokens. The service does not support clearing this value after it is set; removing it from the configuration retains the remote value.
* `request_cost_points` - (Optional, Available since v2.1.0) The cost in points per request. The service does not support clearing this value after it is set; removing it from the configuration retains the remote value.

## Attributes Reference

* `id` - The resource ID in the format `<gateway_id>:<model_service_id>`.
* `model_service_id` - (Available since v2.1.0) The model service ID returned by the service API.
* `status` - (Available since v2.1.0) The model service status.
* `create_time` - (Available since v2.1.0) The creation time.

## Import

The service API does not return the unmasked upstream credential, model category, or cost fields. Configure `api_key`, `model_category`, `input_cost_points_per_million`, `output_cost_points_per_million`, and `request_cost_points` in Terraform after import.

```shell
$ terraform import alicloud_polardb_gateway_model_service.example pg-abc:ms-abc
```
