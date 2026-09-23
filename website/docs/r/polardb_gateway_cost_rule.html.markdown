---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_gateway_cost_rule"
sidebar_current: "docs-alicloud-resource-polardb-gateway-cost-rule"
description: |-
  Provides a PolarDB AI Gateway Cost Rule resource.
---

# alicloud_polardb_gateway_cost_rule

Provides a PolarDB AI Gateway Cost Rule resource.

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
resource "alicloud_polardb_gateway_cost_rule" "example" {
  gateway_id                     = alicloud_polardb_gateway.example.id
  model_name                     = "qwen-plus"
  model_service_id               = alicloud_polardb_gateway_model_service.example.model_service_id
  input_cost_points_per_million  = "10"
  output_cost_points_per_million = "20"
}
```

## Argument Reference

* `gateway_id` - (Required, ForceNew, Available since v1.294.0) The ID of the PolarDB AI gateway.
* `model_name` - (Required, ForceNew, Available since v1.294.0) The model name.
* `model_service_id` - (Required, Available since v1.294.0) The model service ID returned by the service API, without the gateway portion of the Terraform composite ID.
* `input_cost_points_per_million` - (Optional, Available since v1.294.0) Input-token cost points. Default: `0`.
* `output_cost_points_per_million` - (Optional, Available since v1.294.0) Output-token cost points. Default: `0`.
* `cache_cost_points_per_million` - (Optional, Available since v1.294.0) Cached-token cost points. Default: `0`.

## Attributes Reference

* `id` - The resource ID in the format `<gateway_id>:<cost_rule_id>`.
* `cost_rule_id` - (Available since v1.294.0) The cost rule ID returned by the service API.
* `create_time` - (Available since v1.294.0) The creation time.
* `modify_time` - (Available since v1.294.0) The last modification time.

## Import

```shell
$ terraform import alicloud_polardb_gateway_cost_rule.example pg-abc:cost-rule-abc
```
