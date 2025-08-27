---
subcategory: "RabbitMQ (AMQP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_amqp_virtual_host"
description: |-
  Provides a Alicloud RabbitMQ (AMQP) Virtual Host resource.
---

# alicloud_amqp_virtual_host

Amqp Virtual Host.

Provides a RabbitMQ (AMQP) Virtual Host resource.

For information about RabbitMQ (AMQP) Virtual Host and how to use it, see [What is Virtual Host](https://www.alibabacloud.com/help/en/message-queue-for-rabbitmq/latest/createvirtualhost).

-> **NOTE:** Available since v1.126.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_amqp_virtual_host&exampleId=76e02692-42f8-7c52-e927-859a9352ba925dbeb20f&activeTab=example&spm=docs.r.amqp_virtual_host.0.76e0269242&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_amqp_instance" "default" {
  instance_name         = "${var.name}-${random_integer.default.result}"
  instance_type         = "enterprise"
  max_tps               = 3000
  max_connections       = 2000
  queue_capacity        = 200
  payment_type          = "Subscription"
  renewal_status        = "AutoRenewal"
  renewal_duration      = 1
  renewal_duration_unit = "Year"
  support_eip           = true
}

resource "alicloud_amqp_virtual_host" "default" {
  instance_id       = alicloud_amqp_instance.default.id
  virtual_host_name = "${var.name}-${random_integer.default.result}"
}
```

## Argument Reference

The following arguments are supported:
* `instance_id` - (Required, ForceNew) InstanceId.
* `virtual_host_name` - (Required, ForceNew) VirtualHostName.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Sasl User. It formats as `<instance_id>:<virtual_host_name>`.

## Timeouts

-> **NOTE:** Available since v1.259.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Virtual Host.
* `delete` - (Defaults to 5 mins) Used when delete the Virtual Host.

## Import

RabbitMQ (AMQP) Virtual Host can be imported using the id, e.g.

```shell
$ terraform import alicloud_amqp_virtual_host.example <instance_id>:<virtual_host_name>
```
