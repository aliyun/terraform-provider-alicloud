---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_gateway_rate_limit_policy"
sidebar_current: "docs-alicloud-resource-polardb-gateway-rate-limit-policy"
description: |-
  Provides a PolarDB AI Gateway Rate Limit Policy resource.
---

# alicloud_polardb_gateway_rate_limit_policy

Provides a PolarDB AI Gateway Rate Limit Policy resource.

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
resource "alicloud_polardb_gateway_rate_limit_policy" "example" {
  gateway_id     = alicloud_polardb_gateway.example.id
  scope_type     = "Consumer"
  scope_ref_id   = alicloud_polardb_gateway_consumer.example.consumer_id
  rate_limit_rpm = "100"
  rate_limit_tpm = "10000"
}
```

## Argument Reference

* `gateway_id` - (Required, ForceNew, Available since v1.294.0) The ID of the PolarDB AI gateway.
* `scope_type` - (Required, ForceNew, Available since v1.294.0) The rate limit dimension. Valid values: `ConsumerGroup`, `Consumer`.
* `scope_ref_id` - (Required, ForceNew, Available since v1.294.0) The consumer group or consumer ID.
* `rate_limit_rpm` - (Required, Available since v1.294.0) The maximum number of requests per minute.
* `rate_limit_tpm` - (Required, Available since v1.294.0) The maximum number of tokens per minute.

## Attributes Reference

* `id` - The resource ID in the format `<gateway_id>:<policy_id>`.
* `policy_id` - (Available since v1.294.0) The rate limit policy ID.
* `policy_type` - (Available since v1.294.0) The policy type.
* `status` - (Available since v1.294.0) The policy status.
* `create_time` - (Available since v1.294.0) The creation time.
* `modify_time` - (Available since v1.294.0) The last modification time.

## Import

```shell
$ terraform import alicloud_polardb_gateway_rate_limit_policy.example pg-abc:policy-abc
```
