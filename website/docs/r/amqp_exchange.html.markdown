---
subcategory: "RabbitMQ (AMQP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_amqp_exchange"
sidebar_current: "docs-alicloud-resource-amqp-exchange"
description: |-
  Provides a Alicloud RabbitMQ (AMQP) Exchange resource.
---

# alicloud\_amqp\_exchange

Provides a RabbitMQ (AMQP) Exchange resource.

For information about RabbitMQ (AMQP) Exchange and how to use it, see [What is Exchange](https://www.alibabacloud.com/help/product/100989.html).

-> **NOTE:** Available in v1.128.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_amqp_virtual_host" "example" {
  instance_id       = "amqp-abc12345"
  virtual_host_name = "my-VirtualHost"
}
resource "alicloud_amqp_exchange" "example" {
  auto_delete_state = false
  exchange_name     = "my-Exchange"
  exchange_type     = "DIRECT"
  instance_id       = alicloud_amqp_virtual_host.example.instance_id
  internal          = false
  virtual_host_name = alicloud_amqp_virtual_host.example.virtual_host_name
}

```

## Argument Reference

The following arguments are supported:

* `alternate_exchange` - (Optional) The alternate exchange. An alternate exchange is configured for an existing exchange. It is used to receive messages that fail to be routed to queues from the existing exchange.
* `auto_delete_state` - (Required, ForceNew) Specifies whether the Auto Delete attribute is configured. Valid values:
  * true: The Auto Delete attribute is configured. If the last queue that is bound to an exchange is unbound, the exchange is automatically deleted.
  * false: The Auto Delete attribute is not configured. If the last queue that is bound to an exchange is unbound, the exchange is not automatically deleted.

* `exchange_name` - (Required, ForceNew) The name of the exchange. It must be 1 to 255 characters in length, and can contain only letters, digits, hyphens (-), underscores (_), periods (.), and at signs (@).
* `exchange_type` - (Required, ForceNew) The type of the exchange. Valid values:
  * FANOUT: An exchange of this type routes all the received messages to all the queues bound to this exchange. You can use a fanout exchange to broadcast messages.
  * DIRECT: An exchange of this type routes a message to the queue whose binding key is exactly the same as the routing key of the message.
  * TOPIC: This type is similar to the direct exchange type. An exchange of this type routes a message to one or more queues based on the fuzzy match or multi-condition match result between the routing key of the message and the binding keys of the current exchange.
  * HEADERS: Headers Exchange uses the Headers property instead of Routing Key for routing matching. 
    When binding Headers Exchange and Queue, set the key-value pair of the binding property; 
    when sending a message to the Headers Exchange, set the message's Headers property key-value pair and use the message Headers 
    The message is routed to the bound Queue by comparing the attribute key-value pair and the bound attribute key-value pair.
    
* `instance_id` - (Required, ForceNew) The ID of the instance.
* `internal` - (Required) Specifies whether an exchange is an internal exchange. Valid values:
  * false: The exchange is not an internal exchange.
  * true: The exchange is an internal exchange.
  
* `virtual_host_name` - (Required, ForceNew) The name of virtual host where an exchange resides.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Exchange. The value formats as `<instance_id>:<virtual_host_name>:<exchange_name>`.

## Import

RabbitMQ (AMQP) Exchange can be imported using the id, e.g.

```
$ terraform import alicloud_amqp_exchange.example <instance_id>:<virtual_host_name>:<exchange_name>
```
