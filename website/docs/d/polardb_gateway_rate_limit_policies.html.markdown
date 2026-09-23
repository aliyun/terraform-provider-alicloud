---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_gateway_rate_limit_policies"
sidebar_current: "docs-alicloud-datasource-polardb-gateway-rate-limit-policies"
description: |-
  Provides a list of PolarDB AI Gateway Rate Limit Policies.
---

# alicloud_polardb_gateway_rate_limit_policies

Provides a list of PolarDB AI Gateway Rate Limit Policies.

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
data "alicloud_polardb_gateway_rate_limit_policies" "example" {
  gateway_id = alicloud_polardb_gateway.example.id
  scope_type = "Consumer"
}
```

## Argument Reference

* `gateway_id` - (Required, Available since v1.294.0) The ID of the PolarDB AI gateway.
* `ids` - (Optional, Available since v1.294.0) A list of rate limit policy IDs.
* `scope_type` - (Optional, Available since v1.294.0) The rate limit dimension. Valid values: `ConsumerGroup`, `Consumer`.
* `scope_ref_id` - (Optional, Available since v1.294.0) The consumer group or consumer ID.

## Attributes Reference

* `policies` - A list of rate limit policies. Each element contains the following attributes:
  * `id` - The policy ID.
  * `policy_type` - The policy type.
  * `scope_type` - The rate limit dimension.
  * `scope_ref_id` - The dimension object ID.
  * `rate_limit_rpm` - The maximum requests per minute.
  * `rate_limit_tpm` - The maximum tokens per minute.
  * `status` - The policy status.
  * `create_time` - The creation time.
  * `modify_time` - The last modification time.
