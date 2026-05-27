---
subcategory: "RabbitMQ (AMQP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_amqp_open_source_permission"
description: |-
  Provides a Alicloud RabbitMQ (AMQP) Open Source Permission resource.
---

# alicloud_amqp_open_source_permission

Provides a RabbitMQ (AMQP) Open Source Permission resource.

Permissions in the open-source authentication and permission management system.

For information about RabbitMQ (AMQP) Open Source Permission and how to use it, see [What is Open Source Permission](https://next.api.alibabacloud.com/document/amqp-open/2019-12-12/CreateOpenSourcePermission).

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

variable "vhost" {
  default = "/"
}

variable "user_name" {
  default = "Suhao123_WithPer"
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


resource "alicloud_amqp_open_source_permission" "default" {
  write       = ".*"
  read        = ".*"
  vhost       = var.vhost
  user_name   = var.user_name
  instance_id = alicloud_amqp_instance.CreateInstance.id
  configure   = ".*"
}
```

## Argument Reference

The following arguments are supported:
* `configure` - (Optional) Permission configuration, such as .*
* `instance_id` - (Required, ForceNew) Instance ID
* `read` - (Optional) Read permission, such as .*
* `user_name` - (Required, ForceNew) Username
* `vhost` - (Optional, ForceNew, Computed) Vhost of the instance
* `write` - (Optional) Write permission, such as .*

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<user_name>:<vhost>:<instance_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Open Source Permission.
* `delete` - (Defaults to 5 mins) Used when delete the Open Source Permission.
* `update` - (Defaults to 5 mins) Used when update the Open Source Permission.

## Import

RabbitMQ (AMQP) Open Source Permission can be imported using the id, e.g.

```shell
$ terraform import alicloud_amqp_open_source_permission.example <user_name>:<vhost>:<instance_id>
```
