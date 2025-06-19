---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_forwarding_rule"
sidebar_current: "docs-alicloud-resource-ga-forwarding-rule"
description: |-
  Provides a Alicloud Global Accelerator (GA) Forwarding Rule resource.
---

# alicloud_ga_forwarding_rule

Provides a Global Accelerator (GA) Forwarding Rule resource.

For information about Global Accelerator (GA) Forwarding Rule and how to use it, see [What is Forwarding Rule](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-createforwardingrules).

-> **NOTE:** Available since v1.120.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ga_forwarding_rule&exampleId=e41696bb-3f5b-3ed7-4dea-dd501dc13c6c42c46f9f&activeTab=example&spm=docs.r.ga_forwarding_rule.0.e41696bb3f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = var.region
}

variable "region" {
  default = "cn-hangzhou"
}

variable "name" {
  default = "terraform-example"
}

data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_ga_accelerator" "example" {
  duration            = 3
  spec                = "2"
  accelerator_name    = var.name
  auto_use_coupon     = false
  description         = var.name
  auto_renew_duration = "2"
  renewal_status      = "AutoRenewal"
}

resource "alicloud_ga_bandwidth_package" "example" {
  type                   = "Basic"
  bandwidth              = 20
  bandwidth_type         = "Basic"
  duration               = 1
  auto_pay               = true
  payment_type           = "Subscription"
  auto_use_coupon        = false
  bandwidth_package_name = var.name
  description            = var.name
}

resource "alicloud_ga_bandwidth_package_attachment" "example" {
  accelerator_id       = alicloud_ga_accelerator.example.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.example.id
}

resource "alicloud_ga_listener" "example" {
  accelerator_id  = alicloud_ga_bandwidth_package_attachment.example.accelerator_id
  client_affinity = "SOURCE_IP"
  description     = var.name
  name            = var.name
  protocol        = "HTTP"
  proxy_protocol  = true
  port_ranges {
    from_port = 60
    to_port   = 60
  }
}

resource "alicloud_eip_address" "example" {
  bandwidth            = "10"
  internet_charge_type = "PayByBandwidth"
}

resource "alicloud_ga_endpoint_group" "virtual" {
  accelerator_id = alicloud_ga_accelerator.example.id
  endpoint_configurations {
    endpoint                     = alicloud_eip_address.example.ip_address
    type                         = "PublicIp"
    weight                       = "20"
    enable_clientip_preservation = true
  }
  endpoint_group_region         = data.alicloud_regions.default.regions.0.id
  listener_id                   = alicloud_ga_listener.example.id
  description                   = var.name
  endpoint_group_type           = "virtual"
  endpoint_request_protocol     = "HTTPS"
  health_check_interval_seconds = 4
  health_check_path             = "/path"
  name                          = var.name
  threshold_count               = 4
  traffic_percentage            = 20
  port_overrides {
    endpoint_port = 80
    listener_port = 60
  }
}

resource "alicloud_ga_forwarding_rule" "example" {
  accelerator_id = alicloud_ga_accelerator.example.id
  listener_id    = alicloud_ga_listener.example.id
  rule_conditions {
    rule_condition_type = "Path"
    path_config {
      values = ["/testpathconfig"]
    }
  }
  rule_conditions {
    rule_condition_type = "Host"
    host_config {
      values = ["www.test.com"]
    }
  }
  rule_actions {
    order            = "40"
    rule_action_type = "ForwardGroup"
    forward_group_config {
      server_group_tuples {
        endpoint_group_id = alicloud_ga_endpoint_group.virtual.id
      }
    }
  }
  priority             = 2
  forwarding_rule_name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the Global Accelerator instance.
* `listener_id` - (Required, ForceNew) The ID of the listener.
* `priority` - (Optional, Int) Forwarding policy priority.
* `forwarding_rule_name` - (Optional) Forwarding policy name. The length of the name is 2-128 English or Chinese characters. It must start with uppercase and lowercase letters or Chinese characters. It can contain numbers, half width period (.), underscores (_) And dash (-).
* `rule_conditions` - (Required, Set) Forwarding condition list. See [`rule_conditions`](#rule_conditions) below.
* `rule_actions` - (Required, Set) Forward action. See [`rule_actions`](#rule_actions) below.

### `rule_actions`

The rule_actions supports the following:

* `order` - (Required, Int) Forwarding priority.
* `rule_action_type` - (Required) The type of the forwarding action. Valid values: `ForwardGroup`, `Redirect`, `FixResponse`, `Rewrite`, `AddHeader`, `RemoveHeader`, `Drop`.
* `rule_action_value` - (Optional, Available since v1.207.0) The value of the forwarding action type. For more information, see [How to use it](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-createforwardingrules).
* `forward_group_config` - (Optional, Set) Forwarding configuration. See [`forward_group_config`](#rule_actions-forward_group_config) below.
-> **NOTE:** From version 1.207.0, We recommend that you do not use `forward_group_config`, and we recommend that you use the `rule_action_type` and `rule_action_value` to configure forwarding actions.

### `rule_actions-forward_group_config`

The forward_group_config supports the following:

* `server_group_tuples` - (Required, Set) The information about the endpoint group. See [`server_group_tuples`](#rule_actions-forward_group_config-server_group_tuples) below.

### `rule_actions-forward_group_config-server_group_tuples`

The server_group_tuples supports the following:

* `endpoint_group_id` - (Required) The ID of the endpoint group.

### `rule_conditions`

The rule_conditions supports the following:

* `rule_condition_type` - (Required) The type of the forwarding conditions. Valid values: `Host`, `Path`, `RequestHeader`, `Query`, `Method`, `Cookie`, `SourceIP`. **NOTE:** From version 1.231.0, `rule_condition_type` can be set to `RequestHeader`, `Query`, `Method`, `Cookie`, `SourceIP`.
* `rule_condition_value` - (Optional, Available since v1.231.0) The value of the forwarding condition type. For more information, see [How to use it](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-createforwardingrules).
* `path_config` - (Optional, Set) The configuration of the path. See [`path_config`](#rule_conditions-path_config) below.
* `host_config` - (Optional, Set) The configuration of the domain name. See [`host_config`](#rule_conditions-host_config) below.
-> **NOTE:** From version 1.231.0, We recommend that you do not use `path_config` or `host_config`, and we recommend that you use the `rule_condition_type` and `rule_condition_value` to configure forwarding conditions.

### `rule_conditions-path_config`

The path_config supports the following:

* `values` - (Optional, List) The length of the path is 1-128 characters. It must start with a forward slash (/), and can only contain letters, numbers, dollar sign ($), dash (-), and underscores (_) , half width full stop (.), plus sign (+), forward slash (/), and (&), wavy line (~), at (@), half width colon (:), apostrophe ('). It supports asterisk (*) and half width question mark (?) as wildcards.

### `rule_conditions-host_config`

The host_config supports the following:

* `values` - (Optional, List) The domain name is 3-128 characters long, which can contain letters, numbers, dashes (-) and width period (.), and supports the use of asterisk (*) and width question mark (?) as wildcard characters.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Forwarding Rule. It formats as `<accelerator_id>:<listener_id>:<forwarding_rule_id>`.
* `forwarding_rule_id` - The ID of the Forwarding Rule.
* `forwarding_rule_status` - The status of the Forwarding Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Forwarding Rule.
* `update` - (Defaults to 3 mins) Used when update the Forwarding Rule.
* `delete` - (Defaults to 10 mins) Used when delete the Forwarding Rule.

## Import

Ga Forwarding Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_forwarding_rule.example <accelerator_id>:<listener_id>:<forwarding_rule_id>
```
