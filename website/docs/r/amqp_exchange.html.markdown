---
subcategory: "RabbitMQ (AMQP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_amqp_exchange"
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
  <a href="https://api.aliyun.com/terraform?resource=alicloud_amqp_exchange&exampleId=3c0471e9-733b-9de7-7c2a-6bdbc2bd31b52c5640b5&activeTab=example&spm=docs.r.amqp_exchange.0.3c0471e973&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-shanghai"
}

variable "name" {
  default = "tf-example"
}

variable "virtual_host_name" {
  default = "/"
}

resource "alicloud_amqp_instance" "CreateInstance" {
  renewal_duration      = "1"
  max_tps               = "3000"
  period_cycle          = "Month"
  max_connections       = "2000"
  support_eip           = true
  auto_renew            = false
  renewal_status        = "AutoRenewal"
  period                = "12"
  instance_name         = var.name
  support_tracing       = false
  payment_type          = "Subscription"
  renewal_duration_unit = "Month"
  instance_type         = "enterprise"
  queue_capacity        = "200"
  max_eip_tps           = "128"
  storage_size          = "0"
}

resource "alicloud_amqp_exchange" "default" {
  virtual_host_name  = var.virtual_host_name
  instance_id        = alicloud_amqp_instance.CreateInstance.id
  internal           = "true"
  auto_delete_state  = "false"
  exchange_name      = var.name
  exchange_type      = "X_CONSISTENT_HASH"
  alternate_exchange = "bakExchange"
  x_delayed_type     = "DIRECT"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_amqp_exchange&spm=docs.r.amqp_exchange.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `alternate_exchange` - (Optional) The alternate exchange. An alternate exchange is used to receive messages that fail to be routed to queues from the current exchange.
* `auto_delete_state` - (Required, ForceNew) Specifies whether to automatically delete the exchange. Valid values:

  - `true`: If the last queue that is bound to the exchange is unbound, the exchange is automatically deleted.
  - `false`: If the last queue that is bound to the exchange is unbound, the exchange is not automatically deleted.
* `exchange_name` - (Required, ForceNew) The name of the exchange that you want to create. The exchange name must meet the following conventions:

  - The name must be 1 to 255 characters in length, and can contain only letters, digits, hyphens (-), underscores (\_), periods (.), number signs (#), forward slashes (/), and at signs (@).
  - After the exchange is created, you cannot change its name. If you want to change its name, delete the exchange and create another exchange.
* `exchange_type` - (Required, ForceNew) The Exchange type. Value:
  - `DIRECT`: This type of Routing rule routes messages to a Queue whose Binding Key matches the Routing Key.
  - `TOPIC`: This type is similar to the DIRECT type. It uses Routing Key pattern matching and string comparison to route messages to the bound Queue.
  - `FANOUT`: This type of routing rule is very simple. It routes all messages sent to the Exchange to all queues bound to it, which is equivalent to the broadcast function.
  - `HEADERS`: This type is similar to the DIRECT type. Headers Exchange uses the Headers attribute instead of Routing Key for route matching. When binding Headers Exchange and Queue, the Key-value pair of the bound attribute is set. When sending a message to Headers Exchange, the Headers attribute Key-value pair of the message is set, and the message is routed to the bound Queue by comparing the Headers attribute Key-value pair with the bound attribute Key-value pair.
  - `X_delayed_message`: By declaring this type of Exchange, you can customize the Header attribute x-delay of the message to specify the delivery delay time period, in milliseconds. Messages will be delivered to the corresponding Queue after the time period defined in the x-delay according to the routing rules. The routing rule depends on the Exchange route type specified in the x-delayed-type.
  - `X_CONSISTENT_HASH`: The x-consistent-hash Exchange allows you to Hash the Routing Key or Header value and use the consistent hashing algorithm to route messages to different queues.
* `instance_id` - (Required, ForceNew) The ID of the ApsaraMQ for RabbitMQ instance whose exchange you want to delete.
* `internal` - (Required) Specifies whether the exchange is an internal exchange. Valid values:

  - `false`
  - `true`
* `virtual_host_name` - (Required, ForceNew) The name of the vhost to which the exchange that you want to create belongs.
* `x_delayed_type` - (Optional, Available since v1.249.0) RabbitMQ supports the x-delayed-message Exchange. By declaring this type of Exchange, you can customize the x-delay header attribute to specify the delay period for message delivery, measured in milliseconds. The message will be delivered to the corresponding Queue after the period defined in x-delay. The routing rules are determined by the type of Exchange specified in x-delayed-type.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<virtual_host_name>:<exchange_name>`.
* `create_time` - CreateTime

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Exchange.
* `delete` - (Defaults to 5 mins) Used when delete the Exchange.

## Import

RabbitMQ (AMQP) Exchange can be imported using the id, e.g.

```shell
$ terraform import alicloud_amqp_exchange.example <instance_id>:<virtual_host_name>:<exchange_name>
```