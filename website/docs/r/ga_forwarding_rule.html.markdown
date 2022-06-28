---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_forwarding_rule"
sidebar_current: "docs-alicloud-resource-ga-forwarding-rule"
description: |-
  Provides a Alicloud Global Accelerator (GA) Forwarding Rule resource.
---

# alicloud\_ga\_forwarding\_rule

Provides a Global Accelerator (GA) Forwarding Rule resource.

For information about Global Accelerator (GA) Forwarding Rule and how to use it, see [What is Forwarding Rule](https://www.alibabacloud.com/help/zh/doc-detail/205815.htm).

-> **NOTE:** Available in v1.120.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ga_accelerator" "example" {
  duration            = 3
  spec                = "2"
  accelerator_name    = "ga-tf"
  auto_use_coupon     = false
  description         = "ga-tf description"
  auto_renew_duration = "2"
  renewal_status      = "AutoRenewal"
}

resource "alicloud_ga_bandwidth_package" "example" {
  type           = "Basic"
  bandwidth      = 20
  bandwidth_type = "Basic"
  duration       = 1
  timeouts {
    create = "5m"
  }
  auto_pay               = true
  payment_type           = "Subscription"
  billing_type           = "PayByTraffic"
  ratio                  = 40
  auto_use_coupon        = false
  bandwidth_package_name = "bandwidth_package_name_tf"
  description            = "bandwidth_package_name_tf_description"

}

resource "alicloud_ga_bandwidth_package_attachment" "example" {
  accelerator_id       = alicloud_ga_accelerator.example.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.example.id
}


resource "alicloud_ga_listener" "example" {
  depends_on     = [alicloud_ga_bandwidth_package_attachment.example]
  accelerator_id = alicloud_ga_accelerator.example.id
  port_ranges {
    from_port = 60
    to_port   = 60
  }
  client_affinity = "SOURCE_IP"
  description     = "alicloud_ga_listener_description"
  name            = "alicloud_ga_listener_tf"
  protocol        = "HTTP"
  proxy_protocol  = true
}


resource "alicloud_ga_ip_set" "example" {
  depends_on           = [alicloud_ga_bandwidth_package_attachment.example]
  accelerate_region_id = "cn-shanghai"
  accelerator_id       = alicloud_ga_accelerator.example.id
  bandwidth            = "20"
}

resource "alicloud_eip_address" "example" {
  bandwidth            = "10"
  internet_charge_type = "PayByBandwidth"
}

resource "alicloud_ga_endpoint_group" "default" {
  accelerator_id = alicloud_ga_accelerator.example.id
  endpoint_configurations {
    endpoint                     = alicloud_eip_address.example.ip_address
    type                         = "PublicIp"
    weight                       = "20"
    enable_clientip_preservation = true
  }
  endpoint_group_region         = "cn-shanghai"
  listener_id                   = alicloud_ga_listener.example.id
  description                   = "alicloud_ga_endpoint_group_description"
  endpoint_group_type           = "default"
  endpoint_request_protocol     = "HTTPS"
  health_check_interval_seconds = 4
  health_check_path             = "/path"
  name                          = "alicloud_ga_endpoint_group_name"
  threshold_count               = 4
  traffic_percentage            = 20
  port_overrides {
    endpoint_port = 80
    listener_port = 60
  }
}

resource "alicloud_ga_endpoint_group" "virtual" {
  accelerator_id = alicloud_ga_accelerator.example.id
  endpoint_configurations {
    endpoint                     = alicloud_eip_address.example.ip_address
    type                         = "PublicIp"
    weight                       = "20"
    enable_clientip_preservation = true
  }
  endpoint_group_region = "cn-shanghai"
  listener_id           = alicloud_ga_listener.example.id

  description                   = "alicloud_ga_endpoint_group_description"
  endpoint_group_type           = "virtual"
  endpoint_request_protocol     = "HTTPS"
  health_check_interval_seconds = 4
  health_check_path             = "/path"
  name                          = "alicloud_ga_endpoint_group_name"
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
        endpoint_group_id = alicloud_ga_endpoint_group.default.id
      }
    }
  }
  priority             = 2
  forwarding_rule_name = "forwarding_rule_name"
}
```

## Argument Reference

The following arguments are supported:

* `priority` - (Optional) Forwarding policy priority.
* `forwarding_rule_name` - (Optional) Forwarding policy name. The length of the name is 2-128 English or Chinese characters. It must start with uppercase and lowercase letters or Chinese characters. It can contain numbers, half width period (.), underscores (_) And dash (-).
* `rule_conditions` - (Required) Forwarding condition list.
* `rule_actions` - (Required) Forward action.
* `accelerator_id` - (Required, ForceNew) The ID of the Global Accelerator instance.
* `listener_id` - (Required, ForceNew) The ID of the listener.

#### Block rule_conditions

The rule_conditions supports the following:

* `rule_condition_type` (Required) Forwarding condition type. Valid value: `Host`, `Path`. 
* `path_config` (Optional) Path configuration information.
* `host_config` (Optional) Domain name configuration information.

#### Block path_config

The path_config supports the following:

* `values` (Optional) The length of the path is 1-128 characters. It must start with a forward slash (/), and can only contain letters, numbers, dollar sign ($), dash (-), and underscores (_) , half width full stop (.), plus sign (+), forward slash (/), and (&), wavy line (~), at (@), half width colon (:), apostrophe ('). It supports asterisk (*) and half width question mark (?) as wildcards.

#### Block host_config

The host_config supports the following:

* `values` (Optional) The domain name is 3-128 characters long, which can contain letters, numbers, dashes (-) and width period (.), and supports the use of asterisk (*) and width question mark (?) as wildcard characters.

#### Block rule_actions

The rule_actions supports the following:

* `order` (Required) Forwarding priority.
* `rule_action_type` (Required) Forward action type. Default: forwardgroup.
* `forward_group_config` (Required) Forwarding configuration. 

#### Block forward_group_config

The forward_group_config supports the following:

* `server_group_tuples` (Required) Terminal node group configuration.

#### Block server_group_tuples

The server_group_tuples supports the following:

* `endpoint_group_id` (Required) Terminal node group ID.
                                                                        
## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Forwarding Rule. The value is formate as `<accelerator_id>:<listener_id>:<forwarding_rule_id>`.
* `forwarding_rule_id` - Forwarding Policy ID.
* `forwarding_rule_status` - Forwarding Policy Status.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Endpoint Forwarding Rule.
* `delete` - (Defaults to 3 mins) Used when delete the Endpoint Forwarding Rule.
* `update` - (Defaults to 3 mins) Used when update the Endpoint Forwarding Rule.

## Import

Ga Forwarding Rule can be imported using the id, e.g.

```
$ terraform import alicloud_ga_forwarding_rule.example <id>
```
