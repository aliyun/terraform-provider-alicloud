---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_slr"
sidebar_current: "docs-alicloud-resource-event-bridge-slr"
description: |-
  Provides a Alicloud Event Bridge Service Linked Role resource.
---

# alicloud\_event\_bridge\_slr

Provides a Event Bridge Service Linked Role resource.

`alicloud_event_bridge_slr` is an alias of `alicloud_event_bridge_service_linked_role`. The two resource types share the same schema, CRUD behavior and service-linked-role lifecycle; use either name interchangeably.

For information about Event Bridge Service Linked Role and how to use it, see [What is Service Linked Role](https://www.alibabacloud.com/help/en/eventbridge/developer-reference/api-eventbridge-2020-04-01-createservicelinkedroleforproduct).

-> **NOTE:** Available since v1.129.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_event_bridge_slr" "default" {
  product_name = "AliyunServiceRoleForEventBridgeSourceRocketMQ"
}
```

## Argument Reference

The following arguments are supported:

* `product_name` - (Required, ForceNew) The name of the cloud service or the name of the service-linked role with which the cloud service is associated. For more information, see [How to use it](https://www.alibabacloud.com/help/en/eventbridge/developer-reference/api-eventbridge-2020-04-01-createservicelinkedroleforproduct).

## Attributes Reference

* `id` - The resource ID in terraform of Service Linked Role.

## Timeouts

-> **NOTE:** Available since v1.252.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Service Linked Role.
* `delete` - (Defaults to 1 mins) Used when delete the Service Linked Role.

## Import

Event Bridge Service Linked Role can be imported using the id, e.g.

```shell
$ terraform import alicloud_event_bridge_slr.example <product_name>
```
