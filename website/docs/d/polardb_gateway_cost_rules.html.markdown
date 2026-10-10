---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_gateway_cost_rules"
sidebar_current: "docs-alicloud-datasource-polardb-gateway-cost-rules"
description: |-
  Provides a list of PolarDB AI Gateway cost rules.
---

# alicloud_polardb_gateway_cost_rules

Provides a list of PolarDB AI Gateway cost rules.

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
data "alicloud_polardb_gateway_cost_rules" "example" {
  gateway_id       = "pg-abc"
  model_service_id = "ms-abc"
}
```

## Argument Reference

* `gateway_id` - (Required, Available since v1.294.0) The gateway ID.
* `model_name` - (Optional, Available since v1.294.0) The model name.
* `model_service_id` - (Optional, Available since v1.294.0) The model service ID.

## Attributes Reference

* `rules` - A list of cost rules. Each item contains the following attributes:
  * `id` - The cost rule ID.
  * `model_name` - The model name.
  * `model_service_id` - The model service ID.
  * `input_cost_points_per_million` - The input cost points.
  * `output_cost_points_per_million` - The output cost points.
  * `cache_cost_points_per_million` - The cached-token cost points.
  * `create_time` - The creation time.
  * `modify_time` - The last modification time.
