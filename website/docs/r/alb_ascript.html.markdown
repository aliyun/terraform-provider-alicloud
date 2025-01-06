---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_ascript"
description: |-
  Provides a Alicloud Application Load Balancer (ALB) A Script resource.
---

# alicloud_alb_ascript

Provides a Application Load Balancer (ALB) A Script resource.



For information about Application Load Balancer (ALB) A Script and how to use it, see [What is A Script](https://www.alibabacloud.com/help/en/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-createascripts).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_alb_zones" "default" {}
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  count        = 2
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = format("10.4.%d.0/24", count.index + 1)
  zone_id      = data.alicloud_alb_zones.default.zones[count.index].id
  vswitch_name = format("${var.name}_%d", count.index + 1)
}

resource "alicloud_alb_load_balancer" "default" {
  vpc_id                 = alicloud_vpc.default.id
  address_type           = "Internet"
  address_allocated_mode = "Fixed"
  load_balancer_name     = var.name
  load_balancer_edition  = "Standard"
  resource_group_id      = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  load_balancer_billing_config {
    pay_type = "PayAsYouGo"
  }
  tags = {
    Created = "TF"
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default.0.id
    zone_id    = data.alicloud_alb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default.1.id
    zone_id    = data.alicloud_alb_zones.default.zones.1.id
  }
}

resource "alicloud_alb_server_group" "default" {
  protocol          = "HTTP"
  vpc_id            = alicloud_vpc.default.id
  server_group_name = var.name
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  health_check_config {
    health_check_enabled = "false"
  }
  sticky_session_config {
    sticky_session_enabled = "false"
  }
  tags = {
    Created = "TF"
  }
}

resource "alicloud_alb_listener" "default" {
  load_balancer_id     = alicloud_alb_load_balancer.default.id
  listener_protocol    = "HTTP"
  listener_port        = 8081
  listener_description = var.name
  default_actions {
    type = "ForwardGroup"
    forward_group_config {
      server_group_tuples {
        server_group_id = alicloud_alb_server_group.default.id
      }
    }
  }
}

resource "alicloud_alb_ascript" "default" {
  script_content = "time()"
  position       = "RequestHead"
  ascript_name   = var.name
  enabled        = true
  listener_id    = alicloud_alb_listener.default.id
}
```

## Argument Reference

The following arguments are supported:
* `ascript_name` - (Required) AScript name.
* `dry_run` - (Optional, Available since v1.241.0) Whether to PreCheck only this request
* `enabled` - (Optional) Whether AScript is enabled.
* `ext_attribute_enabled` - (Optional) Whether extension parameters are enabled. When ExtAttributeEnabled is true, ExtAttributes must be set.
* `ext_attributes` - (Optional, List) Expand the list of attributes. When ExtAttributeEnabled is true, ExtAttributes must be set. See [`ext_attributes`](#ext_attributes) below.
* `listener_id` - (Required, ForceNew) Listener ID of script attribution
* `position` - (Required, ForceNew) Script execution location.
* `script_content` - (Required) AScript script content.

### `ext_attributes`

The ext_attributes supports the following:
* `attribute_key` - (Optional) Key to extend attribute
* `attribute_value` - (Optional) The value of the extended attribute

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - Script status

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the A Script.
* `delete` - (Defaults to 5 mins) Used when delete the A Script.
* `update` - (Defaults to 5 mins) Used when update the A Script.

## Import

Application Load Balancer (ALB) A Script can be imported using the id, e.g.

```shell
$ terraform import alicloud_alb_ascript.example <id>
```