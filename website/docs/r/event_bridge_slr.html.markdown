---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_slr"
sidebar_current: "docs-alicloud-resource-event-bridge-slr"
description: |-
    Provides a resource to create the Event Bridge service-linked roles(SLR).
---

# alicloud\_event\_bridge\_slr

Using this data source can create Event Bridge service-linked roles(SLR). EventBridge may need to access another Alibaba Cloud service to implement a specific feature. In this case, EventBridge must assume a specific service-linked role, which is a Resource Access Management (RAM) role, to obtain permissions to access another Alibaba Cloud service. 

For information about Event Bridge service-linked roles(SLR) and how to use it, see [What is service-linked roles](https://www.alibabacloud.com/help/doc-detail/181425.htm).

-> **NOTE:** Available in v1.129.0+

## Example Usage

```terraform
resource "alicloud_event_bridge_slr" "slr" {
    product_name = "AliyunServiceRoleForEventBridgeSourceRocketMQ"
}
```

## Argument Reference

The following arguments are supported:

* `product_name` - (Required, ForceNew) The product name for SLR. EventBridge can automatically create the following service-linked roles:
Event source related: `AliyunServiceRoleForEventBridgeSourceRocketMQ`, `AliyunServiceRoleForEventBridgeSourceActionTrail`, `AliyunServiceRoleForEventBridgeSourceRabbitMQ`
Target related: `AliyunServiceRoleForEventBridgeConnectVPC`, `AliyunServiceRoleForEventBridgeSendToFC`, `AliyunServiceRoleForEventBridgeSendToSMS`, `AliyunServiceRoleForEventBridgeSendToDirectMail`, `AliyunServiceRoleForEventBridgeSendToRabbitMQ`, `AliyunServiceRoleForEventBridgeSendToRocketMQ`

## Attributes Reference

* `id` - The ID of the DataSource. The value is same as `product_name`.
