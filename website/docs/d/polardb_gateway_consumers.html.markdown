---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_gateway_consumers"
sidebar_current: "docs-alicloud-datasource-polardb-gateway-consumers"
description: |-
  Provides a list of PolarDB AI Gateway Consumers.
---

# alicloud_polardb_gateway_consumers

Provides a list of PolarDB AI Gateway Consumers. API keys are deliberately omitted.

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
data "alicloud_polardb_gateway_consumers" "example" {
  gateway_id        = alicloud_polardb_gateway.example.id
  consumer_group_id = alicloud_polardb_gateway_consumer_group.example.consumer_group_id
}
```

## Argument Reference

* `gateway_id` - (Required, Available since v1.294.0) The ID of the PolarDB AI gateway.
* `ids` - (Optional, Available since v1.294.0) A list of consumer IDs used to filter results.
* `consumer_group_id` - (Optional, Available since v1.294.0) The consumer group ID used to filter results.

## Attributes Reference

* `consumers` - A list of consumers. Each element contains the following attributes:
  * `id` - The consumer ID.
  * `name` - The consumer name.
  * `consumer_group_id` - The consumer group ID.
  * `consumer_group_name` - The consumer group name.
  * `nickname` - The consumer nickname.
  * `allowed_models` - The models the consumer is allowed to access.
  * `month_to_date_cost_count` - The month-to-date cost count.
  * `lifetime_cost_count` - The lifetime cost count.
  * `month_to_date_token_count` - The month-to-date token count.
  * `lifetime_token_count` - The lifetime token count.
  * `create_time` - The creation time.
  * `modify_time` - The last modification time.
