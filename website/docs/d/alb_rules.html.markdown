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

* `ids` - (Optional, List)  A list of Rule IDs.
* `name_regex` - (Optional) A regex string to filter results by Rule name.
* `load_balancer_ids` - (Optional, List) The load balancer ids.
* `listener_ids` - (Optional, List) The listener ids.
* `rule_ids` - (Optional, List) The rule ids.
* `status` - (Optional) The status of the forwarding rule. Valid values: `Provisioning`, `Configuring`, `Available`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Rule names.
* `rules` - A list of Alb Rules.

---

### `rules`

The rules supports the following attributes:

* `id` - The ID of the Rule.
* `rule_id` - The ID of the Rule.
* `load_balancer_id` - The ID of the Application Load Balancer (ALB) instance to which the forwarding rule belongs.
* `listener_id` - The ID of the listener to which the forwarding rule belongs.
* `rule_name` - The name of the forwarding rule.
* `priority` - The priority of the rule.
* `rule_actions` - The actions of the forwarding rules.
* `rule_conditions` - The conditions of the forwarding rule.
* `status` - The status of the forwarding rule.

---

### `rules.rule_actions`

The rule_actions supports the following attributes:

* `type` - The action type.
* `order` - The order of the forwarding rule actions.
* `fixed_response_config` - The configuration of the fixed response.
* `forward_group_config` - The configurations of the destination server groups.
* `insert_header_config` - The configuration of the inserted header field.
* `redirect_config` - The configuration of the external redirect action.
* `rewrite_config` - The redirect action within ALB.
* `traffic_limit_config` - The Flow speed limit.
* `traffic_mirror_config` - The Traffic mirroring.

---

### `rules.rule_actions.fixed_response_config`

The fixed_response_config supports the following attributes:

* `content` - The fixed response. The response cannot exceed 1 KB in size and can contain only ASCII characters.
* `content_type` - The format of the fixed response.
* `http_code` - The HTTP status code of the response.

---

### `rules.rule_actions.forward_group_config`

The forward_group_config supports the following attributes:

* `server_group_tuples` - The destination server group to which requests are forwarded.

---

### `rules.rule_actions.forward_group_config.server_group_tuples`

The server_group_tuples supports the following attributes:

* `server_group_id` - The ID of the destination server group to which requests are forwarded.
* `weight` - The Weight of server group.

---

### `rules.rule_actions.insert_header_config`

The insert_header_config supports the following attributes:

* `key` - The name of the inserted header field.
* `value` - The content of the inserted header field.
* `value_type` - The value type of the inserted header field.

---

### `rules.rule_actions.redirect_config`

The redirect_config supports the following attributes:

* `host` - The host name of the destination to which requests are directed.
* `http_code` - The redirect method.
* `path` - The path of the destination to which requests are directed.
* `port` - The port of the destination to which requests are redirected.
* `protocol` - The protocol of the requests to be redirected.
* `query` - The query string of the request to be redirected.

---

### `rules.rule_actions.rewrite_config`

The rewrite_config supports the following attributes:

* `host` - The host name of the destination to which requests are redirected within ALB.
* `path` - The path to which requests are to be redirected within ALB.
* `query` - The query string of the request to be redirected within ALB.

---

### `rules.rule_actions.traffic_limit_config`

The traffic_limit_config supports the following attributes:

* `qps` - The Number of requests per second.

---

### `rules.rule_actions.traffic_mirror_config`

The traffic_mirror_config supports the following attributes:

* `target_type` - The Mirror target type.
* `mirror_group_config` - The Traffic is mirrored to the server group.

---

### `rules.rule_actions.traffic_mirror_config.mirror_group_config`

The mirror_group_config supports the following attributes:

* `server_group_tuples` - The destination server group to which requests are forwarded.

---

### `rules.rule_actions.traffic_mirror_config.mirror_group_config.server_group_tuples`

The server_group_tuples supports the following attributes:

* `server_group_id` - The ID of the destination server group to which requests are forwarded.

---

### `rules.rule_conditions`

The rule_conditions supports the following attributes:

* `type` - The type of the forwarding rule.
* `cookie_config` - The configuration of the cookie.
* `header_config` - The configuration of the header field.
* `host_config` - The configuration of the host.
* `method_config` - The configuration of the request method.
* `path_config` - The configuration of the path for the request to be forwarded.
* `query_string_config` - The configuration of the query string.
* `source_ip_config` - The Based on source IP traffic matching.

---

### `rules.rule_conditions.cookie_config`

The cookie_config supports the following attributes:

* `values` - The values of the cookie.

---

### `rules.rule_conditions.cookie_config.values`

The values supports the following attributes:

* `key` - The key of the cookie.
* `value` - The value of the cookie.

---

### `rules.rule_conditions.header_config`

The header_config supports the following attributes:

* `key` - The key of the header field.
* `values` - The value of the header field.

---

### `rules.rule_conditions.host_config`

The host_config supports the following attributes:

* `values` - The name of the host.

---

### `rules.rule_conditions.method_config`

The method_config supports the following attributes:

* `values` - The request method.

---

### `rules.rule_conditions.path_config`

The path_config supports the following attributes:

* `values` - The path of the request to be forwarded.

---

### `rules.rule_conditions.query_string_config`

The query_string_config supports the following attributes:

* `values` - The query string.

---

### `rules.rule_conditions.query_string_config.values`

The values supports the following attributes:

* `key` - The key of the query string.
* `value` - The value of the query string.

---

### `rules.rule_conditions.source_ip_config`

The source_ip_config supports the following attributes:

* `values` - Add one or more IP addresses or IP address segments.
