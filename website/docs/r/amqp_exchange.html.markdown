---
subcategory: "RabbitMQ (AMQP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_amqp_exchange"
sidebar_current: "docs-alicloud-resource-amqp-exchange"
description: |-
  Provides a Alicloud RabbitMQ (AMQP) Exchange resource.
---

# alicloud_amqp_exchange

Provides a RabbitMQ (AMQP) Exchange resource.

For information about RabbitMQ (AMQP) Exchange and how to use it, see [What is Exchange](https://www.alibabacloud.com/help/en/message-queue-for-rabbitmq/latest/createexchange).

-> **NOTE:** Available since v1.128.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_amqp_exchange&exampleId=5de3d9e7-a155-7c82-d84c-5f01dc23c782c0222634&activeTab=example&spm=docs.r.amqp_exchange.0.5de3d9e7a1&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-shanghai"
}

resource "alicloud_amqp_instance" "default" {
  instance_type  = "professional"
  max_tps        = 1000
  queue_capacity = 50
  support_eip    = true
  max_eip_tps    = 128
  payment_type   = "Subscription"
  period         = 1
}

resource "alicloud_amqp_virtual_host" "default" {
  instance_id       = alicloud_amqp_instance.default.id
  virtual_host_name = "tf-example"
}

resource "alicloud_amqp_exchange" "default" {
  auto_delete_state = false
  exchange_name     = "tf-example"
  exchange_type     = "DIRECT"
  instance_id       = alicloud_amqp_instance.default.id
  internal          = false
  virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
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

```shell
$ terraform import alicloud_amqp_exchange.example <instance_id>:<virtual_host_name>:<exchange_name>
```
