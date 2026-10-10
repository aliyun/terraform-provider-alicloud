---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_gateway_consumer_group"
sidebar_current: "docs-alicloud-resource-polardb-gateway-consumer-group"
description: |-
  Provides a PolarDB AI Gateway Consumer Group resource.
---

# alicloud_polardb_gateway_consumer_group

Provides a PolarDB AI Gateway Consumer Group resource.

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
resource "alicloud_polardb_gateway_consumer_group" "example" {
  gateway_id = alicloud_polardb_gateway.example.id
  name       = "production"
  nickname   = "Production consumers"
  is_default = "0"
}
```

## Argument Reference

* `gateway_id` - (Required, ForceNew, Available since v1.294.0) The ID of the PolarDB AI gateway.
* `name` - (Required, ForceNew, Available since v1.294.0) The name of the consumer group.
* `nickname` - (Optional, Available since v1.294.0) The nickname of the consumer group.
* `is_default` - (Optional, Available since v1.294.0) Whether the consumer group is the default group. Valid values: `0`, `1`. Default: `0`.

## Attributes Reference

* `id` - The resource ID in the format `<gateway_id>:<consumer_group_id>`.
* `consumer_group_id` - (Available since v1.294.0) The consumer group ID returned by the service API.
* `allowed_models` - (Available since v1.294.0) The models that the consumer group is allowed to access.
* `create_time` - (Available since v1.294.0) The creation time.
* `modify_time` - (Available since v1.294.0) The last modification time.

## Import

```shell
$ terraform import alicloud_polardb_gateway_consumer_group.example pg-abc:cg-abc
```
