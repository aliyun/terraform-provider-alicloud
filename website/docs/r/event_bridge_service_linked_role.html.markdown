---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_service_linked_role"
sidebar_current: "docs-alicloud-resource-event-bridge-service-linked-role"
description: |-
  Provides a Alicloud Event Bridge Service Linked Role resource.
---

# alicloud_event_bridge_service_linked_role

Provides a Event Bridge Service Linked Role resource.

For information about Event Bridge Service Linked Role and how to use it, see [What is Service Linked Role](https://www.alibabacloud.com/help/en/eventbridge/developer-reference/api-eventbridge-2020-04-01-createservicelinkedroleforproduct).

-> **NOTE:** Available since v1.129.0.

-> **NOTE:** From version 1.142.0, the resource is renamed as `alicloud_event_bridge_service_linked_role`.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_event_bridge_service_linked_role&exampleId=8fa138a7-b303-f1d0-26f5-a744abf8bcfbd42adae9&activeTab=example&spm=docs.r.event_bridge_service_linked_role.0.8fa138a7b3&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_event_bridge_service_linked_role" "default" {
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
$ terraform import alicloud_event_bridge_service_linked_role.example <product_name>
```
