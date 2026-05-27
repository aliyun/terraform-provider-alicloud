---
subcategory: "RabbitMQ (AMQP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_amqp_open_source_account"
description: |-
  Provides a Alicloud RabbitMQ (AMQP) Open Source Account resource.
---

# alicloud_amqp_open_source_account

Provides a RabbitMQ (AMQP) Open Source Account resource.

An account under the open-source authentication and permission management system.  .

For information about RabbitMQ (AMQP) Open Source Account and how to use it, see [What is Open Source Account](https://next.api.alibabacloud.com/document/amqp-open/2019-12-12/CreateOpenSourceAccount).

-> **NOTE:** Available since v1.280.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "instance_name" {
  default = "example开源鉴权实例"
}

variable "user_name" {
  default = "Suhao123_"
}

variable "user_name_update" {
  default = "Suhao456_"
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
  instance_name         = var.instance_name
  support_tracing       = false
  payment_type          = "Subscription"
  renewal_duration_unit = "Month"
  instance_type         = "enterprise"
  queue_capacity        = "200"
  max_eip_tps           = "128"
  vpc_id                = alicloud_vpc.default.id
  vswitch_ids           = [alicloud_vswitch.default_b.id, alicloud_vswitch.default_g.id]
  security_group_id     = alicloud_security_group.default.id
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default_b" {
  vswitch_name = "${var.name}-b"
  cidr_block   = "172.16.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-hangzhou-b"
}

resource "alicloud_vswitch" "default_g" {
  vswitch_name = "${var.name}-g"
  cidr_block   = "172.16.1.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-hangzhou-g"
}

resource "alicloud_security_group" "default" {
  security_group_name = var.name
  vpc_id              = alicloud_vpc.default.id
}

resource "alicloud_amqp_open_source_account" "default" {
  user_name   = var.user_name
  description = var.user_name
  password    = var.user_name
  instance_id = alicloud_amqp_instance.CreateInstance.id
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) Description  
* `instance_id` - (Required, ForceNew) Instance ID  
* `password` - (Required) User password  
* `user_name` - (Required, ForceNew) User name  

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<user_name>:<instance_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Open Source Account.
* `delete` - (Defaults to 5 mins) Used when delete the Open Source Account.
* `update` - (Defaults to 5 mins) Used when update the Open Source Account.

## Import

RabbitMQ (AMQP) Open Source Account can be imported using the id, e.g.

```shell
$ terraform import alicloud_amqp_open_source_account.example <user_name>:<instance_id>
```
