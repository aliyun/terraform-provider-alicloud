---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_ascript"
sidebar_current: "docs-alicloud-resource-alb-ascript"
description: |-
  Provides a Alicloud Alb Ascript resource.
---

# alicloud_alb_ascript

Provides a Alb Ascript resource.

For information about Alb Ascript and how to use it, see [What is AScript](https://www.alibabacloud.com/help/en/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-createascripts).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alb_ascript&exampleId=d74c0c5f-2b1a-bc8c-1003-4229b402a07186b2b332&activeTab=example&spm=docs.r.alb_ascript.0.d74c0c5f2b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
* `listener_id` - (Required, ForceNew) Listener ID of script attribution
* `position` - (Required, ForceNew) Execution location of AScript.
* `ascript_name` - (Required) The name of AScript.
* `script_content` - (Required) The content of AScript.
* `enabled` - (Required) Whether scripts are enabled.
* `ext_attribute_enabled` - (Optional) Whether extension parameters are enabled.
* `ext_attributes` - (Optional) Extended attribute list. See [`ext_attributes`](#ext_attributes) below for details.

### `ext_attributes`

The ext_attributes supports the following:
* `attribute_key` - (Optional) The key of the extended attribute.
* `attribute_value` - (Optional) The value of the extended attribute.


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `load_balancer_id` - The ID of load balancer instance.
* `status` - The status of AScript.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Ascript.
* `delete` - (Defaults to 5 mins) Used when delete the Ascript.
* `update` - (Defaults to 5 mins) Used when update the Ascript.

## Import

Alb AScript can be imported using the id, e.g.

```shell
$terraform import alicloud_alb_ascript.example <id>
```