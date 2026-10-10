---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_gateway_model_apis"
sidebar_current: "docs-alicloud-datasource-polardb-gateway-model-apis"
description: |-
  Provides a list of PolarDB AI Gateway model APIs.
---

# alicloud_polardb_gateway_model_apis

Provides a list of PolarDB AI Gateway model APIs.

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
data "alicloud_polardb_gateway_model_apis" "example" {
  gateway_id = "pg-abc"
  status     = "Enable"
}
```

## Argument Reference

* `gateway_id` - (Required, Available since v1.294.0) The gateway ID.
* `ids` - (Optional, Available since v1.294.0) A list of model API IDs.
* `name` - (Optional, Available since v1.294.0) The model API name.
* `model_category` - (Optional, Available since v1.294.0) The model category.
* `path_prefix` - (Optional, Available since v1.294.0) The path prefix.
* `protocol` - (Optional, Available since v1.294.0) The protocol.
* `status` - (Optional, Available since v1.294.0) The status.

## Attributes Reference

* `apis` - A list of model APIs. Each item contains the following attributes:
  * `id` - The model API ID.
  * `name` - The model API name.
  * `model_category` - The model category.
  * `path_prefix` - The API path prefix.
  * `protocol` - The model API protocol.
  * `status` - The model API status.
  * `record_input` - The input recording configuration.
  * `record_output` - The output recording configuration.
  * `route_rules` - The routing rules as JSON.
  * `force_model` - The forced model.
  * `invoke_endpoint` - The invocation endpoint.
  * `create_time` - The creation time.
