---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_service_linked_role"
sidebar_current: "docs-alicloud-resource-event-bridge-service-linked-role"
description: |-
    Provides a resource to create the Event Bridge service-linked roles(SLR).
---

# alicloud\_event\_bridge\_service\_linked\_role

Using this data source can create Event Bridge service-linked roles(SLR). EventBridge may need to access another Alibaba Cloud service to implement a specific feature. In this case, EventBridge must assume a specific service-linked role, which is a Resource Access Management (RAM) role, to obtain permissions to access another Alibaba Cloud service. 

For information about Event Bridge service-linked roles(SLR) and how to use it, see [What is service-linked roles](https://www.alibabacloud.com/help/doc-detail/181425.htm).

-> **NOTE:** Available in v1.129.0+. After the version 1.142.0, the resource is renamed as `alicloud_event_bridge_service_linked_role`.


## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_event_bridge_service_linked_role&exampleId=a550113e-5065-36cb-7089-d8fe9975ac65fa176683&activeTab=example&spm=docs.r.event_bridge_service_linked_role.0.a550113e50&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_event_bridge_service_linked_role" "service_linked_role" {
  product_name = "AliyunServiceRoleForEventBridgeSourceRocketMQ"
}
```

## Argument Reference

The following arguments are supported:

* `product_name` - (Required, ForceNew) The product name for SLR. EventBridge can automatically create the following service-linked roles:
Event source related: `AliyunServiceRoleForEventBridgeSendToMNS`,`AliyunServiceRoleForEventBridgeSourceRocketMQ`, `AliyunServiceRoleForEventBridgeSourceActionTrail`, `AliyunServiceRoleForEventBridgeSourceRabbitMQ`
Target related: `AliyunServiceRoleForEventBridgeConnectVPC`, `AliyunServiceRoleForEventBridgeSendToFC`, `AliyunServiceRoleForEventBridgeSendToSMS`, `AliyunServiceRoleForEventBridgeSendToDirectMail`, `AliyunServiceRoleForEventBridgeSendToRabbitMQ`, `AliyunServiceRoleForEventBridgeSendToRocketMQ`

## Attributes Reference

* `id` - The ID of the DataSource. The value is same as `product_name`.

## Import

Event Bridge service-linked roles(SLR) can be imported using the id, e.g.

```shell
$ terraform import alicloud_event_bridge_service_linked_role.example <product_name>
```
