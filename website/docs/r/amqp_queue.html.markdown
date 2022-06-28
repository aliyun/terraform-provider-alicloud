---
subcategory: "RabbitMQ (AMQP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_amqp_queue"
sidebar_current: "docs-alicloud-resource-amqp-queue"
description: |-
  Provides a Alicloud RabbitMQ (AMQP) Queue resource.
---

# alicloud\_amqp\_queue

Provides a RabbitMQ (AMQP) Queue resource.

For information about RabbitMQ (AMQP) Queue and how to use it, see [What is Queue](https://www.alibabacloud.com/help/doc-detail/101631.htm).

-> **NOTE:** Available in v1.127.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_amqp_virtual_host" "example" {
  instance_id       = "amqp-abc12345"
  virtual_host_name = "my-VirtualHost"
}
resource "alicloud_amqp_queue" "example" {
  instance_id       = alicloud_amqp_virtual_host.example.instance_id
  queue_name        = "my-Queue"
  virtual_host_name = alicloud_amqp_virtual_host.example.virtual_host_name
}

```

## Argument Reference

The following arguments are supported:

* `auto_delete_state` - (Optional, ForceNew) Specifies whether the Auto Delete attribute is configured. Valid values:
  * true: The Auto Delete attribute is configured. The queue is automatically deleted after the last subscription from consumers to this queue is canceled. 
  * false: The Auto Delete attribute is not configured.
  
* `auto_expire_state` - (Optional) The validity period after which the queue is automatically deleted.
  If the queue is not accessed within a specified period of time, it is automatically deleted.
* `dead_letter_exchange` - (Optional) The dead-letter exchange. A dead-letter exchange is used to receive rejected messages. 
  If a consumer rejects a message that cannot be retried, this message is routed to a specified dead-letter exchange.
  Then, the dead-letter exchange routes the message to the queue that is bound to the dead-letter exchange.
* `dead_letter_routing_key` - (Optional) The dead letter routing key.
* `exclusive_state` - (Optional, ForceNew) Specifies whether the queue is an exclusive queue. Valid values:
  * true: The queue is an exclusive queue. It can be used only for the connection that declares the exclusive queue. After the connection is closed, the exclusive queue is automatically deleted.
  * false: The queue is not an exclusive queue.
  
* `instance_id` - (Required, ForceNew) The ID of the instance.
* `max_length` - (Optional) The maximum number of messages that can be stored in the queue.
  If this threshold is exceeded, the earliest messages that are routed to the queue are discarded.
* `maximum_priority` - (Optional) The highest priority supported by the queue. This parameter is set to a positive integer.
  Valid values: 0 to 255. Recommended values: 1 to 10
* `message_ttl` - (Optional) The message TTL of the queue.
  If the retention period of a message in the queue exceeds the message TTL of the queue, the message expires.
  Message TTL must be set to a non-negative integer, in milliseconds.
  For example, if the message TTL of the queue is 1000, messages survive for at most 1 second in the queue.
* `queue_name` - (Required, ForceNew) The name of the queue.
  The queue name must be 1 to 255 characters in length, and can contain only letters, digits, hyphens (-), underscores (_), periods (.), and at signs (@).
* `virtual_host_name` - (Required, ForceNew) The name of the virtual host.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Queue. The value formats as `<instance_id>:<virtual_host_name>:<queue_name>`.

## Import

RabbitMQ (AMQP) Queue can be imported using the id, e.g.

```
$ terraform import alicloud_amqp_queue.example <instance_id>:<virtual_host_name>:<queue_name>
```
