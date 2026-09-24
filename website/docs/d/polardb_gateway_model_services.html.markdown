---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_gateway_model_services"
sidebar_current: "docs-alicloud-datasource-polardb-gateway-model-services"
description: |-
  Provides a list of PolarDB AI Gateway model services.
---

# alicloud_polardb_gateway_model_services

Provides a list of PolarDB AI Gateway model services. API keys are intentionally not exported.

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
data "alicloud_polardb_gateway_model_services" "example" {
  gateway_id = "pg-abc"
  protocol   = "openai"
}
```

## Argument Reference

* `gateway_id` - (Required, Available since v1.294.0) The gateway ID.
* `ids` - (Optional, Available since v1.294.0) A list of model service IDs.
* `name` - (Optional, Available since v1.294.0) The model service name.
* `model_category` - (Optional, Available since v1.294.0) The model category.
* `protocol` - (Optional, Available since v1.294.0) The protocol.
* `status` - (Optional, Available since v1.294.0) The status.

## Attributes Reference

* `services` - A list of model services. Each item contains the following attributes:
  * `id` - The model service ID.
  * `name` - The model service name.
  * `model_category` - The model category.
  * `protocol` - The upstream protocol.
  * `status` - The model service status.
  * `base_url` - The upstream service URL.
  * `vendor` - The model service vendor.
  * `create_time` - The creation time.
  * `input_cost_points_per_million` - The input cost points.
  * `output_cost_points_per_million` - The output cost points.
  * `request_cost_points` - The per-request cost points.
