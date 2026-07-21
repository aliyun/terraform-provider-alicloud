---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_rules"
sidebar_current: "docs-alicloud-datasource-alb-rules"
description: |-
  Provides a list of Alb Rules to the user.
---

# alicloud_alb_rules

This data source provides the Alb Rules of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.133.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_alb_zones" "default" {
}

data "alicloud_resource_manager_resource_groups" "default" {
}

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
  listener_port        = 80
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

resource "alicloud_alb_rule" "default" {
  rule_name   = var.name
  listener_id = alicloud_alb_listener.default.id
  priority    = "555"
  rule_conditions {
    cookie_config {
      values {
        key   = "created"
        value = "tf"
      }
    }
    type = "Cookie"
  }

  rule_actions {
    forward_group_config {
      server_group_tuples {
        server_group_id = alicloud_alb_server_group.default.id
      }
    }
    order = "9"
    type  = "ForwardGroup"
  }
}

data "alicloud_alb_rules" "ids" {
  ids = [alicloud_alb_rule.default.id]
}

output "alb_rule_id" {
  value = data.alicloud_alb_rules.ids.rules.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List)  A list of Rule IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Rule name.
* `load_balancer_ids` - (Optional, ForceNew, List) The load balancer ids.
* `listener_ids` - (Optional, ForceNew, List) The listener ids.
* `rule_ids` - (Optional, ForceNew, List) The rule ids.
* `status` - (Optional, ForceNew) The status of the forwarding rule. Valid values: `Provisioning`, `Configuring`, `Available`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Rule names.
* `rules` - A list of Alb Rules. Each element contains the following attributes:
  * `id` - The ID of the Rule.
  * `rule_id` - The ID of the Rule.
  * `load_balancer_id` - The ID of the Application Load Balancer (ALB) instance to which the forwarding rule belongs.
  * `listener_id` - The ID of the listener to which the forwarding rule belongs.
  * `rule_name` - The name of the forwarding rule.
  * `priority` - The priority of the rule.
  * `rule_actions` - The actions of the forwarding rules.
    * `type` - The action type.
    * `order` - The order of the forwarding rule actions.
    * `fixed_response_config` - The configuration of the fixed response.
      * `content` - The fixed response. The response cannot exceed 1 KB in size and can contain only ASCII characters.
      * `content_type` - The format of the fixed response.
      * `http_code` - The HTTP status code of the response.
    * `forward_group_config` - The configurations of the destination server groups.
      * `server_group_tuples` - The destination server group to which requests are forwarded.
        * `server_group_id` - The ID of the destination server group to which requests are forwarded.
        * `weight` - The Weight of server group.
    * `insert_header_config` - The configuration of the inserted header field.
      * `key` - The name of the inserted header field.
      * `value` - The content of the inserted header field.
      * `value_type` - The value type of the inserted header field.
    * `redirect_config` - The configuration of the external redirect action.
      * `host` - The host name of the destination to which requests are directed.
      * `http_code` - The redirect method.
      * `path` - The path of the destination to which requests are directed.
      * `port` - The port of the destination to which requests are redirected.
      * `protocol` - The protocol of the requests to be redirected.
      * `query` - The query string of the request to be redirected.
    * `rewrite_config` - The redirect action within ALB.
      * `host` - The host name of the destination to which requests are redirected within ALB.
      * `path` - The path to which requests are to be redirected within ALB.
      * `query` - The query string of the request to be redirected within ALB.
    * `traffic_limit_config` - The Flow speed limit.
      * `qps` - The Number of requests per second.
    * `traffic_mirror_config` - The Traffic mirroring.
      * `target_type` - The Mirror target type.
      * `mirror_group_config` - The Traffic is mirrored to the server group.
        * `server_group_tuples` - The destination server group to which requests are forwarded.
          * `server_group_id` - The ID of the destination server group to which requests are forwarded.
  * `rule_conditions` - The conditions of the forwarding rule.
    * `type` - The type of the forwarding rule.
    * `cookie_config` - The configuration of the cookie.
      * `values` - The values of the cookie.
        * `key` - The key of the cookie.
        * `value` - The value of the cookie.
    * `header_config` - The configuration of the header field.
      * `key` - The key of the header field.
      * `values` - The value of the header field.
    * `host_config` - The configuration of the host.
      * `values` - The name of the host.
    * `method_config` - The configuration of the request method.
      * `values` - The request method.
    * `path_config` - The configuration of the path for the request to be forwarded.
      * `values` - The path of the request to be forwarded.
    * `query_string_config` - The configuration of the query string.
      * `values` - The query string.
        * `key` - The key of the query string.
        * `value` - The value of the query string.
    * `source_ip_config` - The Based on source IP traffic matching.
      * `values` - Add one or more IP addresses or IP address segments.
  * `status` - The status of the forwarding rule.
