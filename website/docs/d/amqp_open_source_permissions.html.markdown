---
subcategory: "RabbitMQ (AMQP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_amqp_open_source_permissions"
sidebar_current: "docs-alicloud-datasource-amqp-open-source-permissions"
description: |-
  Provides a list of RabbitMQ (AMQP) Open Source Permission owned by an Alibaba Cloud account.
---

# alicloud_amqp_open_source_permissions

This data source provides RabbitMQ (AMQP) Open Source Permission available to the user.[What is Open Source Permission](https://next.api.alibabacloud.com/document/amqp-open/2019-12-12/CreateOpenSourcePermission)

-> **NOTE:** Available since v1.280.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "instance_name" {
  default = "测试开源鉴权实例"
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

data "alicloud_amqp_open_source_permissions" "default" {
  ids         = ["${alicloud_amqp_open_source_permission.default.id}"]
  instance_id = alicloud_amqp_instance.CreateInstance.id
  user_name   = var.user_name
}

output "alicloud_amqp_open_source_permission_example_id" {
  value = data.alicloud_amqp_open_source_permissions.default.permissions.0.id
}
```

## Argument Reference

The following arguments are supported:
* `instance_id` - (Required) Instance ID.
* `user_name` - (Required) Username.
* `enable_details` - (Optional) Default to `false`. Set it to `true` to query detailed attributes for each permission.
* `ids` - (Optional, Computed) A list of Open Source Permission IDs. The value is formulated as `<user_name>:<vhost>:<instance_id>`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Open Source Permission IDs.
* `permissions` - A list of Open Source Permission Entries. Each element contains the following attributes:
    * `configure` - Permission configuration, such as .
    * `instance_id` - Instance ID.
    * `read` - Read permission, such as .
    * `user_name` - Username.
    * `vhost` - Vhost of the instance.
    * `write` - Write permission, such as .
    * `id` - The ID of the resource supplied above.
