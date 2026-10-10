---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_gateway_consumer_groups"
sidebar_current: "docs-alicloud-datasource-polardb-gateway-consumer-groups"
description: |-
  Provides a list of PolarDB AI Gateway Consumer Groups.
---

# alicloud_polardb_gateway_consumer_groups

Provides a list of PolarDB AI Gateway Consumer Groups.

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
data "alicloud_polardb_gateway_consumer_groups" "example" {
  gateway_id = alicloud_polardb_gateway.example.id
}
```

## Argument Reference

* `gateway_id` - (Required, Available since v1.294.0) The ID of the PolarDB AI gateway.
* `ids` - (Optional, Available since v1.294.0) A list of consumer group IDs used to filter results.

## Attributes Reference

* `groups` - A list of consumer groups. Each element contains the following attributes:
  * `id` - The consumer group ID.
  * `name` - The consumer group name.
  * `nickname` - The consumer group nickname.
  * `is_default` - Whether the group is the default group.
  * `allowed_models` - The models the group is allowed to access.
  * `create_time` - The creation time.
  * `modify_time` - The last modification time.
